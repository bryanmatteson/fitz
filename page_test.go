package fitz_test

import (
	"image/png"
	"os"
	"testing"

	"go.matteson.dev/fitz"
)

func TestPageRender(t *testing.T) {
	doc, err := fitz.NewDocument("/Users/bryan/Desktop/scratch/stj.pdf")
	if err != nil {
		t.Fatal(err)
	}
	pg, err := doc.LoadPage(0)
	if err != nil {
		t.Fatal(err)
	}

	img, err := pg.RenderImage(fitz.MakeRectWH(50, 50, 160, 30), 4)
	// img, err := pg.RenderImage(fitz.EmptyRect(), 1)
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
