# go-elgato-stream-deck

[![Go Report Card](https://goreportcard.com/badge/github.com/dh1tw/go-elgato-stream-deck)](https://goreportcard.com/report/github.com/dh1tw/go-elgato-stream-deck)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://img.shields.io/badge/license-MIT-blue.svg)
[![GoDoc](https://godoc.org/github.com/dh1tw/go-elgato-stream-deck?status.svg)](https://godoc.org/github.com/dh1tw/go-elgato-stream-deck)

![go-elgato-stream-deck buttons](https://i.imgur.com/tEt3tPr.jpg "go-elgato-stream-deck Buttons")
![go-elgato-stream-deck slide show](https://i.imgur.com/gh1xXiU.jpg "go-elgato-stream-deck Slideshow")



**go-elgato-stream-deck** is a library for interfacing with the [Elgato Stream Deck](https://www.elgato.com/en/gaming/stream-deck)

This library is written in the programing language [Go](https://golang.org).

## Note
This project provides only an API for the Elgato Stream Deck. It is not an
executable program, nor does it replace the OEM's software.

## License

go-elgato-stream-deck is published under the permissive [MIT license](https://github.com/dh1tw/go-elgato-stream-deck/blob/master/LICENSE).

## Dependencies

There are a few go libraries which are needed at compile time. Go-elgato-stream-deck
does not have any runtime dependencies.

## Supported Operating Systems

In principal the library should work on Linux, MacOS and Windows (>=7).

go-elgato-stream-deck works well on SoC boards like the Raspberry / Orange / Banana Pis.

## How to Install

````bash
$ go get github.com/dh1tw/go-elgato-stream-deck
````

On Linux you might have to create an udev rule, to access the streamdeck.

````
sudo vim /etc/udev/rules.d/99-streamdeck.rules

SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", MODE="0664", GROUP="plugdev"
````

After saving the udev rule, unplug and plug the streamdeck again into the USB port.

Make sure that your streamdeck is correctly recognized
by executing:

````bash
$ go run examples/enumerate.go
Found 1 Elgato Stream Deck(s):
	SerialNumber:        AL12H1A07123
````

## Documentation

The auto generated documentation can be found at [godoc.org](https://godoc.org/github.com/dh1tw/go-elgato-stream-deck)

## Examples

There are a couple of examples located in the `examples` folder.

````bash
$ go run examples/led_buttons.go
$ go run examples/slideshow.go
...
````

## Credits

This project would not have been possible without the work of [Alex Van Camp](https://github.com/Lange). In particular his
[notes of the StreamDeck's protocol](https://github.com/Lange/node-elgato-stream-deck/blob/master/NOTES.md)
were very helpful.
Alex has provided a [reference implementation](https://github.com/Lange/node-elgato-stream-deck) in Javascript / Node.js.