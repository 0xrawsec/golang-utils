package main

import (
	"testing"

	"github.com/0xrawsec/golang-utils/datastructs"
)

func TestSets(t *testing.T) {
	s1 := datastructs.NewSyncedSet()
	s2 := datastructs.NewSyncedSet()
	s1.Add("This", "is", "foo", "!!", "!!!!")
	s2.Add("This", "is", "bar", "!!!", "!!!!!!!")
	s1copy := datastructs.NewSyncedSet(&s1)

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
}
