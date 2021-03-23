package fitz

// #include "bridge.h"
import "C"
import (
	"image/color"
	"unsafe"

	"go.matteson.dev/gfx"
)

func rectFromFitz(rect C.fz_rect) gfx.Rect {
	return gfx.MakeRectCorners(float64(rect.x0), float64(rect.y0), float64(rect.x1), float64(rect.y1))
}

func matrixFromFitz(trm C.fz_matrix) gfx.Matrix {
	return gfx.NewMatrix(float64(trm.a), float64(trm.b), float64(trm.c), float64(trm.d), float64(trm.e), float64(trm.f))
}

func isFontDesc(ctx *C.fz_context, obj *C.pdf_obj) bool {
	typ := C.pdf_dict_get(ctx, obj, pdfName(C.PDF_ENUM_NAME_Type))
	return C.pdf_name_eq(ctx, typ, pdfName(C.PDF_ENUM_NAME_FontDescriptor)) != 0
}

func pdfName(ename C.int) *C.pdf_obj {
	return C.pdfname(ename)
}

func getRGBColor(ctx *C.fz_context, col *C.float, colorspace *C.fz_colorspace, alpha C.float, params C.fz_color_params) color.NRGBA {
	var rgb [3]C.float
	if C.fz_colorspace_is_rgb(ctx, colorspace) == 0 {
		C.fz_convert_color(ctx, colorspace, col, C.fz_device_rgb(ctx), (*C.float)(unsafe.Pointer(&rgb[0])), nil, params)
	} else {
		rgb = *(*[3]C.float)(unsafe.Pointer(col))
	}

	return color.NRGBA{
		R: byte(C.fz_clampi(C.int(rgb[0]*255), 0, 255)),
		G: byte(C.fz_clampi(C.int(rgb[1]*255), 0, 255)),
		B: byte(C.fz_clampi(C.int(rgb[2]*255), 0, 255)),
		A: byte(C.fz_clampi(C.int(alpha*255), 0, 255)),
	}
}

func getStroke(stroke *C.fz_stroke_state) *Stroke {
	dashes := make([]float64, int(stroke.dash_len))
	for i := 0; i < int(stroke.dash_len); i++ {
		dashes[i] = float64(stroke.dash_list[i])
	}

	return &Stroke{
		StartCap:   LineCap(stroke.start_cap),
		DashCap:    LineCap(stroke.dash_cap),
		EndCap:     LineCap(stroke.end_cap),
		LineJoin:   LineJoin(stroke.linejoin),
		LineWidth:  float64(stroke.linewidth),
		MiterLimit: float64(stroke.miterlimit),
		DashPhase:  float64(stroke.dash_phase),
		Dashes:     dashes,
	}
}

func getImage(ctx *C.fz_context, ctm C.fz_matrix, img *C.fz_image, colorParams C.fz_color_params) *Image {
	pix := C.fz_get_pixmap_from_image(ctx, img, nil, nil, nil, nil)
	cs := C.fz_pixmap_colorspace(ctx, pix)

	switch C.fz_colorspace_type(ctx, cs) {
	case C.FZ_COLORSPACE_RGB, C.FZ_COLORSPACE_NONE:
		break
	default:
		pix = C.fz_convert_pixmap(ctx, pix, C.fz_device_rgb(ctx), nil, nil, colorParams, 1)
	}

	comp := int(C.fz_pixmap_components(ctx, pix))
	stride := int(C.fz_pixmap_stride(ctx, pix))
	height := int(C.fz_pixmap_height(ctx, pix))
	width := int(C.fz_pixmap_width(ctx, pix))
	x := int(C.fz_pixmap_x(ctx, pix))
	y := int(C.fz_pixmap_y(ctx, pix))
	pixels := C.fz_pixmap_samples(ctx, pix)
	data := C.GoBytes(unsafe.Pointer(pixels), C.int(stride*height))
	bounds := C.fz_transform_rect(C.fz_unit_rect, ctm)

	return &Image{
		Rect:    gfx.MakeRectCorners(float64(x), float64(y), float64(x+width), float64(y+height)),
		Frame:   rectFromFitz(bounds),
		Data:    data,
		Stride:  stride,
		NumComp: comp,
	}
}

func getTextInfo(ctx *C.fz_context, fztext *C.fz_text, ctm C.fz_matrix, col color.Color) (text *Text) {
	text = &Text{}
	for span := fztext.head; span != nil; span = span.next {
		wmode := C.fz_text_span_wmode(span)
		font := getFontInfo(ctx, span)
		bbox := C.fz_font_bbox(ctx, span.font)

		spanMat := matrixFromFitz(span.trm)
		letters := make(Letters, 0, span.len)
		quads := make(gfx.Quads, 0, span.len)

		for i := 0; i < int(span.len); i++ {
			item := (*C.fz_text_item)(unsafe.Pointer(uintptr(unsafe.Pointer(span.items)) + uintptr(i)*unsafe.Sizeof(*span.items)))
			if item.ucs == -1 {
				continue
			}

			trm := spanMat.Translated(float64(item.x), float64(item.y)).Concat(matrixFromFitz(ctm))

			var adv float64 = 0
			if item.gid >= 0 {
				adv = float64(C.fz_advance_glyph(ctx, span.font, item.gid, wmode))
			}

			var dir, p, q, a, d gfx.Point
			if wmode == 0 {
				dir.X, dir.Y = 1, 0
			} else {
				dir.X, dir.Y = 0, -1
			}

			dir = trm.TransformVec(dir)
			size := trm.Expansion()

			if wmode == 0 {
				p.X, p.Y = trm.E, trm.F
				q.X, q.Y = trm.E+adv*dir.X, trm.F+adv*dir.Y
				a.X, a.Y = 0, float64(C.fz_font_ascender(ctx, span.font))
				d.X, d.Y = 0, float64(C.fz_font_descender(ctx, span.font))
			} else {
				q.X, q.Y = trm.E, trm.F
				p.X, p.Y = trm.E-adv*dir.X, trm.F-adv*dir.Y
				a.X, a.Y = float64(bbox.x1), 0
				d.X, d.Y = float64(bbox.x0), 0
			}

			a = trm.TransformVec(a)
			d = trm.TransformVec(d)

			quad := gfx.Quad{
				BottomLeft:  gfx.Point{X: p.X + d.X, Y: p.Y + d.Y},
				TopLeft:     gfx.Point{X: p.X + a.X, Y: p.Y + a.Y},
				BottomRight: gfx.Point{X: q.X + d.X, Y: q.Y + d.Y},
				TopRight:    gfx.Point{X: q.X + a.X, Y: q.Y + a.Y},
			}

			mat := C.fz_matrix{C.float(trm.A), C.float(trm.B), C.float(trm.C), C.float(trm.D), C.float(trm.E), C.float(trm.F)}

			gb := C.fz_bound_glyph(ctx, span.font, item.gid, mat)
			glyphBounds := gfx.MakeRectCorners(
				float64(gb.x0), float64(gb.y0),
				float64(gb.x1), float64(gb.y1),
			)
			glyphPath := C.fz_outline_glyph(ctx, span.font, item.gid, mat)
			quad = gfx.RectToQuad(glyphBounds)
			quads = append(quads, quad)

			letter := Letter{
				Rune:          rune(item.ucs),
				Font:          font,
				GlyphPath:     makePath(ctx, glyphPath),
				Quad:          quad,
				Size:          size,
				Color:         col,
				StartBaseline: p,
				EndBaseline:   q,
				GlyphBounds:   glyphBounds,
			}
			C.fz_drop_path(ctx, glyphPath)

			letters = append(letters, letter)
		}

		sp := &TextSpan{
			Font:    font,
			WMode:   int(wmode),
			Letters: letters,
			Quad:    quads.Union(),
		}

		text.Spans = append(text.Spans, sp)
	}
	return
}

func getFontInfo(ctx *C.fz_context, span *C.fz_text_span) (font *Font) {
	fontFamily := FontFamilySans
	if C.fz_font_is_serif(ctx, span.font) != 0 {
		fontFamily = FontFamilySerif
	} else if C.fz_font_is_monospaced(ctx, span.font) != 0 {
		fontFamily = FontFamilyMono
	}
	fontStyle := FontStyleNormal
	if C.fz_font_is_bold(ctx, span.font) != 0 {
		fontStyle |= FontStyleBold
	} else if C.fz_font_is_italic(ctx, span.font) != 0 {
		fontStyle |= FontStyleItalic
	}

	fontName := C.GoString(C.fz_font_name(ctx, span.font))

	font = GetFont(fontName, fontStyle, fontFamily)
	if font == nil {
		font = &Font{
			Name:   fontName,
			Family: fontFamily,
			Style:  fontStyle,
		}
		RegisterFont(font)
	}
	return
}
