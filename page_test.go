package fitz_test

import (
	"image/png"
	"os"
	"testing"

	"go.matteson.dev/fitz"
	"go.matteson.dev/gfx"
)

func TestPageRender(t *testing.T) {
	doc, err := fitz.NewDocument("/Volumes/SamT5/backup/misc/reports/BIO/bio_quick_view_report/1487927b-ca2e-4121-b0ea-918731053f28.pdf")
	if err != nil {
		t.Fatal(err)
	}
	pg, err := doc.LoadPage(0)
	if err != nil {
		t.Fatal(err)
	}

	img, err := pg.RenderImage(gfx.MakeRectWH(360, 17, 166, 44), 4)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("/Users/bryan/Desktop/output.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
}
