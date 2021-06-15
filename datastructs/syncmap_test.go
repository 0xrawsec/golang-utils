package datastructs

import (
	"testing"
)

func TestSyncedMap(t *testing.T) {
	sm := NewSyncedMap()
	sm.Add("foo", 2)
	sm.Add(2, "foo")
	if !sm.Contains(2) {
		t.Error("Map should contain 2")
	}
	if !sm.Contains("foo") {
		t.Error("Map should contain foo")
	}

	sm.Del(2)
	sm.Del("foo")
	if sm.Contains(2) {
		t.Error("Map should not contain 2")
	}
	if sm.Contains("foo") {
		t.Error("Map should not contain foo")
	}
}
