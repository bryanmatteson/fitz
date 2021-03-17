#ifndef BRIDGE_H
#define BRIDGE_H

#include "mupdf/fitz.h"
#include "mupdf/pdf.h"

extern const char* fz_version;
extern const int fz_lock_max;

typedef const float cfloat_t;
typedef const fz_path cfz_path_t;
typedef const fz_text cfz_text_t;
typedef const char cchar_t;
typedef const fz_stroke_state cfz_stroke_state_t;
typedef const fz_path_walker cfz_path_walker;

fz_context* fzgo_new_context();

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

fz_device* fz_new_go_device(fz_context* ctx, void* user_data);
pdf_obj* pdfname(int typ);

typedef struct fzgo_device {
    fz_device super;
    void* user_data;
} fzgo_device;

extern void gopath_moveto(fz_context* ctx, void* arg, float x, float y);
extern void gopath_lineto(fz_context* ctx, void* arg, float x, float y);
extern void gopath_curveto(fz_context* ctx, void* arg, float x1, float y1, float x2, float y2, float x3, float y3);
extern void gopath_closepath(fz_context* ctx, void* arg);
extern void gopath_quadto(fz_context* ctx, void* arg, float x1, float y1, float x2, float y2);
extern void gopath_curvetov(fz_context* ctx, void* arg, float x2, float y2, float x3, float y3);
extern void gopath_curvetoy(fz_context* ctx, void* arg, float x1, float y1, float x3, float y3);
extern void gopath_rectto(fz_context* ctx, void* arg, float x1, float y1, float x2, float y2);

extern const fz_path_walker go_path_walker;

int fz_text_span_wmode(fz_text_span* span);

#endif