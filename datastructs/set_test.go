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
	s1copy := NewSyncedSet(&s1)

	intersection := s1.Intersect(&s2)
	union := s1.Union(&s2)

	t.Logf("s1.List: %v", s1.List())
	t.Logf("s2.List: %v", s2.List())
	t.Logf("s1.Intersect(s2).List: %v", intersection.List())
	t.Logf("s1.Union(s2).List: %v", union.List())

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
	t.Logf("union after delete: %v", union.List())

	if !s1.Equal(&s1copy) {
		t.Error("equality test failed")
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
