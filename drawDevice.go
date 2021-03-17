package fitz

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

type DrawDevice struct {
	BaseDevice
	image     image.Image
	context   *draw2dimg.GraphicContext
	transform Matrix
}

func NewDrawDevice(transform Matrix, dest *image.RGBA) GoDevice {
	ctx := draw2dimg.NewGraphicContext(dest)
	ctx.SetFillColor(color.White)
	ctx.Clear()
	return &DrawDevice{image: dest, context: ctx, transform: transform}
}

func (dev *DrawDevice) FillPath(path *Path, fillRule FillRule, ctm Matrix, fillColor color.Color) {
	trm := ctm.Concat(dev.transform)
	dev.context.SetMatrixTransform(toDrawMatrix(trm))
	drawPath(dev.context, path)
	dev.context.SetFillColor(fillColor)
	dev.context.SetFillRule(draw2d.FillRule(fillRule))
	dev.context.Fill()
}

func (dev *DrawDevice) StrokePath(path *Path, stroke *Stroke, ctm Matrix, strokeColor color.Color) {
	trm := ctm.Concat(dev.transform)
	dev.context.SetMatrixTransform(toDrawMatrix(trm))
	drawPath(dev.context, path)
	dev.context.SetStrokeColor(strokeColor)
	dev.context.SetLineCap(draw2d.LineCap(stroke.StartCap))
	dev.context.SetLineJoin(draw2d.LineJoin(stroke.LineJoin))
	dev.context.SetLineDash(stroke.Dashes, stroke.DashPhase)
	dev.context.SetLineWidth(stroke.LineWidth)
	dev.context.Stroke()
}

func (dev *DrawDevice) FillShade(shade *Shader, ctm Matrix, alpha float64) {}

func (dev *DrawDevice) FillImage(image *Image, ctm Matrix, alpha float64) {
	trm := ctm.Concat(dev.transform)
	dev.context.Save()

	// ctm is a transform from the image (expressed as a unit rect) to the destination device
	// we need to reverse and scale it to get a mapping from the destination to the source pixels,
	// then finally invert it to end up with a transform from image rect to destination device
	inv := trm.Inverted().PostScaled(image.Rect.Width(), image.Rect.Height()).Inverted()

	dev.context.SetMatrixTransform(toDrawMatrix(inv))
	dev.context.DrawImage(image)
}

func (dev *DrawDevice) FillImageMask(image *Image, ctm Matrix, color color.Color) {}
func (dev *DrawDevice) ClipPath(path *Path, fillRule FillRule, ctm Matrix, scissor Rect) {
}
func (dev *DrawDevice) ClipStrokePath(path *Path, stroke *Stroke, ctm Matrix, scissor Rect) {
}
func (dev *DrawDevice) ClipImageMask(image *Image, ctm Matrix, scissor Rect) {}

func (dev *DrawDevice) FillText(text *Text, ctm Matrix, fillColor color.Color) {
	ctm = ctm.Concat(dev.transform)
	dev.context.SetMatrixTransform(toDrawMatrix(ctm))

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			dev.context.SetFillColor(fillColor)
			drawPath(dev.context, letter.GlyphPath)
			dev.context.Fill()
		}
	}
}

func (dev *DrawDevice) StrokeText(text *Text, stroke *Stroke, ctm Matrix, color color.Color) {
	ctm = ctm.Concat(dev.transform)
	dev.context.SetMatrixTransform(toDrawMatrix(ctm))

	dev.context.SetStrokeColor(color)
	dev.context.SetLineCap(draw2d.LineCap(stroke.StartCap))
	dev.context.SetLineJoin(draw2d.LineJoin(stroke.LineJoin))
	dev.context.SetLineDash(stroke.Dashes, stroke.DashPhase)
	dev.context.SetLineWidth(stroke.LineWidth)

	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			drawPath(dev.context, letter.GlyphPath)
			dev.context.Stroke()
		}
	}
}

func (dev *DrawDevice) ClipText(text *Text, ctm Matrix, scissor Rect) {}
func (dev *DrawDevice) ClipStrokeText(text *Text, stroke *Stroke, ctm Matrix, scissor Rect) {
}

func (dev *DrawDevice) IgnoreText(text *Text, ctm Matrix)                      {}
func (dev *DrawDevice) PopClip()                                               {}
func (dev *DrawDevice) BeginMask(rect Rect, color color.Color, luminosity int) {}
func (dev *DrawDevice) EndMask()                                               {}
func (dev *DrawDevice) BeginGroup(rect Rect, cs *Colorspace, isolated bool, knockout bool, blendmode BlendMode, alpha float64) {
}
func (dev *DrawDevice) EndGroup()                   {}
func (dev *DrawDevice) BeginTile() int              { return 0 }
func (dev *DrawDevice) EndTile()                    {}
func (dev *DrawDevice) BeginLayer(layerName string) {}
func (dev *DrawDevice) EndLayer()                   {}
func (dev *DrawDevice) Close()                      {}

func drawPath(ctx *draw2dimg.GraphicContext, path *Path) {
	var j = 0
	for _, cmd := range path.Components {
		switch cmd {
		case MoveToComp:
			ctx.MoveTo(path.Points[j].X, path.Points[j].Y)
		case LineToComp:
			ctx.LineTo(path.Points[j].X, path.Points[j].Y)
		case QuadCurveToComp:
			ctx.QuadCurveTo(path.Points[j].X, path.Points[j].Y, path.Points[j+1].X, path.Points[j+1].Y)
		case CubicCurveToComp:
			ctx.CubicCurveTo(path.Points[j].X, path.Points[j].Y, path.Points[j+1].X, path.Points[j+1].Y, path.Points[j+2].X, path.Points[j+2].Y)
		case ClosePathComp:
			ctx.Close()
		}
		j += cmd.PointCount()
	}
}

func toDrawMatrix(mat Matrix) draw2d.Matrix {
	return draw2d.Matrix{mat.A, mat.B, mat.C, mat.D, mat.E, mat.F}
}
