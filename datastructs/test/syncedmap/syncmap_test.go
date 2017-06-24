package main

import (
	"testing"

	"github.com/0xrawsec/golang-utils/datastructs"
)

func TestSyncedMap(t *testing.T) {
	sm := datastructs.NewSyncedMap()
	sm.Add("foo", 2)
	sm.Add(2, "foo")
	t.Log(sm.Contains(2))
	t.Log(sm.Contains("foo"))
}
