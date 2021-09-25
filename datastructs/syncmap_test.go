package datastructs

import (
	"math/rand"
	"testing"
)

func TestBasicSyncedMap(t *testing.T) {
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

func TestAdvancedSyncedMap(t *testing.T) {
	size := 10000
	sm := NewSyncedMap()

	for i := 0; i < size; i++ {
		sm.Add(i, "foo")
	}

	// expected size ?
	if sm.Len() != size {
		t.Error("SyncedMap is of wrong size")
	}

	// changing the data in the map
	for i := 0; i < size; i++ {
		sm.Add(i, "bar")
	}

	// re-testing that size did not change
	if sm.Len() != size {
		t.Error("SyncedMap is of wrong size")
	}

	// Testing that modification worked
	for i := 0; i < size; i++ {
		if v, ok := sm.Get(i); !ok {
			t.Error("This key must be there")
			if v.(string) != "bar" {
				t.Error("Wrong value")
			}
		}
	}

	ndel := 0
	for k := range sm.Keys() {
		if rand.Int()%2 == 0 {
			sm.Del(k)
			ndel++
		}
	}

	// testing size after deletions
	if sm.Len() != size-ndel {
		t.Error("Syncedmap has wrong size")
	}

	for _, v := range sm.Values() {
		if v.(string) != "bar" {
			t.Error("Wrong value")
		}
	}

}
