package fitz_test

import (
	"testing"

	"go.matteson.dev/fitz"
)

func TestDocumentSplit(t *testing.T) {
	doc, err := fitz.NewDocument("/Users/bryan/Desktop/scratch/mdt3.pdf")
	if err != nil {
		t.Fatal(err)
	}

	newDoc, err := doc.NewDocumentFromPages(3)
	if err != nil {
		t.Fatal(err)
	}

	if newDoc.NumPages() != 1 {
		t.Fail()
	}
	newDoc.Save("/Users/bryan/Desktop/testoutput.pdf")
}
