package fitz

// #include "bridge.h"
import "C"

type usercontext struct {
	ctx       *C.fz_context
	fontCache *fontCache
}

func newusercontext(fzctx *C.fz_context) *usercontext {
	fc := newfontcache()
	return &usercontext{
		ctx:       fzctx,
		fontCache: fc,
	}
}
