package fitz

import (
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

func NewDrawDevice(transform gfx.Matrix, dest *image.RGBA) GoDevice {
	ctx := gfx.NewImageContext(dest)
	return &DrawDevice{image: dest, context: ctx, transform: transform}
}

func (dev *DrawDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, ctm gfx.Matrix, fillColor color.Color) {
	dev.context.SetTransformationMatrix(ctm.Concat(dev.transform))
	dev.context.SetFillColor(fillColor)
	dev.context.SetFillRule(fillRule)
	dev.context.Fill(path)
}

func (dev *DrawDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, ctm gfx.Matrix, strokeColor color.Color) {
	dev.context.SetTransformationMatrix(ctm.Concat(dev.transform))
	dev.context.SetStrokeColor(strokeColor)
	dev.context.SetLineCap(stroke.StartCap)
	dev.context.SetLineJoin(stroke.LineJoin)
	dev.context.SetLineDash(stroke.Dashes, stroke.DashPhase)
	dev.context.SetLineWidth(stroke.LineWidth)
	dev.context.Stroke(path)
}

func (dev *DrawDevice) FillShade(shade *gfx.Shader, ctm gfx.Matrix, alpha float64) {}

func (dev *DrawDevice) FillImage(image *Image, ctm gfx.Matrix, alpha float64) {
	trm := ctm.Concat(dev.transform)

	// ctm is a transform from the image (expressed as a unit rect) to the destination device
	// we need to reverse and scale it to get a mapping from the destination to the source pixels,
	// then finally invert it to end up with a transform from image rect to destination device
	inv := trm.Inverted().PostScaled(image.Rect.Width(), image.Rect.Height()).Inverted()

	dev.context.SetTransformationMatrix(inv)
	dev.context.DrawImage(image)
}

func (dev *DrawDevice) FillImageMask(image *Image, ctm gfx.Matrix, color color.Color) {}
func (dev *DrawDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, ctm gfx.Matrix, scissor gfx.Rect) {
}
func (dev *DrawDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, ctm gfx.Matrix, scissor gfx.Rect) {
}
func (dev *DrawDevice) ClipImageMask(image *Image, ctm gfx.Matrix, scissor gfx.Rect) {}

func (dev *DrawDevice) FillText(text *Text, ctm gfx.Matrix, fillColor color.Color) {
	ctm = ctm.Concat(dev.transform)
	dev.context.SetTransformationMatrix(ctm)

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			_ = letter
			// dev.context.SetFillColor(fillColor)
			// drawPath(dev.context, letter.GlyphPath)
			// dev.context.Fill()
		}
	}
}

func (dev *DrawDevice) StrokeText(text *Text, stroke *gfx.Stroke, ctm gfx.Matrix, color color.Color) {
	ctm = ctm.Concat(dev.transform)
	dev.context.SetTransformationMatrix(ctm)

	dev.context.SetStrokeColor(color)
	dev.context.SetLineCap(stroke.StartCap)
	dev.context.SetLineJoin(stroke.LineJoin)
	dev.context.SetLineDash(stroke.Dashes, stroke.DashPhase)
	dev.context.SetLineWidth(stroke.LineWidth)

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			_ = letter
			// drawPath(dev.context, letter.GlyphPath)
			// dev.context.Stroke()
		}
	}
}

func (dev *DrawDevice) ClipText(text *Text, ctm gfx.Matrix, scissor gfx.Rect) {}
func (dev *DrawDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, ctm gfx.Matrix, scissor gfx.Rect) {
}

func (dev *DrawDevice) IgnoreText(text *Text, ctm gfx.Matrix)                      {}
func (dev *DrawDevice) PopClip()                                                   {}
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
