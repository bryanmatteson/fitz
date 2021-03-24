package cluster

import (
	"github.com/ahmetb/go-linq"
	"go.matteson.dev/fitz/internal/data"
	"go.matteson.dev/gfx"
)

// InterfaceSlice ...
func InterfaceSlice(slice interface{}) (results []interface{}) {
	linq.From(slice).ToSlice(&results)
	return
}

// MakeTreeNodes ...
func MakeTreeNodes(slice interface{}, pointFn func(interface{}) gfx.Point) (nodes []*KDNode) {
	elements := InterfaceSlice(slice)
	nodes = make([]*KDNode, len(elements))
	for i, elem := range elements {
		nodes[i] = &KDNode{
			Point:   pointFn(elem),
			Element: elem,
			Index:   i,
		}
	}
	return
}

// // NearestNeighbors ...
// func NearestNeighbors(elements []interface{}, k int, config KDTreeConfigurator) []*data.IntegerSet {
// 	indices := make([]int, len(elements))
// 	for i := range indices {
// 		indices[i] = -1
// 	}
// 	tree := NewKDTree(elements)

// 	for i, elem := range elements {
// 		if !config.PivotFilter(elem) {
// 			continue
// 		}

// 		for _, res := range tree.NearestNeighbors(elem, k, config.Distance) {
// 			if !config.CandidateFilter(elem, res.Node.Element) {
// 				continue
// 			}

// 			if res.Distance >= config.MaxDistance(elem, res.Node.Element) {
// 				continue
// 			}
// 			indices[i] = res.Node.Index
// 		}
// 	}

// 	return GroupIndicesNN(indices)
// }

// GroupIndicesNN ...
func GroupIndicesNN(edges []int) (groups []*data.IntegerSet) {
	adjacency := make([][]int, len(edges))
	for i := 0; i < len(edges); i++ {
		matchSet := data.NewIntegerSet(len(edges))
		if edges[i] != -1 {
			matchSet.Add(edges[i])
			for j := 0; j < len(edges); j++ {
				if edges[j] == i {
					matchSet.Add(j)
				}
			}
		}
		adjacency[i] = matchSet.Values()
	}

	isDone := make([]bool, len(edges))

	for i := 0; i < len(edges); i++ {
		if isDone[i] {
			continue
		}
		groups = append(groups, dfsIterative(i, adjacency, isDone))
	}
	return
}

func dfsIterative(s int, adj [][]int, done []bool) (set *data.IntegerSet) {
	set = data.NewIntegerSet(16)
	stack := data.NewStack(16)
	stack.Push(s)
	for stack.Len() > 0 {
		u := stack.Pop().(int)
		if done[u] {
			continue
		}

		set.Add(u)
		done[u] = true
		for _, v := range adj[u] {
			stack.Push(v)
		}
	}
	return
}
