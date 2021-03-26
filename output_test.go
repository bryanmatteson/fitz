package fitz_test

import (
	"os"
	"testing"

	"go.matteson.dev/fitz"
)

func TestOutput(t *testing.T) {
	doc, err := fitz.NewDocument("/Volumes/SamT5/backup/misc/reports/BIO/bio_quick_view_report/1487927b-ca2e-4121-b0ea-918731053f28.pdf")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("/Users/bryan/Desktop/export.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	doc.Write(file, fitz.DefaultWriteOptions())
}
