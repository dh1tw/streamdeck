package Label

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	sd "github.com/dh1tw/streamdeck"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// Label is a basic Element for the StreamDeck.
type Label struct {
	streamDeck *sd.StreamDeck
	text       string
	id         int
	textColor  color.Color
	bgColor    color.Color
}

var font *truetype.Font

// in order to avoid the repetitive loading of the font, we load it once
// during initalization into memory
func init() {
	fontBox := packr.NewBox("./fonts")

	var err error

	// Load the font
	font, err = freetype.ParseFont(fontBox.Bytes("mplus-1m-regular.ttf"))
	if err != nil {
		log.Panic(err)
	}
}

// NewLabel is the constructor method for a Label.
func NewLabel(sd *sd.StreamDeck, btnIndex int, options ...func(*Label)) (*Label, error) {

	l := &Label{
		streamDeck: sd,
		id:         btnIndex,
		text:       "",
		textColor:  image.White,
		bgColor:    image.Black,
	}

	for _, option := range options {
		option(l)
	}

	return l, nil
}

// Draw renders the Label on the designated Button.
func (l *Label) Draw() error {
	img := image.NewRGBA(image.Rect(0, 0, sd.ButtonSize, sd.ButtonSize))
	l.addBgColor(l.bgColor, img)
	if err := l.addText(l.text, img); err != nil {
		return err
	}
	return l.streamDeck.FillImage(l.id, img)
}

// SetText sets the text of the Label.
func (l *Label) SetText(text string) error {
	l.text = text
	return l.Draw()
}

// SetBgColor sets the background color of the Label.
func (l *Label) SetBgColor(color *image.Uniform) error {
	l.bgColor = color
	return l.Draw()
}

func (l *Label) addBgColor(col color.Color, img *image.RGBA) {
	draw.Draw(img, img.Bounds(), image.NewUniform(col), image.ZP, draw.Src)
}

type textParams struct {
	fontSize float64
	posX     int
	posY     int
}

var singleChar = textParams{
	fontSize: 26,
	posX:     30,
	posY:     20,
}

var oneLineTwoChars = textParams{
	fontSize: 26,
	posX:     23,
	posY:     20,
}

var oneLineThreeChars = textParams{
	fontSize: 26,
	posX:     17,
	posY:     20,
}

var oneLineFourChars = textParams{
	fontSize: 26,
	posX:     11,
	posY:     20,
}

var oneLineFiveChars = textParams{
	fontSize: 26,
	posX:     5,
	posY:     20,
}

var oneLine = textParams{
	fontSize: 26,
	posX:     0,
	posY:     20,
}

func (l *Label) addText(text string, img *image.RGBA) error {

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
	c.SetSrc(image.NewUniform(l.textColor))
	pt := freetype.Pt(p.posX, p.posY+int(c.PointToFixed(24)>>6))

	if _, err := c.DrawString(text, pt); err != nil {
		return err
	}

	return nil
}
