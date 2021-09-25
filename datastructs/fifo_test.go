package datastructs

import "testing"

func TestFifoBasic(t *testing.T) {
	f := &Fifo{}
	for i := 0; i < 10; i++ {
		f.Push(i)
		t.Logf("Size: %d", f.Len())
	}
	for p := f.Pop(); p != nil; p = f.Pop() {
		t.Logf("popped: %s", p)
		t.Logf("Size: %d", f.Len())
	}
	// second loops
	for i := 0; i < 10; i++ {
		f.Push(i)
		t.Logf("Size: %d", f.Len())
	}

	t.Logf("Fifo: %s", f)
	for p := f.Pop(); p != nil; p = f.Pop() {
		t.Logf("popped: %s", p)
		t.Logf("Size: %d", f.Len())
	}

	t.Logf("Size: %d", f.Len())
	if !f.Empty() {
		t.Error("Fifo should be empty")
	}
}
