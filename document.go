package fitz

// #include "bridge.h"
import "C"

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/bryanmatteson/gfx"
	"github.com/mattn/go-pointer"
)

type Document struct {
	mut    sync.Mutex
	ctx    *C.fz_context
	native *C.pdf_document
	pages  map[int]*Page
}

func (d *Document) GetFontCache() gfx.FontCache {
	return pointer.Restore(unsafe.Pointer(d.ctx.user)).(*usercontext).fontCache
}

func (d *Document) NumPages() int {
	d.mut.Lock()
	defer d.mut.Unlock()
	return int(C.fz_count_pages(d.ctx, &d.native.super))
}

func (d *Document) LoadPage(num int) (*Page, error) {
	d.mut.Lock()
	defer d.mut.Unlock()

	if int(C.fz_count_pages(d.ctx, &d.native.super)) <= num {
		return nil, ErrInvalidPage
	}

	if _, ok := d.pages[num]; !ok {
		d.pages[num] = newPage(d.native, d.ctx, num)
	}

	return d.pages[num], nil
}

// func (d *Document) ParallelPageProcess(fn func(p *Page)) {
// 	var wg sync.WaitGroup
// 	pageCount := d.NumPages()
// 	wg.Add(pageCount)

// 	for i := 0; i < pageCount; i++ {
// 		p, err := d.LoadPage(i)
// 		if err != nil {
// 			log.Printf("%v", err)
// 			continue
// 		}
// 		go func(p *Page) {
// 			defer wg.Done()
// 			fn(p)
// 		}(p)
// 	}

// 	wg.Wait()
// }

func (d *Document) SequentialPageProcess(fn func(p *Page)) {
	pageCount := d.NumPages()

	for i := 0; i < pageCount; i++ {
		p, err := d.LoadPage(i)
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		fn(p)
	}
}

// Close closes the underlying fitz document.
func (d *Document) Close() {
	for _, pg := range d.pages {
		pg.drop()
	}
	d.pages = nil
	C.pdf_drop_document(d.ctx, d.native)

	if d.ctx.user != nil {
		pointer.Unref(unsafe.Pointer(d.ctx.user))
		d.ctx.user = nil
	}

	C.fz_drop_context(d.ctx)
}

func (d *Document) Save(filePath string, opts WriteOptions) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	file.Close()

	options := opts.fzoptions()
	output := C.CString(filePath)
	defer C.free(unsafe.Pointer(output))

	C.pdf_save_document(d.ctx, d.native, output, &options)
	return nil
}

func (d *Document) Write(w io.Writer, opts WriteOptions) error {
	options := opts.fzoptions()
	output := newOutputForWriter(d.ctx, 8192, w)
	defer C.fz_drop_output(d.ctx, output)

	C.pdf_write_document(d.ctx, d.native, output, &options)
	C.fz_close_output(d.ctx, output)

	return nil
}

func (d *Document) NewDocumentFromPages(pages ...int) (*Document, error) {
	dest := C.pdf_create_document(d.ctx)

	graftMap := C.pdf_new_graft_map(d.ctx, dest)
	defer C.pdf_drop_graft_map(d.ctx, graftMap)

	for _, pg := range pages {
		C.pdf_graft_mapped_page(d.ctx, graftMap, C.int(-1), d.native, C.int(pg))
	}

	return newDocument(C.fz_clone_context(d.ctx), dest), nil
}

func newDocument(ctx *C.fz_context, doc *C.pdf_document) *Document {
	return &Document{
		ctx:    ctx,
		native: doc,
		pages:  make(map[int]*Page),
	}
}

func NewDocumentFromBytes(b []byte) (d *Document, err error) {
	return NewDocument(bytes.NewReader(b))
}

func NewDocument(r io.Reader) (d *Document, err error) {
	ctx := C.fzgo_new_context()

	if err != nil {
		err = ErrCreateContext
		return
	}

	C.fz_register_document_handlers(ctx)

	var stream *C.fz_stream
	if rds, ok := r.(io.ReadSeeker); ok {
		stream, err = newStreamFromReader(ctx, 16384, rds)
	} else {
		var b []byte
		if b, err = ioutil.ReadAll(r); err == nil {
			stream, err = newStreamFromBytes(ctx, 16384, b)
		}
	}

	if err != nil {
		return nil, err
	}

	defer C.fz_drop_stream(ctx, stream)

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

	userCtx := newusercontext()
	userCtx.fontCache.init(ctx, native)
	ctx.user = pointer.Save(userCtx)

	return newDocument(ctx, native), nil
}

func NewDocumentFromFile(fileName string) (d *Document, err error) {
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

	userCtx := newusercontext()
	userCtx.fontCache.init(ctx, native)
	ctx.user = pointer.Save(userCtx)

	return newDocument(ctx, native), nil
}

type WriteOptions struct {
	CompressImages         bool
	CompressFonts          bool
	CompressStreams        bool
	DecompressStreams      bool
	CleanStreams           bool
	SanitizeStreams        bool
	Linearize              bool
	DontRegenerateID       bool
	GarbageCollectionLevel int
}

func DefaultWriteOptions() WriteOptions { return WriteOptions{} }

func (o *WriteOptions) fzoptions() C.pdf_write_options {
	opts := C.pdf_write_options{}
	if o == nil {
		return opts
	}

	if o.DontRegenerateID {
		opts.dont_regenerate_id = 1
	}

	if o.CompressImages {
		opts.do_compress_images = 1
	}

	if o.CompressFonts {
		opts.do_compress_fonts = 1
	}

	if o.CompressStreams {
		opts.do_compress = 1
	}

	if o.DecompressStreams {
		opts.do_decompress = 1
	}

	if o.CleanStreams {
		opts.do_clean = 1
	}

	if o.SanitizeStreams {
		opts.do_sanitize = 1
	}

	if o.Linearize {
		opts.do_linear = 1
	}

	opts.do_garbage = C.int(o.GarbageCollectionLevel)
	return opts
}
