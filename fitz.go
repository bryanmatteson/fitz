package fitz

// #include "bridge.h"
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

// Version returns the version of mupdf
func Version() string {
	return C.GoString(C.fz_version)
}

var _locks []sync.Mutex = make([]sync.Mutex, C.fz_lock_max)

//export lock_mutex
func lock_mutex(user unsafe.Pointer, lock C.int) {
	_locks[int(lock)].Lock()
}

//export unlock_mutex
func unlock_mutex(user unsafe.Pointer, lock C.int) {
	_locks[int(lock)].Unlock()
}

var (
	ErrUnknownSource = errors.New("fitz: unknown source")
	ErrNoSuchFile    = errors.New("fitz: no such file")
	ErrCreateContext = errors.New("fitz: cannot create context")
	ErrOpenDocument  = errors.New("fitz: cannot open document")
	ErrOpenMemory    = errors.New("fitz: cannot open memory")
	ErrOpenReader    = errors.New("fitz: cannot read from reader")
	ErrPageMissing   = errors.New("fitz: page missing")
	ErrCreatePixmap  = errors.New("fitz: cannot create pixmap")
	ErrPixmapSamples = errors.New("fitz: cannot get pixmap samples")
	ErrNeedsPassword = errors.New("fitz: document needs password")
	ErrLoadOutline   = errors.New("fitz: cannot load outline")
	ErrInvalidPage   = errors.New("fitz: cannot load page")
)

//export exception_callback
func exception_callback(code C.int, message *C.char) {
	panic(fmt.Errorf("[PANIC:%d]: %s", code, C.GoString(message)))
}

//export error_callback
func error_callback(userData unsafe.Pointer, message *C.char) {
	panic(fmt.Errorf("[ERR]: %s", C.GoString(message)))
}

//export warn_callback
func warn_callback(userData unsafe.Pointer, message *C.char) {
	fmt.Printf("[WARN]: %s\n", C.GoString(message))
}

func catch(err *error) {
	if r := recover(); r != nil {
		switch t := r.(type) {
		case string:
			*err = errors.New(t)
		case error:
			*err = t
		default:
			*err = fmt.Errorf("%v", t)
		}
	}
}
