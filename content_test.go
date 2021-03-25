package fitz_test

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lucasb-eyer/go-colorful"
	"go.matteson.dev/fitz"
	"go.matteson.dev/gfx"
)

func TestContent(t *testing.T) {
	doc, err := fitz.NewDocument("/Volumes/SamT5/backup/misc/reports/BIO/bio_quick_view_report/1487927b-ca2e-4121-b0ea-918731053f28.pdf")
	if err != nil {
		t.Fatal(err)
	}

	doc.ParallelPageProcess(func(p *fitz.Page, err error) {
		if err != nil {
			t.Fatal(err)
		}

		img, err := p.RenderImage(p.Bounds(), 5)
		if err != nil {
			t.Fatal(err)
		}

		var displayList fitz.DisplayList
		p.RunDevice(fitz.NewListDevice(&displayList))

		var content fitz.PageContent
		displayList.Apply(fitz.NewContentDevice(&content,
			fitz.WithOCR().
				WithMinConfidence(35).
				WithMinLetterWidth(5).
				WithNonImageAreas(p.Bounds(), gfx.MakeRectWH(0, 0, p.Bounds().Width(), 100)),
		))

		ctx := draw2dimg.NewGraphicContext(img)
		ctx.SetMatrixTransform(draw2d.NewScaleMatrix(5, 5))

		for _, block := range content.Blocks {
			ctx.MoveTo(block.Quad.BottomLeft.X, block.Quad.BottomLeft.Y)
			ctx.LineTo(block.Quad.TopLeft.X, block.Quad.TopLeft.Y)
			ctx.LineTo(block.Quad.TopRight.X, block.Quad.TopRight.Y)
			ctx.LineTo(block.Quad.BottomRight.X, block.Quad.BottomRight.Y)
			ctx.Close()
			ctx.SetStrokeColor(colorful.FastHappyColor())
			ctx.SetLineWidth(1.0)
			ctx.Stroke()
		}

		file, err := os.Create(fmt.Sprintf("/Users/bryan/Desktop/output%d.png", p.Number()))
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		if err = png.Encode(file, img); err != nil {
			t.Fatal(err)
		}
	})
}
