package fitz

// #cgo CXXFLAGS: -fno-rtti -fpic -std=c++14
// #cgo darwin linux CFLAGS: -I${SRCDIR}/deps/include/ -I${SRCDIR}/deps/mupdf/thirdparty/freetype/include
// #cgo LDFLAGS: -lmupdf -lm -lmupdf-third -stdlib=libc++ -lstdc++
// #cgo linux LDFLAGS: -L${SRCDIR}/deps/lib/linux_x64
// #cgo darwin LDFLAGS: -L${SRCDIR}/deps/lib/darwin_x64
import "C"

import (
	_ "github.com/bryanmatteson/fitz/deps/include"
	_ "github.com/bryanmatteson/fitz/deps/include/mupdf"
	_ "github.com/bryanmatteson/fitz/deps/include/mupdf/fitz"
	_ "github.com/bryanmatteson/fitz/deps/include/mupdf/helpers"
	_ "github.com/bryanmatteson/fitz/deps/include/mupdf/pdf"
	_ "github.com/bryanmatteson/fitz/deps/lib/darwin_x64"
	_ "github.com/bryanmatteson/fitz/deps/lib/linux_x64"
)
