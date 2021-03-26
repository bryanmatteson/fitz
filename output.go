package fitz

// #include "bridge.h"
import "C"
import (
	"io"
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export gooutput_writer_write
func gooutput_writer_write(ctx *C.fz_context, state *C.void, data C.cvoidptr_t, length C.size_t) {
	output := pointer.Restore(unsafe.Pointer(state)).(io.WriteCloser)
	buffer := C.GoBytes(unsafe.Pointer(data), C.int(length))
	output.Write(buffer)
}

//export gooutput_writer_close
func gooutput_writer_close(ctx *C.fz_context, state *C.void) {
	output := pointer.Restore(unsafe.Pointer(state)).(io.WriteCloser)
	output.Close()
}

//export gooutput_writer_drop
func gooutput_writer_drop(ctx *C.fz_context, state *C.void) {
	pointer.Unref(unsafe.Pointer(state))
}

func newOutputForWriter(ctx *C.fz_context, bufferSize int, w io.WriteCloser) *C.fz_output {
	ref := pointer.Save(w)
	return C.fzgo_new_output_writer(ctx, C.int(bufferSize), ref)
}
