package fitz_test

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"testing"

	"go.matteson.dev/fitz"
	"go.matteson.dev/no/x/urlx"
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
			doc.LoadPage(i)
		}
		doc.Close()

		fmt.Println(item)
	}

	fmt.Println("hi")
}
