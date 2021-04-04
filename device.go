package fitz

import (
	"image/color"
	"strings"

	"go.matteson.dev/gfx"
)

type Device interface {
	FillPathHandler
	StrokePathHandler
	FillShadeHandler
	FillImageHandler
	FillImageMaskHandler
	ClipPathHandler
	ClipStrokePathHandler
	ClipImageMaskHandler
	FillTextHandler
	StrokeTextHandler
	ClipTextHandler
	ClipStrokeTextHandler
	IgnoreTextHandler
	PopClipHandler
	BeginMaskHandler
	EndMaskHandler
	BeginGroupHandler
	EndGroupHandler
	BeginTileHandler
	EndTileHandler
	BeginLayerHandler
	EndLayerHandler

	ShouldCall(CommandKind) bool
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
func (dev *BaseDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, color color.Color) {
}

func (dev *BaseDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
}

func (dev *BaseDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)    {}
func (dev *BaseDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64)         {}
func (dev *BaseDevice) FillImageMask(image *Image, matrix gfx.Matrix, color color.Color) {}
func (dev *BaseDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {}
func (dev *BaseDevice) FillText(text *Text, matrix gfx.Matrix, color color.Color)       {}
func (dev *BaseDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
}

func (dev *BaseDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {}
func (dev *BaseDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) IgnoreText(text *Text, matrix gfx.Matrix)                   {}
func (dev *BaseDevice) PopClip()                                                   {}
func (dev *BaseDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {}
func (dev *BaseDevice) EndMask()                                                   {}
func (dev *BaseDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
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

type FillPathHandler interface {
	FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color)
}

type FillPathHandlerFunc func(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color)

func (fn FillPathHandlerFunc) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
	fn(path, fillRule, matrix, fillColor)
}

type StrokePathHandler interface {
	StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color)
}
type StrokePathHandlerFunc func(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color)

func (fn StrokePathHandlerFunc) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	fn(path, stroke, matrix, strokeColor)
}

type FillShadeHandler interface {
	FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)
}
type FillShadeHandlerFunc func(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)

func (fn FillShadeHandlerFunc) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	fn(shade, matrix, alpha)
}

type FillImageHandler interface {
	FillImage(image *Image, matrix gfx.Matrix, alpha float64)
}
type FillImageHandlerFunc func(image *Image, matrix gfx.Matrix, alpha float64)

func (fn FillImageHandlerFunc) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	fn(image, matrix, alpha)
}

type FillImageMaskHandler interface {
	FillImageMask(image *Image, matrix gfx.Matrix, color color.Color)
}

type FillImageMaskHandlerFunc func(image *Image, matrix gfx.Matrix, color color.Color)

func (fn FillImageMaskHandlerFunc) FillImageMask(image *Image, matrix gfx.Matrix, color color.Color) {
	fn(image, matrix, color)
}

type ClipPathHandler interface {
	ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect)
}
type ClipPathHandlerFunc func(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect)

func (fn ClipPathHandlerFunc) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	fn(path, fillRule, matrix, scissor)
}

type ClipStrokePathHandler interface {
	ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect)
}
type ClipStrokePathHandlerFunc func(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect)

func (fn ClipStrokePathHandlerFunc) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	fn(path, stroke, matrix, scissor)
}

type ClipImageMaskHandler interface {
	ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect)
}
type ClipImageMaskHandlerFunc func(image *Image, matrix gfx.Matrix, scissor gfx.Rect)

func (fn ClipImageMaskHandlerFunc) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	fn(image, matrix, scissor)
}

type FillTextHandler interface {
	FillText(text *Text, matrix gfx.Matrix, color color.Color)
}
type FillTextHandlerFunc func(text *Text, matrix gfx.Matrix, color color.Color)

func (fn FillTextHandlerFunc) FillText(text *Text, matrix gfx.Matrix, color color.Color) {
	fn(text, matrix, color)
}

type StrokeTextHandler interface {
	StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color)
}
type StrokeTextHandlerFunc func(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color)

func (fn StrokeTextHandlerFunc) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	fn(text, stroke, matrix, color)
}

type ClipTextHandler interface {
	ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect)
}
type ClipTextHandlerFunc func(text *Text, matrix gfx.Matrix, scissor gfx.Rect)

func (fn ClipTextHandlerFunc) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	fn(text, matrix, scissor)
}

type ClipStrokeTextHandler interface {
	ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect)
}
type ClipStrokeTextHandlerFunc func(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect)

func (fn ClipStrokeTextHandlerFunc) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	fn(text, stroke, matrix, scissor)
}

type IgnoreTextHandler interface {
	IgnoreText(text *Text, matrix gfx.Matrix)
}
type IgnoreTextHandlerFunc func(text *Text, matrix gfx.Matrix)

func (fn IgnoreTextHandlerFunc) IgnoreText(text *Text, matrix gfx.Matrix) { fn(text, matrix) }

type PopClipHandler interface{ PopClip() }
type PopClipHandlerFunc func()

func (fn PopClipHandlerFunc) PopClip() { fn() }

type BeginMaskHandler interface {
	BeginMask(rect gfx.Rect, color color.Color, luminosity int)
}
type BeginMaskHandlerFunc func(rect gfx.Rect, color color.Color, luminosity int)
type EndMaskHandler interface{ EndMask() }
type EndMaskHandlerFunc func()

func (fn BeginMaskHandlerFunc) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	fn(rect, color, luminosity)
}

func (fn EndMaskHandlerFunc) EndMask() { fn() }

type BeginGroupHandler interface {
	BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64)
}
type BeginGroupHandlerFunc func(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64)

func (fn BeginGroupHandlerFunc) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
	fn(rect, cs, isolated, knockout, blendmode, alpha)
}

type EndGroupHandler interface{ EndGroup() }
type EndGroupHandlerFunc func()

func (fn EndGroupHandlerFunc) EndGroup() { fn() }

type BeginTileHandler interface{ BeginTile() int }
type BeginTileHandlerFunc func() int

func (fn BeginTileHandlerFunc) BeginTile() int { return fn() }

type EndTileHandler interface{ EndTile() }
type EndTileHandlerFunc func()

func (fn EndTileHandlerFunc) EndTile() { fn() }

type BeginLayerHandler interface{ BeginLayer(layerName string) }
type BeginLayerHandlerFunc func(layerName string)

func (fn BeginLayerHandlerFunc) BeginLayer(layerName string) { fn(layerName) }

type EndLayerHandler interface{ EndLayer() }
type EndLayerHandlerFunc func()

func (fn EndLayerHandlerFunc) EndLayer() { fn() }

type CloseHandler interface{ Close() }
type CloseHandlerFunc func()

func (fn CloseHandlerFunc) Close() { fn() }
