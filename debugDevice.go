package fitz

import (
	"image"
	"image/color"

	"go.matteson.dev/gfx"
)

type DebugDevice struct {
	BaseDevice
	image     image.Image
	context   *gfx.ImageContext
	transform gfx.Matrix
}

func NewDebugDevice(transform gfx.Matrix, dest *image.RGBA) GoDevice {
	ctx := gfx.NewImageContext(dest)
	ctx.Clear()
	return &DebugDevice{image: dest, context: ctx, transform: transform}
}

func (dev *DebugDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, ctm gfx.Matrix, fillColor color.Color) {
	dev.context.SetTransformationMatrix(ctm.Concat(dev.transform))
	dev.context.SetFillColor(fillColor)
	dev.context.SetFillRule(fillRule)
	dev.context.Fill(path)
}

func (dev *DebugDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, ctm gfx.Matrix, strokeColor color.Color) {
	dev.context.SetTransformationMatrix(ctm.Concat(dev.transform))
	dev.context.SetStroke(stroke)
	dev.context.SetStrokeColor(strokeColor)
	dev.context.Stroke(path)
}

func (dev *DebugDevice) FillShade(shade *gfx.Shader, ctm gfx.Matrix, alpha float64) {}

func (dev *DebugDevice) FillImage(image *Image, ctm gfx.Matrix, alpha float64) {
	trm := ctm.Concat(dev.transform)

	// ctm is a transform from the image (expressed as a unit rect) to the destination device
	// we need to reverse and scale it to get a mapping from the destination to the source pixels,
	// then finally invert it to end up with a transform from image rect to destination device
	inv := trm.Inverted().PostScaled(image.Rect.Width(), image.Rect.Height()).Inverted()

	dev.context.SetTransformationMatrix(inv)
	dev.context.DrawImage(image)
}

func (dev *DebugDevice) FillImageMask(image *Image, ctm gfx.Matrix, color color.Color) {}
func (dev *DebugDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, ctm gfx.Matrix, scissor gfx.Rect) {
}

func (dev *DebugDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, ctm gfx.Matrix, scissor gfx.Rect) {
}

func (dev *DebugDevice) ClipImageMask(image *Image, ctm gfx.Matrix, scissor gfx.Rect) {}

func (dev *DebugDevice) FillText(text *Text, ctm gfx.Matrix, fillColor color.Color) {
	dev.context.SetTransformationMatrix(gfx.IdentityMatrix)
	dev.context.SetFillColor(fillColor)

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			glyph := span.Font.Glyph(letter.Rune, span.Matrix.Translated(letter.Origin.X, letter.Origin.Y).Compose(ctm, dev.transform))
			dev.context.Fill(glyph.Path)
		}
		q := ctm.Concat(dev.transform).TransformQuad(span.Quad)
		dev.context.SetStrokeColor(color.RGBA{R: 255, A: 255})
		dev.context.DrawQuad(q)
		dev.context.Stroke()
	}

}

func (dev *DebugDevice) StrokeText(text *Text, stroke *gfx.Stroke, ctm gfx.Matrix, color color.Color) {
	dev.context.SetTransformationMatrix(gfx.IdentityMatrix)

	dev.context.SetStrokeColor(color)
	dev.context.SetStroke(stroke)

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			glyph := span.Font.Glyph(letter.Rune, span.Matrix.Translated(letter.Origin.X, letter.Origin.Y).Compose(ctm, dev.transform))
			dev.context.Stroke(glyph.Path)
		}
	}
}

func (dev *DebugDevice) ClipText(text *Text, ctm gfx.Matrix, scissor gfx.Rect) {}
func (dev *DebugDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, ctm gfx.Matrix, scissor gfx.Rect) {
}

func (dev *DebugDevice) IgnoreText(text *Text, ctm gfx.Matrix)                      {}
func (dev *DebugDevice) PopClip()                                                   {}
func (dev *DebugDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {}
func (dev *DebugDevice) EndMask()                                                   {}
func (dev *DebugDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
}
func (dev *DebugDevice) EndGroup()                   {}
func (dev *DebugDevice) BeginTile() int              { return 0 }
func (dev *DebugDevice) EndTile()                    {}
func (dev *DebugDevice) BeginLayer(layerName string) {}
func (dev *DebugDevice) EndLayer()                   {}
func (dev *DebugDevice) Close()                      {}
