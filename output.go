package fitz

// #include "bridge.h"
import "C"

import (
	"bytes"
	"errors"
	"io"
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export gooutput_writer_write
func gooutput_writer_write(ctx *C.fz_context, state unsafe.Pointer, data unsafe.Pointer, length C.size_t) {
	output := pointer.Restore(state).(*outputwriter)
	buffer := C.GoBytes(data, C.int(length))
	output.Write(buffer)
}

//export gooutput_writer_close
func gooutput_writer_close(ctx *C.fz_context, state unsafe.Pointer) {
	output := pointer.Restore(state).(*outputwriter)
	output.Destination.Write(output.Bytes())
}

//export gooutput_writer_tell
func gooutput_writer_tell(ctx *C.fz_context, state unsafe.Pointer) int64 {
	output := pointer.Restore(state).(*outputwriter)
	return int64(output.Position())
}

//export gooutput_writer_seek
func gooutput_writer_seek(ctx *C.fz_context, state unsafe.Pointer, offset C.int64_t, whence C.int) {
	output := pointer.Restore(state).(*outputwriter)
	output.Seek(int64(offset), int(whence))
}

//export gooutput_writer_drop
func gooutput_writer_drop(ctx *C.fz_context, state unsafe.Pointer) {
	pointer.Unref(state)
}

type outputwriter struct {
	*WriterSeeker
	Destination io.Writer
}

func (o *outputwriter) Close() {}

func newOutputForWriter(ctx *C.fz_context, bufferSize int, w io.Writer) *C.fz_output {
	writer := &outputwriter{
		Destination:  w,
		WriterSeeker: &WriterSeeker{},
	}
	return C.fzgo_new_output_writer(ctx, C.int(bufferSize), pointer.Save(writer))
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

func (ws *WriterSeeker) Position() int {
	return ws.pos
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
