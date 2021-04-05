package fitz

// #include "bridge.h"
import "C"
import (
	"bytes"
	"io"
	"log"
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export fzgo_read_stream_next
func fzgo_read_stream_next(ctx *C.fz_context, stm *C.fz_stream, max C.size_t) C.int {
	stream := pointer.Restore(unsafe.Pointer(stm.state)).(*inputstream)
	n, err := stream.rd.Read(stream.buf)

	if err != nil {
		log.Printf("%v", err)
		return -1
	}

	if n < 0 {
		return -1
	}

	stm.rp = (*C.uchar)(unsafe.Pointer(&stream.buf[1]))
	stm.wp = (*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(&stream.buf[0])) + uintptr(n)))
	stm.pos += C.int64_t(n)

	return C.int(stream.buf[0])
}

//export fzgo_read_stream_seek
func fzgo_read_stream_seek(ctx *C.fz_context, stm *C.fz_stream, offset C.int64_t, whence C.int) {
	stream := pointer.Restore(unsafe.Pointer(stm.state)).(*inputstream)
	off, err := stream.rd.Seek(int64(offset), int(whence))
	if err != nil {
		panic(err)
	}

	stm.pos = C.int64_t(off)
	stm.rp = (*C.uchar)(unsafe.Pointer(&stream.buf[0]))
	stm.wp = (*C.uchar)(unsafe.Pointer(&stream.buf[0]))
}

//export fzgo_read_stream_drop
func fzgo_read_stream_drop(ctx *C.fz_context, state unsafe.Pointer) {
	stream := pointer.Restore(state).(*inputstream)
	if closer, ok := stream.rd.(io.Closer); ok {
		closer.Close()
	}
	pointer.Unref(state)
}

func newStreamFromBytes(ctx *C.fz_context, bufferSize int, b []byte) (*C.fz_stream, error) {
	stream := &inputstream{
		rd:  bytes.NewReader(b),
		buf: make([]byte, bufferSize),
	}
	return C.fzgo_new_read_stream(ctx, pointer.Save(stream)), nil
}

func newStreamFromReader(ctx *C.fz_context, bufferSize int, r io.ReadSeeker) (*C.fz_stream, error) {
	stream := &inputstream{
		rd:  r,
		buf: make([]byte, bufferSize),
	}
	return C.fzgo_new_read_stream(ctx, pointer.Save(stream)), nil
}

type inputstream struct {
	buf []byte
	rd  io.ReadSeeker
}
