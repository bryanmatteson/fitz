package fitz_test

import (
	"os"
	"testing"

	"go.matteson.dev/fitz"
)

func TestOutput(t *testing.T) {
	doc, err := fitz.NewDocument("/Users/bryan/Desktop/export.pdf")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("/Users/bryan/Desktop/export1.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	doc.Write(file, fitz.DefaultWriteOptions())
}
