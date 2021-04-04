package fitz

import (
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

type Device interface {
	Should(CommandKind) bool
	Drop()

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

type BaseDevice struct{}

func (dev *BaseDevice) Should(CommandKind) bool { return true }

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
func (dev *BaseDevice) Drop()                       {}

type LoopDevice struct {
	BaseDevice
	done bool
}

func (dev *LoopDevice) Break()                       { dev.done = true }
func (dev *LoopDevice) Should(kind CommandKind) bool { return !dev.done }

type ChainDevice struct {
	ListDevice
	devices []Device
}

func NewChainDevice(devices ...Device) Device {
	return &ChainDevice{devices: devices}
}

func (dev *ChainDevice) Close() {
	for _, device := range dev.devices {
		dev.displayList.Apply(device)
		if device.Should(CloseDeviceCommand) {
			device.Close()
		}
	}
}

func (dev *ChainDevice) Drop() {
	for _, device := range dev.devices {
		device.Drop()
	}
}

type CompositeDevice struct {
	devices []Device
}

func NewCompositeDevice(devices ...Device) Device {
	return &CompositeDevice{devices: devices}
}

func (dev *CompositeDevice) Should(kind CommandKind) bool { return true }

func (dev *CompositeDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		if device.Should(FillPathCommand) {
			device.FillPath(path, fillRule, matrix, fillColor)
		}
	}
}

func (dev *CompositeDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		if device.Should(StrokePathCommand) {
			device.StrokePath(path, stroke, matrix, strokeColor)
		}
	}
}

func (dev *CompositeDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		if device.Should(FillShadeCommand) {
			device.FillShade(shade, matrix, alpha)
		}
	}
}

func (dev *CompositeDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		if device.Should(FillImageCommand) {
			device.FillImage(image, matrix, alpha)
		}
	}
}

func (dev *CompositeDevice) FillImageMask(image *Image, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		if device.Should(FillImageMaskCommand) {
			device.FillImageMask(image, matrix, fillColor)
		}
	}
}

func (dev *CompositeDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.Should(ClipPathCommand) {
			device.ClipPath(path, fillRule, matrix, scissor)
		}
	}
}

func (dev *CompositeDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.Should(ClipStrokePathCommand) {
			device.ClipStrokePath(path, stroke, matrix, scissor)
		}
	}
}

func (dev *CompositeDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.Should(ClipImageMaskCommand) {
			device.ClipImageMask(image, matrix, scissor)
		}
	}
}

func (dev *CompositeDevice) FillText(text *Text, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		if device.Should(FillTextCommand) {
			device.FillText(text, matrix, fillColor)
		}
	}
}

func (dev *CompositeDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		if device.Should(StrokeTextCommand) {
			device.StrokeText(text, stroke, matrix, strokeColor)
		}
	}
}

func (dev *CompositeDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.Should(ClipTextCommand) {
			device.ClipText(text, matrix, scissor)
		}
	}
}

func (dev *CompositeDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.Should(ClipStrokeTextCommand) {
			device.ClipStrokeText(text, stroke, matrix, scissor)
		}
	}
}

func (dev *CompositeDevice) IgnoreText(text *Text, matrix gfx.Matrix) {
	for _, device := range dev.devices {
		if device.Should(IgnoreTextCommand) {
			device.IgnoreText(text, matrix)
		}
	}
}

func (dev *CompositeDevice) PopClip() {
	for _, device := range dev.devices {
		if device.Should(PopClipCommand) {
			device.PopClip()
		}
	}
}

func (dev *CompositeDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	for _, device := range dev.devices {
		if device.Should(BeginMaskCommand) {
			device.BeginMask(rect, color, luminosity)
		}
	}
}

func (dev *CompositeDevice) EndMask() {
	for _, device := range dev.devices {
		if device.Should(EndMaskCommand) {
			device.EndMask()
		}
	}
}

func (dev *CompositeDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendMode gfx.BlendMode, alpha float64) {
	for _, device := range dev.devices {
		if device.Should(BeginGroupCommand) {
			device.BeginGroup(rect, cs, isolated, knockout, blendMode, alpha)
		}
	}
}

func (dev *CompositeDevice) EndGroup() {
	for _, device := range dev.devices {
		if device.Should(EndGroupCommand) {
			device.EndGroup()
		}
	}
}

func (dev *CompositeDevice) BeginTile() int {
	for _, device := range dev.devices {
		if device.Should(BeginTileCommand) {
			device.BeginTile()
		}
	}
	return 0
}

func (dev *CompositeDevice) EndTile() {
	for _, device := range dev.devices {
		if device.Should(EndTileCommand) {
			device.EndTile()
		}
	}
}

func (dev *CompositeDevice) BeginLayer(layerName string) {
	for _, device := range dev.devices {
		if device.Should(BeginLayerCommand) {
			device.BeginLayer(layerName)
		}
	}
}

func (dev *CompositeDevice) EndLayer() {
	for _, device := range dev.devices {
		if device.Should(EndLayerCommand) {
			device.EndLayer()
		}
	}
}

func (dev *CompositeDevice) Close() {
	for _, device := range dev.devices {
		if device.Should(CloseDeviceCommand) {
			device.Close()
		}
	}
}

func (dev *CompositeDevice) Drop() {
	for _, device := range dev.devices {
		device.Drop()
	}
}
