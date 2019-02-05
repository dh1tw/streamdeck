package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"os/signal"

	sdeck "github.com/dh1tw/streamdeck"
	"github.com/gobuffalo/packr"
)

// This example loads icons and places them on buttons in the first row
// of the Stream Deck. The lightbulb icon on button 0 can be toggled.

func main() {

	sd, err := sdeck.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

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

	if err := sd.FillImage(4, warning); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(3, doctor); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(2, tux); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(1, user); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(0, lightbulbOff); err != nil {
		log.Panic(err)
	}

	lightbulb := false

	onPressedCb := func(btnIndex int, state sdeck.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		if btnIndex == 0 && state == sdeck.BtnPressed {
			if lightbulb {
				if err := sd.FillImage(0, lightbulbOff); err != nil {
					log.Panic(err)
				}
				lightbulb = false
			} else {
				if err := sd.FillImage(0, lightbulbOn); err != nil {
					log.Panic(err)
				}
				lightbulb = true
			}
		}
	}

	sd.SetBtnEventCb(onPressedCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
