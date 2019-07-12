# streamdeck

[![Go Report Card](https://goreportcard.com/badge/github.com/dh1tw/streamdeck)](https://goreportcard.com/report/github.com/dh1tw/streamdeck)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://img.shields.io/badge/license-MIT-blue.svg)
[![GoDoc](https://godoc.org/github.com/dh1tw/streamdeck?status.svg)](https://godoc.org/github.com/dh1tw/streamdeck)

![streamdeck buttons](https://i.imgur.com/tEt3tPr.jpg "streamdeck Buttons")
![streamdeck slide show](https://i.imgur.com/gh1xXiU.jpg "streamdeck Slideshow")



**streamdeck** is a library for interfacing with the [Elgato Stream Deck](https://www.elgato.com/en/gaming/stream-deck)

This library is written in the programing language [Go](https://golang.org).

## Note
This project is a golang API for the Elgato/Corsair StreamDeck. This library
unleashes the power of the StreamDeck. It allows you to completely customize
the content of the device, without the need of the OEMs software.

## License

streamdeck is published under the permissive [MIT license](https://github.com/dh1tw/streamdeck/blob/master/LICENSE).

## Dependencies

There are a few go libraries which are needed at compile time. streamdeck
does not have any runtime dependencies.

However compiling this library requires a c compiler since the underlying [HID library](github.com/karalabe/hid) requires cgo for enumerating the
HID devices.

## Supported Operating Systems

In principal the library should work on Linux, MacOS and Windows (>=7).

streamdeck works well on SoC boards like the Raspberry / Orange / Banana Pis.

## How to Install

````bash
$ go get github.com/dh1tw/streamdeck
````

By default the images and fonts are not included in the binary. If you
would like to do so, you can execute:

````
$ go get github.com/gobuffalo/packr/v2/packr2
$ cd $GOPATH/src/github.com/dh1tw/streamdeck
$ packr2
````

[Packr2](github.com/gobuffalo/packr/v2/packr2) will compile all the static
assets into go file while will then be included when you execute `go build`.

On Linux you might have to create an udev rule, to access the streamdeck.

````
sudo vim /etc/udev/rules.d/99-streamdeck.rules

SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", MODE="0664", GROUP="plugdev"
````

After saving the udev rule, unplug and plug the streamdeck again into the USB port.
For the rule above, your user must be a member of the `plugdev` group.

Make sure that your streamdeck is correctly recognized
by executing:

````bash
$ go run examples/enumerate.go
Found 1 Elgato Stream Deck(s):
	SerialNumber:        AL12H1A07123
````

## Documentation

The auto generated documentation can be found at [godoc.org](https://godoc.org/github.com/dh1tw/streamdeck)

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