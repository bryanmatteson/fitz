package fitz

import (
	"image/color"
	"strings"
)

type GoDevice interface {
	ShouldCall(CommandKind) bool

	FillPath(path *Path, fillRule FillRule, ctm Matrix, color color.Color)
	StrokePath(path *Path, stroke *Stroke, ctm Matrix, color color.Color)
	FillShade(shade *Shader, ctm Matrix, alpha float64)
	FillImage(image *Image, ctm Matrix, alpha float64)
	FillImageMask(image *Image, ctm Matrix, color color.Color)
	ClipPath(path *Path, fillRule FillRule, ctm Matrix, scissor Rect)
	ClipStrokePath(path *Path, stroke *Stroke, ctm Matrix, scissor Rect)
	ClipImageMask(image *Image, ctm Matrix, scissor Rect)
	FillText(text *Text, ctm Matrix, color color.Color)
	StrokeText(text *Text, stroke *Stroke, ctm Matrix, color color.Color)
	ClipText(text *Text, ctm Matrix, scissor Rect)
	ClipStrokeText(text *Text, stroke *Stroke, ctm Matrix, scissor Rect)
	IgnoreText(text *Text, ctm Matrix)
	PopClip()
	BeginMask(rect Rect, color color.Color, luminosity int)
	EndMask()
	BeginGroup(rect Rect, cs *Colorspace, isolated bool, knockout bool, blendmode BlendMode, alpha float64)
	EndGroup()
	BeginTile() int
	EndTile()
	BeginLayer(layerName string)
	EndLayer()
	Close()
}

// Commands
const (
	FillPath CommandKind = 1 << iota
	StrokePath
	FillShade
	FillImage
	FillImageMask
	ClipPath
	ClipStrokePath
	ClipImageMask
	FillText
	StrokeText
	ClipText
	ClipStrokeText
	IgnoreText
	PopClip
	BeginMask
	EndMask
	BeginGroup
	EndGroup
	BeginTile
	EndTile
	BeginLayer
	EndLayer
	CloseDevice
)

type BaseDevice struct{}

func (dev *BaseDevice) ShouldCall(CommandKind) bool { return true }
func (dev *BaseDevice) FillPath(path *Path, fillRule FillRule, matrix Matrix, color color.Color) {
}

func (dev *BaseDevice) StrokePath(path *Path, stroke *Stroke, matrix Matrix, color color.Color) {
}

func (dev *BaseDevice) FillShade(shade *Shader, matrix Matrix, alpha float64)        {}
func (dev *BaseDevice) FillImage(image *Image, matrix Matrix, alpha float64)         {}
func (dev *BaseDevice) FillImageMask(image *Image, matrix Matrix, color color.Color) {}
func (dev *BaseDevice) ClipPath(path *Path, fillRule FillRule, matrix Matrix, scissor Rect) {
}

func (dev *BaseDevice) ClipStrokePath(path *Path, stroke *Stroke, matrix Matrix, scissor Rect) {
}

func (dev *BaseDevice) ClipImageMask(image *Image, matrix Matrix, scissor Rect) {}
func (dev *BaseDevice) FillText(txt *Text, matrix Matrix, color color.Color)    {}
func (dev *BaseDevice) StrokeText(txt *Text, stroke *Stroke, matrix Matrix, color color.Color) {
}

func (dev *BaseDevice) ClipText(txt *Text, matrix Matrix, scissor Rect) {}
func (dev *BaseDevice) ClipStrokeText(txt *Text, stroke *Stroke, matrix Matrix, scissor Rect) {
}

func (dev *BaseDevice) IgnoreText(txt *Text, matrix Matrix)                    {}
func (dev *BaseDevice) PopClip()                                               {}
func (dev *BaseDevice) BeginMask(rect Rect, color color.Color, luminosity int) {}
func (dev *BaseDevice) EndMask()                                               {}
func (dev *BaseDevice) BeginGroup(rect Rect, cs *Colorspace, isolated bool, knockout bool, blendmode BlendMode, alpha float64) {
}

func (dev *BaseDevice) EndGroup()                   {}
func (dev *BaseDevice) BeginTile() int              { return 0 }
func (dev *BaseDevice) EndTile()                    {}
func (dev *BaseDevice) BeginLayer(layerName string) {}
func (dev *BaseDevice) EndLayer()                   {}
func (dev *BaseDevice) Close()                      {}

type CommandKind uint32

func (k *CommandKind) Set(flag CommandKind)     { *k = *k | flag }
func (k *CommandKind) Clear(flag CommandKind)   { *k = *k &^ flag }
func (k *CommandKind) Toggle(flag CommandKind)  { *k = *k ^ flag }
func (k CommandKind) Has(flag CommandKind) bool { return k&flag == flag }

var AllCommands = CommandKind(0xFFFFFFFF)

func (k CommandKind) String() string {
	kinds := []string{}

	if k.Has(FillPath) {
		kinds = append(kinds, "FillPath")
	}
	if k.Has(StrokePath) {
		kinds = append(kinds, "StrokePath")
	}
	if k.Has(FillShade) {
		kinds = append(kinds, "FillShade")
	}
	if k.Has(FillImage) {
		kinds = append(kinds, "FillImage")
	}
	if k.Has(FillImageMask) {
		kinds = append(kinds, "FillImageMask")
	}
	if k.Has(ClipPath) {
		kinds = append(kinds, "ClipPath")
	}
	if k.Has(ClipStrokePath) {
		kinds = append(kinds, "ClipStrokePath")
	}
	if k.Has(ClipImageMask) {
		kinds = append(kinds, "ClipImageMask")
	}
	if k.Has(FillText) {
		kinds = append(kinds, "FillText")
	}
	if k.Has(StrokeText) {
		kinds = append(kinds, "StrokeText")
	}
	if k.Has(ClipText) {
		kinds = append(kinds, "ClipText")
	}
	if k.Has(ClipStrokeText) {
		kinds = append(kinds, "ClipStrokeText")
	}
	if k.Has(IgnoreText) {
		kinds = append(kinds, "IgnoreText")
	}
	if k.Has(PopClip) {
		kinds = append(kinds, "PopClip")
	}
	if k.Has(BeginMask) {
		kinds = append(kinds, "BeginMask")
	}
	if k.Has(EndMask) {
		kinds = append(kinds, "EndMask")
	}
	if k.Has(BeginGroup) {
		kinds = append(kinds, "BeginGroup")
	}
	if k.Has(EndGroup) {
		kinds = append(kinds, "EndGroup")
	}
	if k.Has(BeginTile) {
		kinds = append(kinds, "BeginTile")
	}
	if k.Has(EndTile) {
		kinds = append(kinds, "EndTile")
	}
	if k.Has(BeginLayer) {
		kinds = append(kinds, "BeginLayer")
	}
	if k.Has(EndLayer) {
		kinds = append(kinds, "EndLayer")
	}
	if k.Has(CloseDevice) {
		kinds = append(kinds, "Close")
	}

	if len(kinds) == 0 {
		kinds = append(kinds, "unknown")
	}

	return strings.Join(kinds, " ")
}
