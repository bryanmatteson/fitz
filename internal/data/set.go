package data

import "sort"

// IntegerSet ...
type IntegerSet struct {
	data map[int]bool
}

// NewIntegerSet ...
func NewIntegerSet(capacity int) *IntegerSet {
	return &IntegerSet{data: make(map[int]bool, capacity)}
}

// Remove ...
func (set *IntegerSet) Remove(i int) {
	delete(set.data, i)
}

// Len ...
func (set *IntegerSet) Len() int { return len(set.data) }

// Add ...
func (set *IntegerSet) Add(i int) {
	set.data[i] = true
}

// Values ...
func (set *IntegerSet) Values() (vals []int) {
	vals = make([]int, 0, len(set.data))
	for i := range set.data {
		vals = append(vals, i)
	}
	sort.Ints(vals)
	return
}
