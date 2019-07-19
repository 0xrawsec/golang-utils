package datastructs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/0xrawsec/golang-utils/log"
)

type Element struct {
	Value interface{}
	Prev  *Element
	Next  *Element
}

func (e *Element) String() string {
	return fmt.Sprintf("(%T(%v), %p, %p)", e.Value, e.Value, e.Prev, e.Next)
}

type Fifo struct {
	sync.RWMutex
	e    *Element
	last *Element
	size int
}

func (f *Fifo) Push(i interface{}) {
	f.Lock()
	defer f.Unlock()
	e := Element{Value: i}
	if f.e == nil {
		f.e = &e
		f.last = &e
	} else {
		e.Next = f.e
		f.e.Prev = &e
		f.e = &e
	}
	f.size++
}

func (f *Fifo) String() string {
	f.RLock()
	defer f.RUnlock()
	out := make([]string, 0)
	for e := f.e; e != nil; e = e.Next {
		log.Info(e)
		out = append(out, e.String())
	}
	return strings.Join(out, "->")
}

func (f *Fifo) Empty() bool {
	f.RLock()
	defer f.RUnlock()
	return f.size == 0
}

func (f *Fifo) Pop() *Element {
	f.Lock()
	defer f.Unlock()
	if f.last == nil {
		return nil
	}

	popped := f.last
	f.last = f.last.Prev
	if f.last != nil {
		f.last.Next = nil
	}
	// we have to nil out f.e if we pop
	// the last element of the Fifo
	if f.size == 1 {
		f.e = nil
	}
	f.size--
	return popped
}

func (f *Fifo) Len() int {
	f.RLock()
	defer f.RUnlock()
	return f.size
}
