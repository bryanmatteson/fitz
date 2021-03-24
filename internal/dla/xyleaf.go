package dla

import (
	"github.com/ahmetb/go-linq"
	"go.matteson.dev/gfx"
)

// type xynode interface {
// 	Bounds() gfx.Quad
// 	IsLeaf() bool
// 	WordCount() int
// 	GetLeaves() []*xyleaf
// }

// type xycontainer struct {
// 	children []xynode
// 	bounds   gfx.Quad
// }

// func newXYContainer(children ...xynode) xynode {
// 	quads := make(gfx.Quads, len(children))
// 	node := &xycontainer{children: children, bounds: quads.Union()}
// 	return node
// }

// // Bounds implements xynode interface
// func (c *xycontainer) Bounds() gfx.Quad { return c.bounds }

// // IsLeaf implements xynode interface
// func (c *xycontainer) IsLeaf() bool { return false }

// // GetLeaves implements xynode interface
// func (c *xycontainer) GetLeaves() (leaves []*xyleaf) {
// 	leaves = c.recursiveGetLeaves(0)
// 	return
// }

// func (c *xycontainer) recursiveGetLeaves(level int) (leaves []*xyleaf) {
// 	vertical := level%2 == 0
// 	containers := []*xycontainer{}
// 	for _, node := range c.children {
// 		if node.IsLeaf() {
// 			leaves = append(leaves, node.(*xyleaf))
// 		} else {
// 			containers = append(containers, node.(*xycontainer))
// 		}
// 	}

// 	sort.Sort(&byOrientation{vertical: vertical, nodes: containers})
// 	for _, node := range containers {
// 		leaves = append(leaves, node.recursiveGetLeaves(level+1)...)
// 	}
// 	return
// }

// WordCount ...
// func (c *xycontainer) WordCount() (count int) {
// 	for _, child := range c.children {
// 		count += child.WordCount()
// 	}
// 	return
// }

type xyleaf struct {
	bounds gfx.Quad
	words  gfx.TextWords
}

func newXYLeaf(words gfx.TextWords) *xyleaf {
	quads := make(gfx.Quads, len(words))
	for _, word := range words {
		quads = append(quads, word.Quad)
	}
	return &xyleaf{words: words, bounds: quads.Normalize()}
}

// Bounds implements xynode interface
func (l *xyleaf) Bounds() gfx.Quad { return l.bounds }

// IsLeaf implements xynode interface
func (l *xyleaf) IsLeaf() bool { return true }

// GetLeaves implements xynode interface
func (l *xyleaf) GetLeaves() []*xyleaf {
	return nil
}

// WordCount ...
func (l *xyleaf) WordCount() int {
	return len(l.words)
}

func (l *xyleaf) GetLines(wordSep string) (lines gfx.TextLines) {
	groupedWords := make([]gfx.TextWords, 0)
	linq.From(l.words).
		GroupBy(func(i interface{}) interface{} { return i.(gfx.TextWord).Quad.Bottom() }, func(i interface{}) interface{} { return i.(gfx.TextWord) }).
		Select(func(i interface{}) interface{} {
			words := make(gfx.TextWords, len(i.(linq.Group).Group))
			linq.From(i.(linq.Group).Group).ToSlice(&words)
			return words
		}).ToSlice(&groupedWords)

	for _, words := range groupedWords {
		lines = append(lines, gfx.MakeTextLine(words.OrderByReadingOrder(), wordSep))
	}

	return lines.OrderByReadingOrder()
}

// //
// //
// //
// //
// // byOrientation
// type byOrientation struct {
// 	vertical bool
// 	nodes    []*xycontainer
// }

// // Len is the number of elements in the collection.
// func (b *byOrientation) Len() int {
// 	return len(b.nodes)
// }

// // Less reports whether the element with
// // index i should sort before the element with index j.
// func (b *byOrientation) Less(i, j int) bool {
// 	if b.vertical {
// 		return b.nodes[i].bounds.Left() < b.nodes[j].bounds.Left()
// 	}
// 	return b.nodes[j].bounds.Top() < b.nodes[i].bounds.Top()
// }

// // Swap swaps the elements with indexes i and j.
// func (b *byOrientation) Swap(i, j int) {
// 	b.nodes[i], b.nodes[j] = b.nodes[j], b.nodes[i]
// }
