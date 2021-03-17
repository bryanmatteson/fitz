package fitz

import (
	"sync"
)

type FontStyle byte

// Styles
const (
	FontStyleNormal FontStyle = iota
	FontStyleBold   FontStyle = 1 << iota
	FontStyleItalic
)

type FontFamily byte

// Families
const (
	FontFamilySans FontFamily = iota
	FontFamilySerif
	FontFamilyMono
)

type Font struct {
	Name   string
	Style  FontStyle
	Family FontFamily
}

func (f Font) IsBold() bool { return f.Style&FontStyleBold == FontStyleBold }

func (f Font) IsItalic() bool { return f.Style&FontStyleItalic == FontStyleItalic }

type FontFileNamer interface {
	GetFontName(name string, style FontStyle, family FontFamily) string
}

type FontFileNamerFunc func(name string, style FontStyle, family FontFamily) string

func (fn FontFileNamerFunc) GetFontName(name string, style FontStyle, family FontFamily) string {
	return fn(name, style, family)
}

type FontCacher struct {
	sync.RWMutex
	fonts map[string]*Font
	namer FontFileNamer
}

func NewFontCache(namer FontFileNamer) *FontCacher {
	return &FontCacher{namer: namer, fonts: make(map[string]*Font)}
}

// Store a font to this cache
func (cache *FontCacher) Store(font *Font) {
	cache.Lock()
	cache.fonts[cache.namer.GetFontName(font.Name, font.Style, font.Family)] = font
	cache.Unlock()
}

// Load a font from cache if exists otherwise it will load the font from file
func (cache *FontCacher) Load(name string, style FontStyle, family FontFamily) (font *Font) {
	cache.RLock()
	font = cache.fonts[cache.namer.GetFontName(name, style, family)]
	cache.RUnlock()

	return font
}

func defaultFontNamer(name string, style FontStyle, family FontFamily) string {
	switch family {
	case FontFamilySans:
		name += "-Sans"
	case FontFamilySerif:
		name += "-Serif"
	case FontFamilyMono:
		name += "-Mono"
	}

	if style&FontStyleBold != 0 {
		name += "-Bold"
	} else {
		name += "-Regular"
	}

	if style&FontStyleItalic != 0 {
		name += "Italic"
	}
	return name
}

var fontCache *FontCacher = NewFontCache(FontFileNamerFunc(defaultFontNamer))

func RegisterFont(font *Font) {
	if fontCache.Load(font.Name, font.Style, font.Family) == nil {
		fontCache.Store(font)
	}
}

func GetFont(name string, style FontStyle, family FontFamily) *Font {
	return fontCache.Load(name, style, family)
}
