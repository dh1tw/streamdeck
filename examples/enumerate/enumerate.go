package main

import (
	"html/template"
	"log"
	"os"

	sd "github.com/dh1tw/streamdeck"
	"github.com/karalabe/hid"
)

// This program enumerates all instances of Elgato Stream Deck connected
// to this computer.

var tmpl = template.Must(template.New("").Parse(
	`Found {{. | len}} Elgato Stream Deck(s): {{range .}}
	SerialNumber:        {{.Serial}}
	{{end}}`,
))

// Enumerate shows all connected Elgato Stream Decks
func main() {

	devices := hid.Enumerate(sd.VendorID, sd.ProductID)
	if err := tmpl.Execute(os.Stdout, devices); err != nil {
		log.Fatal(err)
	}
}
