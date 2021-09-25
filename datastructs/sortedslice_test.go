package datastructs

import (
	"math/rand"
	"testing"
	"time"
)

type MyInt int

func (m MyInt) Less(other Sortable) bool {
	return m < other.(MyInt)
}

type MyTime struct {
	t time.Time
}

func (m MyTime) Less(other Sortable) bool {
	return m.t.Before(other.(MyTime).t)
}

func (m MyTime) String() string {
	return m.t.Format(time.RFC3339Nano)
}

var (
	ints = [...]int{10, 13, 1, 2, 12, 99, 100, 102, 103, 103, 100, 1100, -2, -4}
)

func TestInsert(t *testing.T) {
	s := NewSortedSlice()
	for _, i := range ints {
		s.Insert(MyInt(i))
		t.Log(s)
		if !s.Control() {
			t.Fail()
		}
	}
	t.Log(s)
}

func TestRandom(t *testing.T) {
	s := NewSortedSlice()
	for i := 0; i < 1000; i++ {
		s.Insert(MyInt(rand.Int()))
	}
	if !s.Control() {
		t.Fail()
	}
}

func TestFail(t *testing.T) {
	s := NewSortedSlice()
	fail := [...]int{937, 821, 551, 410, 51, 320}
	for _, i := range fail {
		s.Insert(MyInt(i))
	}
	t.Log(s)
}

func TestTime(t *testing.T) {
	now := time.Now()
	s := NewSortedSlice()
	for i := 0; i < 50; i++ {
		mt := now.Add(time.Minute * time.Duration((rand.Int63() % 60)))
		mt.Add(time.Second * time.Duration((rand.Int63() % 60)))

		s.Insert(MyTime{mt})
	}
	if !s.Control() {
		t.Fail()
	}
	for mt := range s.ReversedIter() {
		t.Logf("%s", mt)
	}
}

func TestSearchRange(t *testing.T) {
	s := NewSortedSlice()
	for _, i := range ints {
		s.Insert(MyInt(i))
		t.Log(s)
		if !s.Control() {
			t.Fail()
		}
	}
	t.Log(s)
	t.Log(s.RangeLessThan(MyInt(0)))
}

func TestIter(t *testing.T) {
	s := NewSortedSlice()
	for _, i := range ints {
		s.Insert(MyInt(i))
		t.Log(s)
		if !s.Control() {
			t.Fail()
		}
	}
	t.Log(s)
	t.Log("Iter")
	for m := range s.Iter(0, 3) {
		t.Log(m)
	}
	t.Log("ReversedIter")
	for m := range s.ReversedIter(0, 3) {
		t.Log(m)
	}
}
