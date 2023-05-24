package fitz_test

import (
	"os"
	"testing"

	"github.com/bryanmatteson/fitz"
)

func TestOutput(t *testing.T) {
	doc, err := fitz.NewDocumentFromFile("/Volumes/SamT5/backup/misc/reports/BIO/bio_quick_view_report/1487927b-ca2e-4121-b0ea-918731053f28.pdf")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("/Users/bryan/Desktop/export.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	newDoc, err := doc.NewDocumentFromPages(1, 3, 5, 7)
	if err != nil {
		t.Fatal(err)
	}
	defer newDoc.Close()

	newDoc.Write(file, fitz.WriteOptions{
		DontRegenerateID: true,
		CleanStreams:     true,
		SanitizeStreams:  true,
		Linearize:        true,

		GarbageCollectionLevel: 4,
	})
}
