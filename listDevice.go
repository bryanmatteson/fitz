package fitz

import (
	"image/color"

	"go.matteson.dev/gfx"
)

type ListDevice struct {
	BaseDevice
	displayList *DisplayList
}

func NewListDevice(displayList *DisplayList) Device { return &ListDevice{displayList: displayList} }

// FillPath implements the GoDevice interface
func (dev *ListDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, color color.Color) {
	cmd := fillpathcmd{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Color:    color,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// StrokePath implements the GoDevice interface
func (dev *ListDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	cmd := strokepathcmd{
		Matrix: matrix,
		Path:   path,
		Stroke: stroke,
		Color:  color,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// FillShade implements the GoDevice interface
func (dev *ListDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	cmd := fillshadecmd{
		Matrix: matrix,
		Alpha:  float64(alpha),
		Shader: shade,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// FillImage implements the GoDevice interface
func (dev *ListDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	cmd := fillimagecmd{
		Matrix: matrix,
		Image:  image,
		Alpha:  float64(alpha),
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// FillImageMask implements the GoDevice interface
func (dev *ListDevice) FillImageMask(image *Image, matrix gfx.Matrix, color color.Color) {
	cmd := fillimagemaskcmd{
		Matrix: matrix,
		Image:  image,
		Color:  color,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// ClipPath implements the GoDevice interface
func (dev *ListDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := clippathcmd{
		Matrix:   matrix,
		Path:     path,
		FillRule: fillRule,
		Scissor:  scissor,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// ClipStrokePath implements the GoDevice interface
func (dev *ListDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := clipstrokepathcmd{
		Matrix:  matrix,
		Path:    path,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// ClipImageMask implements the GoDevice interface
func (dev *ListDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := clipimagemaskcmd{
		Matrix:  matrix,
		Image:   image,
		Scissor: scissor,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// FillText implements the GoDevice interface
func (dev *ListDevice) FillText(txt *Text, matrix gfx.Matrix, color color.Color) {
	cmd := filltextcmd{
		Matrix: matrix,
		Text:   txt,
		Color:  color,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// StrokeText implements the GoDevice interface
func (dev *ListDevice) StrokeText(txt *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	cmd := stroketextcmd{
		Matrix: matrix,
		Text:   txt,
		Stroke: stroke,
		Color:  color,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// ClipText implements the GoDevice interface
func (dev *ListDevice) ClipText(txt *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := cliptextcmd{
		Matrix:  matrix,
		Text:    txt,
		Scissor: scissor,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// ClipStrokeText implements the GoDevice interface
func (dev *ListDevice) ClipStrokeText(txt *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	cmd := clipstroketextcmd{
		Matrix:  matrix,
		Text:    txt,
		Stroke:  stroke,
		Scissor: scissor,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// IgnoreText implements the GoDevice interface
func (dev *ListDevice) IgnoreText(txt *Text, matrix gfx.Matrix) {
	cmd := ignoretextcmd{
		Matrix: matrix,
		Text:   txt,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// PopClip implements the GoDevice interface
func (dev *ListDevice) PopClip() {
	dev.displayList.commands = append(dev.displayList.commands, &popclipcmd{})
}

// BeginMask implements the GoDevice interface
func (dev *ListDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	cmd := beginmaskcmd{
		Rect:       rect,
		Luminosity: luminosity,
		Color:      color,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// EndMask implements the GoDevice interface
func (dev *ListDevice) EndMask() {
	dev.displayList.commands = append(dev.displayList.commands, &endmaskcmd{})
}

// BeginGroup implements the GoDevice interface
func (dev *ListDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
	cmd := begingroupcmd{
		Rect:      rect,
		Isolated:  isolated,
		Knockout:  knockout,
		BlendMode: blendmode,
		Alpha:     alpha,
	}
	dev.displayList.commands = append(dev.displayList.commands, &cmd)
}

// EndGroup implements the GoDevice interface
func (dev *ListDevice) EndGroup() {
	dev.displayList.commands = append(dev.displayList.commands, &endgroupcmd{})
}

// BeginTile implements the GoDevice interface
func (dev *ListDevice) BeginTile() int {
	dev.displayList.commands = append(dev.displayList.commands, &begintilecmd{})
	return 0
}

// EndTile implements the GoDevice interface
func (dev *ListDevice) EndTile() {
	dev.displayList.commands = append(dev.displayList.commands, &endtilecmd{})
}

// BeginLayer implements the GoDevice interface
func (dev *ListDevice) BeginLayer(layerName string) {
	dev.displayList.commands = append(dev.displayList.commands, &beginlayercmd{Name: layerName})
}

// EndLayer implements the GoDevice interface
func (dev *ListDevice) EndLayer() {
	dev.displayList.commands = append(dev.displayList.commands, &endlayercmd{})
}

// Close implements the GoDevice interface
func (dev *ListDevice) Close() {
	dev.displayList.commands = append(dev.displayList.commands, &closecmd{})
}

func (dev *ListDevice) Drop() {}

type fillpathcmd struct {
	Matrix   gfx.Matrix
	Path     *gfx.Path
	FillRule gfx.FillRule
	Color    color.Color
}

type strokepathcmd struct {
	Matrix gfx.Matrix
	Path   *gfx.Path
	Stroke *gfx.Stroke
	Color  color.Color
}

type fillshadecmd struct {
	Matrix gfx.Matrix
	Shader *gfx.Shader
	Alpha  float64
}

type fillimagecmd struct {
	Matrix gfx.Matrix
	Image  *Image
	Alpha  float64
}

type fillimagemaskcmd struct {
	Matrix gfx.Matrix
	Image  *Image
	Color  color.Color
}

type clippathcmd struct {
	Matrix   gfx.Matrix
	Path     *gfx.Path
	FillRule gfx.FillRule
	Scissor  gfx.Rect
}

type clipstrokepathcmd struct {
	Matrix  gfx.Matrix
	Path    *gfx.Path
	Stroke  *gfx.Stroke
	Scissor gfx.Rect
}

type clipimagemaskcmd struct {
	Matrix  gfx.Matrix
	Image   *Image
	Scissor gfx.Rect
}

type filltextcmd struct {
	Matrix gfx.Matrix
	Text   *Text
	Color  color.Color
}

type stroketextcmd struct {
	Matrix gfx.Matrix
	Text   *Text
	Stroke *gfx.Stroke
	Color  color.Color
}

type cliptextcmd struct {
	Matrix  gfx.Matrix
	Text    *Text
	Scissor gfx.Rect
}

type clipstroketextcmd struct {
	Matrix  gfx.Matrix
	Text    *Text
	Stroke  *gfx.Stroke
	Scissor gfx.Rect
}

type ignoretextcmd struct {
	Matrix gfx.Matrix
	Text   *Text
}

type popclipcmd struct{}

type beginmaskcmd struct {
	Rect       gfx.Rect
	Color      color.Color
	Luminosity int
}

type endmaskcmd struct{}

type begingroupcmd struct {
	Rect       gfx.Rect
	Colorspace *gfx.Colorspace
	Isolated   bool
	Knockout   bool
	BlendMode  gfx.BlendMode
	Alpha      float64
}

type endgroupcmd struct{}
type begintilecmd struct{}
type endtilecmd struct{}
type beginlayercmd struct{ Name string }
type endlayercmd struct{}
type closecmd struct{}
