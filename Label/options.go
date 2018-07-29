package Label

import "image"

// Text is a functional option for providing the initial text on the label.
// Max 5 characters.
func Text(text string) func(*Label) {
	return func(l *Label) {
		l.text = text
	}
}

// TextColor is a functional option which sets the text color.
func TextColor(c image.Uniform) func(*Label) {
	return func(l *Label) {
		l.textColor = &c
	}
}

// BgColor is a functional option which sets the background color of the label.
func BgColor(c image.Uniform) func(*Label) {
	return func(l *Label) {
		l.bgColor = &c
	}
}
