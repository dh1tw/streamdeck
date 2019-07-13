//go:generate stringer -type=BtnState

package StreamDeck

import (
	"fmt"
	"image"
	"os"
	"sync"

	"github.com/disintegration/gift"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"github.com/dh1tw/hid"

	"image/color"
	"image/draw"
	_ "image/gif"  // support gif
	_ "image/jpeg" // support jpeg
	_ "image/png"  // support png
)

// VendorID is the USB VendorID assigned to Elgato (0x0fd9)
const VendorID = 4057

// ProductID is the USB ProductID assigned to Elgato's Stream Deck (0x0060)
const ProductID = 96

// NumButtons is the total amount of Buttons located on the Stream Deck.
const NumButtons = 15

// numFirstMsgPixels is the amount of pixels which have to be sent to the
// Stream Deck in the first message.
const numFirstMsgPixels = 2583

// numSecondMsgPixels is the amount of pixels which have to be send to the
// Stream Deck in the second message.
const numSecondMsgPixels = 2601

// ButtonSize is the size of a button (in pixel).
const ButtonSize = 72

// NumButtonColumns is the number of columns on the Stream Deck.
const NumButtonColumns = 5

// NumButtonRows is the number of button rows on the Stream Deck.
const NumButtonRows = 3

// Spacer is the spacing distance (in pixel) of two buttons on the Stream Deck.
const Spacer = 19

// PanelWidth is the total screen width of the Stream Deck (including spacers).
const PanelWidth = NumButtonColumns*ButtonSize + Spacer*(NumButtonColumns-1)

// PanelHeight is the total screen height of the stream deck (including spacers).
const PanelHeight = NumButtonRows*ButtonSize + Spacer*(NumButtonRows-1)

// BtnEvent is a callback which gets executed when the state of a button changes,
// so whenever it get's pressed or released.
type BtnEvent func(btnIndex int, newBtnState BtnState)

// BtnState is a type representing the button state.
type BtnState int

const (
	// BtnPressed button pressed
	BtnPressed BtnState = iota
	// BtnReleased button released
	BtnReleased
)

// ReadErrorCb is a callback which gets executed in case reading from the
// Stream Deck fails (e.g. the cable get's disconnected).
type ReadErrorCb func(err error)

// StreamDeck is the object representing the Elgato Stream Deck.
type StreamDeck struct {
	sync.Mutex
	device     *hid.Device
	btnEventCb BtnEvent
	btnState   []BtnState
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
func NewStreamDeck(serial ...string) (*StreamDeck, error) {

	if len(serial) > 1 {
		return nil, fmt.Errorf("only <= 1 serial numbers must be provided")
	}

	devices := hid.Enumerate(VendorID, ProductID)

	if len(devices) == 0 {
		return nil, fmt.Errorf("no stream deck device found")
	}

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

	device, err := devices[id].Open()
	if err != nil {
		return nil, err
	}

	sd := &StreamDeck{
		device:   device,
		btnState: make([]BtnState, NumButtons),
	}

	// initialize buttons to state BtnReleased
	for i := range sd.btnState {
		sd.btnState[i] = BtnReleased
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
// It is typcially executed in a dedicated go routine.
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
			if sd.btnState[i] != itob(int(b)) {
				sd.btnState[i] = itob(int(b))
				if sd.btnEventCb != nil {
					btnState := sd.btnState[i]
					go sd.btnEventCb(i, btnState)
				}
			}
		}
		sd.Unlock()
	}
}

// Close the connection to the Elgato Stream Deck
func (sd *StreamDeck) Close() error {
	sd.Lock()
	sd.Unlock()
	return sd.Close()
}

// ClearBtn fills a particular key with the color black
func (sd *StreamDeck) ClearBtn(btnIndex int) error {

	if err := checkValidKeyIndex(btnIndex); err != nil {
		return err
	}
	return sd.FillColor(btnIndex, 0, 0, 0)
}

// ClearAllBtns fills all keys with the color black
func (sd *StreamDeck) ClearAllBtns() {
	for i := 14; i >= 0; i-- {
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

	img := image.NewRGBA(image.Rect(0, 0, ButtonSize, ButtonSize))
	color := color.RGBA{uint8(r), uint8(g), uint8(b), 0}
	draw.Draw(img, img.Bounds(), image.NewUniform(color), image.Point{0, 0}, draw.Src)

	return sd.FillImage(btnIndex, img)
}

// FillImage fills the given key with an image. For best performance, provide
// the image in the size of 72x72 pixels. Otherwise it will be automatically
// resized.
func (sd *StreamDeck) FillImage(btnIndex int, img image.Image) error {
	if err := checkValidKeyIndex(btnIndex); err != nil {
		return err
	}

	// if necessary, rescale the picture
	rect := img.Bounds()
	if rect.Dx() != ButtonSize {
		img = resize(img, ButtonSize, ButtonSize)
	}

	imgBuf := make([]byte, 0, ButtonSize*ButtonSize*3)

	for row := 0; row < ButtonSize; row++ {
		for line := ButtonSize - 1; line >= 0; line-- {
			r, g, b, _ := img.At(line, row).RGBA()
			imgBuf = append(imgBuf, byte(r), byte(b), byte(g))
		}
	}

	page1 := imgBuf[0 : numFirstMsgPixels*3]
	page2 := imgBuf[numFirstMsgPixels*3:]

	sd.Lock()
	sd.writeMsg1(btnIndex, page1)
	sd.writeMsg2(btnIndex, page2)
	sd.Unlock()

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
	if rect.Dx() != PanelWidth {
		newWidthRatio := float32(rect.Dx()) / float32((PanelWidth))
		img = resize(img, PanelWidth, int(float32(rect.Dy())/newWidthRatio))
	}

	// if the Canvas is larger than PanelWidth x PanelHeight then we crop
	// the Center match PanelWidth x PanelHeight
	rect = img.Bounds()
	if rect.Dx() > PanelWidth || rect.Dy() > PanelHeight {
		img = cropCenter(img, PanelWidth, PanelHeight)
	}

	counter := 0

	for row := 0; row < NumButtonRows; row++ {
		for col := 0; col < NumButtonColumns; col++ {
			rect := image.Rectangle{
				Min: image.Point{
					PanelWidth - ButtonSize - col*ButtonSize - col*Spacer,
					row*ButtonSize + row*Spacer,
				},
				Max: image.Point{
					PanelWidth - 1 - col*ButtonSize - col*Spacer,
					ButtonSize - 1 + row*ButtonSize + row*Spacer,
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

	if err := checkValidKeyIndex(btnIndex); err != nil {
		return err
	}

	img := image.NewRGBA(image.Rect(0, 0, ButtonSize, ButtonSize))
	bg := image.NewUniform(textBtn.BgColor)
	// fill button with Background color
	draw.Draw(img, img.Bounds(), bg, image.Point{0, 0}, draw.Src)

	for _, line := range textBtn.Lines {
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

// writeMsg1 writes the first part of a button's content to the stream deck.
func (sd *StreamDeck) writeMsg1(btnIndex int, c []byte) {
	prefix := []byte{'\x02', '\x01', '\x01', '\x00', '\x00', byte(btnIndex + 1), '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x42', '\x4D', '\xF6', '\x3C', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x36', '\x00', '\x00', '\x00', '\x28', '\x00', '\x00', '\x00', '\x48', '\x00',
		'\x00', '\x00', '\x48', '\x00', '\x00', '\x00', '\x01', '\x00', '\x18', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\xC0', '\x3C', '\x00', '\x00', '\xC4', '\x0E', '\x00', '\x00', '\xC4', '\x0E', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00'}
	merged := append(prefix, c...)
	sd.device.Write(merged)
}

// writeMsg2 writes the second part of a button's content to the stream deck.
func (sd *StreamDeck) writeMsg2(btnIndex int, c []byte) {
	prefix := []byte{'\x02', '\x01', '\x02', '\x00', '\x01', byte(btnIndex + 1), '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00'}
	merged := append(prefix, c...)
	sd.device.Write(merged)
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

// checkValidKeyIndex checks that the keyIndex is valid
func checkValidKeyIndex(keyIndex int) error {
	if keyIndex < 0 || keyIndex > 15 {
		return fmt.Errorf("invalid key index")
	}
	return nil
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
