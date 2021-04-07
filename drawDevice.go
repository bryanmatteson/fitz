package fitz

import (
	"fmt"
	"image"
	"image/color"

	"go.matteson.dev/gfx"
)

type DrawDevice struct {
	BaseDevice
	image     image.Image
	context   *gfx.ImageContext
	transform gfx.Matrix
}

func NewDrawDevice(transform gfx.Matrix, dest *image.RGBA) Device {
	ctx := gfx.NewContextForImage(dest)
	ctx.Clear(color.White)
	return &DrawDevice{image: dest, context: ctx, transform: transform}
}

func (dev *DrawDevice) Should(kind CommandKind) bool { return true }

func (dev *DrawDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, ctm gfx.Matrix, fillColor color.Color) {
	dev.context.Save()
	defer dev.context.Restore()

	dev.context.SetTransformationMatrix(ctm.Concat(dev.transform))
	dev.context.SetFillColor(fillColor)
	dev.context.SetFillRule(fillRule)
	dev.context.Fill(path)
}

func (dev *DrawDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, ctm gfx.Matrix, strokeColor color.Color) {
	dev.context.Save()
	defer dev.context.Restore()

	dev.context.SetTransformationMatrix(ctm.Concat(dev.transform))
	dev.context.SetStroke(stroke)
	dev.context.SetStrokeColor(strokeColor)
	dev.context.Stroke(path)
}

func (dev *DrawDevice) FillShade(shade *gfx.Shader, ctm gfx.Matrix, alpha float64) {
	fmt.Println("FillShade")
}

func (dev *DrawDevice) FillImage(img image.Image, ctm gfx.Matrix, alpha float64) {
	fmt.Println("FillImage")
	dev.context.Save()
	defer dev.context.Restore()

	trm := ctm.Concat(dev.transform)

	// ctm is a transform from the image (expressed as a unit rect) to the destination device
	// we need to reverse and scale it to get a mapping from the destination to the source pixels,
	// then finally invert it to end up with a transform from image rect to destination device
	inv := trm.Inverted().PostScaled(float64(img.Bounds().Dx()), float64(img.Bounds().Dy())).Inverted()

	dev.context.SetTransformationMatrix(inv)
	dev.context.DrawImage(img)
}

func (dev *DrawDevice) FillImageMask(img image.Image, ctm gfx.Matrix, color color.Color) {
	fmt.Println("FillImageMask")
	dev.context.Save()
	defer dev.context.Restore()

	trm := ctm.Concat(dev.transform)
	inv := trm.Inverted().PostScaled(float64(img.Bounds().Dx()), float64(img.Bounds().Dy())).Inverted()
	dev.context.SetTransformationMatrix(inv)
}

func (dev *DrawDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, ctm gfx.Matrix, scissor gfx.Rect) {
	trm := ctm.Concat(dev.transform)

	clip := dev.transform.TransformRect(scissor)
	clip = clip.Intersection(gfx.ImageRect(dev.image.Bounds()))
	clip = clip.Intersection(trm.TransformRect(path.Bounds()))

	dev.context.Save()
	dev.context.DrawRect(clip)
	dev.context.Clip()
}

func (dev *DrawDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, ctm gfx.Matrix, scissor gfx.Rect) {
	trm := ctm.Concat(dev.transform)

	bounds := path.Bounds().Expanded(gfx.MakePoint(stroke.LineWidth*0.5, stroke.LineWidth*0.5))
	clip := trm.TransformRect(scissor.Intersection(bounds))

	dev.context.Save()
	dev.context.DrawRect(clip)
	dev.context.Clip()
}

func (dev *DrawDevice) ClipImageMask(img image.Image, ctm gfx.Matrix, scissor gfx.Rect) {
	inv := ctm.Concat(dev.transform).Inverted().PostScaled(float64(img.Bounds().Dx()), float64(img.Bounds().Dy())).Inverted()
	frame := inv.TransformRect(gfx.ImageRect(img.Bounds()))

	scissor = dev.transform.TransformRect(scissor)
	clip := scissor.Intersection(frame)

	dev.context.Save()
	dev.context.DrawRect(clip)
	dev.context.Clip()
	dev.context.DrawMask(img, inv)
}

func (dev *DrawDevice) FillText(text *Text, ctm gfx.Matrix, fillColor color.Color) {
	dev.context.Save()
	defer dev.context.Restore()

	dev.context.SetFillColor(fillColor)

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			glyph := span.Font.Glyph(letter.Rune, span.Matrix.Translated(letter.Origin.X, letter.Origin.Y).Compose(ctm, dev.transform))
			dev.context.Fill(glyph.Path)
		}
	}
}

func (dev *DrawDevice) StrokeText(text *Text, stroke *gfx.Stroke, ctm gfx.Matrix, color color.Color) {
	dev.context.Save()
	defer dev.context.Restore()

	dev.context.SetStrokeColor(color)
	dev.context.SetStroke(stroke)

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			glyph := span.Font.Glyph(letter.Rune, span.Matrix.Translated(letter.Origin.X, letter.Origin.Y).Compose(ctm, dev.transform))
			dev.context.Stroke(glyph.Path)
		}
	}
}

func (dev *DrawDevice) ClipText(text *Text, ctm gfx.Matrix, scissor gfx.Rect) {
	trm := ctm.Concat(dev.transform)
	dev.context.Save()
	dev.context.DrawRect(trm.TransformRect(scissor))
	dev.context.Clip()
}

func (dev *DrawDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, ctm gfx.Matrix, scissor gfx.Rect) {
	trm := ctm.Concat(dev.transform)
	clip := trm.TransformRect(scissor)

	dev.context.Save()
	dev.context.DrawRect(clip)
	dev.context.Clip()
}

func (dev *DrawDevice) PopClip() { dev.context.Restore() }

func (dev *DrawDevice) IgnoreText(text *Text, ctm gfx.Matrix)                      {}
func (dev *DrawDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {}
func (dev *DrawDevice) EndMask()                                                   {}
func (dev *DrawDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
}
func (dev *DrawDevice) EndGroup()                   {}
func (dev *DrawDevice) BeginTile() int              { return 0 }
func (dev *DrawDevice) EndTile()                    {}
func (dev *DrawDevice) BeginLayer(layerName string) {}
func (dev *DrawDevice) EndLayer()                   {}
func (dev *DrawDevice) Close()                      {}
