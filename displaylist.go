package fitz

import (
	"log"
)

type DisplayList struct {
	Commands []GraphicsCommand
}

func (list *DisplayList) Apply(device GoDevice) {
	kinds := device.Methods()

	for _, command := range list.Commands {
		if !kinds.Has(command.Kind()) {
			continue
		}

		switch cmd := command.(type) {
		case (*FillPathCommand):
			device.FillPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Color)
		case (*StrokePathCommand):
			device.StrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Color)
		case (*FillShadeCommand):
			device.FillShade(cmd.Shader, cmd.Matrix, cmd.Alpha)
		case (*FillImageCommand):
			device.FillImage(cmd.Image, cmd.Matrix, cmd.Alpha)
		case (*FillImageMaskCommand):
			device.FillImageMask(cmd.Image, cmd.Matrix, cmd.Color)
		case (*ClipPathCommand):
			device.ClipPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Scissor)
		case (*ClipStrokePathCommand):
			device.ClipStrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Scissor)
		case (*ClipImageMaskCommand):
			device.ClipImageMask(cmd.Image, cmd.Matrix, cmd.Scissor)
		case (*FillTextCommand):
			device.FillText(cmd.Text, cmd.Matrix, cmd.Color)
		case (*StrokeTextCommand):
			device.StrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Color)
		case (*ClipTextCommand):
			device.ClipText(cmd.Text, cmd.Matrix, cmd.Scissor)
		case (*ClipStrokeTextCommand):
			device.ClipStrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Scissor)
		case (*IgnoreTextCommand):
			device.IgnoreText(cmd.Text, cmd.Matrix)
		case (*PopClipCommand):
			device.PopClip()
		case (*BeginMaskCommand):
			device.BeginMask(cmd.Rect, cmd.Color, cmd.Luminosity)
		case (*EndMaskCommand):
			device.EndMask()
		case (*BeginGroupCommand):
			device.BeginGroup(cmd.Rect, cmd.Colorspace, cmd.Isolated, cmd.Knockout, cmd.BlendMode, cmd.Alpha)
		case (*EndGroupCommand):
			device.EndGroup()
		case (*BeginTileCommand):
			device.BeginTile()
		case (*EndTileCommand):
			device.EndTile()
		case (*BeginLayerCommand):
			device.BeginLayer(cmd.Name)
		case (*EndLayerCommand):
			device.EndLayer()
		case (*CloseCommand):
			device.Close()
		default:
			log.Printf("skipping unknown command in display list: %v\n", cmd.Kind())
		}
	}
}
