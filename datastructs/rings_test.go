package datastructs

import (
	"encoding/json"
	"math/rand"
	"testing"
)

func TestRingSlice(t *testing.T) {
	r := NewRingSlice(10)
	for i := 0; i < 11; i++ {
		r.Add(i)
	}
	if r.GetItem(0).(int) != 10 {
		t.Error("Bad item at index 0")
	}
	if r.GetItem(r.Len()-1) != 9 {
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
	if r.GetItem(0).(float64) != 1 {
		t.Error("Bad item at index 0")
	}
	if r.GetItem(r.Len()-1).(float64) != 10 {
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

func TestRingSet(t *testing.T) {
	r := NewRingSet(10)
	for i := 0; i < 100; i++ {
		if r.Contains(i) {
			t.Errorf("RingSet should not contain value before being added: %d", i)
		}
		r.Add(i)
		if !r.Contains(i) {
			t.Errorf("RingSet should contain value just after being added: %d", i)
		}
	}
	t.Log(r)

	// we do some modifications on copies that should trigger error
	// if we modify original structures
	rs := r.RingSlice()
	set := r.Set()
	for i := 0; i < 10; i++ {
		rs.Add(i)
		set.Add(i)
	}

	for i := 0; i < 100; i++ {
		if i < 90 {
			// ring set should not contain those values
			if r.Contains(i) {
				t.Errorf("RingSet should not contain value: %d", i)
			}
		} else {
			if !r.Contains(i) {
				t.Errorf("RingSet should contain value: %d", i)
			}
		}
	}

	if r.rslice.Len() != r.set.Len() {
		t.Errorf("RingSlice and Set must have the same size")
	}

	if r.set.Len() != r.Len() {
		t.Errorf("Inconsistent size")
	}

	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("Failed to marshal RingSet: %s", err)
		t.FailNow()
	}
	t.Log(string(b))

	new := NewRingSet(0)
	if err = json.Unmarshal(b, &new); err != nil {
		t.Errorf("Failed to unmarshall RingSet: %s", err)
		t.FailNow()
	}
	t.Log(new)

	// json unmarshal integers as float64
	for i := float64(0); i < 100; i++ {
		if i < 90 {
			// ring set should not contain those values
			if new.Contains(i) {
				t.Errorf("RingSet should not contain value: %f", i)
			}
		} else {
			if !new.Contains(i) {
				t.Errorf("RingSet should contain value: %f", i)
			}
		}
	}

	if new.rslice.Len() != new.set.Len() {
		t.Errorf("RingSlice and Set must have the same size, even after json un/marshalling")
	}
}

func TestRingSetNestedJSON(t *testing.T) {
	type T struct {
		R *RingSet `json:"r"`
	}

	ts := T{NewRingSet(10)}
	data, err := json.Marshal(&ts)
	if err != nil {
		t.Errorf("Failed to marshal nested structure: %s", err)
	}
	t.Log(string(data))
}

func TestRingGetSet(t *testing.T) {
	size := 100
	r := NewRingSet(100)
	for i := 0; i < r.Len(); i++ {
		r.SetItem(i, rand.Int())
	}

	s := r.Slice()
	for i := 0; i < size; i++ {
		item := r.GetItem(i)
		if item.(int) != s[i].(int) {
			t.Error("Wrong value returned by GetItem")
		}
	}
}

func TestRingSetCopy(t *testing.T) {
	size := 100
	r := NewRingSet(100)
	for i := 0; i < r.Len(); i++ {
		r.SetItem(i, rand.Int())
	}

	copy := r.Copy()
	for i := 0; i < size; i++ {
		if r.GetItem(i).(int) != copy.GetItem(i).(int) {
			t.Error("Wrong value returned by GetItem")
		}
	}
}
