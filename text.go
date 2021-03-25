package fitz

import (
	"image/color"
	"strings"
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

func (s *TextSpan) String() string {
	var builder strings.Builder
	for _, letter := range s.Letters {
		builder.WriteRune(letter.Rune)
	}
	return builder.String()
}

type Text struct {
	Spans []*TextSpan
}
