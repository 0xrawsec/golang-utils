package datastructs

import (
	"testing"
)

func TestSyncedMap(t *testing.T) {
	sm := NewSyncedMap()
	sm.Add("foo", 2)
	sm.Add(2, "foo")
	t.Log(sm.Contains(2))
	t.Log(sm.Contains("foo"))
}
