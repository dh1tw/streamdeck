//go:generate stringer -type=BtnState

package streamdeck

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"sync"
	"time"

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

// VendorID is the USB VendorID assigned to Elgato (0x0fd9)
const VendorID = 4057

// BtnEvent is a callback which gets executed when the state of a button changes,
// so whenever it gets pressed or released.
type BtnEvent func(btnIndex int, newBtnState BtnState)

// BtnState is a type representing the button state.
type BtnState int

const (
	// BtnPressed button pressed
	BtnPressed BtnState = iota
	// BtnReleased button released
	BtnReleased
	BtnLongPressed
)

// ReadErrorCb is a callback which gets executed in case reading from the
// Stream Deck fails (e.g. the cable get's disconnected).
type ReadErrorCb func(err error)

// StreamDeck is the object representing the Elgato Stream Deck.

type btn struct {
	state          BtnState
	longPressTimer *time.Timer
	stopTimer      chan struct{}
}

type StreamDeck struct {
	sync.Mutex
	device           *hid.Device
	btnEventCb       BtnEvent
	buttons          map[int]*btn
	longPressTimeout time.Duration
	config           Config
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

// Page contains the configuration of one particular page of buttons. Pages
// can be nested to an arbitrary depth.
type Page interface {
	Set(btnIndex int, state BtnState) Page
	Parent() Page
	Draw()
	SetActive(bool)
}

// NewStreamDeck is the constructor of the StreamDeck object. If several StreamDecks
// are connected to this PC, the Streamdeck can be selected by supplying
// the optional serial number of the Device. In the examples folder there is
// a small program which enumerates all available Stream Decks. If no serial number
// is supplied, the first StreamDeck found will be selected.
func NewStreamDeck(logger *log.Logger, c Config, serial ...string) (*StreamDeck, error) {

	if len(serial) > 1 {
		return nil, fmt.Errorf("only <= 1 serial numbers must be provided")
	}

	devices := hid.Enumerate(VendorID, c.ProductID)

	if len(devices) == 0 {
		return nil, fmt.Errorf("no stream deck device found")
	}

	log.Printf("Found %d StreamDecks", len(devices))

	id := 0
	if len(serial) == 1 {
		found := false
		for i, d := range devices {
			if d.Serial == serial[0] {
				id = i
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("no stream deck device found with serial number %s", serial[0])
		}
	}

	log.Printf("Connecting to StreamDeck: %v", devices[id])

	device, err := devices[id].Open()
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to StreamDeck: %v", devices[id])

	sd := &StreamDeck{
		device:           device,
		buttons:          make(map[int]*btn),
		longPressTimeout: time.Second,
		config:           c,
	}

	// initialize buttons to state BtnReleased
	for i := 0; i < c.NumButtons(); i++ {
		newBtn := &btn{
			state:          BtnReleased,
			longPressTimer: time.NewTimer(sd.longPressTimeout),
			stopTimer:      make(chan struct{}),
		}
		newBtn.longPressTimer.Stop()
		sd.buttons[i] = newBtn
	}

	sd.ClearAllBtns()

	go sd.read()

	return sd, nil
}

// SetBtnEventCb sets the BtnEvent callback which get's executed whenever
// a Button event (pressed/released) occures.
func (sd *StreamDeck) SetBtnEventCb(ev BtnEvent) {
	sd.Lock()
	defer sd.Unlock()
	sd.btnEventCb = ev
}

// Read will listen in a for loop for incoming messages from the Stream Deck.
// It is typically executed in a dedicated go routine.
func (sd *StreamDeck) read() {

	for {
		data := make([]byte, 16)
		_, err := sd.device.Read(data)
		if err != nil {
			fmt.Println(err)
		}

		data = data[1:] // strip off the first byte; usage unknown, but it is always '\x01'

		sd.Lock()
		// we have to iterate over all 15 buttons and check if the state
		// has changed. If it has changed, execute the callback.
		for i, b := range data {
			myBtn, exists := sd.buttons[i]
			if !exists {
				fmt.Println("unknown button ", i)
			}
			// if state didn't change then move on
			if myBtn.state == itob(int(b)) {
				continue
			}
			// button pressed
			if itob(int(b)) == BtnPressed {
				myBtn.longPressTimer.Reset(sd.longPressTimeout)
				myBtn.stopTimer = make(chan struct{})
				myBtn.state = BtnPressed
				if sd.btnEventCb != nil {
					go sd.btnEventCb(i, BtnPressed)
				}
				go func(index int) {
					select {
					case <-myBtn.longPressTimer.C:
						if sd.btnEventCb != nil {
							go sd.btnEventCb(index, BtnLongPressed)
						}
					case <-myBtn.stopTimer:
						myBtn.longPressTimer.Stop()
					}
				}(i)
				// continue
			} else if itob(int(b)) == BtnReleased {
				myBtn.state = BtnReleased
				close(myBtn.stopTimer)
				if sd.btnEventCb != nil {
					go sd.btnEventCb(i, BtnReleased)
				}
			} else {
				fmt.Println("** Should never arrive here!**")
			}
		}
		sd.Unlock()
	}
}

// Close the connection to the Elgato Stream Deck
func (sd *StreamDeck) Close() error {
	sd.Lock()
	defer sd.Unlock()
	return sd.device.Close()
}

// Serial returns the Serial number of this Elgato Stream Deck
func (sd *StreamDeck) Serial() string {
	sd.Lock()
	defer sd.Unlock()
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
func (sd *StreamDeck) ClearAllBtns() {
	for i := sd.config.NumButtons() - 1; i >= 0; i-- {
		sd.ClearBtn(i)
	}
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

	img := image.NewRGBA(image.Rect(0, 0, sd.config.ButtonSize, sd.config.ButtonSize))
	color := color.RGBA{uint8(r), uint8(g), uint8(b), 1}
	draw.Draw(img, img.Bounds(), image.NewUniform(color), image.Point{0, 0}, draw.Src)

	return sd.FillImage(btnIndex, img)
}

func (sd *StreamDeck) encodeImage(img image.Image) ([]byte, error) {
	if sd.config.ImageFormat == "bmp" {
		imgBuf := make([]byte, 0, sd.config.ButtonSize*sd.config.ButtonSize*3)

		for row := 0; row < sd.config.ButtonSize; row++ {
			for line := sd.config.ButtonSize - 1; line >= 0; line-- {
				r, g, b, _ := img.At(line, row).RGBA()
				imgBuf = append(imgBuf, byte(r), byte(b), byte(g))
			}
		}
		return imgBuf, nil
	}

	if sd.config.ImageFormat == "jpg" {
		buf := bytes.Buffer{}
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unknown image format [%s]", sd.config.ImageFormat)

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
	if rect.Dx() != sd.config.ButtonSize {
		img = resize(img, sd.config.ButtonSize, sd.config.ButtonSize)
	}

	imgBuf, err := sd.encodeImage(img)
	if err != nil {
		return err
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
		buf[2] = byte(btnIndex)
		if bytesLeft == 0 {
			buf[3] = 1
		} else {
			buf[3] = 0
		}
		binary.LittleEndian.PutUint16(buf[4:], uint16(imgToSend))
		log.Printf("x %x %x", buf[4], buf[5])
		binary.LittleEndian.PutUint16(buf[6:], pageNumber)

		copy(buf[8:], imgBuf[pos:(pos+imgToSend)])

		log.Printf("going to Write len(buf): %d imgToSend: %d bytesLeft: %d pageNumber: %d len(imgBuf): %d", len(buf), imgToSend, bytesLeft, pageNumber, len(imgBuf))

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
	if rect.Dx() != sd.config.PanelWidth() {
		newWidthRatio := float32(rect.Dx()) / float32((sd.config.PanelWidth()))
		img = resize(img, sd.config.PanelWidth(), int(float32(rect.Dy())/newWidthRatio))
	}

	// if the Canvas is larger than sd.config.PanelWidth() x sd.config.PanelHeight() then we crop
	// the Center match sd.config.PanelWidth() x sd.config.PanelHeight()
	rect = img.Bounds()
	if rect.Dx() > sd.config.PanelWidth() || rect.Dy() > sd.config.PanelHeight() {
		img = cropCenter(img, sd.config.PanelWidth(), sd.config.PanelHeight())
	}

	counter := 0

	for row := 0; row < sd.config.NumButtonRows; row++ {
		for col := 0; col < sd.config.NumButtonColumns; col++ {
			rect := image.Rectangle{
				Min: image.Point{
					sd.config.PanelWidth() - sd.config.ButtonSize - col*sd.config.ButtonSize - col*sd.config.Spacer,
					row*sd.config.ButtonSize + row*sd.config.Spacer,
				},
				Max: image.Point{
					sd.config.PanelWidth() - 1 - col*sd.config.ButtonSize - col*sd.config.Spacer,
					sd.config.ButtonSize - 1 + row*sd.config.ButtonSize + row*sd.config.Spacer,
				},
			}
			sd.FillImage(counter, img.(*image.RGBA).SubImage(rect))
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

	img := image.NewRGBA(image.Rect(0, 0, sd.config.ButtonSize, sd.config.ButtonSize))
	bg := image.NewUniform(textBtn.BgColor)
	// fill button with Background color
	draw.Draw(img, img.Bounds(), bg, image.Point{0, 0}, draw.Src)

	for _, line := range textBtn.Lines {
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

	sd.FillImage(btnIndex, img)
	return nil
}

// checkValidKeyIndex checks that the keyIndex is valid
func (sd *StreamDeck) checkValidKeyIndex(keyIndex int) error {
	if keyIndex < 0 || keyIndex > sd.config.NumButtons() {
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
func resize(img image.Image, width, height int) image.Image {
	g := gift.New(
		gift.Resize(width, height, gift.LanczosResampling),
		gift.UnsharpMask(1, 1, 0),
	)
	res := image.NewRGBA(g.Bounds(img.Bounds()))
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

// int to ButtonState
func itob(i int) BtnState {
	if i == 0 {
		return BtnReleased
	}
	return BtnPressed
}
