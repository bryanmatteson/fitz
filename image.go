package fitz

import (
	"image"
	"image/color"
)

type Image struct {
	Rect  Rect
	Frame Rect

	Data    []byte
	Stride  int
	NumComp int
}

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Image) Bounds() image.Rectangle {
	return image.Rect(int(img.Rect.X.Min), int(img.Rect.Y.Min), int(img.Rect.X.Max), int(img.Rect.Y.Max))
}

func (img Image) At(x, y int) color.Color {
	bounds := img.Bounds()
	if !(image.Point{x, y}.In(bounds)) {
		return color.RGBA{}
	}

	i := (y-bounds.Min.Y)*img.Stride + (x-bounds.Min.X)*img.NumComp

	c := img.NumComp
	s := img.Data[i : i+c : i+c]

	switch c {
	case 4:
		return color.RGBA{s[0], s[1], s[2], s[3]}
	case 3:
		return color.RGBA{s[0], s[1], s[2], 255}
	case 1:
		return color.Gray{Y: s[0]}
	default:
		return color.Transparent
	}
}
