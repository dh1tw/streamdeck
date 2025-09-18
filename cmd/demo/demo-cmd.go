package main

import (
	"image/color"
	"log"
	"time"

	"github.com/dh1tw/streamdeck"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	sd, err := streamdeck.NewStreamDeck()
	if err != nil {
		return err
	}
	defer sd.Close()

	err = sd.SetBrightness(100)
	if err != nil {
		return err
	}

	err = sd.FillColor(0, 255, 0, 0)
	if err != nil {
		return err
	}

	err = sd.WriteText(1, streamdeck.TextButton{
		Lines: []streamdeck.TextLine{
			{Text: "foo", PosX: 10, PosY: 10, FontSize: 20, FontColor: color.RGBA{255, 0, 0, 255}},
			{Text: "bar", PosX: 10, PosY: 40, FontSize: 20, FontColor: color.RGBA{0, 0, 255, 255}},
		},
		BgColor: color.RGBA{0, 255, 0, 255},
	})
	if err != nil {
		return err
	}

	sd.SetBtnEventCb(func(s streamdeck.State, e streamdeck.Event) {
		log.Printf("got event: %v state: %v", e, s)
	})

	log.Printf("sleeping")

	time.Sleep(10 * time.Second)

	err = sd.SetBrightness(50)
	if err != nil {
		return err
	}

	return nil
}
