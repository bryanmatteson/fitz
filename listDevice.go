package fitz

import (
	"image/color"

	"go.matteson.dev/gfx"
)

type ListDevice struct {
	displayList *DisplayList
}

func NewListDevice(displayList *DisplayList) Device {
	return &ListDevice{
		displayList: displayList,
	}
}

// FillPath implements the GoDevice interface
func (dev *ListDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, color color.Color) {
	cmd := FillPathCommand{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Color:    color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// StrokePath implements the GoDevice interface
func (dev *ListDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	cmd := StrokePathCommand{
		Matrix: matrix,
		Path:   path,
		Stroke: stroke,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillShade implements the GoDevice interface
func (dev *ListDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	cmd := FillShadeCommand{
		Matrix: matrix,
		Alpha:  float64(alpha),
		Shader: shade,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillImage implements the GoDevice interface
func (dev *ListDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	cmd := FillImageCommand{
		Matrix: matrix,
		Image:  image,
		Alpha:  float64(alpha),
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillImageMask implements the GoDevice interface
func (dev *ListDevice) FillImageMask(image *Image, matrix gfx.Matrix, color color.Color) {
	cmd := FillImageMaskCommand{
		Matrix: matrix,
		Image:  image,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipPath implements the GoDevice interface
func (dev *ListDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipPathCommand{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Scissor:  scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipStrokePath implements the GoDevice interface
func (dev *ListDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipStrokePathCommand{
		Matrix:  matrix,
		Path:    path,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipImageMask implements the GoDevice interface
func (dev *ListDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipImageMaskCommand{
		Matrix:  matrix,
		Image:   image,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// FillText implements the GoDevice interface
func (dev *ListDevice) FillText(txt *Text, matrix gfx.Matrix, color color.Color) {
	cmd := FillTextCommand{
		Matrix: matrix,
		Text:   txt,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// StrokeText implements the GoDevice interface
func (dev *ListDevice) StrokeText(txt *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	cmd := StrokeTextCommand{
		Matrix: matrix,
		Text:   txt,
		Stroke: stroke,
		Color:  color,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipText implements the GoDevice interface
func (dev *ListDevice) ClipText(txt *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipTextCommand{
		Matrix:  matrix,
		Text:    txt,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// ClipStrokeText implements the GoDevice interface
func (dev *ListDevice) ClipStrokeText(txt *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := ClipStrokeTextCommand{
		Matrix:  matrix,
		Text:    txt,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.displayList.Commands = append(dev.displayList.Commands, &cmd)
}

// IgnoreText implements the GoDevice interface
func (dev *ListDevice) IgnoreText(txt *Text, matrix gfx.Matrix) {
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
func (dev *ListDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
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
func (dev *ListDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
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
func (dev *ListDevice) Done() {
	dev.displayList.Commands = append(dev.displayList.Commands, &DoneCommand{})
}

type FillPathCommand struct {
	Matrix   gfx.Matrix
	Path     *gfx.Path
	FillRule gfx.FillRule
	Color    color.Color
}

type StrokePathCommand struct {
	Matrix gfx.Matrix
	Path   *gfx.Path
	Stroke *gfx.Stroke
	Color  color.Color
}

type FillShadeCommand struct {
	Matrix gfx.Matrix
	Shader *gfx.Shader
	Alpha  float64
}

type FillImageCommand struct {
	Matrix gfx.Matrix
	Image  *Image
	Alpha  float64
}

type FillImageMaskCommand struct {
	Matrix gfx.Matrix
	Image  *Image
	Color  color.Color
}

type ClipPathCommand struct {
	Matrix   gfx.Matrix
	Path     *gfx.Path
	FillRule gfx.FillRule
	Scissor  gfx.Rect
}

type ClipStrokePathCommand struct {
	Matrix  gfx.Matrix
	Path    *gfx.Path
	Stroke  *gfx.Stroke
	Scissor gfx.Rect
}

type ClipImageMaskCommand struct {
	Matrix  gfx.Matrix
	Image   *Image
	Scissor gfx.Rect
}

type FillTextCommand struct {
	Matrix gfx.Matrix
	Text   *Text
	Color  color.Color
}

type StrokeTextCommand struct {
	Matrix gfx.Matrix
	Text   *Text
	Stroke *gfx.Stroke
	Color  color.Color
}

type ClipTextCommand struct {
	Matrix  gfx.Matrix
	Text    *Text
	Scissor gfx.Rect
}

type ClipStrokeTextCommand struct {
	Matrix  gfx.Matrix
	Text    *Text
	Stroke  *gfx.Stroke
	Scissor gfx.Rect
}

type IgnoreTextCommand struct {
	Matrix gfx.Matrix
	Text   *Text
}

type PopClipCommand struct{}

type BeginMaskCommand struct {
	Rect       gfx.Rect
	Color      color.Color
	Luminosity int
}

type EndMaskCommand struct{}

type BeginGroupCommand struct {
	Rect       gfx.Rect
	Colorspace *gfx.Colorspace
	Isolated   bool
	Knockout   bool
	BlendMode  gfx.BlendMode
	Alpha      float64
}

type EndGroupCommand struct{}
type BeginTileCommand struct{}
type EndTileCommand struct{}
type BeginLayerCommand struct{ Name string }
type EndLayerCommand struct{}
type DoneCommand struct{}
