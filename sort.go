package collectjs

import "sort"

// A couple of type definitions to make the units clear.
type earthMass float64
type au float64

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *Collect) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(collects []*Collect) {
	ps := &collectSorter{
		collects: collects,
		by:       by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type collectSorter struct {
	collects []*Collect
	by       func(p1, p2 *Collect) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *collectSorter) Len() int {
	return len(s.collects)
}

// Swap is part of sort.Interface.
func (s *collectSorter) Swap(i, j int) {
	s.collects[i], s.collects[j] = s.collects[j], s.collects[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *collectSorter) Less(i, j int) bool {
	return s.by(s.collects[i], s.collects[j])
}
