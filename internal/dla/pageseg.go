package dla

import "go.matteson.dev/gfx"

// BasicPageSegmenter ...
type BasicPageSegmenter struct {
	options *PageSegmenterOptions
}

// NewBasicPageSegmenter ...
func NewBasicPageSegmenter(options ...PageSegmenterOptionsFunc) PageSegmenter {
	opts := DefaultPageSegmenterOptions()
	for _, o := range options {
		o(opts)
	}

	return &BasicPageSegmenter{options: opts}
}

// GetBlocks implements the PageSegmenter interface
func (s *BasicPageSegmenter) GetBlocks(words gfx.TextWords) (results gfx.TextBlocks) {
	results = append(results, gfx.MakeTextBlock(newXYLeaf(words).GetLines(s.options.WordSeparator), s.options.LineSeparator))
	return
}

// PageSegmenterOptions ...
type PageSegmenterOptions struct {
	WordSeparator string
	LineSeparator string
}

// DefaultPageSegmenterOptions ...
func DefaultPageSegmenterOptions() *PageSegmenterOptions {
	return &PageSegmenterOptions{
		WordSeparator: " ",
		LineSeparator: "\n",
	}
}

// PageSegmenterOptionsFunc ...
type PageSegmenterOptionsFunc func(*PageSegmenterOptions)

// WithBasicLineSeparator ...
func WithBasicLineSeparator(ls string) PageSegmenterOptionsFunc {
	return func(o *PageSegmenterOptions) {
		o.LineSeparator = ls
	}
}

// WithBasicWordSeparator ...
func WithBasicWordSeparator(ws string) PageSegmenterOptionsFunc {
	return func(o *PageSegmenterOptions) {
		o.WordSeparator = ws
	}
}
