package main

import (
	"image/color"
	"log"
	"os"
	"os/signal"

	esd "github.com/dh1tw/go-elgato-stream-deck"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var monoFont *truetype.Font

func main() {

	fontBox := packr.NewBox("fonts")

	var err error

	// Load the font
	monoFont, err = freetype.ParseFont(fontBox.Bytes("mplus-1m-regular.ttf"))
	if err != nil {
		log.Panic(err)
	}

	lineLabel := esd.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 255, 0, 0}, // Yellow
		FontSize:  22,
		PosX:      10,
		PosY:      5,
		Text:      "STATE",
	}

	linePressed := esd.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 255, 255, 0}, // White
		FontSize:  14,
		PosX:      12,
		PosY:      30,
		Text:      "PRESSED",
	}

	lineReleased := esd.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 0, 0, 0}, // Red
		FontSize:  14,
		PosX:      9,
		PosY:      30,
		Text:      "RELEASED",
	}

	pressedText := esd.TextButton{
		BgColor: color.RGBA{0, 0, 0, 0},
		Lines:   []esd.TextLine{lineLabel, linePressed},
	}

	releasedText := esd.TextButton{
		BgColor: color.RGBA{0, 0, 0, 0},
		Lines:   []esd.TextLine{lineLabel, lineReleased},
	}

	sd, err := esd.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < 15; i++ {
		sd.WriteText(i, releasedText)
	}

	btnEvtCb := func(btnIndex int, state esd.BtnState) {
		if state == esd.BtnPressed {
			sd.WriteText(btnIndex, pressedText)
		} else {
			sd.WriteText(btnIndex, releasedText)
		}
	}

	sd.SetBtnEventCb(btnEvtCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
