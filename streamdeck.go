package streamdeck

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"sync"

	"github.com/disintegration/gift"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"github.com/bearsh/hid"

	"image/color"
	"image/draw"
	_ "image/gif"  // support gif
	_ "image/jpeg" // support jpeg
	_ "image/png"  // support png
)

var Debug = false

func debug(f string, args ...interface{}) {
	if !Debug {
		return
	}
	log.Printf(f, args...)
}

// VendorID is the USB VendorID assigned to Elgato (0x0fd9)
const VendorID = 4057

// BtnEvent is a callback which gets executed when the state of a button changes,
// so whenever it gets pressed or released.
type BtnEvent func(s State, e Event)

// StreamDeck is the object representing the Elgato Stream Deck.
type StreamDeck struct {
	lock       sync.Mutex
	device     *hid.Device
	btnEventCb BtnEvent
	Config     *Config

	waitGroup sync.WaitGroup
	cancel    context.CancelFunc
}

// TextButton holds the lines to be written to a button and the desired
// Background color.
type TextButton struct {
	Lines   []TextLine
	BgColor color.Color
}

// TextLine holds the content of one text line.
type TextLine struct {
	Text      string
	PosX      int
	PosY      int
	Font      *truetype.Font
	FontSize  float64
	FontColor color.Color
}

// NewStreamDeck is the constructor of the StreamDeck object. If several StreamDecks
// are connected to this PC, the Streamdeck can be selected by supplying
// the optional serial number of the Device. In the examples folder there is
// a small program which enumerates all available Stream Decks. If no serial number
// is supplied, the first StreamDeck found will be selected.
func NewStreamDeck(serial ...string) (*StreamDeck, error) {

	s := ""
	if len(serial) > 0 {
		s = serial[0]
	}

	return NewStreamDeckWithConfig(nil, s)
}

// NewStreamDeckWithConfig is the constructor for a custom config.
func NewStreamDeckWithConfig(c *Config, serial string) (*StreamDeck, error) {

	if c == nil {
		cc, found := FindConnectedConfig()
		if !found {
			return nil, fmt.Errorf("no streamdeck device found with any config")
		}
		c = &cc
	}

	devices := hid.Enumerate(VendorID, c.ProductID)

	if len(devices) == 0 {
		return nil, fmt.Errorf("no stream deck device found")
	}

	log.Printf("Found %d StreamDecks", len(devices))

	id := 0
	if serial != "" {
		found := false
		for i, d := range devices {
			if d.Serial == serial {
				id = i
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("no stream deck device found with serial number %s", serial)
		}
	}

	log.Printf("Connecting to StreamDeck: %v", devices[id])

	device, err := devices[id].Open()
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to StreamDeck: %v", devices[id])

	sd := &StreamDeck{
		device: device,
		Config: c,
	}

	sd.ClearAllBtns()

	cancelCtx, cancel := context.WithCancel(context.Background())
	sd.cancel = cancel

	sd.waitGroup.Add(1)
	go sd.read(cancelCtx)

	return sd, nil
}

// SetBtnEventCb sets the BtnEvent callback which get's executed whenever
// a Button event (pressed/released) occures.
func (sd *StreamDeck) SetBtnEventCb(ev BtnEvent) {
	sd.lock.Lock()
	defer sd.lock.Unlock()
	sd.btnEventCb = ev
}

// Read will listen in a for loop for incoming messages from the Stream Deck.
// It is typically executed in a dedicated go routine.
func (sd *StreamDeck) read(ctx context.Context) {
	defer sd.waitGroup.Done()
	myState := State{}
	for ctx.Err() == nil {
		data := make([]byte, 24)
		_, err := sd.device.Read(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		event, err := myState.Update(sd.Config, data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var cb BtnEvent

		sd.lock.Lock()
		cb = sd.btnEventCb
		sd.lock.Unlock()

		if cb != nil {
			go func() {
				cb(myState, event)
			}()
		}
	}
}

// Close the connection to the Elgato Stream Deck
func (sd *StreamDeck) Close() error {
	sd.cancel()
	err := sd.device.Close()
	sd.waitGroup.Wait()
	return err
}

// Serial returns the Serial number of this Elgato Stream Deck
func (sd *StreamDeck) Serial() string {
	return sd.device.Serial
}

// ClearBtn fills a particular key with the color black
func (sd *StreamDeck) ClearBtn(btnIndex int) error {

	if err := sd.checkValidKeyIndex(btnIndex); err != nil {
		return err
	}
	return sd.FillColor(btnIndex, 0, 0, 0)
}

// ClearAllBtns fills all keys with the color black
func (sd *StreamDeck) ClearAllBtns() error {
	for i := sd.Config.NumButtons() - 1; i >= 0; i-- {
		err := sd.ClearBtn(i)
		if err != nil {
			return err
		}
	}
	return nil
}

// FillColor fills the given button with a solid color.
func (sd *StreamDeck) FillColor(btnIndex, r, g, b int) error {
	if err := checkRGB(r); err != nil {
		return err
	}
	if err := checkRGB(g); err != nil {
		return err
	}
	if err := checkRGB(b); err != nil {
		return err
	}

	img := image.NewRGBA(image.Rect(0, 0, sd.Config.ButtonSize, sd.Config.ButtonSize))
	color := color.RGBA{uint8(r), uint8(g), uint8(b), 1}
	draw.Draw(img, img.Bounds(), image.NewUniform(color), image.Point{0, 0}, draw.Src)

	return sd.FillImage(btnIndex, img)
}

func (sd *StreamDeck) encodeImage(img image.Image) ([]byte, error) {
	if sd.Config.ImageRotate {
		b := img.Bounds()
		newImage := image.NewRGBA(image.Rect(0, 0, sd.Config.ButtonSize, sd.Config.ButtonSize))
		for x := 0; x < sd.Config.ButtonSize; x++ {
			for y := 0; y < sd.Config.ButtonSize; y++ {
				newImage.Set(x, y, img.At(b.Min.X+sd.Config.ButtonSize-x, b.Min.Y+sd.Config.ButtonSize-y))
			}
		}
		img = newImage
	}

	if sd.Config.ImageFormat == "bmp" {
		return encodeBMP(sd.Config, img)
	}

	if sd.Config.ImageFormat == "jpg" {
		buf := bytes.Buffer{}
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unknown image format [%s]", sd.Config.ImageFormat)

}

func encodeBMP(c *Config, img image.Image) ([]byte, error) {
	imgBuf := []byte{
		'\x42', '\x4D', '\xF6', '\x3C', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x36', '\x00', '\x00', '\x00', '\x28', '\x00',
		'\x00', '\x00', '\x48', '\x00', '\x00', '\x00', '\x48', '\x00',
		'\x00', '\x00', '\x01', '\x00', '\x18', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\xC0', '\x3C', '\x00', '\x00', '\xC4', '\x0E',
		'\x00', '\x00', '\xC4', '\x0E', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00'}

	for row := 0; row < c.ButtonSize; row++ {
		for line := c.ButtonSize - 1; line >= 0; line-- {
			r, g, b, _ := img.At(line, row).RGBA()
			imgBuf = append(imgBuf, byte(b), byte(g), byte(r))
		}
	}
	return imgBuf, nil
}

// FillImage fills the given key with an image. For best performance, provide
// the image in the size of 72x72 pixels. Otherwise it will be automatically
// resized.
func (sd *StreamDeck) FillImage(btnIndex int, img image.Image) error {
	if err := sd.checkValidKeyIndex(btnIndex); err != nil {
		return err
	}

	// if necessary, rescale the picture
	rect := img.Bounds()
	if rect.Dx() != sd.Config.ButtonSize {
		img = resize(img, sd.Config.ButtonSize, sd.Config.ButtonSize)
	}

	imgBuf, err := sd.encodeImage(img)
	if err != nil {
		return err
	}

	sd.lock.Lock()
	defer sd.lock.Unlock()

	if sd.Config.ImageFormat == "bmp" {
		splitPoint := 7803
		err := sd.sendOriginalSingleMsgInLock(btnIndex, 1, imgBuf[0:splitPoint])
		if err != nil {
			return err
		}

		return sd.sendOriginalSingleMsgInLock(btnIndex, 2, imgBuf[splitPoint:])
	}

	headerSize := 8
	bytesLeft := len(imgBuf)
	pos := 0
	pageNumber := uint16(0)

	for bytesLeft > 0 {
		imgToSend := min(bytesLeft, 1024-headerSize)

		buf := make([]byte, 1024)
		bytesLeft -= imgToSend

		buf[0] = 0x02
		buf[1] = 0x07
		buf[2] = byte(sd.Config.fixKey(btnIndex))
		if bytesLeft == 0 {
			buf[3] = 1
		} else {
			buf[3] = 0
		}
		binary.LittleEndian.PutUint16(buf[4:], uint16(imgToSend))
		debug("x %x %x", buf[4], buf[5])
		binary.LittleEndian.PutUint16(buf[6:], pageNumber)

		copy(buf[8:], imgBuf[pos:(pos+imgToSend)])

		debug("going to Write len(buf): %d imgToSend: %d bytesLeft: %d pageNumber: %d len(imgBuf): %d", len(buf), imgToSend, bytesLeft, pageNumber, len(imgBuf))

		n, err := sd.device.Write(buf)
		if err != nil {
			return err
		}
		if n != len(buf) {
			return fmt.Errorf("only wrote %d of %d", n, len(buf))
		}

		pageNumber++
		pos += imgToSend

	}

	return nil
}

func (sd *StreamDeck) sendOriginalSingleMsgInLock(btnIndex int, pageNumber uint16, data []byte) error {
	buf := make([]byte, 8191)
	buf[0] = 0x02
	buf[1] = 0x01
	binary.LittleEndian.PutUint16(buf[2:], pageNumber)
	if pageNumber == 2 {
		buf[4] = 1
	}
	buf[5] = byte(sd.Config.fixKey(btnIndex))
	copy(buf[16:], data)

	n, err := sd.device.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("only wrote %d of %d", n, len(buf))
	}
	return nil
}

// FillImageFromFile fills the given key with an image from a file.
func (sd *StreamDeck) FillImageFromFile(keyIndex int, path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	return sd.FillImage(keyIndex, img)
}

// FillPanel fills the whole panel witn an image. The image is scaled to fit
// and then center-cropped (if necessary). The native picture size is 360px x 216px.
func (sd *StreamDeck) FillPanel(img image.Image) error {

	// resize if the picture width is larger or smaller than panel
	rect := img.Bounds()
	if rect.Dx() != sd.Config.PanelWidth() {
		newWidthRatio := float32(rect.Dx()) / float32((sd.Config.PanelWidth()))
		img = resize(img, sd.Config.PanelWidth(), int(float32(rect.Dy())/newWidthRatio))
	}

	// if the Canvas is larger than sd.Config.PanelWidth() x sd.Config.PanelHeight() then we crop
	// the Center match sd.Config.PanelWidth() x sd.Config.PanelHeight()
	rect = img.Bounds()
	if rect.Dx() > sd.Config.PanelWidth() || rect.Dy() > sd.Config.PanelHeight() {
		img = cropCenter(img, sd.Config.PanelWidth(), sd.Config.PanelHeight())
	}

	counter := 0

	for row := 0; row < sd.Config.NumButtonRows; row++ {
		for col := 0; col < sd.Config.NumButtonColumns; col++ {
			rect := image.Rectangle{
				Min: image.Point{
					X: col*sd.Config.ButtonSize + col*sd.Config.Spacer,
					Y: row*sd.Config.ButtonSize + row*sd.Config.Spacer,
				},
				Max: image.Point{
					X: (1+col)*sd.Config.ButtonSize + col*sd.Config.Spacer,
					Y: (1+row)*sd.Config.ButtonSize + row*sd.Config.Spacer,
				},
			}
			sub := img.(*image.RGBA).SubImage(rect)
			sd.FillImage(counter, sub)
			counter++
		}
	}

	return nil
}

// FillPanelFromFile fills the entire panel with an image from a file.
func (sd *StreamDeck) FillPanelFromFile(path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	return sd.FillPanel(img)
}

// WriteText can write several lines of Text to a button. It is up to the
// user to ensure that the lines fit properly on the button.
func (sd *StreamDeck) WriteText(btnIndex int, textBtn TextButton) error {

	if err := sd.checkValidKeyIndex(btnIndex); err != nil {
		return err
	}

	img := image.NewRGBA(image.Rect(0, 0, sd.Config.ButtonSize, sd.Config.ButtonSize))
	bg := image.NewUniform(textBtn.BgColor)
	// fill button with Background color
	draw.Draw(img, img.Bounds(), bg, image.Point{0, 0}, draw.Src)

	return sd.WriteTextOnImage(btnIndex, img, textBtn.Lines)
}

// WriteText can write several lines of Text to a button. It is up to the
// user to ensure that the lines fit properly on the button.
func (sd *StreamDeck) WriteTextOnImage(btnIndex int, imgIn image.Image, lines []TextLine) error {
	img := resize(imgIn, sd.Config.ButtonSize, sd.Config.ButtonSize)

	for _, line := range lines {
		if line.Font == nil {
			line.Font = MonoRegular
		}
		fontColor := image.NewUniform(line.FontColor)
		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(line.Font)
		c.SetFontSize(line.FontSize)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(fontColor)
		pt := freetype.Pt(line.PosX, line.PosY+int(c.PointToFixed(24)>>6))

		if _, err := c.DrawString(line.Text, pt); err != nil {
			return err
		}
	}

	return sd.FillImage(btnIndex, img)
}

// checkValidKeyIndex checks that the keyIndex is valid
func (sd *StreamDeck) checkValidKeyIndex(keyIndex int) error {
	if keyIndex < 0 || keyIndex > sd.Config.NumButtons() {
		return fmt.Errorf("invalid key index")
	}
	return nil
}

// b 0 -> 100
func (sd *StreamDeck) SetBrightness(b uint16) error {

	buf := []byte{0x03, 0x08, 0xFF, 0xFF}
	binary.LittleEndian.PutUint16(buf[2:], b)

	_, err := sd.device.SendFeatureReport(buf)
	return err
}

// resize returns a resized copy of the supplied image with the given width and height.
func resize(img image.Image, width, height int) *image.RGBA {
	g := gift.New(
		gift.Resize(width, height, gift.LanczosResampling),
		gift.UnsharpMask(1, 1, 0),
	)
	res := image.NewRGBA(g.Bounds(image.Rect(0, 0, width, height)))
	g.Draw(res, img)
	return res
}

// crop center will extract a sub image with the given width and height
// from the center of the supplied picture.
func cropCenter(img image.Image, width, height int) image.Image {
	g := gift.New(
		gift.CropToSize(width, height, gift.CenterAnchor),
	)
	res := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(res, img)
	return res
}

// checkRGB returns an error in case of an invalid color (8 bit)
func checkRGB(value int) error {
	if value < 0 || value > 255 {
		return fmt.Errorf("invalid color range")
	}
	return nil
}
