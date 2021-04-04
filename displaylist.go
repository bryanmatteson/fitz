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
			device.FillPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Color)
		case (*strokepathcmd):
			device.StrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Color)
		case (*fillshadecmd):
			device.FillShade(cmd.Shader, cmd.Matrix, cmd.Alpha)
		case (*fillimagecmd):
			device.FillImage(cmd.Image, cmd.Matrix, cmd.Alpha)
		case (*fillimagemaskcmd):
			device.FillImageMask(cmd.Image, cmd.Matrix, cmd.Color)
		case (*clippathcmd):
			device.ClipPath(cmd.Path, cmd.FillRule, cmd.Matrix, cmd.Scissor)
		case (*clipstrokepathcmd):
			device.ClipStrokePath(cmd.Path, cmd.Stroke, cmd.Matrix, cmd.Scissor)
		case (*clipimagemaskcmd):
			device.ClipImageMask(cmd.Image, cmd.Matrix, cmd.Scissor)
		case (*filltextcmd):
			device.FillText(cmd.Text, cmd.Matrix, cmd.Color)
		case (*stroketextcmd):
			device.StrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Color)
		case (*cliptextcmd):
			device.ClipText(cmd.Text, cmd.Matrix, cmd.Scissor)
		case (*clipstroketextcmd):
			device.ClipStrokeText(cmd.Text, cmd.Stroke, cmd.Matrix, cmd.Scissor)
		case (*ignoretextcmd):
			device.IgnoreText(cmd.Text, cmd.Matrix)
		case (*popclipcmd):
			device.PopClip()
		case (*beginmaskcmd):
			device.BeginMask(cmd.Rect, cmd.Color, cmd.Luminosity)
		case (*endmaskcmd):
			device.EndMask()
		case (*begingroupcmd):
			device.BeginGroup(cmd.Rect, cmd.Colorspace, cmd.Isolated, cmd.Knockout, cmd.BlendMode, cmd.Alpha)
		case (*endgroupcmd):
			device.EndGroup()
		case (*begintilecmd):
			device.BeginTile()
		case (*endtilecmd):
			device.EndTile()
		case (*beginlayercmd):
			device.BeginLayer(cmd.Name)
		case (*endlayercmd):
			device.EndLayer()
		case (*closecmd):
			device.Close()
		default:
			panic(fmt.Sprintf("unknown command in display list: %v\n", cmd))
		}
	}
}
