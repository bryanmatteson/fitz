package dla

import (
	"go.matteson.dev/gfx"
)

// WordExtractor ...
type WordExtractor interface {
	GetWords(letters gfx.Chars) gfx.TextWords
}

// PageSegmenter ...
type PageSegmenter interface {
	GetBlocks(words gfx.TextWords) gfx.TextBlocks
}
