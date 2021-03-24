package fitz

import (
	"image/color"

	"go.matteson.dev/fitz/internal/dla"
	"go.matteson.dev/gfx"
	"go.matteson.dev/tess"
)

type PageContent struct {
	Blocks  gfx.TextBlocks
	Paths   []*gfx.Path
	Strokes []*gfx.Path
	Images  []*Image
}

type ContentDevice struct {
	BaseDevice
	content *PageContent
	ocr     *tess.Client
	letters gfx.Chars
	words   gfx.TextWords
}

func NewContentDevice(content *PageContent) GoDevice {
	return &ContentDevice{content: content, ocr: tess.NewClient()}
}

func (dev *ContentDevice) ShouldCall(kind CommandKind) bool {
	handles := FillText | FillImage | CloseDevice
	return handles.Has(kind)
}

func (dev *ContentDevice) FillText(text *Text, ctm gfx.Matrix, fillColor color.Color) {
	for _, span := range text.Spans {
		for _, letter := range span.Letters {
			dev.letters = append(dev.letters, gfx.Char{
				Rune:          letter.Rune,
				Quad:          letter.Quad,
				StartBaseline: letter.StartBaseline,
				EndBaseline:   letter.EndBaseline,
				Confidence:    100,
				Orientation:   letter.Quad.Orientation(),
				DeskewAngle:   letter.Quad.T(),
			})
		}
	}
}

func (dev *ContentDevice) FillPath(path *gfx.Path, fillRule FillRule, matrix gfx.Matrix, color color.Color) {
	dev.content.Paths = append(dev.content.Paths, path)
}

func (dev *ContentDevice) StrokePath(path *gfx.Path, stroke *Stroke, matrix gfx.Matrix, color color.Color) {
	dev.content.Strokes = append(dev.content.Strokes, path)
}

func (dev *ContentDevice) FillImage(img *Image, matrix gfx.Matrix, alpha float64) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	mat := matrix.Inverted().PostScaled(float64(width), float64(height)).Inverted()
	dev.ocr.SetImageFromFileData(img.PngBytes())

	words, _ := dev.ocr.GetWords()
	for _, word := range words {
		if word.Confidence < 35 || word.Quad.Width() < 5 {
			continue
		}
		word.Quad = mat.TransformQuad(word.Quad)
		word.StartBaseline = mat.TransformPoint(word.StartBaseline)
		word.EndBaseline = mat.TransformPoint(word.EndBaseline)
		dev.words = append(dev.words, word)
	}
	dev.content.Images = append(dev.content.Images, img)
}

func (dev *ContentDevice) Close() {
	extractor := dla.NewNearestNeighborWordExtractor()
	words := append(extractor.GetWords(dla.RemoveOverlappingLetters(dev.letters)), dev.words...)
	docstrum := dla.NewDocstrumBoundingBoxPageSegmenter()
	blocks := docstrum.GetBlocks(words)
	dev.content.Blocks = append(dev.content.Blocks, blocks...)
	dev.ocr.Close()
}
