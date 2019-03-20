package datastructs

import (
	"math/rand"
	"testing"
)

func TestBitSetBasic(t *testing.T) {
	bs := NewBitSet(255)
	offset := 10

	bs.Set(offset)
	if !bs.Get(offset) {
		t.Logf("Failed to retrieve bit at offset: %d", offset)
		t.FailNow()
	}
}

func TestBitSetRookie(t *testing.T) {
	bs := NewBitSet(1013)
	for i := 0; i < bs.Len(); i++ {
		if i%2 == 0 {
			bs.Set(i)
		}
	}

	for i := 0; i < bs.Len(); i++ {
		if i%2 == 0 {
			if !bs.Get(i) {
				t.Logf("Failed to retrieve bit at offset: %d", i)
				t.FailNow()
			}
		}
	}
}

func TestBitSetHardcore(t *testing.T) {
	for i := 0; i < 10000; i++ {
		size := rand.Uint32() % (10 * 1024)
		bs := NewBitSet(int(size))
		for i := 0; i < bs.Len(); i++ {
			if i%2 == 0 {
				bs.Set(i)
			}
		}

		for i := 0; i < bs.Len(); i++ {
			if i%2 == 0 {
				if !bs.Get(i) {
					t.Logf("Failed to retrieve bit at offset: %d", i)
					t.FailNow()
				}
			}
		}
	}
}
