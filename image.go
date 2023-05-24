package fitz

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/bryanmatteson/gfx"
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

func (img Image) Cropped(color color.Color) gfx.Rect {
	var minx, miny, maxx, maxy float64
	minx, miny = math.Inf(1), math.Inf(1)
	maxx, maxy = math.Inf(-1), math.Inf(-1)

	cr, cg, cb, _ := color.RGBA()

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if (r == cr && g == cg && b == cb) || a == 0 {
				continue
			}

			minx = math.Min(minx, float64(x))
			miny = math.Min(miny, float64(y))
			maxx = math.Max(maxx, float64(x))
			maxy = math.Max(maxy, float64(y))
		}
	}

	return gfx.MakeRect(minx, miny, maxx, maxy)
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
