package fitz

// #include "bridge.h"
import "C"
import (
	"bytes"
	"errors"
	"io"
	"log"
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export gooutput_writer_write
func gooutput_writer_write(ctx *C.fz_context, state *C.void, data C.cvoidptr_t, length C.size_t) {
	output := pointer.Restore(unsafe.Pointer(state)).(io.Writer)
	buffer := C.GoBytes(unsafe.Pointer(data), C.int(length))
	output.Write(buffer)
}

//export gooutput_writer_close
func gooutput_writer_close(ctx *C.fz_context, state *C.void) {
	output := pointer.Restore(unsafe.Pointer(state)).(io.Writer)
	if closer, ok := output.(io.Closer); ok {
		closer.Close()
	}
}

//export gooutput_writer_tell
func gooutput_writer_tell(ctx *C.fz_context, state *C.void) int64 {
	output := pointer.Restore(unsafe.Pointer(state)).(io.WriteSeeker)
	cur, err := output.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Printf("%v", err)
	}
	return cur
}

//export gooutput_writer_seek
func gooutput_writer_seek(ctx *C.fz_context, state *C.void, offset C.int64_t, whence C.int) {
	output := pointer.Restore(unsafe.Pointer(state)).(io.WriteSeeker)
	_, err := output.Seek(int64(offset), int(whence))
	if err != nil {
		log.Printf("%v", err)
	}
}

//export gooutput_writer_drop
func gooutput_writer_drop(ctx *C.fz_context, state *C.void) {
	pointer.Unref(unsafe.Pointer(state))
}

func newOutputForWriter(ctx *C.fz_context, bufferSize int, w io.WriteSeeker) *C.fz_output {
	ref := pointer.Save(w)
	return C.fzgo_new_output_writer(ctx, C.int(bufferSize), ref)
}

type WriterSeeker struct {
	buf bytes.Buffer
	pos int
}

func (ws *WriterSeeker) Bytes() []byte { return ws.buf.Bytes() }

func (ws *WriterSeeker) Write(p []byte) (n int, err error) {
	if extra := ws.pos - ws.buf.Len(); extra > 0 {
		if _, err := ws.buf.Write(make([]byte, extra)); err != nil {
			return n, err
		}
	}

	if ws.pos < ws.buf.Len() {
		n = copy(ws.buf.Bytes()[ws.pos:], p)
		p = p[n:]
	}

	if len(p) > 0 {
		var bn int
		bn, err = ws.buf.Write(p)
		n += bn
	}

	ws.pos += n
	return n, err
}

func (ws *WriterSeeker) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = ws.pos + offs
	case io.SeekEnd:
		newPos = ws.buf.Len() + offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	ws.pos = newPos
	return int64(newPos), nil
}

func (ws *WriterSeeker) Reader() io.Reader {
	return bytes.NewReader(ws.buf.Bytes())
}

func (ws *WriterSeeker) Close() error {
	return nil
}

func (ws *WriterSeeker) BytesReader() *bytes.Reader {
	return bytes.NewReader(ws.buf.Bytes())
}
