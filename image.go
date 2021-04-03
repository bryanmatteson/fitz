package fitz

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"go.matteson.dev/gfx"
)

type Image struct {
	Rect    gfx.Rect
	Data    []byte
	Stride  int
	NumComp int
}

func (img Image) ColorModel() color.Model { return color.RGBAModel }
func (img Image) Bounds() image.Rectangle { return img.Rect.ImageRect() }

func (img Image) PngBytes() []byte {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
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
		return color.RGBA{}
	}
}
