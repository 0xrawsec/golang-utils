package datastructs

import "testing"

func shouldPanic(t *testing.T, f func()) {
	defer func() { recover() }()
	f()
	t.Errorf("should have panicked")
}

func TestToInterfaceSlice(t *testing.T) {
	size := 1000
	intSlice := make([]int, 0)

	for i := 0; i < size; i++ {
		intSlice = append(intSlice, i)
	}

	is := ToInterfaceSlice(intSlice)
	if len(is) != size {
		t.Error("Interface slice has wrong size")
	}

	shouldPanic(t, func() { ToInterfaceSlice(42) })
}
