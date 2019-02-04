package main

import (
	"image/color"
	"log"
	"os"
	"os/signal"

	sd "github.com/dh1tw/streamdeck"
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

	lineLabel := sd.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 255, 0, 0}, // Yellow
		FontSize:  22,
		PosX:      10,
		PosY:      5,
		Text:      "STATE",
	}

	linePressed := sd.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 255, 255, 0}, // White
		FontSize:  14,
		PosX:      12,
		PosY:      30,
		Text:      "PRESSED",
	}

	lineReleased := sd.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 0, 0, 0}, // Red
		FontSize:  14,
		PosX:      9,
		PosY:      30,
		Text:      "RELEASED",
	}

	pressedText := sd.TextButton{
		BgColor: color.RGBA{0, 0, 0, 0},
		Lines:   []sd.TextLine{lineLabel, linePressed},
	}

	releasedText := sd.TextButton{
		BgColor: color.RGBA{0, 0, 0, 0},
		Lines:   []sd.TextLine{lineLabel, lineReleased},
	}

	sdeck, err := sd.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sdeck.ClearAllBtns()

	for i := 0; i < 15; i++ {
		sdeck.WriteText(i, releasedText)
	}

	btnEvtCb := func(btnIndex int, state sd.BtnState) {
		if state == sd.BtnPressed {
			sdeck.WriteText(btnIndex, pressedText)
		} else {
			sdeck.WriteText(btnIndex, releasedText)
		}
	}

	sdeck.SetBtnEventCb(btnEvtCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
