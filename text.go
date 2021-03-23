package fitz

import (
	"image/color"
	"unicode"

	"github.com/ahmetb/go-linq"
	"go.matteson.dev/gfx"
)

// Writing modes
const (
	WModeHorizontal int = iota
	WModeVertical
)

type Letters []Letter

func (l Letters) IsWhitespace() bool {
	return linq.From(l).All(func(i interface{}) bool { return i.(*Letter).IsWhitespace() })
}

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
	Letters Letters
	Quad    gfx.Quad
}

type Text struct {
	Spans []*TextSpan
}
