package fitz

// #include "bridge.h"
import "C"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"unsafe"
)

type PageRange struct {
	Start, End int
}

type Document struct {
	sync.Mutex
	ctx    *C.fz_context
	native *C.pdf_document
	pages  map[int]*Page
}

func NewDocument(source interface{}) (*Document, error) {
	switch t := source.(type) {
	case string:
		return newDocumentFromFile(t)
	case []byte:
		return newDocumentFromBytes(t)
	case io.Reader:
		return newDocumentFromReader(t)
	default:
		return nil, ErrUnknownSource
	}
}

func (d *Document) loadFont(num int) {
	ref := C.pdf_new_indirect(d.ctx, d.native, C.int(num), 0)
	defer C.pdf_drop_obj(d.ctx, ref)

	if isFontDesc(d.ctx, ref) {
		var stream *C.pdf_obj
		var ext string

		obj := C.pdf_dict_get(d.ctx, ref, pdfName(C.PDF_ENUM_NAME_FontFile))
		if obj != nil {
			stream = obj
			ext = "pfa"
		}

		obj = C.pdf_dict_get(d.ctx, ref, pdfName(C.PDF_ENUM_NAME_FontFile2))
		if obj != nil {
			stream = obj
			ext = "ttf"
		}

		obj = C.pdf_dict_get(d.ctx, ref, pdfName(C.PDF_ENUM_NAME_FontFile3))
		if obj != nil {
			stream = obj
			obj = C.pdf_dict_get(d.ctx, obj, pdfName(C.PDF_ENUM_NAME_Subtype))
			if obj != nil && C.pdf_is_name(d.ctx, obj) != 0 {
				log.Printf("invalid font descriptor subtype")
				return
			}

			if C.pdf_name_eq(d.ctx, obj, pdfName(C.PDF_ENUM_NAME_Type1C)) != 0 {
				ext = "cff"
			} else if C.pdf_name_eq(d.ctx, obj, pdfName(C.PDF_ENUM_NAME_CIDFontType0C)) != 0 {
				ext = "cid"
			} else if C.pdf_name_eq(d.ctx, obj, pdfName(C.PDF_ENUM_NAME_OpenType)) != 0 {
				ext = "otf"
			} else {
				log.Printf("unhandled font type %s", C.GoString(C.pdf_to_name(d.ctx, obj)))
				return
			}
		}

		if stream == nil {
			return
		}

		var data *C.uchar
		buf := C.pdf_load_stream(d.ctx, stream)
		defer C.fz_drop_buffer(d.ctx, buf)

		buflen := C.fz_buffer_storage(d.ctx, buf, &data)
		num := C.pdf_to_num(d.ctx, ref)

		fontData := C.GoBytes(unsafe.Pointer(data), C.int(buflen))
		name := fmt.Sprintf("font-%04d.%s", num, ext)

		fmt.Printf("font %s: %d", name, len(fontData))
	}

	C.fz_empty_store(d.ctx)
}

func (d *Document) LoadFonts() {
	numObj := int(C.pdf_count_objects(d.ctx, d.native))
	for i := 1; i < numObj; i++ {
		d.loadFont(i)
	}
}

func (d *Document) NumPages() int {
	d.Lock()
	defer d.Unlock()
	return int(C.fz_count_pages(d.ctx, &d.native.super))
}

func (d *Document) LoadPage(num int) (*Page, error) {
	d.Lock()
	defer d.Unlock()

	if int(C.fz_count_pages(d.ctx, &d.native.super)) <= num {
		return nil, ErrInvalidPage
	}

	if _, ok := d.pages[num]; !ok {
		pg := C.fz_load_page(d.ctx, &d.native.super, C.int(num))
		defer C.fz_drop_page(d.ctx, pg)

		list := C.fz_new_display_list_from_page(d.ctx, pg)
		bounds := C.fz_bound_page(d.ctx, pg)
		number := pg.number
		d.pages[num] = newPage(C.fz_clone_context(d.ctx), number, bounds, list)
	}

	return d.pages[num], nil
}

func (d *Document) ParallelPageProcess(fn func(p *Page, err error)) {
	var wg sync.WaitGroup
	pageCount := d.NumPages()
	wg.Add(pageCount)

	for i := 0; i < pageCount; i++ {
		go func(p *Page, err error) {
			defer wg.Done()
			fn(p, err)
		}(d.LoadPage(i))
	}

	wg.Wait()
}

func (d *Document) SequentialPageProcess(fn func(p *Page, err error)) {
	pageCount := d.NumPages()

	for i := 0; i < pageCount; i++ {
		fn(d.LoadPage(i))
	}
}

// Close closes the underlying fitz document.
func (d *Document) Close() {
	for _, pg := range d.pages {
		pg.drop()
	}
	d.pages = nil
	C.pdf_drop_document(d.ctx, d.native)
	C.fz_drop_context(d.ctx)
}

func (d *Document) ExtractPageRanges(filePath string, ranges ...PageRange) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	file.Close()

	dest := C.pdf_create_document(d.ctx)
	defer C.pdf_drop_document(d.ctx, dest)

	graftMap := C.pdf_new_graft_map(d.ctx, dest)
	defer C.pdf_drop_graft_map(d.ctx, graftMap)

	pageCount := d.NumPages()

	for _, rng := range ranges {
		if rng.Start < 0 || rng.End > pageCount {
			return errors.New("invalid page range")
		}

		C.pdf_graft_mapped_page(d.ctx, graftMap, C.int(rng.End), d.native, C.int(rng.Start))
	}

	options := C.pdf_write_options{}
	output := C.CString(filePath)
	defer C.free(unsafe.Pointer(output))

	C.pdf_save_document(d.ctx, dest, output, &options)
	return nil
}

func newDocument(ctx *C.fz_context, doc *C.pdf_document) *Document {
	return &Document{
		ctx:    ctx,
		native: doc,
		pages:  make(map[int]*Page),
	}
}

func newDocumentFromReader(r io.Reader) (*Document, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, ErrOpenReader
	}

	return newDocumentFromBytes(b)
}

func newDocumentFromBytes(b []byte) (d *Document, err error) {
	ctx := C.fzgo_new_context()
	if err != nil {
		err = ErrCreateContext
		return
	}

	C.fz_register_document_handlers(ctx)

	data := (*C.uchar)(C.CBytes(b))
	stream := C.fz_open_memory(ctx, data, C.size_t(len(b)))
	if stream == nil {
		err = ErrOpenMemory
		return
	}

	native := C.pdf_open_document_with_stream(ctx, stream)
	if native == nil {
		err = ErrOpenDocument
		return
	}

	ret := C.pdf_needs_password(ctx, native)
	if bool(int(ret) != 0) {
		err = ErrNeedsPassword
		C.pdf_drop_document(ctx, native)
		C.fz_drop_context(ctx)
		return nil, err
	}

	return newDocument(ctx, native), nil
}

func newDocumentFromFile(fileName string) (d *Document, err error) {
	fileName, err = filepath.Abs(fileName)
	if err != nil {
		return
	}

	if _, e := os.Stat(fileName); e != nil {
		err = ErrNoSuchFile
		return
	}

	ctx := C.fzgo_new_context()
	if err != nil {
		err = ErrCreateContext
		return
	}

	C.fz_register_document_handlers(ctx)

	fname := C.CString(fileName)
	defer C.free(unsafe.Pointer(fname))

	native := C.pdf_open_document(ctx, fname)
	if native == nil {
		err = ErrOpenDocument
		return
	}

	ret := C.pdf_needs_password(ctx, native)
	if bool(int(ret) != 0) {
		err = ErrNeedsPassword
		C.pdf_drop_document(ctx, native)
		C.fz_drop_context(ctx)
		return nil, err
	}

	return newDocument(ctx, native), nil
}
