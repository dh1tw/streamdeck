package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"os/signal"

	sd "github.com/dh1tw/streamdeck"
	"github.com/gobuffalo/packr"
)

// This example loads icons and places them on buttons in the first row
// of the Stream Deck. The lightbulb icon on button 0 can be toggled.

func main() {

	sdeck, err := sd.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sdeck.ClearAllBtns()

	imgBox := packr.NewBox("images")

	user, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("user.png")))
	if err != nil {
		log.Panic(err)
	}

	tux, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("tux.png")))
	if err != nil {
		log.Panic(err)
	}

	warning, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("warning.png")))
	if err != nil {
		log.Panic(err)
	}

	doctor, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("doctor.png")))
	if err != nil {
		log.Panic(err)
	}
	lightbulbOn, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("lightbulb_on.png")))
	if err != nil {
		log.Panic(err)
	}

	lightbulbOff, _, err := image.Decode(bytes.NewBuffer(imgBox.Bytes("lightbulb_off.png")))
	if err != nil {
		log.Panic(err)
	}

	if err := sdeck.FillImage(4, warning); err != nil {
		log.Panic(err)
	}
	if err := sdeck.FillImage(3, doctor); err != nil {
		log.Panic(err)
	}
	if err := sdeck.FillImage(2, tux); err != nil {
		log.Panic(err)
	}
	if err := sdeck.FillImage(1, user); err != nil {
		log.Panic(err)
	}
	if err := sdeck.FillImage(0, lightbulbOff); err != nil {
		log.Panic(err)
	}

	lightbulb := false

	onPressedCb := func(btnIndex int, state sd.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		if btnIndex == 0 && state == sd.BtnPressed {
			if lightbulb {
				if err := sdeck.FillImage(0, lightbulbOff); err != nil {
					log.Panic(err)
				}
				lightbulb = false
			} else {
				if err := sdeck.FillImage(0, lightbulbOn); err != nil {
					log.Panic(err)
				}
				lightbulb = true
			}
		}
	}

	sdeck.SetBtnEventCb(onPressedCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
