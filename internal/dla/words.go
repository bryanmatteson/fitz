package dla

import (
	"math"
	"unicode"

	"go.matteson.dev/fitz/internal/cluster"
	"go.matteson.dev/gfx"
)

// NearestNeighborWordExtractor ...
type NearestNeighborWordExtractor struct {
	options *NearestNeighborWordExtractorOptions
}

// NewNearestNeighborWordExtractor ...
func NewNearestNeighborWordExtractor(options ...NearestNeighborWordExtractorOptionFunc) WordExtractor {
	opts := DefaultNearestNeighborWordExtractorOptions()
	for _, o := range options {
		o(opts)
	}

	return &NearestNeighborWordExtractor{
		options: opts,
	}
}

// GetWords ...
func (e *NearestNeighborWordExtractor) GetWords(letters gfx.Chars) (words gfx.TextWords) {
	if len(letters) == 0 {
		return
	}

	if e.options.GroupByOrientation {
		var horiz, left, right, upsideDown, other gfx.Chars
		for _, l := range letters {
			switch l.Orientation {
			case gfx.PageUp:
				horiz = append(horiz, l)
			case gfx.PageLeft:
				left = append(left, l)
			case gfx.PageRight:
				right = append(right, l)
			case gfx.PageDown:
				upsideDown = append(upsideDown, l)
			default:
				other = append(other, l)
			}
		}

		words = append(words, e.getWords(horiz, e.options.DistanceMeasureAxisAlignedFn)...)
		words = append(words, e.getWords(left, e.options.DistanceMeasureAxisAlignedFn)...)
		words = append(words, e.getWords(right, e.options.DistanceMeasureAxisAlignedFn)...)
		words = append(words, e.getWords(upsideDown, e.options.DistanceMeasureAxisAlignedFn)...)
		words = append(words, e.getWords(other, e.options.DistanceMeasureFn)...)
		return
	}

	return e.getWords(letters, e.options.DistanceMeasureFn)
}

func (e *NearestNeighborWordExtractor) getWords(letters gfx.Chars, measureFn cluster.DistanceMeasureFunc) (words gfx.TextWords) {
	if len(letters) == 0 {
		return
	}

	pivotFn := func(i interface{}) gfx.Point { return i.(gfx.Char).EndBaseline }
	distanceMeasureFn := func(i1, i2 interface{}) float64 { return measureFn(i1.(gfx.Point), i2.(gfx.Point)) }

	nodes := cluster.MakeTreeNodes(letters, func(i interface{}) gfx.Point { return i.(gfx.Char).StartBaseline })
	indices := make([]int, len(nodes))
	for i := range indices {
		indices[i] = -1
	}

	tree := cluster.NewKDTree(nodes)

	for i, letter := range letters {
		if !e.options.FilterPivotFn(letter) {
			continue
		}

		node, dist := tree.NearestNeighbor(letter, pivotFn, distanceMeasureFn)
		if !e.options.FilterFn(letter, node.Element) {
			continue
		}

		if dist >= e.options.MaxDistanceFn(letter, node.Element) {
			continue
		}

		indices[i] = node.Index
	}

	groups := cluster.GroupIndicesNN(indices)

	words = make(gfx.TextWords, len(groups))
	for grpIdx, grp := range groups {
		wordLetters := make(gfx.Chars, grp.Len())
		for i, val := range grp.Values() {
			wordLetters[i] = letters[val]
		}
		words[grpIdx] = gfx.MakeWord(wordLetters)
	}

	return
}

func maxDistance(l1 interface{}, l2 interface{}) float64 {
	letter1, letter2 := l1.(gfx.Char), l2.(gfx.Char)

	quadWidth := math.Max(letter1.Quad.Width(), letter2.Quad.Width())
	return quadWidth * 0.2
}

// NearestNeighborWordExtractorOptions ...
type NearestNeighborWordExtractorOptions struct {
	GroupByOrientation bool

	MaxDistanceFn                cluster.MaxDistanceFunc
	DistanceMeasureFn            cluster.DistanceMeasureFunc
	DistanceMeasureAxisAlignedFn cluster.DistanceMeasureFunc
	FilterFn                     cluster.CandidateFilterFunc
	FilterPivotFn                cluster.PivotFilterFunc
}

// DefaultNearestNeighborWordExtractorOptions ...
func DefaultNearestNeighborWordExtractorOptions() *NearestNeighborWordExtractorOptions {
	return &NearestNeighborWordExtractorOptions{
		GroupByOrientation:           true,
		MaxDistanceFn:                maxDistance,
		DistanceMeasureFn:            cluster.EuclideanDistance,
		DistanceMeasureAxisAlignedFn: cluster.ManhattanDistance,
		FilterFn:                     func(l1 interface{}, l2 interface{}) bool { return !unicode.IsSpace(l2.(gfx.Char).Rune) },
		FilterPivotFn:                func(l interface{}) bool { return !unicode.IsSpace(l.(gfx.Char).Rune) },
	}
}

// NearestNeighborWordExtractorOptionFunc ...
type NearestNeighborWordExtractorOptionFunc func(*NearestNeighborWordExtractorOptions)

// WithGroupByOrientation ...
func WithGroupByOrientation(group bool) NearestNeighborWordExtractorOptionFunc {
	return func(e *NearestNeighborWordExtractorOptions) {
		e.GroupByOrientation = group
	}
}

// WithMaxDistanceFunction ...
func WithMaxDistanceFunction(fn cluster.MaxDistanceFunc) NearestNeighborWordExtractorOptionFunc {
	return func(e *NearestNeighborWordExtractorOptions) {
		e.MaxDistanceFn = fn
	}
}

// WithDistanceMeasureFunction ...
func WithDistanceMeasureFunction(fn cluster.DistanceMeasureFunc) NearestNeighborWordExtractorOptionFunc {
	return func(e *NearestNeighborWordExtractorOptions) {
		e.DistanceMeasureFn = fn
	}
}

// WithDistanceMeasureAxisAlignedFunction ...
func WithDistanceMeasureAxisAlignedFunction(fn cluster.DistanceMeasureFunc) NearestNeighborWordExtractorOptionFunc {
	return func(e *NearestNeighborWordExtractorOptions) {
		e.DistanceMeasureAxisAlignedFn = fn
	}
}

// WithFilterFunction ...
func WithFilterFunction(fn cluster.CandidateFilterFunc) NearestNeighborWordExtractorOptionFunc {
	return func(e *NearestNeighborWordExtractorOptions) {
		e.FilterFn = fn
	}
}

// WithFilterPivotFunction ...
func WithFilterPivotFunction(fn cluster.PivotFilterFunc) NearestNeighborWordExtractorOptionFunc {
	return func(e *NearestNeighborWordExtractorOptions) {
		e.FilterPivotFn = fn
	}
}
