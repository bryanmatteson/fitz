package fitz

import "fmt"

type DisplayList struct {
	PageNumber int
	commands   []interface{}
}

func (list *DisplayList) Apply(device Device) {
	for _, command := range list.commands {
		switch cmd := command.(type) {
		case (*fillpathcmd):
			if device.Should(FillPathCommand) {
				device.FillPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Color)
			}
		case (*strokepathcmd):
			if device.Should(StrokePathCommand) {
				device.StrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Color)
			}
		case (*fillshadecmd):
			if device.Should(FillShadeCommand) {
				device.FillShade(cmd.Shader, cmd.Matrix, cmd.Alpha)
			}
		case (*fillimagecmd):
			if device.Should(FillImageCommand) {
				device.FillImage(cmd.Image, cmd.Matrix, cmd.Alpha)
			}
		case (*fillimagemaskcmd):
			if device.Should(FillImageMaskCommand) {
				device.FillImageMask(cmd.Image, cmd.Matrix, cmd.Color)
			}
		case (*clippathcmd):
			if device.Should(ClipPathCommand) {
				device.ClipPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Scissor)
			}
		case (*clipstrokepathcmd):
			if device.Should(ClipStrokePathCommand) {
				device.ClipStrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Scissor)
			}
		case (*clipimagemaskcmd):
			if device.Should(ClipImageMaskCommand) {
				device.ClipImageMask(cmd.Image, cmd.Matrix, cmd.Scissor)
			}
		case (*filltextcmd):
			if device.Should(FillTextCommand) {
				device.FillText(cmd.Text, cmd.Matrix, cmd.Color)
			}
		case (*stroketextcmd):
			if device.Should(StrokeTextCommand) {
				device.StrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Color)
			}
		case (*cliptextcmd):
			if device.Should(ClipTextCommand) {
				device.ClipText(cmd.Text, cmd.Matrix, cmd.Scissor)
			}
		case (*clipstroketextcmd):
			if device.Should(ClipStrokeTextCommand) {
				device.ClipStrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Scissor)
			}
		case (*ignoretextcmd):
			if device.Should(IgnoreTextCommand) {
				device.IgnoreText(cmd.Text, cmd.Matrix)
			}
		case (*popclipcmd):
			if device.Should(PopClipCommand) {
				device.PopClip()
			}
		case (*beginmaskcmd):
			if device.Should(BeginMaskCommand) {
				device.BeginMask(cmd.Rect, cmd.Color, cmd.Luminosity)
			}
		case (*endmaskcmd):
			if device.Should(EndMaskCommand) {
				device.EndMask()
			}
		case (*begingroupcmd):
			if device.Should(BeginGroupCommand) {
				device.BeginGroup(cmd.Rect, cmd.Colorspace, cmd.Isolated, cmd.Knockout, cmd.BlendMode, cmd.Alpha)
			}
		case (*endgroupcmd):
			if device.Should(EndGroupCommand) {
				device.EndGroup()
			}
		case (*begintilecmd):
			if device.Should(BeginTileCommand) {
				device.BeginTile()
			}
		case (*endtilecmd):
			if device.Should(EndTileCommand) {
				device.EndTile()
			}
		case (*beginlayercmd):
			if device.Should(BeginLayerCommand) {
				device.BeginLayer(cmd.Name)
			}
		case (*endlayercmd):
			if device.Should(EndLayerCommand) {
				device.EndLayer()
			}
		case (*closecmd):
			if device.Should(CloseDeviceCommand) {
				device.Close()
			}
		case (*dropcmd):
			device.Drop()
		default:
			panic(fmt.Sprintf("unknown command in display list: %v\n", cmd))
		}
	}
}
