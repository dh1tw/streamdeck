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
	// pick the first Stream Deck we find
	sd, err := streamdeck.NewStreamDeck()
	if err != nil {
		return err
	}
	defer sd.Close()

	err = sd.SetBrightness(100)
	if err != nil {
		return err
	}

	// First button will be plain red
	err = sd.FillColor(0, 255, 0, 0)
	if err != nil {
		return err
	}

	// Second button is green with red and blue text
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

	// If this is a Stream Deck Plus, demonstrate the LCD panel
	if sd.Config.HasLCD() {
		log.Printf("Stream Deck Plus detected - demonstrating LCD panel")

		// Fill each LCD segment with a different color and label
		segments := []struct {
			color color.RGBA
			label string
		}{
			{color.RGBA{255, 0, 0, 255}, "Red"},
			{color.RGBA{0, 255, 0, 255}, "Green"},
			{color.RGBA{0, 0, 255, 255}, "Blue"},
			{color.RGBA{255, 255, 0, 255}, "Yellow"},
		}

		for i, seg := range segments {
			err = sd.WriteTextToLCDSegment(i, seg.color, []streamdeck.LCDTextLine{
				{Text: seg.label, PosX: 50, PosY: 40, FontSize: 24, FontColor: color.RGBA{255, 255, 255, 255}},
			})
			if err != nil {
				return err
			}
		}
	}

	// capture Streamdeck events
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
