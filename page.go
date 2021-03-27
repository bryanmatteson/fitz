package fitz

// #include "bridge.h"
import "C"
import (
	"image"
	"sync"
	"unsafe"

	"github.com/mattn/go-pointer"
	"go.matteson.dev/gfx"
)

type Page struct {
	number int
	mut    sync.Mutex
	ctx    *C.fz_context
	list   *C.fz_display_list
	bounds C.fz_rect
}

func newPage(ctx *C.fz_context, number int, bounds C.fz_rect, list *C.fz_display_list) *Page {
	ctx = C.fz_clone_context(ctx)
	return &Page{ctx: ctx, number: number, bounds: bounds, list: list}
}

func (p *Page) drop() {
	C.fz_drop_display_list(p.ctx, p.list)
	C.fz_drop_context(p.ctx)
	p.list = nil
	p.ctx = nil
}

func (p *Page) Number() int      { return p.number }
func (p *Page) Bounds() gfx.Rect { return rectFromFitz(p.bounds) }

func (p *Page) RenderImage(region gfx.Rect, scale float64) (img *image.RGBA, err error) {
	p.mut.Lock()
	defer p.mut.Unlock()
	img = &image.RGBA{}

	bounds := p.bounds
	ctm := C.fz_identity

	if scale > 0 {
		ctm = C.fz_scale(C.float(scale), C.float(scale))
	}

	if !region.IsEmpty() {
		region = region.Intersection(p.Bounds())
		bounds = C.fz_make_rect(C.float(region.X.Min), C.float(region.Y.Min), C.float(region.X.Max), C.float(region.Y.Max))
	}

	ctm = C.fz_concat(C.fz_translate(-bounds.x0, -bounds.y0), ctm)
	bounds = C.fz_transform_rect(bounds, ctm)
	bbox := C.fz_round_rect(bounds)

	width := bbox.x1 - bbox.x0
	height := bbox.y1 - bbox.y0
	pixBbox := C.fz_make_irect(0, 0, width, height)

	pixmap := C.fz_new_pixmap_with_bbox(p.ctx, C.fz_device_rgb(p.ctx), pixBbox, nil, 1)
	if pixmap == nil {
		return nil, ErrCreatePixmap
	}

	C.fz_clear_pixmap_with_value(p.ctx, pixmap, C.int(0xff))
	defer C.fz_drop_pixmap(p.ctx, pixmap)

	device := C.fz_new_draw_device_with_bbox(p.ctx, ctm, pixmap, &pixBbox)
	defer C.fz_drop_device(p.ctx, device)

	C.fz_enable_device_hints(p.ctx, device, C.FZ_NO_CACHE)

	C.fz_run_display_list(p.ctx, p.list, device, C.fz_identity, bounds, nil)
	C.fz_close_device(p.ctx, device)

	pixels := C.fz_pixmap_samples(p.ctx, pixmap)
	if pixels == nil {
		return nil, ErrPixmapSamples
	}

	img.Pix = C.GoBytes(unsafe.Pointer(pixels), C.int(4*width*height))
	img.Rect = image.Rect(0, 0, int(width), int(height))
	img.Stride = 4 * img.Rect.Max.X

	return img, nil
}

func (p *Page) RunDevice(device GoDevice) {
	ref := pointer.Save(device)
	defer pointer.Unref(ref)

	fzdev := C.fz_new_go_device(p.ctx, ref)
	defer C.fz_drop_device(p.ctx, fzdev)

	C.fz_run_display_list(p.ctx, p.list, fzdev, C.fz_identity, C.fz_infinite_rect, nil)
	C.fz_close_device(p.ctx, fzdev)
}

// RenderSVG returns svg document for given page number.
func (p *Page) RenderSVG(scale float64) (string, error) {
	p.mut.Lock()
	defer p.mut.Unlock()
	bounds := p.bounds

	ctm := C.fz_scale(C.float(scale), C.float(scale))
	bounds = C.fz_transform_rect(bounds, ctm)

	buf := C.fz_new_buffer(p.ctx, 1024)
	defer C.fz_drop_buffer(p.ctx, buf)

	out := C.fz_new_output_with_buffer(p.ctx, buf)
	defer C.fz_drop_output(p.ctx, out)

	device := C.fz_new_svg_device(p.ctx, out, bounds.x1-bounds.x0, bounds.y1-bounds.y0, C.FZ_SVG_TEXT_AS_PATH, 1)
	C.fz_enable_device_hints(p.ctx, device, C.FZ_NO_CACHE)
	defer C.fz_drop_device(p.ctx, device)

	var cookie C.fz_cookie
	C.fz_run_display_list(p.ctx, p.list, device, C.fz_identity, bounds, &cookie)

	C.fz_close_device(p.ctx, device)

	str := C.GoString(C.fz_string_from_buffer(p.ctx, buf))
	return str, nil
}

// GetText returns text for page
func (p *Page) GetText() string {
	p.mut.Lock()
	defer p.mut.Unlock()
	bounds := p.bounds

	text := C.fz_new_stext_page(p.ctx, bounds)
	defer C.fz_drop_stext_page(p.ctx, text)

	opts := C.fz_stext_options{}
	opts.flags = 0

	device := C.fz_new_stext_device(p.ctx, text, &opts)
	C.fz_enable_device_hints(p.ctx, device, C.FZ_NO_CACHE)
	defer C.fz_drop_device(p.ctx, device)

	var cookie C.fz_cookie
	C.fz_run_display_list(p.ctx, p.list, device, C.fz_identity, bounds, &cookie)
	C.fz_close_device(p.ctx, device)

	buf := C.fz_new_buffer_from_stext_page(p.ctx, text)
	defer C.fz_drop_buffer(p.ctx, buf)

	return C.GoString(C.fz_string_from_buffer(p.ctx, buf))
}
