#ifndef BRIDGE_H
#define BRIDGE_H

#include "mupdf/fitz.h"
#include "mupdf/pdf.h"

extern const char* fz_version;
extern const fz_path_walker go_path_walker;

typedef const float cfloat_t;
typedef const fz_path cfz_path_t;
typedef const fz_text cfz_text_t;
typedef const char cchar_t;
typedef const fz_stroke_state cfz_stroke_state_t;
typedef const fz_path_walker cfz_path_walker;
typedef const void* cvoidptr_t;

fz_stream* fzgo_new_read_stream(fz_context* ctx, void* state);
fz_context* fzgo_new_context();
fz_context* fzgo_new_user_context(void* user);
fz_device* fz_new_go_device(fz_context* ctx, void* user_data);
pdf_obj* pdfname(int typ);
int fz_text_span_wmode(fz_text_span* span);
fz_output* fzgo_new_output_writer(fz_context* ctx, int bufsize, void* iowriter);

typedef struct fzgo_device {
    fz_device super;
    void* user_data;
} fzgo_device;

#endif