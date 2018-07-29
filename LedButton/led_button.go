package LedButton

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"log"

	"github.com/gobuffalo/packr"

	esd "github.com/dh1tw/go-elgato-stream-deck"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type LedButton struct {
	streamDeck *esd.StreamDeck
	ledColor   LEDColor
	text       string
	textColor  *image.Uniform
	id         int
	state      bool
}

type LEDColor int

const (
	LEDRed LEDColor = iota
	LEDGreen
	LEDYellow
	LEDOff
)

var ledOff image.Image
var ledGreen image.Image
var ledYellow image.Image
var ledRed image.Image
var font *truetype.Font

func init() {
	fontBox := packr.NewBox("./fonts")
	imgBox := packr.NewBox("./images")

	var err error

	// Load the font
	font, err = freetype.ParseFont(fontBox.Bytes("mplus-1m-regular.ttf"))
	if err != nil {
		log.Panic(err)
	}

	// Load the LED images
	ledOff, _, err = image.Decode(bytes.NewBuffer(imgBox.Bytes("led_off.png")))
	if err != nil {
		log.Panic(err)
	}
	ledGreen, _, err = image.Decode(bytes.NewBuffer(imgBox.Bytes("led_green_on.png")))
	if err != nil {
		log.Panic(err)
	}
	ledYellow, _, err = image.Decode(bytes.NewBuffer(imgBox.Bytes("led_yellow_on.png")))
	if err != nil {
		log.Panic(err)
	}
	ledRed, _, err = image.Decode(bytes.NewBuffer(imgBox.Bytes("led_red_on.png")))
	if err != nil {
		log.Panic(err)
	}
}

// NewLedButton is the constructor for a new Led Button. Functional
// arguments can be supplied to modify it's default characteristics
func NewLedButton(sd *esd.StreamDeck, id int, options ...func(*LedButton)) (*LedButton, error) {

	if sd == nil {
		return nil, fmt.Errorf("stream deck must not be nil")
	}

	btn := &LedButton{
		streamDeck: sd,
		id:         id,
		ledColor:   LEDGreen,
		text:       "",
		textColor:  image.White,
		state:      false,
	}

	for _, option := range options {
		option(btn)
	}

	return btn, nil
}

func (btn *LedButton) State() bool {
	return btn.state
}

func (btn *LedButton) SetState(state bool) error {
	btn.state = state
	return btn.Draw()
}

func (btn *LedButton) Draw() error {

	img := image.NewRGBA(image.Rect(0, 0, esd.ButtonSize, esd.ButtonSize))
	btn.addLED(btn.ledColor, img)
	if err := btn.addText(btn.text, img); err != nil {
		return err
	}
	return btn.streamDeck.FillImage(btn.id, img)
}

func (btn *LedButton) SetText(text string) error {
	btn.text = text
	return btn.Draw()
}

func (btn *LedButton) addLED(color LEDColor, img *image.RGBA) {

	if !btn.state {
		draw.Draw(img, img.Bounds(), ledOff, image.Point{0, 0}, draw.Src)
		return
	}

	switch color {
	case LEDRed:
		draw.Draw(img, img.Bounds(), ledRed, image.Point{0, 0}, draw.Src)
	case LEDGreen:
		draw.Draw(img, img.Bounds(), ledGreen, image.Point{0, 0}, draw.Src)
	case LEDYellow:
		draw.Draw(img, img.Bounds(), ledYellow, image.Point{0, 0}, draw.Src)
	}

}

type textParams struct {
	fontSize float64
	posX     int
	posY     int
}

var singleChar = textParams{
	fontSize: 26,
	posX:     30,
	posY:     27,
}

var oneLineTwoChars = textParams{
	fontSize: 26,
	posX:     23,
	posY:     27,
}

var oneLineThreeChars = textParams{
	fontSize: 26,
	posX:     17,
	posY:     27,
}

var oneLineFourChars = textParams{
	fontSize: 26,
	posX:     11,
	posY:     27,
}

var oneLineFiveChars = textParams{
	fontSize: 26,
	posX:     5,
	posY:     27,
}

var oneLine = textParams{
	fontSize: 26,
	posX:     0,
	posY:     27,
}

func (btn *LedButton) addText(text string, img *image.RGBA) error {

	var p textParams

	switch len(text) {
	case 1:
		p = singleChar
	case 2:
		p = oneLineTwoChars
	case 3:
		p = oneLineThreeChars
	case 4:
		p = oneLineFourChars
	case 5:
		p = oneLineFiveChars
	default:
		return fmt.Errorf("text line contains more than 5 characters")
	}

	// create Context
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(p.fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(btn.textColor)
	pt := freetype.Pt(p.posX, p.posY+int(c.PointToFixed(24)>>6))

	if _, err := c.DrawString(text, pt); err != nil {
		return err
	}

	return nil
}
