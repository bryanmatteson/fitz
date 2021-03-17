package fitz

import (
	"image/color"
)

type ListDevice struct {
	BaseDevice
	displayList *DisplayList
}

func NewListDevice(displayList *DisplayList) GoDevice {
	return &ListDevice{
		displayList: displayList,
	}
}

// FillPath implements the GoDevice interface
func (dev *ListDevice) FillPath(path *Path, fillRule FillRule, matrix Matrix, color color.Color) {
	cmd := FillPathCommand{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Color:    color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// StrokePath implements the GoDevice interface
func (dev *ListDevice) StrokePath(path *Path, stroke *Stroke, matrix Matrix, color color.Color) {
	cmd := StrokePathCommand{
		Matrix: matrix,
		Path:   path,
		Stroke: stroke,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillShade implements the GoDevice interface
func (dev *ListDevice) FillShade(shade *Shader, matrix Matrix, alpha float64) {
	cmd := FillShadeCommand{
		Matrix: matrix,
		Alpha:  float64(alpha),
		Shader: shade,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillImage implements the GoDevice interface
func (dev *ListDevice) FillImage(image *Image, matrix Matrix, alpha float64) {
	cmd := FillImageCommand{
		Matrix: matrix,
		Image:  image,
		Alpha:  float64(alpha),
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillImageMask implements the GoDevice interface
func (dev *ListDevice) FillImageMask(image *Image, matrix Matrix, color color.Color) {
	cmd := FillImageMaskCommand{
		Matrix: matrix,
		Image:  image,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipPath implements the GoDevice interface
func (dev *ListDevice) ClipPath(path *Path, fillRule FillRule, matrix Matrix, scissor Rect) {
	cmd := ClipPathCommand{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Scissor:  scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipStrokePath implements the GoDevice interface
func (dev *ListDevice) ClipStrokePath(path *Path, stroke *Stroke, matrix Matrix, scissor Rect) {
	cmd := ClipStrokePathCommand{
		Matrix:  matrix,
		Path:    path,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipImageMask implements the GoDevice interface
func (dev *ListDevice) ClipImageMask(image *Image, matrix Matrix, scissor Rect) {
	cmd := ClipImageMaskCommand{
		Matrix:  matrix,
		Image:   image,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillText implements the GoDevice interface
func (dev *ListDevice) FillText(txt *Text, matrix Matrix, color color.Color) {
	cmd := FillTextCommand{
		Matrix: matrix,
		Text:   txt,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// StrokeText implements the GoDevice interface
func (dev *ListDevice) StrokeText(txt *Text, stroke *Stroke, matrix Matrix, color color.Color) {
	cmd := StrokeTextCommand{
		Matrix: matrix,
		Text:   txt,
		Stroke: stroke,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipText implements the GoDevice interface
func (dev *ListDevice) ClipText(txt *Text, matrix Matrix, scissor Rect) {
	cmd := ClipTextCommand{
		Matrix:  matrix,
		Text:    txt,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipStrokeText implements the GoDevice interface
func (dev *ListDevice) ClipStrokeText(txt *Text, stroke *Stroke, matrix Matrix, scissor Rect) {
	cmd := ClipStrokeTextCommand{
		Matrix:  matrix,
		Text:    txt,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// IgnoreText implements the GoDevice interface
func (dev *ListDevice) IgnoreText(txt *Text, matrix Matrix) {
	cmd := IgnoreTextCommand{
		Matrix: matrix,
		Text:   txt,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// PopClip implements the GoDevice interface
func (dev *ListDevice) PopClip() {
	dev.displayList.Commands = append(dev.displayList.Commands, &PopClipCommand{})
}

// BeginMask implements the GoDevice interface
func (dev *ListDevice) BeginMask(rect Rect, color color.Color, luminosity int) {
	cmd := BeginMaskCommand{
		Rect:       rect,
		Luminosity: luminosity,
		Color:      color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// EndMask implements the GoDevice interface
func (dev *ListDevice) EndMask() {
	dev.displayList.Commands = append(dev.displayList.Commands, &EndMaskCommand{})
}

// BeginGroup implements the GoDevice interface
func (dev *ListDevice) BeginGroup(rect Rect, cs *Colorspace, isolated bool, knockout bool, blendmode BlendMode, alpha float64) {
	cmd := BeginGroupCommand{
		Rect:      rect,
		Isolated:  isolated,
		Knockout:  knockout,
		BlendMode: blendmode,
		Alpha:     alpha,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// EndGroup implements the GoDevice interface
func (dev *ListDevice) EndGroup() {
	dev.displayList.Commands = append(dev.displayList.Commands, &EndGroupCommand{})
}

// BeginTile implements the GoDevice interface
func (dev *ListDevice) BeginTile() int {
	dev.displayList.Commands = append(dev.displayList.Commands, &BeginTileCommand{})
	return 0
}

// EndTile implements the GoDevice interface
func (dev *ListDevice) EndTile() {
	dev.displayList.Commands = append(dev.displayList.Commands, &EndTileCommand{})
}

// BeginLayer implements the GoDevice interface
func (dev *ListDevice) BeginLayer(layerName string) {
	dev.displayList.Commands = append(dev.displayList.Commands, &BeginLayerCommand{Name: layerName})
}

// EndLayer implements the GoDevice interface
func (dev *ListDevice) EndLayer() {
	dev.displayList.Commands = append(dev.displayList.Commands, &EndLayerCommand{})
}

// Close implements the GoDevice interface
func (dev *ListDevice) Close() {
	dev.displayList.Commands = append(dev.displayList.Commands, &CloseCommand{})
}

type GraphicsCommand interface {
	Kind() CommandKind
}

type FillPathCommand struct {
	Matrix   Matrix
	Path     *Path
	FillRule FillRule
	Color    color.Color
}

func (c FillPathCommand) Kind() CommandKind { return FillPath }

type StrokePathCommand struct {
	Matrix Matrix
	Path   *Path
	Stroke *Stroke
	Color  color.Color
}

func (c StrokePathCommand) Kind() CommandKind { return StrokePath }

type FillShadeCommand struct {
	Matrix Matrix
	Shader *Shader
	Alpha  float64
}

func (c FillShadeCommand) Kind() CommandKind { return FillShade }

type FillImageCommand struct {
	Matrix Matrix
	Image  *Image
	Alpha  float64
}

func (c FillImageCommand) Kind() CommandKind { return FillImage }

type FillImageMaskCommand struct {
	Matrix Matrix
	Image  *Image
	Color  color.Color
}

func (c FillImageMaskCommand) Kind() CommandKind { return FillImageMask }

type ClipPathCommand struct {
	Matrix   Matrix
	Path     *Path
	FillRule FillRule
	Scissor  Rect
}

func (c ClipPathCommand) Kind() CommandKind { return ClipPath }

type ClipStrokePathCommand struct {
	Matrix  Matrix
	Path    *Path
	Stroke  *Stroke
	Scissor Rect
}

func (c ClipStrokePathCommand) Kind() CommandKind { return ClipStrokePath }

type ClipImageMaskCommand struct {
	Matrix  Matrix
	Image   *Image
	Scissor Rect
}

func (c ClipImageMaskCommand) Kind() CommandKind { return ClipImageMask }

type FillTextCommand struct {
	Matrix Matrix
	Text   *Text
	Color  color.Color
}

func (c FillTextCommand) Kind() CommandKind { return FillText }

type StrokeTextCommand struct {
	Matrix Matrix
	Text   *Text
	Stroke *Stroke
	Color  color.Color
}

func (c StrokeTextCommand) Kind() CommandKind { return StrokeText }

type ClipTextCommand struct {
	Matrix  Matrix
	Text    *Text
	Scissor Rect
}

func (c ClipTextCommand) Kind() CommandKind { return ClipText }

type ClipStrokeTextCommand struct {
	Matrix  Matrix
	Text    *Text
	Stroke  *Stroke
	Scissor Rect
}

func (c ClipStrokeTextCommand) Kind() CommandKind { return ClipStrokeText }

type IgnoreTextCommand struct {
	Matrix Matrix
	Text   *Text
}

func (c IgnoreTextCommand) Kind() CommandKind { return IgnoreText }

type PopClipCommand struct{}

func (c PopClipCommand) Kind() CommandKind { return PopClip }

type BeginMaskCommand struct {
	Rect       Rect
	Color      color.Color
	Luminosity int
}

func (c BeginMaskCommand) Kind() CommandKind { return BeginMask }

type EndMaskCommand struct{}

func (c EndMaskCommand) Kind() CommandKind { return EndMask }

type BeginGroupCommand struct {
	Rect       Rect
	Colorspace *Colorspace
	Isolated   bool
	Knockout   bool
	BlendMode  BlendMode
	Alpha      float64
}

func (c BeginGroupCommand) Kind() CommandKind { return BeginGroup }

type EndGroupCommand struct{}

func (c EndGroupCommand) Kind() CommandKind { return EndGroup }

type BeginTileCommand struct{}

func (c BeginTileCommand) Kind() CommandKind { return BeginTile }

type EndTileCommand struct{}

func (c EndTileCommand) Kind() CommandKind { return EndTile }

type BeginLayerCommand struct {
	Name string
}

func (c BeginLayerCommand) Kind() CommandKind { return BeginLayer }

type EndLayerCommand struct{}

func (c EndLayerCommand) Kind() CommandKind { return EndLayer }

type CloseCommand struct{}

func (c CloseCommand) Kind() CommandKind { return CloseDevice }
