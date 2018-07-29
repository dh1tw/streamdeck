package main

import (
	"bytes"
	"image"
	"log"
	"os"
	"os/signal"
	"time"

	esd "github.com/dh1tw/go-elgato-stream-deck"
	"github.com/gobuffalo/packr"
)

// This example creates a slideshow on the Stream Deck, across all buttons.
// Images of different formats (png, jpeg, gif) are loaded, resized to match
// the panel and if necessary cropped to the center.

func main() {

	sd, err := esd.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

	imgBox := packr.NewBox("images")

	dices, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("dices.png")))
	if err != nil {
		log.Panic(err)
	}

	dna, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("dna.gif")))
	if err != nil {
		log.Panic(err)
	}

	octocat, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("Octocat.jpg")))
	if err != nil {
		log.Panic(err)
	}

	// start drawing octocat
	if err := sd.FillPanel(octocat); err != nil {
		log.Panic(err)
	}

	images := []image.Image{dices, dna, octocat}

	// launch a ticker for the slideshow
	ticker := time.NewTicker(time.Second * 3)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	position := 0
	for {
		select {
		case <-ticker.C:
			if err := sd.FillPanel(images[position]); err != nil {
				log.Panic(err)
			}
			if position == len(images)-1 {
				position = 0
				break
			}
			if position < len(images)-1 {
				position++
			}
		case <-c:
			return
		}
	}
}
