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
	Close()
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
func (dev *NullDevice) Close()                      {}

type BreakDevice struct {
	Device
	done bool
}

func (dev *BreakDevice) Break() { dev.done = true }

func (dev *BreakDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
	if dev.done {
		return
	}
	dev.Device.FillPath(path, fillRule, matrix, fillColor)
}

func (dev *BreakDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	if dev.done {
		return
	}
	dev.Device.StrokePath(path, stroke, matrix, strokeColor)
}

func (dev *BreakDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	if dev.done {
		return
	}
	dev.Device.FillShade(shade, matrix, alpha)
}

func (dev *BreakDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	if dev.done {
		return
	}
	dev.Device.FillImage(image, matrix, alpha)
}

func (dev *BreakDevice) FillImageMask(image *Image, matrix gfx.Matrix, fillColor color.Color) {
	if dev.done {
		return
	}
	dev.Device.FillImageMask(image, matrix, fillColor)
}

func (dev *BreakDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	if dev.done {
		return
	}
	dev.Device.ClipPath(path, fillRule, matrix, scissor)
}

func (dev *BreakDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	if dev.done {
		return
	}
	dev.Device.ClipStrokePath(path, stroke, matrix, scissor)
}

func (dev *BreakDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	if dev.done {
		return
	}
	dev.Device.ClipImageMask(image, matrix, scissor)
}

func (dev *BreakDevice) FillText(text *Text, matrix gfx.Matrix, color color.Color) {
	if dev.done {
		return
	}
	dev.Device.FillText(text, matrix, color)
}

func (dev *BreakDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	if dev.done {
		return
	}
	dev.Device.StrokeText(text, stroke, matrix, strokeColor)
}

func (dev *BreakDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	if dev.done {
		return
	}
	dev.Device.ClipText(text, matrix, scissor)
}

func (dev *BreakDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	if dev.done {
		return
	}
	dev.Device.ClipStrokeText(text, stroke, matrix, scissor)
}

func (dev *BreakDevice) IgnoreText(text *Text, matrix gfx.Matrix) {
	if dev.done {
		return
	}
	dev.Device.IgnoreText(text, matrix)
}

func (dev *BreakDevice) PopClip() {
	if dev.done {
		return
	}
	dev.Device.PopClip()
}

func (dev *BreakDevice) BeginMask(rect gfx.Rect, maskColor color.Color, luminosity int) {
	if dev.done {
		return
	}
	dev.Device.BeginMask(rect, maskColor, luminosity)
}

func (dev *BreakDevice) EndMask() {
	if dev.done {
		return
	}
	dev.EndMask()
}

func (dev *BreakDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
	if dev.done {
		return
	}
	dev.BeginGroup(rect, cs, isolated, knockout, blendmode, alpha)
}

func (dev *BreakDevice) EndGroup() {
	if dev.done {
		return
	}
	dev.Device.EndGroup()
}

func (dev *BreakDevice) BeginTile() int {
	if dev.done {
		return 0
	}
	return dev.Device.BeginTile()
}

func (dev *BreakDevice) EndTile() {
	if dev.done {
		return
	}
	dev.Device.EndTile()
}

func (dev *BreakDevice) BeginLayer(layerName string) {
	if dev.done {
		return
	}
	dev.Device.BeginLayer(layerName)
}

func (dev *BreakDevice) EndLayer() {
	if dev.done {
		return
	}
	dev.Device.EndLayer()
}

func (dev *BreakDevice) Close() {
	dev.Device.Close()
}

type CompositeDevice struct {
	devices []Device
}

func NewCompositeDevice(devices ...Device) Device {
	return &CompositeDevice{devices: devices}
}

func (dev *CompositeDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillPath(path, fillRule, matrix, fillColor)
	}
}

func (dev *CompositeDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		device.StrokePath(path, stroke, matrix, strokeColor)
	}
}

func (dev *CompositeDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		device.FillShade(shade, matrix, alpha)
	}
}

func (dev *CompositeDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.devices {
		device.FillImage(image, matrix, alpha)
	}
}

func (dev *CompositeDevice) FillImageMask(image *Image, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillImageMask(image, matrix, fillColor)
	}
}

func (dev *CompositeDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipPath(path, fillRule, matrix, scissor)
	}
}

func (dev *CompositeDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipStrokePath(path, stroke, matrix, scissor)
	}
}

func (dev *CompositeDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipImageMask(image, matrix, scissor)
	}
}

func (dev *CompositeDevice) FillText(text *Text, matrix gfx.Matrix, fillColor color.Color) {
	for _, device := range dev.devices {
		device.FillText(text, matrix, fillColor)
	}
}

func (dev *CompositeDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, strokeColor color.Color) {
	for _, device := range dev.devices {
		device.StrokeText(text, stroke, matrix, strokeColor)
	}
}

func (dev *CompositeDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipText(text, matrix, scissor)
	}
}

func (dev *CompositeDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.devices {
		device.ClipStrokeText(text, stroke, matrix, scissor)
	}
}

func (dev *CompositeDevice) IgnoreText(text *Text, matrix gfx.Matrix) {
	for _, device := range dev.devices {
		device.IgnoreText(text, matrix)
	}
}

func (dev *CompositeDevice) PopClip() {
	for _, device := range dev.devices {
		device.PopClip()
	}
}

func (dev *CompositeDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	for _, device := range dev.devices {
		device.BeginMask(rect, color, luminosity)
	}
}

func (dev *CompositeDevice) EndMask() {
	for _, device := range dev.devices {
		device.EndMask()
	}
}

func (dev *CompositeDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendMode gfx.BlendMode, alpha float64) {
	for _, device := range dev.devices {
		device.BeginGroup(rect, cs, isolated, knockout, blendMode, alpha)
	}
}

func (dev *CompositeDevice) EndGroup() {
	for _, device := range dev.devices {
		device.EndGroup()
	}
}

func (dev *CompositeDevice) BeginTile() int {
	for _, device := range dev.devices {
		device.BeginTile()
	}
	return 0
}

func (dev *CompositeDevice) EndTile() {
	for _, device := range dev.devices {
		device.EndTile()
	}
}

func (dev *CompositeDevice) BeginLayer(layerName string) {
	for _, device := range dev.devices {
		device.BeginLayer(layerName)
	}
}

func (dev *CompositeDevice) EndLayer() {
	for _, device := range dev.devices {
		device.EndLayer()
	}
}

func (dev *CompositeDevice) Close() {
	for _, device := range dev.devices {
		device.Close()
	}
}
