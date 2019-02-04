package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	sd "github.com/dh1tw/streamdeck"
	label "github.com/dh1tw/streamdeck/Label"
)

// This example will instanciate 15 labels on the streamdeck. Each Label
// is setup as a counter which will increment every 50ms. If a button is
// pressed it will be colored blue until it is released.

func main() {

	sd, err := sd.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

	labels := make(map[int]*label.Label)

	for i := 0; i < 15; i++ {
		label, err := label.NewLabel(sd, i, label.Text(strconv.Itoa(i)))
		if err != nil {
			fmt.Println(err)
		}
		label.Draw()
		labels[i] = label
	}

	handleBtnEvents := func(btnIndex int, state sd.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		if state == sd.BtnPressed {
			col := color.RGBA{0, 0, 153, 0}
			labels[btnIndex].SetBgColor(image.NewUniform(col))
		} else { // must be BtnReleased
			col := color.RGBA{0, 0, 0, 255}
			labels[btnIndex].SetBgColor(image.NewUniform(col))
		}
	}

	sd.SetBtnEventCb(handleBtnEvents)

	ticker := time.NewTicker(time.Millisecond * 50)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	counter := 0

	for {
		select {
		case <-c:
			return
		case <-ticker.C:
			for i := 0; i < 15; i++ {
				if err := labels[i].SetText(fmt.Sprintf("%03d", counter)); err != nil {
					log.Println(err)
				}
			}
			counter++
		}
	}
}
