package fitz

import (
	"fmt"
	"image/color"

	"go.matteson.dev/gfx"
)

type ReplayList struct {
	commands []interface{}
}

func (list *ReplayList) ReplayAndClose(device Device) error {
	defer device.Close()
	return list.Replay(device)
}

func (list *ReplayList) Replay(device Device) error {
	if device.Error() != nil {
		return device.Error()
	}

	for _, command := range list.commands {
		rundevcmd(device, command)
		if device.Error() != nil {
			break
		}
	}

	return device.Error()
}

func rundevcmd(device Device, command interface{}) {
	switch cmd := command.(type) {
	case (*FillPathCmd):
		device.FillPath(cmd.Path, cmd.FillRule, cmd.Transform, cmd.Color)
	case (*StrokePathCmd):
		device.StrokePath(cmd.Path, cmd.Stroke, cmd.Transform, cmd.Color)
	case (*FillShadeCmd):
		device.FillShade(cmd.Shader, cmd.Transform, cmd.Alpha)
	case (*FillImageCmd):
		device.FillImage(cmd.Image, cmd.Transform, cmd.Alpha)
	case (*FillImageMaskCmd):
		device.FillImageMask(cmd.Image, cmd.Transform, cmd.Color)
	case (*ClipPathCmd):
		device.ClipPath(cmd.Path, cmd.FillRule, cmd.Transform, cmd.Scissor)
	case (*ClipStrokePathCmd):
		device.ClipStrokePath(cmd.Path, cmd.Stroke, cmd.Transform, cmd.Scissor)
	case (*ClipImageMaskCmd):
		device.ClipImageMask(cmd.Image, cmd.Transform, cmd.Scissor)
	case (*FillTextCmd):
		device.FillText(cmd.Text, cmd.Transform, cmd.Color)
	case (*StrokeTextCmd):
		device.StrokeText(cmd.Text, cmd.Stroke, cmd.Transform, cmd.Color)
	case (*ClipTextCmd):
		device.ClipText(cmd.Text, cmd.Transform, cmd.Scissor)
	case (*ClipStrokeTextCmd):
		device.ClipStrokeText(cmd.Text, cmd.Stroke, cmd.Transform, cmd.Scissor)
	case (*IgnoreTextCmd):
		device.IgnoreText(cmd.Text, cmd.Transform)
	case (*PopClipCmd):
		device.PopClip()
	case (*BeginMaskCmd):
		device.BeginMask(cmd.Rect, cmd.Color, cmd.Luminosity)
	case (*EndMaskCmd):
		device.EndMask()
	case (*BegingGoupCmd):
		device.BeginGroup(cmd.Rect, cmd.Colorspace, cmd.Isolated, cmd.Knockout, cmd.BlendMode, cmd.Alpha)
	case (*EndGroupCmd):
		device.EndGroup()
	case (*BeginTileCmd):
		device.BeginTile()
	case (*EndTileCmd):
		device.EndTile()
	case (*BeginLayerCmd):
		device.BeginLayer(cmd.Name)
	case (*EndLayerCmd):
		device.EndLayer()
	default:
		panic(fmt.Sprintf("unknown command in display list: %v\n", cmd))
	}
}

type ReplayDevice struct {
	BaseDevice
	replayList *ReplayList
}

func NewReplayDevice(replayList *ReplayList) Device { return &ReplayDevice{replayList: replayList} }

// FillPath implements the GoDevice interface
func (dev *ReplayDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, trm gfx.Matrix, color color.Color) {
	cmd := FillPathCmd{
		Transform: trm,
		Path:      path,
		FillRule:  fillRule,
		Color:     color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// StrokePath implements the GoDevice interface
func (dev *ReplayDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, trm gfx.Matrix, color color.Color) {
	cmd := StrokePathCmd{
		Transform: trm,
		Path:      path,
		Stroke:    stroke,
		Color:     color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillShade implements the GoDevice interface
func (dev *ReplayDevice) FillShade(shade *gfx.Shader, trm gfx.Matrix, alpha float64) {
	cmd := FillShadeCmd{
		Transform: trm,
		Alpha:     float64(alpha),
		Shader:    shade,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillImage implements the GoDevice interface
func (dev *ReplayDevice) FillImage(image *Image, trm gfx.Matrix, alpha float64) {
	cmd := FillImageCmd{
		Transform: trm,
		Image:     image,
		Alpha:     float64(alpha),
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillImageMask implements the GoDevice interface
func (dev *ReplayDevice) FillImageMask(image *Image, trm gfx.Matrix, color color.Color) {
	cmd := FillImageMaskCmd{
		Transform: trm,
		Image:     image,
		Color:     color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipPath implements the GoDevice interface
func (dev *ReplayDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipPathCmd{
		Transform: trm,
		Path:      path,
		FillRule:  fillRule,
		Scissor:   scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipStrokePath implements the GoDevice interface
func (dev *ReplayDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipStrokePathCmd{
		Transform: trm,
		Path:      path,
		Stroke:    stroke,
		Scissor:   scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipImageMask implements the GoDevice interface
func (dev *ReplayDevice) ClipImageMask(image *Image, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipImageMaskCmd{
		Transform: trm,
		Image:     image,
		Scissor:   scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillText implements the GoDevice interface
func (dev *ReplayDevice) FillText(text *Text, trm gfx.Matrix, color color.Color) {
	cmd := FillTextCmd{
		Transform: trm,
		Text:      text,
		Color:     color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// StrokeText implements the GoDevice interface
func (dev *ReplayDevice) StrokeText(text *Text, stroke *gfx.Stroke, trm gfx.Matrix, color color.Color) {
	cmd := StrokeTextCmd{
		Transform: trm,
		Text:      text,
		Stroke:    stroke,
		Color:     color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipText implements the GoDevice interface
func (dev *ReplayDevice) ClipText(text *Text, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipTextCmd{
		Transform: trm,
		Text:      text,
		Scissor:   scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipStrokeText implements the GoDevice interface
func (dev *ReplayDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipStrokeTextCmd{
		Transform: trm,
		Text:      text,
		Stroke:    stroke,
		Scissor:   scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// IgnoreText implements the GoDevice interface
func (dev *ReplayDevice) IgnoreText(text *Text, trm gfx.Matrix) {
	cmd := IgnoreTextCmd{
		Transform: trm,
		Text:      text,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// PopClip implements the GoDevice interface
func (dev *ReplayDevice) PopClip() {
	dev.replayList.commands = append(dev.replayList.commands, &PopClipCmd{})
}

// BeginMask implements the GoDevice interface
func (dev *ReplayDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	cmd := BeginMaskCmd{
		Rect:       rect,
		Luminosity: luminosity,
		Color:      color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// EndMask implements the GoDevice interface
func (dev *ReplayDevice) EndMask() {
	dev.replayList.commands = append(dev.replayList.commands, &EndMaskCmd{})
}

// BeginGroup implements the GoDevice interface
func (dev *ReplayDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
	cmd := BegingGoupCmd{
		Rect:      rect,
		Isolated:  isolated,
		Knockout:  knockout,
		BlendMode: blendmode,
		Alpha:     alpha,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// EndGroup implements the GoDevice interface
func (dev *ReplayDevice) EndGroup() {
	dev.replayList.commands = append(dev.replayList.commands, &EndGroupCmd{})
}

// BeginTile implements the GoDevice interface
func (dev *ReplayDevice) BeginTile() int {
	dev.replayList.commands = append(dev.replayList.commands, &BeginTileCmd{})
	return 0
}

// EndTile implements the GoDevice interface
func (dev *ReplayDevice) EndTile() {
	dev.replayList.commands = append(dev.replayList.commands, &EndTileCmd{})
}

// BeginLayer implements the GoDevice interface
func (dev *ReplayDevice) BeginLayer(layerName string) {
	dev.replayList.commands = append(dev.replayList.commands, &BeginLayerCmd{Name: layerName})
}

// EndLayer implements the GoDevice interface
func (dev *ReplayDevice) EndLayer() {
	dev.replayList.commands = append(dev.replayList.commands, &EndLayerCmd{})
}
