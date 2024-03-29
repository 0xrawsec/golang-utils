package datastructs

import (
	"encoding/json"
	"testing"
)

func TestSets(t *testing.T) {
	s1 := NewSyncedSet()
	s2 := NewSyncedSet()
	s1.Add("This", "is", "foo", "!!", "!!!!")
	s2.Add("This", "is", "bar", "!!!", "!!!!!!!")
	s1copy := NewSyncedSet(s1)

	intersection := s1.Intersect(s2)
	union := s1.Union(s2)

	t.Logf("s1.Slice: %v", s1.Slice())
	t.Logf("s2.Slice: %v", s2.Slice())
	t.Logf("s1.Intersect(s2).Slice: %v", intersection.Slice())
	t.Logf("s1.Union(s2).Slice: %v", union.Slice())

	if !s1.Contains("This", "is", "foo") {
		t.Error("string missing")
	}
	if !intersection.Contains("This", "is") {
		t.Error("string missing")
	}
	if !union.Contains("This", "is", "foo", "bar", "!!", "!!!", "!!!!", "!!!!!!!") {
		t.Error("string missing")
	}
	union.Del("This", "is", "foo", "!!", "!!!!")
	if union.Contains("This") || union.Contains("is") || union.Contains("foo") {
		t.Error("string should be missing")
	}
	t.Logf("union after delete: %v", union.Slice())

	if !s1.Equal(s1copy) {
		t.Error("equality test failed")
	}

	if s1.Len() != s1copy.Len() {
		t.Error("length is not equal between original and copy")
	}

	for it := range union.Items() {
		t.Logf("%v", it)
	}
}
func TestSetJSON(t *testing.T) {
	var data []byte
	var err error

	s1 := NewSyncedSet()
	s1.Add("This", "is", "bar", "!!!", "!!!!!!!")

	if data, err = json.Marshal(&s1); err != nil {
		t.Error("Failed to marshal JSON")
	} else {
		t.Log(string(data))
	}

	s2 := NewSyncedSet()
	if err = json.Unmarshal(data, &s2); err != nil {
		t.Errorf("Failed to unmarshal JSON: %s", err)
		t.FailNow()
	}

	if !s2.Contains("This", "is", "bar", "!!!", "!!!!!!!") {
		t.Error("Set does not contain expected data")
	}
}

func TestSetOrder(t *testing.T) {
	size := 10000
	s := NewSet()
	for i := 0; i < size; i++ {
		s.Add(i)
	}

	ss := s.SortSlice()
	for i := 0; i < size; i++ {
		if ss[i].(int) != i {
			t.Error("Bad set order")
		}
	}
}

func TestJSONMarshalStability(t *testing.T) {
	// aims at testing that order of elements in serialization is stable
	var data, prev []byte
	var err error

	s1 := NewSyncedSet()
	s1.Add("This", "is", "bar", "!!!", "!!!!!!!")

	for i := 0; i < 100; i++ {
		if data, err = json.Marshal(&s1); err != nil {
			t.Error("Failed to marshal JSON")
		}
		if prev == nil {
			goto copy
		}
		if string(data) != string(prev) {
			t.Error("JSON serialization is not stable")
		}
	copy:
		prev = make([]byte, len(data))
		copy(prev, data)
	}
}
