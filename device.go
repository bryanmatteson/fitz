package fitz

import (
	"image/color"

	"go.matteson.dev/gfx"
)

type Device interface {
	FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color)
	StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color)
	FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)
	FillImage(image *Image, matrix gfx.Matrix, alpha float64)
	FillImageMask(image *Image, matrix gfx.Matrix, color color.Color)
	ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect)
	ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect)
	ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect)
	FillText(text *Text, matrix gfx.Matrix, color color.Color)
	StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color)
	ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect)
	ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect)
	IgnoreText(text *Text, matrix gfx.Matrix)
	PopClip()
	BeginMask(rect gfx.Rect, color color.Color, luminosity int)
	BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64)
	EndMask()
	EndGroup()
	BeginTile() int
	EndTile()
	BeginLayer(layerName string)
	EndLayer()
	Done()
}

type NullDevice struct{}

func (dev *NullDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
}

func (dev *NullDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
}

func (dev *NullDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)        {}
func (dev *NullDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64)             {}
func (dev *NullDevice) FillImageMask(image *Image, matrix gfx.Matrix, fillColor color.Color) {}
func (dev *NullDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *NullDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *NullDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {}
func (dev *NullDevice) FillText(text *Text, matrix gfx.Matrix, color color.Color)       {}
func (dev *NullDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
}

func (dev *NullDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {}
func (dev *NullDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *NullDevice) IgnoreText(text *Text, matrix gfx.Matrix)                       {}
func (dev *NullDevice) PopClip()                                                       {}
func (dev *NullDevice) BeginMask(rect gfx.Rect, maskColor color.Color, luminosity int) {}
func (dev *NullDevice) EndMask()                                                       {}
func (dev *NullDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
}

func (dev *NullDevice) EndGroup()                   {}
func (dev *NullDevice) BeginTile() int              { return 0 }
func (dev *NullDevice) EndTile()                    {}
func (dev *NullDevice) BeginLayer(layerName string) {}
func (dev *NullDevice) EndLayer()                   {}
func (dev *NullDevice) Done()                       {}
