package fitz

import (
	"image/color"

	"go.matteson.dev/gfx"
)

type CompositionDevice struct {
	BaseDevice
	devices []GoDevice
}

func NewCompositionCropper(devices ...GoDevice) GoDevice {
	return &CompositionDevice{devices: devices}
}

func (dev *CompositionDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.devices {
		if device.ShouldCall(FillPath) {
			device.FillPath(path, fillRule, matrix, color)
		}
	}
}

func (dev *CompositionDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.devices {
		if device.ShouldCall(StrokePath) {
			device.StrokePath(path, stroke, matrix, color)
		}
	}
}

func (dev *CompositionDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		if device.ShouldCall(FillShade) {
			device.FillShade(shade, matrix, alpha)
		}
	}
}

func (dev *CompositionDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		if device.ShouldCall(FillImage) {
			device.FillImage(image, matrix, alpha)
		}
	}
}

func (dev *CompositionDevice) FillImageMask(image *Image, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.devices {
		if device.ShouldCall(FillImageMask) {
			device.FillImageMask(image, matrix, color)
		}
	}
}

func (dev *CompositionDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.ShouldCall(ClipPath) {
			device.ClipPath(path, fillRule, matrix, scissor)
		}
	}
}

func (dev *CompositionDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.ShouldCall(ClipStrokePath) {
			device.ClipStrokePath(path, stroke, matrix, scissor)
		}
	}
}

func (dev *CompositionDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.ShouldCall(ClipImageMask) {
			device.ClipImageMask(image, matrix, scissor)
		}
	}
}

func (dev *CompositionDevice) FillText(text *Text, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.devices {
		if device.ShouldCall(FillText) {
			device.FillText(text, matrix, color)
		}
	}
}

func (dev *CompositionDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.devices {
		if device.ShouldCall(StrokeText) {
			device.StrokeText(text, stroke, matrix, color)
		}
	}
}

func (dev *CompositionDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.ShouldCall(ClipText) {
			device.ClipText(text, matrix, scissor)
		}
	}
}

func (dev *CompositionDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		if device.ShouldCall(ClipStrokeText) {
			device.ClipStrokeText(text, stroke, matrix, scissor)
		}
	}
}

func (dev *CompositionDevice) IgnoreText(text *Text, matrix gfx.Matrix) {
	for _, device := range dev.devices {
		if device.ShouldCall(IgnoreText) {
			device.IgnoreText(text, matrix)
		}
	}
}

func (dev *CompositionDevice) PopClip() {
	for _, device := range dev.devices {
		if device.ShouldCall(PopClip) {
			device.PopClip()
		}
	}
}

func (dev *CompositionDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	for _, device := range dev.devices {
		if device.ShouldCall(BeginMask) {
			device.BeginMask(rect, color, luminosity)
		}
	}
}

func (dev *CompositionDevice) EndMask() {
	for _, device := range dev.devices {
		if device.ShouldCall(EndMask) {
			device.EndMask()
		}
	}
}

func (dev *CompositionDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendMode gfx.BlendMode, alpha float64) {
	for _, device := range dev.devices {
		if device.ShouldCall(BeginGroup) {
			device.BeginGroup(rect, cs, isolated, knockout, blendMode, alpha)
		}
	}
}

func (dev *CompositionDevice) EndGroup() {
	for _, device := range dev.devices {
		if device.ShouldCall(EndGroup) {
			device.EndGroup()
		}
	}
}

func (dev *CompositionDevice) BeginTile() int {
	for _, device := range dev.devices {
		if device.ShouldCall(BeginTile) {
			device.BeginTile()
		}
	}
	return 0
}

func (dev *CompositionDevice) EndTile() {
	for _, device := range dev.devices {
		if device.ShouldCall(EndTile) {
			device.EndTile()
		}
	}
}

func (dev *CompositionDevice) BeginLayer(layerName string) {
	for _, device := range dev.devices {
		if device.ShouldCall(BeginLayer) {
			device.BeginLayer(layerName)
		}
	}
}

func (dev *CompositionDevice) EndLayer() {
	for _, device := range dev.devices {
		if device.ShouldCall(EndLayer) {
			device.EndLayer()
		}
	}
}

func (dev *CompositionDevice) Close() {
	for _, device := range dev.devices {
		device.Close()
	}
}
