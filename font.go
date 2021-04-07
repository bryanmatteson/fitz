package fitz

// #include "bridge.h"
import "C"
import (
	"log"
	"sync"
	"unsafe"

	"go.matteson.dev/gfx"
)

type fitzFont struct {
	mut  sync.Mutex
	ctx  *C.fz_context
	font *C.fz_font
	info gfx.FontData
}

func newfitzfont(ctx *C.fz_context, font *C.fz_font) gfx.Font {
	fontFamily := gfx.FontFamilySans
	if C.fz_font_is_serif(ctx, font) != 0 {
		fontFamily = gfx.FontFamilySerif
	} else if C.fz_font_is_monospaced(ctx, font) != 0 {
		fontFamily = gfx.FontFamilyMono
	}

	fontStyle := gfx.FontStyleNormal
	if C.fz_font_is_bold(ctx, font) != 0 {
		fontStyle |= gfx.FontStyleBold
	} else if C.fz_font_is_italic(ctx, font) != 0 {
		fontStyle |= gfx.FontStyleItalic
	}

	fontName := C.GoString(C.fz_font_name(ctx, font))
	info := gfx.FontData{
		Name:   fontName,
		Style:  fontStyle,
		Family: fontFamily,
	}

	return &fitzFont{
		ctx:  ctx,
		font: font,
		info: info,
	}
}

func (f *fitzFont) Info() gfx.FontData { return f.info }
func (f *fitzFont) Name() string       { return f.info.Name }

func (f *fitzFont) BoundingBox() gfx.Rect {
	f.mut.Lock()
	defer f.mut.Unlock()

	bbox := C.fz_font_bbox(f.ctx, f.font)
	return gfx.MakeRect(float64(bbox.x0), float64(bbox.y0), float64(bbox.x1), float64(bbox.y1))
}

func (f *fitzFont) Glyph(chr rune, trm gfx.Matrix) *gfx.Glyph {
	f.mut.Lock()
	defer f.mut.Unlock()

	mat := C.fz_matrix{C.float(trm.A), C.float(trm.B), C.float(trm.C), C.float(trm.D), C.float(trm.E), C.float(trm.F)}
	gid := int(C.fz_encode_character(f.ctx, f.font, C.int(chr)))
	hi := int(*(*C.ushort)(unsafe.Pointer(uintptr(unsafe.Pointer(f.font.encoding_cache[0])) + uintptr(chr))))
	_ = hi

	width := float64(C.fz_advance_glyph(f.ctx, f.font, C.int(gid), 0))
	glyphPath := C.fz_outline_glyph(f.ctx, f.font, C.int(gid), mat)

	defer C.fz_drop_path(f.ctx, glyphPath)
	path := convertPath(f.ctx, glyphPath)

	return &gfx.Glyph{
		Path:  path,
		Width: width,
	}
}

func (f *fitzFont) Advance(chr rune, mode int) float64 {
	f.mut.Lock()
	defer f.mut.Unlock()

	gid := int(C.fz_encode_character(f.ctx, f.font, C.int(chr)))
	return float64(C.fz_advance_glyph(f.ctx, f.font, C.int(gid), C.int(mode)))
}

type fontCache struct {
	mut   sync.Mutex
	fonts map[string]gfx.Font
}

func newfontcache() *fontCache {
	return &fontCache{fonts: make(map[string]gfx.Font)}
}

func (fc *fontCache) init(ctx *C.fz_context, doc *C.pdf_document, page *C.pdf_page) {
	fc.mut.Lock()
	defer fc.mut.Unlock()

	var fonts []*C.pdf_obj
	rsrc := C.pdf_page_resources(ctx, page)
	fontObj := C.pdf_dict_get(ctx, rsrc, pdfName(C.PDF_ENUM_NAME_Font))

	if fontObj == nil {
		return
	}

	n := int(C.pdf_dict_len(ctx, fontObj))
	for i := 0; i < n; i++ {
		fontDict := C.pdf_dict_get_val(ctx, fontObj, C.int(i))
		if C.pdf_is_dict(ctx, fontDict) == 0 {
			log.Printf("not a font dict (%d 0 R)", int(C.pdf_to_num(ctx, fontDict)))
			continue
		}

		found := false
		for _, f := range fonts {
			if C.pdf_objcmp(ctx, f, fontDict) == 0 {
				found = true
				break
			}
		}
		if found {
			continue
		}

		fonts = append(fonts, fontDict)

		desc := C.pdf_load_font(ctx, doc, rsrc, fontDict)
		if desc == nil {
			desc = C.pdf_load_hail_mary_font(ctx, doc)
		}

		font := newfitzfont(ctx, desc.font)
		key := font.Info().String()
		if _, ok := fc.fonts[key]; !ok {
			fc.fonts[key] = font
		}
	}
}

func (fc *fontCache) Load(fontData gfx.FontData) (gfx.Font, error) {
	fc.mut.Lock()
	defer fc.mut.Unlock()

	key := fontData.String()
	if font, ok := fc.fonts[key]; ok {
		return font, nil
	}

	return nil, nil
}

func (fc *fontCache) Store(font gfx.Font) {
	fc.mut.Lock()
	defer fc.mut.Unlock()

	key := font.Info().String()
	if _, ok := fc.fonts[key]; !ok {
		fc.fonts[key] = font
	}
}
