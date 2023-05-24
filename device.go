package fitz

import (
	"errors"
	"image"
	"image/color"

	"github.com/bryanmatteson/gfx"
)

type CommandKind int

const (
	FillPathCommand CommandKind = 1 << iota
	StrokePathCommand
	FillShadeCommand
	FillImageCommand
	FillImageMaskCommand
	ClipPathCommand
	ClipStrokePathCommand
	ClipImageMaskCommand
	FillTextCommand
	StrokeTextCommand
	ClipTextCommand
	ClipStrokeTextCommand
	IgnoreTextCommand
	PopClipCommand
	BeginMaskCommand
	BeginGroupCommand
	EndMaskCommand
	EndGroupCommand
	BeginTileCommand
	EndTileCommand
	BeginLayerCommand
	EndLayerCommand
)

var ErrBreak = errors.New("break")

type Device interface {
	Error() error

	FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color)
	StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color)
	FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)
	FillImage(img image.Image, matrix gfx.Matrix, alpha float64)
	FillImageMask(img image.Image, matrix gfx.Matrix, color color.Color)
	ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect)
	ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect)
	ClipImageMask(img image.Image, matrix gfx.Matrix, scissor gfx.Rect)
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
	Close()
}

type BaseDevice struct {
	Err error
}

func (dev *BaseDevice) Error() error { return dev.Err }

func (dev *BaseDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
}

func (dev *BaseDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
}

func (dev *BaseDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)           {}
func (dev *BaseDevice) FillImage(img image.Image, matrix gfx.Matrix, alpha float64)             {}
func (dev *BaseDevice) FillImageMask(img image.Image, matrix gfx.Matrix, fillColor color.Color) {}
func (dev *BaseDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) ClipImageMask(img image.Image, matrix gfx.Matrix, scissor gfx.Rect) {}
func (dev *BaseDevice) FillText(text *Text, matrix gfx.Matrix, color color.Color)          {}
func (dev *BaseDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
}

func (dev *BaseDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {}
func (dev *BaseDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) IgnoreText(text *Text, matrix gfx.Matrix)                       {}
func (dev *BaseDevice) PopClip()                                                       {}
func (dev *BaseDevice) BeginMask(rect gfx.Rect, maskColor color.Color, luminosity int) {}
func (dev *BaseDevice) EndMask()                                                       {}
func (dev *BaseDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
}

func (dev *BaseDevice) EndGroup()                   {}
func (dev *BaseDevice) BeginTile() int              { return 0 }
func (dev *BaseDevice) EndTile()                    {}
func (dev *BaseDevice) BeginLayer(layerName string) {}
func (dev *BaseDevice) EndLayer()                   {}
func (dev *BaseDevice) Close()                      {}
