package fitz_test

import (
	"bytes"
	"fmt"
	"image"
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

		ctx := draw2dimg.NewGraphicContext(img)
		ctx.SetMatrixTransform(draw2d.NewScaleMatrix(5, 5))

		drawRect(ctx, gfx.MakeRectWH(100, 0, p.Bounds().Width()-100, 100))
		ctx.SetFillColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		ctx.Fill()

		var buf bytes.Buffer
		png.Encode(&buf, img)

		ioutil.WriteFile(fmt.Sprintf("/Users/bryan/Desktop/test%d.png", p.Number()), buf.Bytes(), os.ModePerm)
	})
}

func TestRects(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 500, 100))
	ctx := draw2dimg.NewGraphicContext(img)
	drawRect(ctx, gfx.MakeRectWH(0, 0, 500, 100))
	ctx.SetFillColor(color.Black)
	ctx.Fill()

	trm := gfx.NewTranslationMatrix(50, 50).Inverted()
	ctx.SetMatrixTransform(draw2d.Matrix{trm.A, trm.B, trm.C, trm.D, trm.E, trm.F})

	drawRect(ctx, gfx.MakeRectWH(50, 50, 100, 50))
	ctx.SetFillColor(color.White)
	ctx.Fill()

	var buf bytes.Buffer
	png.Encode(&buf, img)
	ioutil.WriteFile("/Users/bryan/Desktop/out.png", buf.Bytes(), os.ModePerm)
}

func drawRect(ctx *draw2dimg.GraphicContext, rect gfx.Rect) {
	ctx.MoveTo(rect.X.Min, rect.Y.Min)
	ctx.LineTo(rect.X.Min, rect.Y.Max)
	ctx.LineTo(rect.X.Max, rect.Y.Max)
	ctx.LineTo(rect.X.Max, rect.Y.Min)
	ctx.Close()
}
