package fitz

// #include "bridge.h"
import "C"
import (
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export fzgo_fill_path
func fzgo_fill_path(ctx *C.fz_context, dev *C.fz_device, path *C.cfz_path_t, evenOdd C.int, ctm C.fz_matrix, colorspace *C.fz_colorspace, color *C.cfloat_t, alpha C.float, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(FillPath) {
		return
	}

	p := makePath(ctx, path)
	rgb := getRGBColor(ctx, color, colorspace, alpha, colorParams)
	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))

	fillRule := FillRuleWinding
	if evenOdd != 0 {
		fillRule = FillRuleEvenOdd
	}

	device.FillPath(p, fillRule, matrix, rgb)
}

//export fzgo_stroke_path
func fzgo_stroke_path(ctx *C.fz_context, dev *C.fz_device, path *C.cfz_path_t, stroke *C.cfz_stroke_state_t, ctm C.fz_matrix, colorspace *C.fz_colorspace, color *C.cfloat_t, alpha C.float, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(StrokePath) {
		return
	}

	p := makePath(ctx, path)
	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	rgb := getRGBColor(ctx, color, colorspace, alpha, colorParams)
	s := getStroke(stroke)

	device.StrokePath(p, s, matrix, rgb)
}

//export fzgo_fill_shade
func fzgo_fill_shade(ctx *C.fz_context, dev *C.fz_device, shade *C.fz_shade, ctm C.fz_matrix, alpha C.float, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(FillShade) {
		return
	}

	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	shaderMatrix := NewMatrix(float64(shade.matrix.a), float64(shade.matrix.b), float64(shade.matrix.c), float64(shade.matrix.d), float64(shade.matrix.e), float64(shade.matrix.f))
	bounds := rectFromFitz(C.fz_bound_shade(ctx, shade, C.fz_identity))

	shader := &Shader{
		Matrix: shaderMatrix,
		Bounds: bounds,
	}

	switch shade._type {
	case C.FZ_FUNCTION_BASED:
		shader.Kind = FunctionShaderKind
	case C.FZ_RADIAL:
		shader.Kind = RadialShaderKind
	case C.FZ_MESH_TYPE4, C.FZ_MESH_TYPE5, C.FZ_MESH_TYPE6, C.FZ_MESH_TYPE7:
		shader.Kind = MeshShaderKind
	}

	device.FillShade(shader, matrix, float64(alpha))
}

//export fzgo_fill_image
func fzgo_fill_image(ctx *C.fz_context, dev *C.fz_device, image *C.fz_image, ctm C.fz_matrix, alpha C.float, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)

	if !device.Methods().Has(FillImage) {
		return
	}

	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	im := getImage(ctx, ctm, image, colorParams)

	device.FillImage(im, matrix, float64(alpha))
}

//export fzgo_fill_image_mask
func fzgo_fill_image_mask(ctx *C.fz_context, dev *C.fz_device, image *C.fz_image, ctm C.fz_matrix, colorspace *C.fz_colorspace, color *C.cfloat_t, alpha C.float, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)

	if !device.Methods().Has(FillImageMask) {
		return
	}

	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	rgb := getRGBColor(ctx, color, colorspace, alpha, colorParams)
	im := getImage(ctx, ctm, image, colorParams)

	device.FillImageMask(im, matrix, rgb)
}

//export fzgo_clip_path
func fzgo_clip_path(ctx *C.fz_context, dev *C.fz_device, path *C.cfz_path_t, evenOdd C.int, ctm C.fz_matrix, scissor C.fz_rect) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(ClipPath) {
		return
	}

	p := makePath(ctx, path)
	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	sci := rectFromFitz(scissor)
	fillRule := FillRuleWinding
	if evenOdd != 0 {
		fillRule = FillRuleEvenOdd
	}

	device.ClipPath(p, fillRule, matrix, sci)
}

//export fzgo_clip_stroke_path
func fzgo_clip_stroke_path(ctx *C.fz_context, dev *C.fz_device, path *C.cfz_path_t, stroke *C.cfz_stroke_state_t, ctm C.fz_matrix, scissor C.fz_rect) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(ClipStrokePath) {
		return
	}

	p := makePath(ctx, path)
	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	sci := rectFromFitz(scissor)
	s := getStroke(stroke)

	device.ClipStrokePath(p, s, matrix, sci)
}

//export fzgo_fill_text
func fzgo_fill_text(ctx *C.fz_context, dev *C.fz_device, text *C.cfz_text_t, ctm C.fz_matrix, colorspace *C.fz_colorspace, color *C.cfloat_t, alpha C.float, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(FillText) {
		return
	}

	rgb := getRGBColor(ctx, color, colorspace, alpha, colorParams)
	txt := getTextInfo(ctx, text, ctm, rgb)
	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))

	device.FillText(txt, matrix, rgb)
}

//export fzgo_stroke_text
func fzgo_stroke_text(ctx *C.fz_context, dev *C.fz_device, text *C.cfz_text_t, stroke *C.cfz_stroke_state_t, ctm C.fz_matrix, colorspace *C.fz_colorspace, color *C.cfloat_t, alpha C.float, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(StrokeText) {
		return
	}

	rgb := getRGBColor(ctx, color, colorspace, alpha, colorParams)
	txt := getTextInfo(ctx, text, ctm, rgb)
	s := getStroke(stroke)
	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))

	device.StrokeText(txt, s, matrix, rgb)
}

//export fzgo_clip_text
func fzgo_clip_text(ctx *C.fz_context, dev *C.fz_device, text *C.cfz_text_t, ctm C.fz_matrix, scissor C.fz_rect) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(ClipText) {
		return
	}

	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	sci := rectFromFitz(scissor)
	txt := getTextInfo(ctx, text, ctm, nil)

	device.ClipText(txt, matrix, sci)
}

//export fzgo_clip_stroke_text
func fzgo_clip_stroke_text(ctx *C.fz_context, dev *C.fz_device, text *C.cfz_text_t, stroke *C.cfz_stroke_state_t, ctm C.fz_matrix, scissor C.fz_rect) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(ClipStrokeText) {
		return
	}

	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	sci := rectFromFitz(scissor)
	s := getStroke(stroke)
	txt := getTextInfo(ctx, text, ctm, nil)

	device.ClipStrokeText(txt, s, matrix, sci)
}

//export fzgo_ignore_text
func fzgo_ignore_text(ctx *C.fz_context, dev *C.fz_device, text *C.cfz_text_t, ctm C.fz_matrix) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(IgnoreText) {
		return
	}

	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	txt := getTextInfo(ctx, text, ctm, nil)

	device.IgnoreText(txt, matrix)
}

//export fzgo_clip_image_mask
func fzgo_clip_image_mask(ctx *C.fz_context, dev *C.fz_device, image *C.fz_image, ctm C.fz_matrix, scissor C.fz_rect) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(ClipImageMask) {
		return
	}

	matrix := NewMatrix(float64(ctm.a), float64(ctm.b), float64(ctm.c), float64(ctm.d), float64(ctm.e), float64(ctm.f))
	sci := rectFromFitz(scissor)
	im := getImage(ctx, ctm, image, C.fz_default_color_params)

	device.ClipImageMask(im, matrix, sci)
}

//export fzgo_pop_clip
func fzgo_pop_clip(ctx *C.fz_context, dev *C.fz_device) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(PopClip) {
		return
	}
	device.PopClip()
}

//export fzgo_begin_mask
func fzgo_begin_mask(ctx *C.fz_context, dev *C.fz_device, rect C.fz_rect, luminosity C.int, colorspace *C.fz_colorspace, color *C.cfloat_t, colorParams C.fz_color_params) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(BeginMask) {
		return
	}

	rgb := getRGBColor(ctx, color, colorspace, 1.0, colorParams)
	r := rectFromFitz(rect)

	device.BeginMask(r, rgb, int(luminosity))
}

//export fzgo_end_mask
func fzgo_end_mask(ctx *C.fz_context, dev *C.fz_device) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(EndMask) {
		return
	}
	device.EndMask()
}

//export fzgo_begin_group
func fzgo_begin_group(ctx *C.fz_context, dev *C.fz_device, rect C.fz_rect, cs *C.fz_colorspace, isolated C.int, knockout C.int, blendmode C.int, alpha C.float) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(BeginGroup) {
		return
	}
	r := rectFromFitz(rect)
	colorspace := &Colorspace{
		Kind:          ColorspaceKind(cs._type),
		Name:          C.GoString(C.fz_colorspace_name(ctx, cs)),
		ColorantCount: int(C.fz_colorspace_n(ctx, cs)),
		Flags:         uint32(cs.flags),
	}
	device.BeginGroup(r, colorspace, isolated != 0, knockout != 0, BlendMode(blendmode), float64(alpha))
}

//export fzgo_end_group
func fzgo_end_group(ctx *C.fz_context, dev *C.fz_device) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(EndGroup) {
		return
	}
	device.EndGroup()
}

//export fzgo_begin_tile
func fzgo_begin_tile(ctx *C.fz_context, dev *C.fz_device, area C.fz_rect, view C.fz_rect, xstep C.float, ystep C.float, ctm C.fz_matrix, ID C.int) C.int {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(BeginTile) {
		return 0
	}
	return C.int(device.BeginTile())
}

//export fzgo_end_tile
func fzgo_end_tile(ctx *C.fz_context, dev *C.fz_device) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(EndTile) {
		return
	}
	device.EndTile()
}

//export fzgo_begin_layer
func fzgo_begin_layer(ctx *C.fz_context, dev *C.fz_device, layerName *C.cchar_t) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(BeginLayer) {
		return
	}
	device.BeginLayer(C.GoString(layerName))
}

//export fzgo_end_layer
func fzgo_end_layer(ctx *C.fz_context, dev *C.fz_device) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	if !device.Methods().Has(EndLayer) {
		return
	}
	device.EndLayer()
}

//export fzgo_close_device
func fzgo_close_device(ctx *C.fz_context, dev *C.fz_device) {
	device := pointer.Restore(((*C.fzgo_device)(unsafe.Pointer(dev))).user_data).(GoDevice)
	device.Close()
}

//export gopath_moveto
func gopath_moveto(ctx *C.fz_context, arg *C.void, x C.float, y C.float) {
	walker := pointer.Restore(unsafe.Pointer(arg)).(PathBuilder)
	walker.MoveTo(float64(x), float64(y))
}

//export gopath_lineto
func gopath_lineto(ctx *C.fz_context, arg *C.void, x C.float, y C.float) {
	walker := pointer.Restore(unsafe.Pointer(arg)).(PathBuilder)
	walker.LineTo(float64(x), float64(y))
}

//export gopath_curveto
func gopath_curveto(ctx *C.fz_context, arg *C.void, x1 C.float, y1 C.float, x2 C.float, y2 C.float, x3 C.float, y3 C.float) {
	walker := pointer.Restore(unsafe.Pointer(arg)).(PathBuilder)
	walker.CubicCurveTo(float64(x1), float64(y1), float64(x2), float64(y2), float64(x3), float64(y3))
}

//export gopath_quadto
func gopath_quadto(ctx *C.fz_context, arg *C.void, x1 C.float, y1 C.float, x2 C.float, y2 C.float) {
	walker := pointer.Restore(unsafe.Pointer(arg)).(PathBuilder)
	walker.QuadCurveTo(float64(x1), float64(y1), float64(x2), float64(y2))
}

//export gopath_closepath
func gopath_closepath(ctx *C.fz_context, arg *C.void) {
	walker := pointer.Restore(unsafe.Pointer(arg)).(PathBuilder)
	walker.ClosePath()
}

func makePath(ctx *C.fz_context, path *C.fz_path) *Path {
	p := &Path{}

	ref := pointer.Save(p)
	defer pointer.Unref(ref)

	C.fz_walk_path(ctx, path, &C.go_path_walker, ref)
	return p
}
