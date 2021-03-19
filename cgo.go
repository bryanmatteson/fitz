package fitz

// #cgo CXXFLAGS: -fno-rtti -fpic -std=c++14
// #cgo darwin linux CXXFLAGS: -I${SRCDIR}/deps/include
// #cgo LDFLAGS: -lmupdf -lm -lmupdf-third -stdlib=libc++ -lstdc++
// #cgo linux LDFLAGS: -L${SRCDIR}/deps/lib/linux_x64
// #cgo darwin LDFLAGS: -L${SRCDIR}/deps/lib/darwin_x64
import "C"

import (
	_ "go.matteson.dev/fitz/deps/include"
	_ "go.matteson.dev/fitz/deps/include/mupdf"
	_ "go.matteson.dev/fitz/deps/include/mupdf/fitz"
	_ "go.matteson.dev/fitz/deps/include/mupdf/helpers"
	_ "go.matteson.dev/fitz/deps/include/mupdf/pdf"
	_ "go.matteson.dev/fitz/deps/lib/darwin_x64"
	_ "go.matteson.dev/fitz/deps/lib/linux_x64"
)
