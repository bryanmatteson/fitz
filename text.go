package fitz

import (
	"image/color"
	"unicode"

	"github.com/ahmetb/go-linq"
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
	GlyphPath     *Path
	GlyphBounds   Rect
	Font          *Font
	Size          float64
	Color         color.Color
	Quad          Quad
	StartBaseline Point
	EndBaseline   Point
}

func (l Letter) IsWhitespace() bool { return unicode.IsSpace(l.Rune) }

type TextSpan struct {
	Font    *Font
	WMode   int
	Letters Letters
	Quad    Quad
}

type Text struct {
	Spans []*TextSpan
}
