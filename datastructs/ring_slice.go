package datastructs

import (
	"encoding/json"
	"fmt"
)

type RingSlice struct {
	ring   []interface{}
	cursor int
}

func NewRingSlice(len int) RingSlice {
	return RingSlice{make([]interface{}, len), 0}
}

func (r RingSlice) String() string {
	return fmt.Sprintf("%v", r.ring)
}

func (r *RingSlice) Add(item interface{}) {
	if r.cursor < len(r.ring) {
		r.ring[r.cursor] = item
		if r.cursor < len(r.ring)-1 {
			r.cursor++
		} else {
			r.cursor = 0
		}
	}
}

func (r *RingSlice) Len() int {
	return len(r.ring)
}

func (r *RingSlice) Get(i int) interface{} {
	return r.ring[i]
}

func (r *RingSlice) Set(i int, item interface{}) {
	r.ring[i] = item
}

func (r *RingSlice) List() []interface{} {
	return r.ring
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
