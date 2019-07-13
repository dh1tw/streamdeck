package label

import (
	"image/color"

	sd "github.com/dh1tw/streamdeck"
)

// Text is a functional option for providing the initial text on the label.
// Max 5 characters.
func Text(text string) func(*Label) {
	return func(l *Label) {
		l.text = text
	}
}

// TextColor is a functional option which sets the text color.
func TextColor(c color.Color) func(*Label) {
	return func(l *Label) {
		l.textColor = c
	}
}

// BgColor is a functional option which sets the background color of the label.
func BgColor(c color.Color) func(*Label) {
	return func(l *Label) {
		l.bgColor = c
	}
}

//
func Callback(cb func(int, sd.BtnState)) func(*Label) {
	return func(l *Label) {
		l.cb = cb
	}
}
