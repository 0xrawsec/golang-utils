package main

import (
	"testing"

	"github.com/0xrawsec/golang-utils/datastructs"
)

type TestHashable struct {
	ID  string
	Map map[string]interface{}
}

func (t TestHashable) Hash() string {
	return t.ID
}

func TestHashMap(t *testing.T) {
	hm := datastructs.NewHashMap()
	altered := TestHashable{"altered", make(map[string]interface{})}
	hm.Set(TestHashable{"it", make(map[string]interface{})}, "blop")
	hm.Set(TestHashable{"works", make(map[string]interface{})}, 42.0)
	hm.Set(TestHashable{"very", make(map[string]interface{})}, int64(42))
	hm.Set(TestHashable{"nice", make(map[string]interface{})}, uint(2))
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
