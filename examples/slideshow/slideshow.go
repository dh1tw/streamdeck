package main

import (
	"bytes"
	"image"
	"log"
	"os"
	"os/signal"
	"time"

	sdeck "github.com/dh1tw/streamdeck"
	"github.com/gobuffalo/packr/v2"
)

// This example creates a slideshow on the Stream Deck, across all buttons.
// Images of different formats (png, jpeg, gif) are loaded, resized to match
// the panel and if necessary cropped to the center.

func main() {

	sd, err := sdeck.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

	imgBox := packr.New("slideshow-images", "../assets/images")

	_dices, err := imgBox.Find("dices.png")
	if err != nil {
		log.Fatal(err)
	}
	dices, _, err := image.Decode(bytes.NewBuffer(_dices))
	if err != nil {
		log.Panic(err)
	}

	_dna, err := imgBox.Find("dna.gif")
	if err != nil {
		log.Fatal(err)
	}
	dna, _, err := image.Decode(bytes.NewBuffer(_dna))
	if err != nil {
		log.Panic(err)
	}

	_octocat, err := imgBox.Find("octocat.jpg")
	if err != nil {
		log.Fatal(err)
	}
	octocat, _, err := image.Decode(bytes.NewBuffer(_octocat))
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
