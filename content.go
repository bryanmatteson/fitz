package fitz

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/disintegration/imaging"
	"go.matteson.dev/dla"
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
	ocrEnabled bool
	ocrOpts    ocroptions
	content    *PageContent
	ocr        *tess.Client
	letters    gfx.Chars
	words      gfx.TextWords
	drawDevice GoDevice
	img        *image.RGBA
}

func NewContentDevice(content *PageContent, opts ...ContentOption) GoDevice {
	options := contentopts{}
	for _, opt := range opts {
		opt.Apply(&options)
	}

	var ocr *tess.Client
	if options.ocrEnabled {
		ocr = tess.NewClient()
	}

	var drawdev GoDevice
	var img *image.RGBA
	if options.ocropts.NonImageAreas {
		bounds := options.ocropts.PageBounds
		img = image.NewRGBA(image.Rect(int(bounds.X.Min), int(bounds.Y.Min), int(bounds.X.Max), int(bounds.Y.Max)))
		drawdev = NewDrawDevice(gfx.NewScaleMatrix(1, 1), img)
	}

	return &ContentDevice{
		content:    content,
		ocrEnabled: options.ocrEnabled,
		ocrOpts:    options.ocropts,
		ocr:        ocr,
		img:        img,
		drawDevice: drawdev,
	}
}

func (dev *ContentDevice) ShouldCall(kind CommandKind) bool {
	handles := FillText | FillImage | FillPath | StrokePath | CloseDevice
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
	if dev.ocrEnabled && dev.ocrOpts.NonImageAreas {
		dev.drawDevice.FillPath(path, fillRule, matrix, color)
	}
	dev.content.Paths = append(dev.content.Paths, path)
}

func (dev *ContentDevice) StrokePath(path *gfx.Path, stroke *Stroke, matrix gfx.Matrix, color color.Color) {
	if dev.ocrEnabled && dev.ocrOpts.NonImageAreas {
		dev.drawDevice.StrokePath(path, stroke, matrix, color)
	}
	dev.content.Strokes = append(dev.content.Strokes, path)
}

func (dev *ContentDevice) FillImage(img *Image, matrix gfx.Matrix, alpha float64) {
	if dev.ocrEnabled && img.Rect.Width() > dev.ocrOpts.MinImageSize.X && img.Rect.Height() > dev.ocrOpts.MinImageSize.Y {
		width, height := img.Bounds().Dx(), img.Bounds().Dy()
		mat := matrix.Inverted().PostScaled(float64(width), float64(height)).Inverted()
		dev.ocr.SetImageFromFileData(img.PngBytes())

		words, _ := dev.ocr.GetWords()
		for _, word := range words {
			if word.Confidence < dev.ocrOpts.MinConfidence || word.Quad.Width() < dev.ocrOpts.MinLetterWidth {
				continue
			}
			word.Quad = mat.TransformQuad(word.Quad)
			word.StartBaseline = mat.TransformPoint(word.StartBaseline)
			word.EndBaseline = mat.TransformPoint(word.EndBaseline)
			dev.words = append(dev.words, word)
		}
	}

	dev.content.Images = append(dev.content.Images, img)
}

func (dev *ContentDevice) Close() {
	if dev.ocrEnabled && dev.ocrOpts.NonImageAreas {
		dev.drawDevice.Close()
		dev.words = append(dev.words, dev.doNonImageOCR()...)
	}

	extractor := dla.NewNearestNeighborWordExtractor()
	words := append(extractor.GetWords(dla.RemoveOverlappingLetters(dev.letters)), dev.words...)

	docstrum := dla.NewDocstrumBoundingBoxPageSegmenter()
	blocks := docstrum.GetBlocks(words)
	dev.content.Blocks = append(dev.content.Blocks, blocks...)

	if dev.ocr != nil {
		dev.ocr.Close()
		dev.ocr = nil
	}
}

func (dev *ContentDevice) doNonImageOCR() gfx.TextWords {
	words := make(gfx.TextWords, 0)
	inv := gfx.NewScaleMatrix(4, 4).Inverted()

	for _, area := range dev.ocrOpts.AdditionalAreas {
		subArea := image.Rect(int(area.X.Min), int(area.Y.Min), int(area.X.Max), int(area.Y.Max))
		subimg := dev.img.SubImage(subArea)
		resized := imaging.Resize(subimg, subArea.Dx()*4, subArea.Dy()*4, imaging.Lanczos)
		var buf bytes.Buffer
		if err := png.Encode(&buf, resized); err != nil {
			log.Printf("%v", err)
			return nil
		}

		dev.ocr.SetImageFromFileData(buf.Bytes())
		ocrWords, err := dev.ocr.GetWords()
		if err != nil {
			log.Printf("%v", err)
			return nil
		}

		for _, word := range ocrWords {
			if word.Confidence < dev.ocrOpts.MinConfidence || word.Quad.Width() < dev.ocrOpts.MinLetterWidth {
				continue
			}

			w := word
			w.Quad = inv.TransformQuad(w.Quad)
			w.StartBaseline = inv.TransformPoint(w.StartBaseline)
			w.EndBaseline = inv.TransformPoint(w.EndBaseline)
			words = append(words, w)
		}
	}

	return words
}

type contentopts struct {
	ocrEnabled bool
	ocropts    ocroptions
}

type ContentOption interface{ Apply(*contentopts) }
type ContentOptionFunc func(*contentopts)

func (fn ContentOptionFunc) Apply(o *contentopts) { fn(o) }

type ocroptions struct {
	MinImageSize    gfx.Point
	MinConfidence   float64
	MinLetterWidth  float64
	NonImageAreas   bool
	PageBounds      gfx.Rect
	AdditionalAreas []gfx.Rect
}

type OCROptionBuilder struct{ options ocroptions }

func (b *OCROptionBuilder) WithMinImageSize(w, h float64) *OCROptionBuilder {
	b.options.MinImageSize = gfx.Point{X: w, Y: h}
	return b
}
func (b *OCROptionBuilder) WithMinConfidence(conf float64) *OCROptionBuilder {
	b.options.MinConfidence = conf
	return b
}
func (b *OCROptionBuilder) WithMinLetterWidth(width float64) *OCROptionBuilder {
	b.options.MinLetterWidth = width
	return b
}
func (b *OCROptionBuilder) WithNonImageAreas(page gfx.Rect, areas ...gfx.Rect) *OCROptionBuilder {
	b.options.NonImageAreas = true
	b.options.PageBounds = page
	b.options.AdditionalAreas = append(b.options.AdditionalAreas, areas...)
	return b
}

func (b *OCROptionBuilder) Apply(o *contentopts) {
	o.ocrEnabled = true
	o.ocropts = b.options
}

func WithOCR() *OCROptionBuilder {
	return &OCROptionBuilder{
		options: ocroptions{
			MinImageSize:   gfx.Point{X: 0, Y: 0},
			MinConfidence:  70,
			MinLetterWidth: 5,
		},
	}
}
