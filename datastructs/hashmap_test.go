package datastructs

import (
	"fmt"
	"math/rand"
	"testing"
)

type IntHashable int

func (i IntHashable) Hash() string {
	return fmt.Sprintf("%d", i)
}

type TestHashable struct {
	ID  string
	Map map[string]interface{}
}

func (t TestHashable) Hash() string {
	return t.ID
}

func TestBasicHashMap(t *testing.T) {
	hm := NewSyncedHashMap()
	altered := TestHashable{"altered", make(map[string]interface{})}
	hm.Add(TestHashable{"it", make(map[string]interface{})}, "blop")
	hm.Add(TestHashable{"works", make(map[string]interface{})}, 42.0)
	hm.Add(TestHashable{"very", make(map[string]interface{})}, int64(42))
	hm.Add(TestHashable{"nice", make(map[string]interface{})}, uint(2))
	t.Log("Printing values refered by keys")
	for k := range hm.Keys() {
		t.Log(hm.Get(k))
		k = altered
	}

	t.Log("Look for altered key")
	if hm.Contains(altered) {
		t.Log("Keys are modifiable and it is not good")
		t.Fail()
	}

	t.Log("Printing only values")
	for v := range hm.Values() {
		t.Log(v)
	}
}

func TestStressHashMap(t *testing.T) {
	size := 1000
	hm := NewSyncedHashMap()

	for i := 0; i < size; i++ {
		hm.Add(IntHashable(i), i)
	}

	if hm.Len() != size {
		t.Error("Hashmap has wrong size")
	}

	del := 0
	for item := range hm.Items() {
		if item.Key.(IntHashable) != IntHashable(item.Value.(int)) {
			t.Error("Wrong item")
		}
		// deleting item
		if rand.Int()%2 == 0 {
			hm.Del(item.Key)
			del++
		}
	}

	if hm.Len() != size-del {
		t.Error("Hashmap has wrong size after deletions")
	}

}
