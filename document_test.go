package fitz_test

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"go.matteson.dev/fitz"
	"go.matteson.dev/gfx"
	"go.matteson.dev/no/x/urlx"
)

func TestDocumentFont(t *testing.T) {
	doc, err := fitz.NewDocument("/Users/bryan/Desktop/scratch/mdt3.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer doc.Close()

	pg, _ := doc.LoadPage(0)
	trm := gfx.NewScaleMatrix(3, 3)
	bounds := trm.TransformRect(pg.Bounds())
	img := image.NewRGBA(bounds.ImageRect())

	pg.RunDevice(fitz.NewDebugDevice(gfx.NewScaleMatrix(3, 3), img))

	var buf bytes.Buffer
	png.Encode(&buf, img)
	ioutil.WriteFile("/Users/bryan/Desktop/test.png", buf.Bytes(), os.ModePerm)
}

func TestDocumentSplit(t *testing.T) {
	doc, err := fitz.NewDocument("/Users/bryan/Desktop/scratch/mdt3.pdf")
	if err != nil {
		t.Fatal(err)
	}

	newDoc, err := doc.NewDocumentFromPages(3)
	if err != nil {
		t.Fatal(err)
	}
	defer newDoc.Close()

	if newDoc.NumPages() != 1 {
		t.Fail()
	}
	newDoc.Save("/Users/bryan/Desktop/testoutput.pdf", fitz.DefaultWriteOptions())
}

func TestDocumentMemory(t *testing.T) {
	uri, err := url.Parse("/Volumes/SamT5/reports")
	if err != nil {
		t.Fatal(err)
	}

	items, err := urlx.GetFiles(context.Background(), uri, true, func(s string) bool { return filepath.Ext(s) == ".pdf" })
	if err != nil {
		t.Fatal(err)
	}

	for item := range items {
		doc, err := fitz.NewDocument(item)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < doc.NumPages(); i++ {
			pg, _ := doc.LoadPage(i)
			pg.RenderImage(pg.Bounds(), 5)
			pg.GetText()
			var displayList fitz.DisplayList
			pg.RunDevice(fitz.NewListDevice(&displayList))
		}
		doc.Close()

		fmt.Println(item)
	}

	fmt.Println("hi")
}
