package fitz

import (
	"image/color"
	"unicode"

	"go.matteson.dev/gfx"
)

// Writing modes
const (
	WModeHorizontal int = iota
	WModeVertical
)

type Letter struct {
	Rune          rune
	GlyphPath     *gfx.Path
	GlyphBounds   gfx.Rect
	Font          *Font
	Size          float64
	Color         color.Color
	Quad          gfx.Quad
	StartBaseline gfx.Point
	EndBaseline   gfx.Point
}

func (l Letter) IsWhitespace() bool { return unicode.IsSpace(l.Rune) }

type TextSpan struct {
	Font    *Font
	WMode   int
	Letters []Letter
	Quad    gfx.Quad
}

type Text struct {
	Spans []*TextSpan
}
