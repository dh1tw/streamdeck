package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	sd "github.com/dh1tw/streamdeck"
	ledBtn "github.com/dh1tw/streamdeck/LedButton"
)

// This example shows how to use the 'streamdeck/LedButton‘. It will
// enumerate all the buttons on the panel with their ID and with a green LED
// which can be activated / deactivated with a button press.

func main() {

	sd, err := sd.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}

	defer sd.ClearAllBtns()

	btns := make(map[int]*ledBtn.LedButton)

	// Red Buttons
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("%03d", i)
		btn, err := ledBtn.NewLedButton(sd, i, ledBtn.Text(text), ledBtn.LedColor(ledBtn.LEDRed))
		if err != nil {
			fmt.Println(err)
		}
		btn.Draw()
		btns[i] = btn
	}

	// Yellow Buttons
	for i := 5; i < 10; i++ {
		text := fmt.Sprintf("%03d", i)
		btn, err := ledBtn.NewLedButton(sd, i, ledBtn.Text(text), ledBtn.LedColor(ledBtn.LEDYellow))
		if err != nil {
			fmt.Println(err)
		}
		btn.Draw()
		btns[i] = btn
	}

	// Green Buttons
	for i := 10; i < 15; i++ {
		text := fmt.Sprintf("%03d", i)
		btn, err := ledBtn.NewLedButton(sd, i, ledBtn.Text(text), ledBtn.LedColor(ledBtn.LEDGreen))
		if err != nil {
			fmt.Println(err)
		}
		btn.Draw()
		btns[i] = btn
	}

	btnChangedCb := func(btnIndex int, state sd.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		if state == sd.BtnPressed {
			btn := btns[btnIndex]
			btn.SetState(!btn.State())
		}
	}
	sd.SetBtnEventCb(btnChangedCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
