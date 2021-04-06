package fitz

import (
	"errors"
	"image/color"

	"go.matteson.dev/gfx"
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
	CloseDeviceCommand
)

var ErrBreak = errors.New("break")

type Device interface {
	Error() error

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

func (dev *BaseDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64)        {}
func (dev *BaseDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64)             {}
func (dev *BaseDevice) FillImageMask(image *Image, matrix gfx.Matrix, fillColor color.Color) {}
func (dev *BaseDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
}

func (dev *BaseDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {}
func (dev *BaseDevice) FillText(text *Text, matrix gfx.Matrix, color color.Color)       {}
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

type ChainDevice struct {
	ReplayDevice
	devices []Device
}

func NewChainDevice(devices ...Device) Device {
	return &ChainDevice{
		ReplayDevice: ReplayDevice{
			replayList: &ReplayList{},
		},
		devices: devices,
	}
}

func (dev *ChainDevice) Close() {
	for _, device := range dev.devices {
		if dev.Err == nil {
			dev.Err = dev.replayList.Apply(device)
		}
		device.Close()
	}
}

type CompositeDevice struct {
	BaseDevice
	devices []Device
}

func NewCompositeDevice(devices ...Device) Device {
	return &CompositeDevice{devices: devices}
}

func (dev *CompositeDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillPath(path, fillRule, matrix, fillColor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		device.StrokePath(path, stroke, matrix, strokeColor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		device.FillShade(shade, matrix, alpha)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		device.FillImage(image, matrix, alpha)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) FillImageMask(image *Image, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillImageMask(image, matrix, fillColor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipPath(path, fillRule, matrix, scissor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipStrokePath(path, stroke, matrix, scissor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipImageMask(image, matrix, scissor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) FillText(text *Text, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillText(text, matrix, fillColor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		device.StrokeText(text, stroke, matrix, strokeColor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipText(text, matrix, scissor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipStrokeText(text, stroke, matrix, scissor)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) IgnoreText(text *Text, matrix gfx.Matrix) {
	for _, device := range dev.devices {
		device.IgnoreText(text, matrix)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) PopClip() {
	for _, device := range dev.devices {
		device.PopClip()
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	for _, device := range dev.devices {
		device.BeginMask(rect, color, luminosity)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) EndMask() {
	for _, device := range dev.devices {
		device.EndMask()
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendMode gfx.BlendMode, alpha float64) {
	for _, device := range dev.devices {
		device.BeginGroup(rect, cs, isolated, knockout, blendMode, alpha)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) EndGroup() {
	for _, device := range dev.devices {
		device.EndGroup()
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) BeginTile() int {
	for _, device := range dev.devices {
		device.BeginTile()
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
	return 0
}

func (dev *CompositeDevice) EndTile() {
	for _, device := range dev.devices {
		device.EndTile()
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) BeginLayer(layerName string) {
	for _, device := range dev.devices {
		device.BeginLayer(layerName)
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) EndLayer() {
	for _, device := range dev.devices {
		device.EndLayer()
		if device.Error() != nil {
			dev.Err = device.Error()
			break
		}
	}
}

func (dev *CompositeDevice) Close() {
	for _, device := range dev.devices {
		device.Close()
	}
}
