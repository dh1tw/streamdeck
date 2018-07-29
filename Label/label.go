package Label

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	esd "github.com/dh1tw/go-elgato-stream-deck"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type Label struct {
	streamDeck *esd.StreamDeck
	text       string
	id         int
	textColor  color.Color
	bgColor    color.Color
}

var font *truetype.Font

func init() {
	fontBox := packr.NewBox("./fonts")

	var err error

	// Load the font
	font, err = freetype.ParseFont(fontBox.Bytes("mplus-1m-regular.ttf"))
	if err != nil {
		log.Panic(err)
	}
}

func NewLabel(sd *esd.StreamDeck, keyIndex int, options ...func(*Label)) (*Label, error) {

	l := &Label{
		streamDeck: sd,
		id:         keyIndex,
		text:       "",
		textColor:  image.White,
		bgColor:    image.Black,
	}

	for _, option := range options {
		option(l)
	}

	return l, nil
}

func (l *Label) Draw() error {
	img := image.NewRGBA(image.Rect(0, 0, esd.ButtonSize, esd.ButtonSize))
	l.addBgColor(l.bgColor, img)
	if err := l.addText(l.text, img); err != nil {
		return err
	}
	return l.streamDeck.FillImage(l.id, img)
}

func (l *Label) SetText(text string) error {
	l.text = text
	return l.Draw()
}

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
