package fitz

import (
	"strings"

	"go.matteson.dev/gfx"
)

type TextSpan struct {
	FontData *gfx.FontData
	WMode    int
	Letters  gfx.Letters
	Quad     gfx.Quad
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
