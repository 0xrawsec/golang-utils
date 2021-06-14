package datastructs

import (
	"encoding/json"
	"testing"
)

func TestRingSlice(t *testing.T) {
	r := NewRingSlice(10)
	for i := 0; i < 11; i++ {
		r.Add(i)
	}
	if r.Get(0).(int) != 10 {
		t.Error("Bad item at index 0")
	}
	if r.Get(r.Len()-1) != 9 {
		t.Error("Bad last item 0")
	}
	t.Log(r)
}

func TestRingSliceJSON(t *testing.T) {
	r := NewRingSlice(10)
	for i := 0; i < 11; i++ {
		r.Add(i)
	}

	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("JSON marshalling failed")
	} else {
		t.Log(string(b))
	}

	r = NewRingSlice(0)
	if err := json.Unmarshal(b, &r); err != nil {
		t.Errorf("JSON unmarshalling failed")
	}
	// items have been reordered by json marshaling
	// and casted to another type float64
	if r.Get(0).(float64) != 1 {
		t.Error("Bad item at index 0")
	}
	if r.Get(r.Len()-1).(float64) != 10 {
		t.Error("Bad last item 0")
	}
	t.Log(r)
}

func TestRingSliceEmpty(t *testing.T) {
	for size := 0; size < 2; size++ {
		r := NewRingSlice(size)
		for i := 0; i < 100; i++ {
			r.Add(i)
		}

		b, err := json.Marshal(&r)
		if err != nil {
			t.Errorf("JSON marshalling failed")
		} else {
			t.Log(string(b))
		}

		if err := json.Unmarshal(b, &r); err != nil {
			t.Errorf("JSON unmarshalling failed")
		}
		t.Log(r)
	}
}
