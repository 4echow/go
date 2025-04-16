package main

import (
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

type Image struct {
	width  int
	height int
}

func (im Image) Bounds() image.Rectangle {
	return image.Rect(600, 300, im.width, im.height)
}

func (im Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (im Image) At(x, y int) color.Color {
	v := uint8((x + y) / 2 * (x ^ y))
	return color.RGBA{v, v, 89, 255}
}

func main() {
	m := Image{}
	pic.ShowImage(m)
}
