package fitz

import (
	"image/color"

	"go.matteson.dev/gfx"
)

type CompositionDevice struct {
	devices []Device
}

func NewCompositionDevice(devices ...Device) Device {
	return &CompositionDevice{devices: devices}
}

func (dev *CompositionDevice) AddDevice(devices ...Device) {
	dev.devices = append(dev.devices, devices...)
}

func (dev *CompositionDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillPath(path, fillRule, matrix, fillColor)
	}
}

func (dev *CompositionDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		device.StrokePath(path, stroke, matrix, strokeColor)
	}
}

func (dev *CompositionDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		device.FillShade(shade, matrix, alpha)
	}
}

func (dev *CompositionDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		device.FillImage(image, matrix, alpha)
	}
}

func (dev *CompositionDevice) FillImageMask(image *Image, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillImageMask(image, matrix, fillColor)
	}
}

func (dev *CompositionDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipPath(path, fillRule, matrix, scissor)
	}
}

func (dev *CompositionDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipStrokePath(path, stroke, matrix, scissor)
	}
}

func (dev *CompositionDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipImageMask(image, matrix, scissor)
	}
}

func (dev *CompositionDevice) FillText(text *Text, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillText(text, matrix, fillColor)
	}
}

func (dev *CompositionDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		device.StrokeText(text, stroke, matrix, strokeColor)
	}
}

func (dev *CompositionDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipText(text, matrix, scissor)
	}
}

func (dev *CompositionDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipStrokeText(text, stroke, matrix, scissor)
	}
}

func (dev *CompositionDevice) IgnoreText(text *Text, matrix gfx.Matrix) {
	for _, device := range dev.devices {
		device.IgnoreText(text, matrix)
	}
}

func (dev *CompositionDevice) PopClip() {
	for _, device := range dev.devices {
		device.PopClip()
	}
}

func (dev *CompositionDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	for _, device := range dev.devices {
		device.BeginMask(rect, color, luminosity)
	}
}

func (dev *CompositionDevice) EndMask() {
	for _, device := range dev.devices {
		device.EndMask()
	}
}

func (dev *CompositionDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendMode gfx.BlendMode, alpha float64) {
	for _, device := range dev.devices {
		device.BeginGroup(rect, cs, isolated, knockout, blendMode, alpha)
	}
}

func (dev *CompositionDevice) EndGroup() {
	for _, device := range dev.devices {
		device.EndGroup()
	}
}

func (dev *CompositionDevice) BeginTile() int {
	for _, device := range dev.devices {
		device.BeginTile()
	}
	return 0
}

func (dev *CompositionDevice) EndTile() {
	for _, device := range dev.devices {
		device.EndTile()
	}
}

func (dev *CompositionDevice) BeginLayer(layerName string) {
	for _, device := range dev.devices {
		device.BeginLayer(layerName)
	}
}

func (dev *CompositionDevice) EndLayer() {
	for _, device := range dev.devices {
		device.EndLayer()
	}
}

func (dev *CompositionDevice) Done() {
	for _, device := range dev.devices {
		device.Done()
	}
}
