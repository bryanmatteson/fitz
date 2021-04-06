package fitz

import (
	"fmt"
	"image/color"

	"go.matteson.dev/gfx"
)

type ReplayList struct {
	commands []interface{}
}

func (list *ReplayList) ApplyAndClose(device Device) (err error) {
	defer device.Close()
	return list.Apply(device)
}

func (list *ReplayList) Apply(device Device) error {
	if device.Error() != nil {
		return device.Error()
	}

	for _, command := range list.commands {
		switch cmd := command.(type) {
		case (*FillPathCmd):
			device.FillPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Color)
		case (*StrokePathCmd):
			device.StrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Color)
		case (*FillShadeCmd):
			device.FillShade(cmd.Shader, cmd.Matrix, cmd.Alpha)
		case (*FillImageCmd):
			device.FillImage(cmd.Image, cmd.Matrix, cmd.Alpha)
		case (*FillImageMaskCmd):
			device.FillImageMask(cmd.Image, cmd.Matrix, cmd.Color)
		case (*ClipPathCmd):
			device.ClipPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Scissor)
		case (*ClipStrokePathCmd):
			device.ClipStrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Scissor)
		case (*ClipImageMaskCmd):
			device.ClipImageMask(cmd.Image, cmd.Matrix, cmd.Scissor)
		case (*FillTextCmd):
			device.FillText(cmd.Text, cmd.Matrix, cmd.Color)
		case (*StrokeTextCmd):
			device.StrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Color)
		case (*ClipTextCmd):
			device.ClipText(cmd.Text, cmd.Matrix, cmd.Scissor)
		case (*ClipStrokeTextCmd):
			device.ClipStrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Scissor)
		case (*IgnoreTextCmd):
			device.IgnoreText(cmd.Text, cmd.Matrix)
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

		if device.Error() != nil {
			break
		}
	}

	return device.Error()
}

type ReplayDevice struct {
	BaseDevice
	replayList *ReplayList
}

func NewReplayDevice(replayList *ReplayList) Device { return &ReplayDevice{replayList: replayList} }

// FillPath implements the GoDevice interface
func (dev *ReplayDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, color color.Color) {
	cmd := FillPathCmd{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Color:    color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// StrokePath implements the GoDevice interface
func (dev *ReplayDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	cmd := StrokePathCmd{
		Matrix: matrix,
		Path:   path,
		Stroke: stroke,
		Color:  color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillShade implements the GoDevice interface
func (dev *ReplayDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	cmd := FillShadeCmd{
		Matrix: matrix,
		Alpha:  float64(alpha),
		Shader: shade,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillImage implements the GoDevice interface
func (dev *ReplayDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	cmd := FillImageCmd{
		Matrix: matrix,
		Image:  image,
		Alpha:  float64(alpha),
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillImageMask implements the GoDevice interface
func (dev *ReplayDevice) FillImageMask(image *Image, matrix gfx.Matrix, color color.Color) {
	cmd := FillImageMaskCmd{
		Matrix: matrix,
		Image:  image,
		Color:  color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipPath implements the GoDevice interface
func (dev *ReplayDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipPathCmd{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Scissor:  scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipStrokePath implements the GoDevice interface
func (dev *ReplayDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipStrokePathCmd{
		Matrix:  matrix,
		Path:    path,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipImageMask implements the GoDevice interface
func (dev *ReplayDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipImageMaskCmd{
		Matrix:  matrix,
		Image:   image,
		Scissor: scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// FillText implements the GoDevice interface
func (dev *ReplayDevice) FillText(txt *Text, matrix gfx.Matrix, color color.Color) {
	cmd := FillTextCmd{
		Matrix: matrix,
		Text:   txt,
		Color:  color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// StrokeText implements the GoDevice interface
func (dev *ReplayDevice) StrokeText(txt *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	cmd := StrokeTextCmd{
		Matrix: matrix,
		Text:   txt,
		Stroke: stroke,
		Color:  color,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipText implements the GoDevice interface
func (dev *ReplayDevice) ClipText(txt *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipTextCmd{
		Matrix:  matrix,
		Text:    txt,
		Scissor: scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// ClipStrokeText implements the GoDevice interface
func (dev *ReplayDevice) ClipStrokeText(txt *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipStrokeTextCmd{
		Matrix:  matrix,
		Text:    txt,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.replayList.commands = append(dev.replayList.commands, &cmd)
}

// IgnoreText implements the GoDevice interface
func (dev *ReplayDevice) IgnoreText(txt *Text, matrix gfx.Matrix) {
	cmd := IgnoreTextCmd{
		Matrix: matrix,
		Text:   txt,
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
