package fitz

import (
	"strings"

	"github.com/bryanmatteson/gfx"
)

type TextSpan struct {
	Letters Letters
	Font    gfx.Font
	Matrix  gfx.Matrix
	WMode   int
	Quad    gfx.Quad
}

func (s *TextSpan) String() string {
	var builder strings.Builder
	for _, letter := range s.Letters {
		builder.WriteRune(letter.Rune)
	}
	return builder.String()
}

type Text struct {
	FontCache gfx.FontCache
	Spans     []*TextSpan
}

// Writing modes
const (
	WModeHorizontal int = iota
	WModeVertical
)

type Letters []Letter
type Letter struct {
	Rune    rune
	GlyphID int
	Origin  gfx.Point
	Quad    gfx.Quad
}
