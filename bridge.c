
#include "bridge.h"

#include <memory.h>
#include <pthread.h>

extern void lock_mutex(void* user, int lock);
extern void unlock_mutex(void* user, int lock);

extern void exception_callback(int code, char* message);
extern void error_callback(void* user_data, const char* message);
extern void warn_callback(void* user_data, const char* message);

extern void fzgo_fill_path(fz_context* ctx, fz_device* dev, const fz_path* path, int even_odd, fz_matrix ctm, fz_colorspace* colorspace, const float* color, float alpha, fz_color_params color_params);
extern void fzgo_stroke_path(fz_context* ctx, fz_device* dev, const fz_path* path, const fz_stroke_state* stroke, fz_matrix ctm, fz_colorspace* colorspace, const float* color, float alpha, fz_color_params color_params);
extern void fzgo_fill_shade(fz_context* ctx, fz_device* dev, fz_shade* shade, fz_matrix ctm, float alpha, fz_color_params color_params);
extern void fzgo_fill_image(fz_context* ctx, fz_device* dev, fz_image* image, fz_matrix ctm, float alpha, fz_color_params color_params);
extern void fzgo_fill_image_mask(fz_context* ctx, fz_device* dev, fz_image* image, fz_matrix ctm, fz_colorspace* colorspace, const float* color, float alpha, fz_color_params color_params);
extern void fzgo_clip_path(fz_context* ctx, fz_device* dev, const fz_path* path, int even_odd, fz_matrix ctm, fz_rect scissor);
extern void fzgo_clip_stroke_path(fz_context* ctx, fz_device* dev, const fz_path* path, const fz_stroke_state* stroke, fz_matrix ctm, fz_rect scissor);
extern void fzgo_fill_text(fz_context* ctx, fz_device* dev, const fz_text* text, fz_matrix ctm, fz_colorspace* colorspace, const float* color, float alpha, fz_color_params color_params);
extern void fzgo_stroke_text(fz_context* ctx, fz_device* dev, const fz_text* text, const fz_stroke_state* stroke, fz_matrix ctm, fz_colorspace* colorspace, const float* color, float alpha, fz_color_params color_params);
extern void fzgo_clip_text(fz_context* ctx, fz_device* dev, const fz_text* text, fz_matrix ctm, fz_rect scissor);
extern void fzgo_clip_stroke_text(fz_context* ctx, fz_device* dev, const fz_text* text, const fz_stroke_state* stroke, fz_matrix ctm, fz_rect scissor);
extern void fzgo_ignore_text(fz_context* ctx, fz_device* dev, const fz_text* text, fz_matrix ctm);
extern void fzgo_clip_image_mask(fz_context* ctx, fz_device* dev, fz_image* image, fz_matrix ctm, fz_rect scissor);
extern void fzgo_pop_clip(fz_context* ctx, fz_device* dev);
extern void fzgo_begin_mask(fz_context* ctx, fz_device* dev, fz_rect rect, int luminosity, fz_colorspace* colorspace, const float* color, fz_color_params color_params);
extern void fzgo_end_mask(fz_context* ctx, fz_device* dev);
extern void fzgo_begin_group(fz_context* ctx, fz_device* dev, fz_rect rect, fz_colorspace* cs, int isolated, int knockout, int blendmode, float alpha);
extern void fzgo_end_group(fz_context* ctx, fz_device* dev);
extern int fzgo_begin_tile(fz_context* ctx, fz_device* dev, fz_rect area, fz_rect view, float xstep, float ystep, fz_matrix ctm, int id);
extern void fzgo_end_tile(fz_context* ctx, fz_device* dev);
extern void fzgo_begin_layer(fz_context* ctx, fz_device* dev, const char* layer_name);
extern void fzgo_end_layer(fz_context* ctx, fz_device* dev);
extern void fzgo_close_device(fz_context* ctx, fz_device* dev);

extern void gopath_moveto(fz_context* ctx, void* arg, float x, float y);
extern void gopath_lineto(fz_context* ctx, void* arg, float x, float y);
extern void gopath_curveto(fz_context* ctx, void* arg, float x1, float y1, float x2, float y2, float x3, float y3);
extern void gopath_closepath(fz_context* ctx, void* arg);
extern void gopath_quadto(fz_context* ctx, void* arg, float x1, float y1, float x2, float y2);
extern void gopath_curvetov(fz_context* ctx, void* arg, float x2, float y2, float x3, float y3);
extern void gopath_curvetoy(fz_context* ctx, void* arg, float x1, float y1, float x3, float y3);
extern void gopath_rectto(fz_context* ctx, void* arg, float x1, float y1, float x2, float y2);

extern void gooutput_writer_write(fz_context* ctx, void* state, const void* data, size_t n);
extern void gooutput_writer_close(fz_context* ctx, void* state);
extern void gooutput_writer_drop(fz_context* ctx, void* state);
extern void gooutput_writer_seek(fz_context* ctx, void* state, int64_t offset, int whence);
extern int64_t gooutput_writer_tell(fz_context* ctx, void* state);

fz_output* fzgo_new_output_writer(fz_context* ctx, int bufsize, void* iowriter) {
    fz_output* output = fz_new_output(ctx, bufsize, iowriter, gooutput_writer_write, gooutput_writer_close, gooutput_writer_drop);
    output->tell = gooutput_writer_tell;
    output->seek = gooutput_writer_seek;
    return output;
}

const char* fz_version = FZ_VERSION;

fz_context* fzgo_new_context() {
    fz_locks_context locks;
    locks.lock = lock_mutex;
    locks.unlock = unlock_mutex;

    fz_context* ctx = fz_new_context(NULL, &locks, FZ_STORE_UNLIMITED);
    fz_set_error_callback(ctx, error_callback, NULL);
    fz_set_warning_callback(ctx, warn_callback, NULL);
    return ctx;
}

const fz_path_walker go_path_walker = {gopath_moveto, gopath_lineto, gopath_curveto, gopath_closepath, gopath_quadto, NULL, NULL, NULL};

fz_device* fz_new_go_device(fz_context* ctx, void* user_data) {
    fzgo_device* dev = (fzgo_device*)fz_new_derived_device(ctx, fzgo_device);

    dev->super.fill_path = fzgo_fill_path;
    dev->super.stroke_path = fzgo_stroke_path;
    dev->super.clip_path = fzgo_clip_path;
    dev->super.clip_stroke_path = fzgo_clip_stroke_path;

    dev->super.fill_shade = fzgo_fill_shade;
    dev->super.fill_image = fzgo_fill_image;
    dev->super.fill_image_mask = fzgo_fill_image_mask;
    dev->super.clip_image_mask = fzgo_clip_image_mask;

    dev->super.fill_text = fzgo_fill_text;
    dev->super.stroke_text = fzgo_stroke_text;
    dev->super.clip_text = fzgo_clip_text;
    dev->super.clip_stroke_text = fzgo_clip_stroke_text;
    dev->super.ignore_text = fzgo_ignore_text;

    dev->super.pop_clip = fzgo_pop_clip;

    dev->super.begin_mask = fzgo_begin_mask;
    dev->super.end_mask = fzgo_end_mask;
    dev->super.begin_group = fzgo_begin_group;
    dev->super.end_group = fzgo_end_group;

    dev->super.begin_tile = fzgo_begin_tile;
    dev->super.end_tile = fzgo_end_tile;

    dev->super.begin_layer = fzgo_begin_layer;
    dev->super.end_layer = fzgo_end_layer;
    dev->super.close_device = fzgo_close_device;

    dev->user_data = user_data;

    return (fz_device*)dev;
}

int fz_text_span_wmode(fz_text_span* span) {
    return span->wmode;
}

pdf_obj* pdfname(int typ) {
    return (pdf_obj*)((intptr_t)typ);
}