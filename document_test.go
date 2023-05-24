package fitz_test

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bryanmatteson/fitz"

	"github.com/bryanmatteson/gfx"
)

func TestDocumentReader(t *testing.T) {
	file, err := os.Open("/Users/bryan/Desktop/scratch/mdt3.pdf")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := fitz.NewDocument(file)
	if err != nil {
		t.Fatal(err)
	}
	doc.LoadPage(0)
	doc.Close()
}

func TestDocumentFont(t *testing.T) {
	doc, err := fitz.NewDocumentFromFile("/Users/bryan/Desktop/scratch/mdt3.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer doc.Close()

	doc.SequentialPageProcess(func(pg *fitz.Page) {
		trm := gfx.NewScaleMatrix(3, 3)
		bounds := trm.TransformRect(pg.Bounds())
		img := image.NewRGBA(bounds.ImageRect())

		pg.RunDevice(fitz.NewDrawDevice(gfx.NewScaleMatrix(3, 3), img))
		var buf bytes.Buffer
		png.Encode(&buf, img)
		ioutil.WriteFile(fmt.Sprintf("/Users/bryan/Desktop/mdt3/test%d.png", pg.Number()), buf.Bytes(), os.ModePerm)
	})
}

func TestDocumentSplit(t *testing.T) {
	doc, err := fitz.NewDocumentFromFile("/Users/bryan/Desktop/scratch/mdt3.pdf")
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
