package datastructs

import (
	"encoding/json"
	"fmt"
)

type RingSet struct {
	rslice *RingSlice
	set    *Set
}

func NewRingSet(len int) *RingSet {
	rs := RingSet{NewRingSlice(len), NewSet()}
	return &rs
}

func (r RingSet) String() string {
	return r.rslice.String()
}

func (r *RingSet) Contains(item ...interface{}) bool {
	return r.set.Contains(item...)
}

func (r *RingSet) Add(item interface{}) {
	// we add item only if not already there
	if !r.Contains(item) {
		// we delete item only if RingSet is full
		if r.rslice.full {
			// delete the item which is going to be erased
			r.set.Del(r.rslice.GetItem(r.rslice.cursor))
		}
		r.rslice.Add(item)
		r.set.Add(item)
	}
}

func (r *RingSet) Len() int {
	return r.rslice.Len()
}

func (r *RingSet) GetItem(i int) interface{} {
	return r.rslice.GetItem(i)
}

func (r *RingSet) SetItem(i int, item interface{}) {
	r.rslice.SetItem(i, item)
	r.set.Add(item)
}

func (r *RingSet) Slice() []interface{} {
	return r.rslice.Slice()
}

func (r *RingSet) Copy() *RingSet {
	new := NewRingSet(r.Len())
	new.rslice = r.rslice.Copy()
	new.set = r.set.Copy()
	return new
}

func (r *RingSet) RingSlice() *RingSlice {
	return r.rslice.Copy()
}

func (r *RingSet) Set() *Set {
	return r.set.Copy()
}

func (r *RingSet) UnmarshalJSON(data []byte) (err error) {
	if err = json.Unmarshal(data, &r.rslice); err != nil {
		return
	}
	r.set = NewInitSet(r.rslice.Slice()...)
	return
}

func (r *RingSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&(r.rslice))
}

type RingSlice struct {
	ring   []interface{}
	cursor int
	full   bool
}

func NewRingSlice(len int) *RingSlice {
	return &RingSlice{make([]interface{}, len), 0, false}
}

func (r *RingSlice) incCursor() {
	r.cursor = r.nextCursor()
	if r.cursor == 0 {
		r.full = true
	}
}

func (r *RingSlice) nextCursor() int {
	if r.cursor < len(r.ring)-1 {
		return r.cursor + 1
	}
	return 0
}

func (r RingSlice) String() string {
	return fmt.Sprintf("%v", r.ring)
}

func (r *RingSlice) Add(item interface{}) {
	if r.cursor < len(r.ring) {
		r.ring[r.cursor] = item
		r.incCursor()
	}
}

func (r *RingSlice) Len() int {
	return len(r.ring)
}

func (r *RingSlice) GetItem(i int) interface{} {
	return r.ring[i]
}

func (r *RingSlice) SetItem(i int, item interface{}) {
	r.ring[i] = item
}

func (r *RingSlice) Slice() []interface{} {
	l := make([]interface{}, len(r.ring))
	copy(l, r.ring)
	return l
}

func (r *RingSlice) Copy() *RingSlice {
	new := NewRingSlice(r.Len())
	copy(new.ring, r.ring)
	new.cursor = r.cursor
	return new
}

func (r *RingSlice) UnmarshalJSON(data []byte) (err error) {
	r.cursor = 0
	return json.Unmarshal(data, &r.ring)
}

func (r *RingSlice) MarshalJSON() ([]byte, error) {
	s := make([]interface{}, len(r.ring))
	if len(s) > 1 {
		p1 := r.ring[r.cursor:len(r.ring)]
		p2 := r.ring[0:r.cursor]
		copy(s, p1)
		copy(s[len(p1):], p2)
	} else {
		copy(s, r.ring)
	}
	return json.Marshal(s)
}
