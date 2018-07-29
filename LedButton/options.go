package LedButton

import "image"

// TextColor is a functional option which sets the text color.
func TextColor(c image.Uniform) func(*LedButton) {
	return func(btn *LedButton) {
		btn.textColor = &c
	}
}

// LedColor is a functional option to set the color of the LED.
func LedColor(color LEDColor) func(*LedButton) {
	return func(btn *LedButton) {
		btn.ledColor = color
	}
}

// Text is a functional option for providing the initial text on the LED Button.
// Max 5 characters.
func Text(text string) func(*LedButton) {
	return func(btn *LedButton) {
		btn.text = text
	}
}
