package fitz

// #include "bridge.h"
import "C"

type usercontext struct {
	fontCache *fontCache
}

func newusercontext() *usercontext {
	fc := newfontcache()
	return &usercontext{
		fontCache: fc,
	}
}
