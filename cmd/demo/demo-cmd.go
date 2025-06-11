package main

import (
	"log"
	"time"

	"github.com/dh1tw/streamdeck"
)

func main() {
	logger := log.Default()
	err := realMain(logger)
	if err != nil {
		log.Fatal(err)
	}
}

func realMain(logger *log.Logger) error {
	sd, err := streamdeck.NewStreamDeck(logger, streamdeck.Plus)
	if err != nil {
		return err
	}
	defer sd.Close()

	err = sd.SetBrightness(100)
	if err != nil {
		return err
	}

	logger.Printf("going to FillColor")
	err = sd.FillColor(0, 255, 0, 0)
	if err != nil {
		return err
	}

	logger.Printf("sleeping")

	time.Sleep(5 * time.Second)

	err = sd.SetBrightness(50)
	if err != nil {
		return err
	}

	return nil
}
