package datastructs

import (
	"fmt"
	"reflect"
)

// Sortable interface definition
type Sortable interface {
	Less(*Sortable) bool
}

// SortedSlice structure
// by convention the smallest value is at the end
type SortedSlice struct {
	s []Sortable
}

// NewSortedSlice returns an empty initialized slice. Opts takes len and cap in
// order to initialize the underlying slice
func NewSortedSlice(opts ...int) *SortedSlice {
	l, c := 0, 0
	if len(opts) >= 1 {
		l = opts[0]
	}
	if len(opts) >= 2 {
		c = opts[1]
	}
	return &SortedSlice{make([]Sortable, l, c)}
}

// Recursive function to search for the next index less than Sortable
func (ss *SortedSlice) searchLessThan(e *Sortable, i, j int) int {
	pivot := ((j + 1 - i) / 2) + i
	if j-i == 1 {
		if ss.s[i].Less(e) {
			return i
		}
		return j
	}
	if ss.s[pivot].Less(e) {
		return ss.searchLessThan(e, i, pivot)
	}
	return ss.searchLessThan(e, pivot, j)
}

// RangeLessThan returns the indexes of the objects Less than Sortable
func (ss *SortedSlice) RangeLessThan(e Sortable) (int, int) {
	i := ss.searchLessThan(&e, 0, len(ss.s)-1)
	return i, len(ss.s) - 1
}

// Insertion method in the slice for a structure implementing Sortable
func (ss *SortedSlice) Insert(e Sortable) {
	switch {
	// Particular cases
	case len(ss.s) == 0, !ss.s[len(ss.s)-1].Less(&e):
		ss.s = append(ss.s, e)
	case len(ss.s) == 1 && ss.s[0].Less(&e):
		ss.s = append(ss.s, e)
		ss.s[1] = ss.s[0]
		ss.s[0] = e
	default:
		//log.Printf("want to insert v=%v into %v", e, ss.s)
		i := ss.searchLessThan(&e, 0, len(ss.s)-1)
		//log.Printf("insert v=%v @ i=%d in ss=%v", e, i, ss.s)
		// Avoid creating intermediary slices
		ss.s = append(ss.s, e)
		copy(ss.s[i+1:], ss.s[i:])
		ss.s[i] = e
	}
}

// Iter returns a chan of Sortable in the slice. Start and Stop indexes can be
// specified via optional parameters
func (ss *SortedSlice) Iter(idx ...int) (c chan Sortable) {
	c = make(chan Sortable)
	i, j := 0, len(ss.s)-1
	if len(idx) >= 1 {
		i = idx[0]
	}
	if len(idx) >= 2 {
		j = idx[1]
	}
	if i < len(ss.s) && j < len(ss.s) && i <= j && i >= 0 {
		go func() {
			defer close(c)
			//for _, v := range ss.s {
			for ; i <= j; i++ {
				v := ss.s[i]
				c <- v
			}
		}()
	} else {
		close(c)
	}
	return c
}

// Iter returns a chan of Sortable in the slice but in reverse order. Start and
// Stop indexes can be specified via optional parameters
func (ss *SortedSlice) ReversedIter(idx ...int) (c chan Sortable) {
	c = make(chan Sortable)
	i, j := 0, len(ss.s)-1
	if len(idx) >= 1 {
		i = idx[0]
	}
	if len(idx) >= 2 {
		j = idx[1]
	}
	if i < len(ss.s) && j < len(ss.s) && i <= j && i >= 0 {
		go func() {
			defer close(c)
			for k := len(ss.s) - 1 - i; k >= len(ss.s)-1-j; k-- {
				v := ss.s[k]
				c <- v
			}
		}()
	} else {
		close(c)
	}
	return c
}

// Slice returns the underlying slice
func (ss *SortedSlice) Slice() []Sortable {
	return ss.s
}

// Control controls if the slice has been properly ordered. A return value of
// true means it is in good order
func (ss *SortedSlice) Control() bool {
	v := ss.s[0]
	for _, tv := range ss.s {
		if !reflect.DeepEqual(v, tv) && !tv.Less(&v) {
			return false
		}
	}
	return true
}

// String fmt helper
func (ss *SortedSlice) String() string {
	return fmt.Sprintf("%v", ss.s)
}
