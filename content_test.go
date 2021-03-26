package fitz_test

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"go.matteson.dev/fitz"
	"go.matteson.dev/gfx"
)

func TestContent(t *testing.T) {
	doc, err := fitz.NewDocument("/Volumes/SamT5/backup/misc/reports/BIO/bio_quick_view_report/1487927b-ca2e-4121-b0ea-918731053f28.pdf")
	if err != nil {
		t.Fatal(err)
	}

	doc.SequentialPageProcess(func(p *fitz.Page, err error) {
		if err != nil {
			t.Fatal(err)
		}

		if p.Number() > 1 {
			return
		}

		img, err := p.RenderImage(p.Bounds(), 5)
		if err != nil {
			t.Fatal(err)
		}

		var displayList fitz.DisplayList
		p.RunDevice(fitz.NewListDevice(&displayList))

		var content fitz.PageContent
		displayList.Apply(fitz.NewContentDevice(&content))

		ctx := draw2dimg.NewGraphicContext(img)
		ctx.SetMatrixTransform(draw2d.NewScaleMatrix(5, 5))

		for _, stroke := range content.Paths {
			drawRect(ctx, stroke.Bounds())
			ctx.SetStrokeColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
			ctx.Stroke()
		}

		var buf bytes.Buffer
		png.Encode(&buf, img)

		ioutil.WriteFile(fmt.Sprintf("/Users/bryan/Desktop/test%d.png", p.Number()), buf.Bytes(), os.ModePerm)
	})
}

func drawRect(ctx *draw2dimg.GraphicContext, rect gfx.Rect) {
	ctx.MoveTo(rect.X.Min, rect.Y.Min)
	ctx.LineTo(rect.X.Min, rect.Y.Max)
	ctx.LineTo(rect.X.Max, rect.Y.Max)
	ctx.LineTo(rect.X.Max, rect.Y.Min)
	ctx.Close()
}
