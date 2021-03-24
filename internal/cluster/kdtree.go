package cluster

import (
	"math"
	"sort"

	"go.matteson.dev/gfx"
)

// KDNode ...
type KDNode struct {
	Element interface{}
	Point   gfx.Point
	Index   int
}

// KDTree ...
type KDTree struct {
	root  *kdtreenode
	count int
}

// NewKDTree ...
func NewKDTree(nodes []*KDNode) *KDTree {
	return &KDTree{
		root:  buildtree(nodes, 0),
		count: len(nodes),
	}
}

func buildtree(nodes []*KDNode, depth int) *kdtreenode {
	if len(nodes) == 0 {
		return nil
	}

	if len(nodes) == 1 {
		return newleafnode(nodes[0], depth)
	}

	sort.Sort(&byAxis{axis: depth, nodes: nodes})

	if len(nodes) == 2 {
		return newtreenode(newleafnode(nodes[0], depth+1), nil, nodes[1], depth)
	}

	median := len(nodes) / 2

	left := buildtree(nodes[:median], depth+1)
	right := buildtree(nodes[median+1:], depth+1)
	return newtreenode(left, right, nodes[median], depth)
}

// NearestNeighbor ...
func (tree *KDTree) NearestNeighbor(pivot interface{}, pivotPointFn func(interface{}) gfx.Point, distance func(interface{}, interface{}) float64) (*KDNode, float64) {
	queue := newkdnnqueue(1)
	tree.knn(tree.root, 1, pivot, pivotPointFn, distance, queue)
	return queue.lastElement.KDNode, queue.lastDistance
}

// KDNResult ...
type KDNResult struct {
	Node     *KDNode
	Distance float64
}

// NearestNeighbors ...
func (tree *KDTree) NearestNeighbors(pivot interface{}, k int, pivotPointFn func(interface{}) gfx.Point, distance func(interface{}, interface{}) float64) (nodes []*KDNResult) {
	queue := newkdnnqueue(k)

	tree.knn(tree.root, k, pivot, pivotPointFn, distance, queue)
	nodes = make([]*KDNResult, 0, k)
	for _, dist := range queue.sortedKeys {
		for _, node := range queue.data[dist].items {
			nodes = append(nodes, &KDNResult{Node: node.KDNode, Distance: dist})
		}
	}
	return
}

// func (tree *KDTree) knnalt(node *kdtreenode, pivot interface{}, k int, currentAxis int, pivotPointFn func(interface{}) gfx.Point, distance func(interface{}, interface{}) float64, queue *PriorityQueue) {
// 	if k == 0 || node == nil {
// 		return
// 	}

// 	p := pivotPointFn(pivot)

// 	var path []*kdtreenode
// 	currentNode := node

// 	// 1. move down
// 	for currentNode != nil {
// 		path = append(path, currentNode)
// 		if ptdim(p, currentAxis) < ptdim(currentNode.KDNode.Point, currentAxis) {
// 			currentNode = currentNode.Left
// 		} else {
// 			currentNode = currentNode.Right
// 		}
// 		currentAxis = (currentAxis + 1) % 2
// 	}

// 	// 2. move up
// 	currentAxis = (currentAxis + 1) % 2
// 	for path, currentNode = popLast(path); currentNode != nil; path, currentNode = popLast(path) {
// 		currentDistance := distance(p, currentNode.KDNode.Point)
// 		checkedDistance := getKthOrLastDistance(queue, k-1)
// 		if currentDistance < checkedDistance {
// 			queue.Insert(currentNode, currentDistance)
// 			checkedDistance = getKthOrLastDistance(queue, k-1)
// 		}

// 		// check other side of plane
// 		if planeDistance(p, ptdim(currentNode.KDNode.Point, currentAxis), currentAxis) < checkedDistance {
// 			var next *kdtreenode
// 			if ptdim(p, currentAxis) < ptdim(currentNode.KDNode.Point, currentAxis) {
// 				next = currentNode.Right
// 			} else {
// 				next = currentNode.Left
// 			}
// 			tree.knnalt(next, k, pivot, (currentAxis+1)%2, pivotPointFn, distance, queue)
// 		}
// 		currentAxis = (currentAxis + 1) % 2
// 	}
// }

// func planeDistance(p gfx.Point, planePosition float64, dim int) float64 {
// 	return math.Abs(planePosition - ptdim(p, dim))
// }

// func popLast(arr []*kdtreenode) ([]*kdtreenode, *kdtreenode) {
// 	l := len(arr) - 1
// 	if l < 0 {
// 		return arr, nil
// 	}
// 	return arr[:l], arr[l]
// }

// func getKthOrLastDistance(queue *PriorityQueue, i int) float64 {
// 	if queue.Len() <= i {
// 		return math.MaxFloat64
// 	}
// 	_, prio := queue.Get(i)
// 	return prio
// }

func (tree *KDTree) knn(node *kdtreenode, k int, pivot interface{}, pivotPointFn func(interface{}) gfx.Point, distance func(interface{}, interface{}) float64, queue *kdnnqueue) (*kdtreenode, float64) {
	if node == nil {
		return nil, math.NaN()
	}

	point := pivotPointFn(pivot)
	if node.Leaf {
		if node.Element == pivot {
			return nil, math.NaN()
		}

		currentDistance := distance(node.Point, point)
		currentNearestNode := node

		if !queue.isFull() || currentDistance <= queue.lastDistance {
			queue.addNode(currentDistance, currentNearestNode)
			currentDistance = queue.lastDistance
			currentNearestNode = queue.lastElement
		}

		return currentNearestNode, currentDistance
	}

	currentNearestNode := node
	currentDistance := distance(node.Point, point)

	if (!queue.isFull() || currentDistance <= queue.lastDistance) && node.Element != pivot {
		queue.addNode(currentDistance, currentNearestNode)
		currentDistance = queue.lastDistance
		currentNearestNode = queue.lastElement
	}

	var newDist float64
	var newNode *kdtreenode

	pointValue := point.Y
	if node.isAxisCutX() {
		pointValue = point.X
	}

	if pointValue < node.splitValue() {
		newNode, newDist = tree.knn(node.Left, k, pivot, pivotPointFn, distance, queue)
		if !math.IsNaN(newDist) && newDist <= currentDistance && newNode.Element != pivot {
			queue.addNode(newDist, newNode)
			currentDistance = queue.lastDistance
			currentNearestNode = queue.lastElement
		}

		if node.Right != nil && pointValue+currentDistance >= node.splitValue() {
			newNode, newDist = tree.knn(node.Right, k, pivot, pivotPointFn, distance, queue)
		}
	} else {
		newNode, newDist = tree.knn(node.Right, k, pivot, pivotPointFn, distance, queue)
		if !math.IsNaN(newDist) && newDist <= currentDistance && newNode.Element != pivot {
			queue.addNode(newDist, newNode)
			currentDistance = queue.lastDistance
			currentNearestNode = queue.lastElement
		}

		if node.Left != nil && pointValue-currentDistance <= node.splitValue() {
			newNode, newDist = tree.knn(node.Left, k, pivot, pivotPointFn, distance, queue)
		}
	}

	if !math.IsNaN(newDist) && newDist <= currentDistance && newNode.Element != pivot {
		queue.addNode(newDist, newNode)
		currentDistance = queue.lastDistance
		currentNearestNode = queue.lastElement
	}

	return currentNearestNode, currentDistance
}

type kdnnqueue struct {
	k int

	lastDistance float64
	lastElement  *kdtreenode
	sortedKeys   []float64
	data         map[float64]*kdnodeset
}

func newkdnnqueue(k int) *kdnnqueue {
	return &kdnnqueue{
		k: k,

		data:         make(map[float64]*kdnodeset, k),
		sortedKeys:   make([]float64, 0, k),
		lastDistance: math.Inf(1),
	}
}

func (q *kdnnqueue) addNode(dist float64, node *kdtreenode) {
	if dist > q.lastDistance && q.isFull() {
		return
	}
	if _, ok := q.data[dist]; !ok {
		q.add(dist, newnodeset(0))
		if len(q.data) > q.k {
			q.removeAt(len(q.data) - 1)
		}
	}
	if q.data[dist].add(node) {
		lastKey := q.sortedKeys[len(q.sortedKeys)-1]
		lastSet := q.data[lastKey]
		q.lastElement = lastSet.items[lastSet.size()-1]
		q.lastDistance = lastKey
	}
}

func (q *kdnnqueue) add(dist float64, set *kdnodeset) {
	q.data[dist] = set
	q.sortedKeys = append(q.sortedKeys, dist)
	sort.Float64s(q.sortedKeys)
}

func (q *kdnnqueue) removeAt(idx int) {
	dist := q.sortedKeys[idx]
	q.sortedKeys = append(q.sortedKeys[:idx], q.sortedKeys[idx+1:]...)
	delete(q.data, dist)
}

// func (q *kdnnqueue) keys() []float64 { return q.sortedKeys }
func (q *kdnnqueue) isFull() bool { return len(q.sortedKeys) >= q.k }

// func (q *kdnnqueue) values() (nodes []*kdtreenode) {
// 	for _, tn := range q.data {
// 		for _, node := range tn.items {
// 			nodes = append(nodes, node)
// 		}
// 	}
// 	return
// }

type kdtreenode struct {
	*KDNode
	Left  *kdtreenode
	Right *kdtreenode
	Depth int
	Leaf  bool
}

func newtreenode(left, right *kdtreenode, node *KDNode, depth int) *kdtreenode {
	return &kdtreenode{
		KDNode: node,
		Left:   left,
		Right:  right,
		Depth:  depth,
		Leaf:   false,
	}
}

func newleafnode(node *KDNode, depth int) *kdtreenode {
	return &kdtreenode{
		KDNode: node,
		Depth:  depth,
		Leaf:   true,
	}
}

func (node *kdtreenode) isAxisCutX() bool { return node.Depth%2 == 0 }
func (node *kdtreenode) splitValue() float64 {
	if node.isAxisCutX() {
		return node.Point.X
	}
	return node.Point.Y
}

// func (node *kdtreenode) getLeaves() (leaves []*kdtreenode) {
// 	leaves = append(leaves, node.Left.recursiveGetLeaves()...)
// 	leaves = append(leaves, node.Right.recursiveGetLeaves()...)
// 	return
// }

// func (node *kdtreenode) recursiveGetLeaves() (leaves []*kdtreenode) {
// 	if node == nil {
// 		return
// 	}

// 	if node.Leaf {
// 		leaves = append(leaves, node)
// 		return
// 	}

// 	leaves = append(leaves, node.Left.recursiveGetLeaves()...)
// 	leaves = append(leaves, node.Right.recursiveGetLeaves()...)
// 	return
// }

type kdnodeset struct {
	data  map[*kdtreenode]bool
	items []*kdtreenode
}

func newnodeset(capacity int) *kdnodeset {
	return &kdnodeset{data: make(map[*kdtreenode]bool, capacity), items: make([]*kdtreenode, capacity)}
}

// func (set *kdnodeset) remove(node *kdtreenode) {
// 	items := make([]*kdtreenode, 0, len(set.items))
// 	for i := range set.items {
// 		if set.items[i] == node {
// 			continue
// 		}
// 		items = append(items, set.items[i])
// 	}
// 	delete(set.data, node)
// }

func (set *kdnodeset) size() int { return len(set.data) }
func (set *kdnodeset) add(node *kdtreenode) bool {
	if _, ok := set.data[node]; ok {
		return false
	}
	set.data[node] = true
	set.items = append(set.items, node)
	return true
}

// func (set *kdnodeset) values() []*kdtreenode {
// 	return set.items
// }

// byAxis
//
//
//
type byAxis struct {
	axis  int
	nodes []*KDNode
}

// Len is the number of elements in the collection.
func (b *byAxis) Len() int {
	return len(b.nodes)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (b *byAxis) Less(i, j int) bool {
	if b.axis%2 == 0 {
		return b.nodes[i].Point.X < b.nodes[j].Point.X
	}
	return b.nodes[i].Point.Y < b.nodes[j].Point.Y
}

// Swap swaps the elements with indexes i and j.
func (b *byAxis) Swap(i, j int) {
	b.nodes[i], b.nodes[j] = b.nodes[j], b.nodes[i]
}

// func ptdim(p gfx.Point, dim int) float64 {
// 	if dim%2 == 0 {
// 		return p.X
// 	}
// 	return p.Y
// }
