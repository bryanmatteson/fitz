package fitz

import (
	"image"
	"image/color"

	"github.com/bryanmatteson/gfx"
)

type CommandList struct {
	FillPaths       []*FillPathCmd
	StrokePaths     []*StrokePathCmd
	FillShades      []*FillShadeCmd
	FillImages      []*FillImageCmd
	FillImageMasks  []*FillImageMaskCmd
	ClipPaths       []*ClipPathCmd
	ClipStrokePaths []*ClipStrokePathCmd
	ClipImageMasks  []*ClipImageMaskCmd
	FillTexts       []*FillTextCmd
	StrokeTexts     []*StrokeTextCmd
	ClipTexts       []*ClipTextCmd
	ClipStrokeTexts []*ClipStrokeTextCmd
	IgnoreTexts     []*IgnoreTextCmd
	BeginMasks      []*BeginMaskCmd
	BeginGroups     []*BegingGoupCmd
	BeginLayers     []*BeginLayerCmd
}

type CommandDevice struct {
	BaseDevice
	commandList *CommandList
}

func NewCommandDevice(cl *CommandList) Device {
	return &CommandDevice{commandList: cl}
}

func (dev *CommandDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, trm gfx.Matrix, color color.Color) {
	cmd := &FillPathCmd{Transform: trm, Path: path, FillRule: fillRule, Color: color}
	dev.commandList.FillPaths = append(dev.commandList.FillPaths, cmd)
}

func (dev *CommandDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, trm gfx.Matrix, color color.Color) {
	cmd := &StrokePathCmd{Transform: trm, Path: path, Stroke: stroke, Color: color}
	dev.commandList.StrokePaths = append(dev.commandList.StrokePaths, cmd)
}

func (dev *CommandDevice) FillShade(shade *gfx.Shader, trm gfx.Matrix, alpha float64) {
	cmd := &FillShadeCmd{Transform: trm, Alpha: float64(alpha), Shader: shade}
	dev.commandList.FillShades = append(dev.commandList.FillShades, cmd)
}

func (dev *CommandDevice) FillImage(img image.Image, trm gfx.Matrix, alpha float64) {
	cmd := &FillImageCmd{Transform: trm, Image: img, Alpha: float64(alpha)}
	dev.commandList.FillImages = append(dev.commandList.FillImages, cmd)
}

func (dev *CommandDevice) FillImageMask(img image.Image, trm gfx.Matrix, color color.Color) {
	cmd := &FillImageMaskCmd{Transform: trm, Image: img, Color: color}
	dev.commandList.FillImageMasks = append(dev.commandList.FillImageMasks, cmd)
}

func (dev *CommandDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := &ClipPathCmd{Transform: trm, Path: path, FillRule: fillRule, Scissor: scissor}
	dev.commandList.ClipPaths = append(dev.commandList.ClipPaths, cmd)
}

func (dev *CommandDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := &ClipStrokePathCmd{Transform: trm, Path: path, Stroke: stroke, Scissor: scissor}
	dev.commandList.ClipStrokePaths = append(dev.commandList.ClipStrokePaths, cmd)
}

func (dev *CommandDevice) ClipImageMask(img image.Image, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := &ClipImageMaskCmd{Transform: trm, Image: img, Scissor: scissor}
	dev.commandList.ClipImageMasks = append(dev.commandList.ClipImageMasks, cmd)
}

func (dev *CommandDevice) FillText(text *Text, trm gfx.Matrix, color color.Color) {
	cmd := &FillTextCmd{Transform: trm, Text: text, Color: color}
	dev.commandList.FillTexts = append(dev.commandList.FillTexts, cmd)
}

func (dev *CommandDevice) StrokeText(text *Text, stroke *gfx.Stroke, trm gfx.Matrix, color color.Color) {
	cmd := &StrokeTextCmd{Transform: trm, Text: text, Stroke: stroke, Color: color}
	dev.commandList.StrokeTexts = append(dev.commandList.StrokeTexts, cmd)
}

func (dev *CommandDevice) ClipText(text *Text, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := &ClipTextCmd{Transform: trm, Text: text, Scissor: scissor}
	dev.commandList.ClipTexts = append(dev.commandList.ClipTexts, cmd)
}

func (dev *CommandDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, trm gfx.Matrix, scissor gfx.Rect) {
	cmd := &ClipStrokeTextCmd{Transform: trm, Text: text, Stroke: stroke, Scissor: scissor}
	dev.commandList.ClipStrokeTexts = append(dev.commandList.ClipStrokeTexts, cmd)
}

func (dev *CommandDevice) IgnoreText(text *Text, trm gfx.Matrix) {
	cmd := &IgnoreTextCmd{Transform: trm, Text: text}
	dev.commandList.IgnoreTexts = append(dev.commandList.IgnoreTexts, cmd)
}

func (dev *CommandDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	cmd := &BeginMaskCmd{Rect: rect, Luminosity: luminosity, Color: color}
	dev.commandList.BeginMasks = append(dev.commandList.BeginMasks, cmd)
}

func (dev *CommandDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendmode gfx.BlendMode, alpha float64) {
	cmd := &BegingGoupCmd{Rect: rect, Isolated: isolated, Knockout: knockout, BlendMode: blendmode, Alpha: alpha}
	dev.commandList.BeginGroups = append(dev.commandList.BeginGroups, cmd)
}

func (dev *CommandDevice) BeginLayer(layerName string) {
	dev.commandList.BeginLayers = append(dev.commandList.BeginLayers, &BeginLayerCmd{Name: layerName})
}

type FillPathCmd struct {
	Transform gfx.Matrix
	Path      *gfx.Path
	FillRule  gfx.FillRule
	Color     color.Color
}

type StrokePathCmd struct {
	Transform gfx.Matrix
	Path      *gfx.Path
	Stroke    *gfx.Stroke
	Color     color.Color
}

type FillShadeCmd struct {
	Transform gfx.Matrix
	Shader    *gfx.Shader
	Alpha     float64
}

type FillImageCmd struct {
	Transform gfx.Matrix
	Image     image.Image
	Alpha     float64
}

type FillImageMaskCmd struct {
	Transform gfx.Matrix
	Image     image.Image
	Color     color.Color
}

type ClipPathCmd struct {
	Transform gfx.Matrix
	Path      *gfx.Path
	FillRule  gfx.FillRule
	Scissor   gfx.Rect
}

type ClipStrokePathCmd struct {
	Transform gfx.Matrix
	Path      *gfx.Path
	Stroke    *gfx.Stroke
	Scissor   gfx.Rect
}

type ClipImageMaskCmd struct {
	Transform gfx.Matrix
	Image     image.Image
	Scissor   gfx.Rect
}

type FillTextCmd struct {
	Transform gfx.Matrix
	Text      *Text
	Color     color.Color
}

type StrokeTextCmd struct {
	Transform gfx.Matrix
	Text      *Text
	Stroke    *gfx.Stroke
	Color     color.Color
}

type ClipTextCmd struct {
	Transform gfx.Matrix
	Text      *Text
	Scissor   gfx.Rect
}

type ClipStrokeTextCmd struct {
	Transform gfx.Matrix
	Text      *Text
	Stroke    *gfx.Stroke
	Scissor   gfx.Rect
}

type IgnoreTextCmd struct {
	Transform gfx.Matrix
	Text      *Text
}

type BeginMaskCmd struct {
	Rect       gfx.Rect
	Color      color.Color
	Luminosity int
}

type BegingGoupCmd struct {
	Rect       gfx.Rect
	Colorspace *gfx.Colorspace
	Isolated   bool
	Knockout   bool
	BlendMode  gfx.BlendMode
	Alpha      float64
}

type PopClipCmd struct{}
type EndMaskCmd struct{}
type EndGroupCmd struct{}
type BeginTileCmd struct{}
type EndTileCmd struct{}
type BeginLayerCmd struct{ Name string }
type EndLayerCmd struct{}
