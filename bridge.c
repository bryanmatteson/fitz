
#include "bridge.h"

#include <memory.h>
#include <pthread.h>

extern void lock_mutex(void* user, int lock);
extern void unlock_mutex(void* user, int lock);

extern void exception_callback(int code, char* message);
extern void error_callback(void* user_data, const char* message);
extern void warn_callback(void* user_data, const char* message);

const char* fz_version = FZ_VERSION;
const int fz_lock_max = (int)FZ_LOCK_MAX;

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