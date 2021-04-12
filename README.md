# streamdeck

[![Go Report Card](https://goreportcard.com/badge/github.com/dh1tw/streamdeck)](https://goreportcard.com/report/github.com/dh1tw/streamdeck)
[![Go Reference](https://pkg.go.dev/badge/github.com/dh1tw/streamdeck.svg)](https://pkg.go.dev/github.com/dh1tw/streamdeck)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://img.shields.io/badge/license-MIT-blue.svg)

![streamdeck buttons](https://i.imgur.com/tEt3tPr.jpg "streamdeck Buttons")
![streamdeck slide show](https://i.imgur.com/gh1xXiU.jpg "streamdeck Slideshow")

**streamdeck** is a library for interfacing with the [Elgato/Corsair Stream Deck](https://www.elgato.com/en/gaming/stream-deck)

This library is written in the programing language [Go](https://golang.org).

## Note
This project is a golang API for the Elgato/Corsair Stream Deck. This library
unleashes the power of the StreamDeck. It allows you to completely customize
the content of the device, without the need of the OEM's software.

## License

streamdeck is published under the permissive [MIT license](https://github.com/dh1tw/streamdeck/blob/master/LICENSE).

## Dependencies

There are a few go libraries which are needed at compile time. streamdeck
does not have any runtime dependencies.

## CGO

Compiling this library requires unfortunately a C compiler since the underlying [HID library](github.com/karalabe/hid) requires cgo for enumerating the HID devices.

## Supported Operating Systems

The library should work on Linux, MacOS and Windows (>=7).

streamdeck works well on SoC boards like the Raspberry / Orange / Banana Pis.

### Linux Device rules

On Linux you might have to create an udev rule, to access the streamdeck.

````
sudo vim /etc/udev/rules.d/99-streamdeck.rules

SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", MODE="0664", GROUP="plugdev"
````

After saving the udev rule, unplug and plug the streamdeck again into the USB port.
For the rule above, your user must be a member of the `plugdev` group.

Make sure that your streamdeck is correctly recognized
by executing the following command:

````bash
$ go run examples/enumerate/enumerate.go
Found 1 Elgato Stream Deck(s):
	SerialNumber:        AL12H1A07123
````

## Documentation

The auto generated documentation can be found at [godoc.org](https://godoc.org/github.com/dh1tw/streamdeck)

## Examples

Please checkout the dedicated repository [streamdeck-examples](https://github.com/dh1tw/streamdeck-examples) for examples.

My personal library of streamdeck elements / buttons can be found in the [streamdeck-buttons](https://github.com/dh1tw/streamdeck-buttons) repository.

## Credits

This project would not have been possible without the work of [Alex Van Camp](https://github.com/Lange). In particular his
[notes of the StreamDeck's protocol](https://github.com/Lange/node-elgato-stream-deck/blob/master/NOTES.md)
were very helpful.
Alex has provided a [reference implementation](https://github.com/Lange/node-elgato-stream-deck) in Javascript / Node.js.