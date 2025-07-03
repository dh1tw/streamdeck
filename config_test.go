package streamdeck

import (
	"go.viam.com/test"
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestFixKey(t *testing.T) {
	test.That(t, Plus.fixKey(5), test.ShouldEqual, 5)
	test.That(t, Original2.fixKey(5), test.ShouldEqual, 5)

	test.That(t, Original.fixKey(4), test.ShouldEqual, 1)
	test.That(t, Original.fixKey(5), test.ShouldEqual, 10)
}

func TestRGB(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, Original.ButtonSize, Original.ButtonSize))
	color := color.RGBA{uint8(255), uint8(0), uint8(0), 1}
	draw.Draw(img, img.Bounds(), image.NewUniform(color), image.Point{0, 0}, draw.Src)

	data, err := encodeBMP(Original, img)
	test.That(t, err, test.ShouldBeNil)

	test.That(t, len(data), test.ShouldEqual, 54+(72*72*3))

}
