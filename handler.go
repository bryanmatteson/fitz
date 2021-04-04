package fitz

import (
	"image/color"

	"go.matteson.dev/gfx"
)

type HandlerFunc interface{}

type HandlerDevice struct {
	BaseDevice
	fillPathHandlers       []FillPathHandler
	strokePathHandlers     []StrokePathHandler
	fillShadeHandlers      []FillShadeHandler
	fillImageHandlers      []FillImageHandler
	fillImageMaskHandlers  []FillImageMaskHandler
	clipPathHandlers       []ClipPathHandler
	clipStrokePathHandlers []ClipStrokePathHandler
	clipImageMaskHandlers  []ClipImageMaskHandler
	fillTextHandlers       []FillTextHandler
	strokeTextHandlers     []StrokeTextHandler
	clipTextHandlers       []ClipTextHandler
	clipStrokeTextHandlers []ClipStrokeTextHandler
	ignoreTextHandlers     []IgnoreTextHandler
	popClipHandlers        []PopClipHandler
	beginMaskHandlers      []BeginMaskHandler
	endMaskHandlers        []EndMaskHandler
	beginGroupHandlers     []BeginGroupHandler
	endGroupHandlers       []EndGroupHandler
	beginTileHandlers      []BeginTileHandler
	endTileHandlers        []EndTileHandler
	beginLayerHandlers     []BeginLayerHandler
	endLayerHandlers       []EndLayerHandler
	closeHandlers          []CloseHandler
}

func (dev *HandlerDevice) AddHandler(handlers ...HandlerFunc) {
	for _, handler := range handlers {
		switch h := handler.(type) {
		case FillPathHandler:
			dev.fillPathHandlers = append(dev.fillPathHandlers, h)
		case StrokePathHandler:
			dev.strokePathHandlers = append(dev.strokePathHandlers, h)
		case FillShadeHandler:
			dev.fillShadeHandlers = append(dev.fillShadeHandlers, h)
		case FillImageHandler:
			dev.fillImageHandlers = append(dev.fillImageHandlers, h)
		case FillImageMaskHandler:
			dev.fillImageMaskHandlers = append(dev.fillImageMaskHandlers, h)
		case ClipPathHandler:
			dev.clipPathHandlers = append(dev.clipPathHandlers, h)
		case ClipStrokePathHandler:
			dev.clipStrokePathHandlers = append(dev.clipStrokePathHandlers, h)
		case ClipImageMaskHandler:
			dev.clipImageMaskHandlers = append(dev.clipImageMaskHandlers, h)
		case FillTextHandler:
			dev.fillTextHandlers = append(dev.fillTextHandlers, h)
		case StrokeTextHandler:
			dev.strokeTextHandlers = append(dev.strokeTextHandlers, h)
		case ClipTextHandler:
			dev.clipTextHandlers = append(dev.clipTextHandlers, h)
		case ClipStrokeTextHandler:
			dev.clipStrokeTextHandlers = append(dev.clipStrokeTextHandlers, h)
		case IgnoreTextHandler:
			dev.ignoreTextHandlers = append(dev.ignoreTextHandlers, h)
		case PopClipHandler:
			dev.popClipHandlers = append(dev.popClipHandlers, h)
		case BeginMaskHandler:
			dev.beginMaskHandlers = append(dev.beginMaskHandlers, h)
		case EndMaskHandler:
			dev.endMaskHandlers = append(dev.endMaskHandlers, h)
		case BeginGroupHandler:
			dev.beginGroupHandlers = append(dev.beginGroupHandlers, h)
		case EndGroupHandler:
			dev.endGroupHandlers = append(dev.endGroupHandlers, h)
		case BeginTileHandler:
			dev.beginTileHandlers = append(dev.beginTileHandlers, h)
		case EndTileHandler:
			dev.endTileHandlers = append(dev.endTileHandlers, h)
		case BeginLayerHandler:
			dev.beginLayerHandlers = append(dev.beginLayerHandlers, h)
		case EndLayerHandler:
			dev.endLayerHandlers = append(dev.endLayerHandlers, h)
		case CloseHandler:
			dev.closeHandlers = append(dev.closeHandlers, h)
		default:
			panic("unknown handler, must be a handler")
		}
	}
}

func (dev *HandlerDevice) ShouldCall(kind CommandKind) bool {
	switch kind {
	case FillPath:
		return len(dev.fillPathHandlers) > 0
	case StrokePath:
		return len(dev.strokePathHandlers) > 0
	case FillShade:
		return len(dev.fillShadeHandlers) > 0
	case FillImage:
		return len(dev.fillImageHandlers) > 0
	case FillImageMask:
		return len(dev.fillImageMaskHandlers) > 0
	case ClipPath:
		return len(dev.clipPathHandlers) > 0
	case ClipStrokePath:
		return len(dev.clipStrokePathHandlers) > 0
	case ClipImageMask:
		return len(dev.clipImageMaskHandlers) > 0
	case FillText:
		return len(dev.fillTextHandlers) > 0
	case StrokeText:
		return len(dev.strokeTextHandlers) > 0
	case ClipText:
		return len(dev.clipTextHandlers) > 0
	case ClipStrokeText:
		return len(dev.clipStrokeTextHandlers) > 0
	case IgnoreText:
		return len(dev.ignoreTextHandlers) > 0
	case PopClip:
		return len(dev.popClipHandlers) > 0
	case BeginMask:
		return len(dev.beginMaskHandlers) > 0
	case EndMask:
		return len(dev.endMaskHandlers) > 0
	case BeginGroup:
		return len(dev.beginGroupHandlers) > 0
	case EndGroup:
		return len(dev.endGroupHandlers) > 0
	case BeginTile:
		return len(dev.beginTileHandlers) > 0
	case EndTile:
		return len(dev.endTileHandlers) > 0
	case BeginLayer:
		return len(dev.beginLayerHandlers) > 0
	case EndLayer:
		return len(dev.endLayerHandlers) > 0
	case CloseDevice:
		return len(dev.closeHandlers) > 0
	default:
		return false
	}
}

func (dev *HandlerDevice) FillPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.fillPathHandlers {
		device.FillPath(path, fillRule, matrix, color)
	}
}

func (dev *HandlerDevice) StrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.strokePathHandlers {
		device.StrokePath(path, stroke, matrix, color)
	}
}

func (dev *HandlerDevice) FillShade(shade *gfx.Shader, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.fillShadeHandlers {
		device.FillShade(shade, matrix, alpha)
	}
}

func (dev *HandlerDevice) FillImage(image *Image, matrix gfx.Matrix, alpha float64) {
	for _, device := range dev.fillImageHandlers {
		device.FillImage(image, matrix, alpha)
	}
}

func (dev *HandlerDevice) FillImageMask(image *Image, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.fillImageMaskHandlers {
		device.FillImageMask(image, matrix, color)
	}
}

func (dev *HandlerDevice) ClipPath(path *gfx.Path, fillRule gfx.FillRule, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.clipPathHandlers {
		device.ClipPath(path, fillRule, matrix, scissor)
	}
}

func (dev *HandlerDevice) ClipStrokePath(path *gfx.Path, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.clipStrokePathHandlers {
		device.ClipStrokePath(path, stroke, matrix, scissor)
	}
}

func (dev *HandlerDevice) ClipImageMask(image *Image, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.clipImageMaskHandlers {
		device.ClipImageMask(image, matrix, scissor)
	}
}

func (dev *HandlerDevice) FillText(text *Text, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.fillTextHandlers {
		device.FillText(text, matrix, color)
	}
}

func (dev *HandlerDevice) StrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, color color.Color) {
	for _, device := range dev.strokeTextHandlers {
		device.StrokeText(text, stroke, matrix, color)
	}
}

func (dev *HandlerDevice) ClipText(text *Text, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.clipTextHandlers {
		device.ClipText(text, matrix, scissor)
	}
}

func (dev *HandlerDevice) ClipStrokeText(text *Text, stroke *gfx.Stroke, matrix gfx.Matrix, scissor gfx.Rect) {
	for _, device := range dev.clipStrokeTextHandlers {
		device.ClipStrokeText(text, stroke, matrix, scissor)
	}
}

func (dev *HandlerDevice) IgnoreText(text *Text, matrix gfx.Matrix) {
	for _, device := range dev.ignoreTextHandlers {
		device.IgnoreText(text, matrix)
	}
}

func (dev *HandlerDevice) PopClip() {
	for _, device := range dev.popClipHandlers {
		device.PopClip()
	}
}

func (dev *HandlerDevice) BeginMask(rect gfx.Rect, color color.Color, luminosity int) {
	for _, device := range dev.beginMaskHandlers {
		device.BeginMask(rect, color, luminosity)
	}
}

func (dev *HandlerDevice) EndMask() {
	for _, device := range dev.endMaskHandlers {
		device.EndMask()
	}
}

func (dev *HandlerDevice) BeginGroup(rect gfx.Rect, cs *gfx.Colorspace, isolated bool, knockout bool, blendMode gfx.BlendMode, alpha float64) {
	for _, device := range dev.beginGroupHandlers {
		device.BeginGroup(rect, cs, isolated, knockout, blendMode, alpha)
	}
}

func (dev *HandlerDevice) EndGroup() {
	for _, device := range dev.endGroupHandlers {
		device.EndGroup()
	}
}

func (dev *HandlerDevice) BeginTile() int {
	for _, device := range dev.beginTileHandlers {
		device.BeginTile()
	}
	return 0
}

func (dev *HandlerDevice) EndTile() {
	for _, device := range dev.endTileHandlers {
		device.EndTile()
	}
}

func (dev *HandlerDevice) BeginLayer(layerName string) {
	for _, device := range dev.beginLayerHandlers {
		device.BeginLayer(layerName)
	}
}

func (dev *HandlerDevice) EndLayer() {
	for _, device := range dev.endLayerHandlers {
		device.EndLayer()
	}
}

func (dev *HandlerDevice) Close() {
	for _, device := range dev.closeHandlers {
		device.Close()
	}
}
