package datastructs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/0xrawsec/golang-utils/log"
)

type Element struct {
	value interface{}
	prev  *Element
	next  *Element
}

func (e *Element) String() string {
	return fmt.Sprintf("(%T(%v), %p, %p)", e.value, e.value, e.prev, e.next)
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
	e := Element{value: i}
	if f.e == nil {
		f.e = &e
		f.last = &e
	} else {
		e.next = f.e
		f.e.prev = &e
		f.e = &e
	}
	f.size++
}

func (f *Fifo) String() string {
	f.RLock()
	defer f.RUnlock()
	out := make([]string, 0)
	for e := f.e; e != nil; e = e.next {
		log.Info(e)
		out = append(out, e.String())
	}
	return strings.Join(out, "->")
}

func (f *Fifo) Empty() bool {
	f.RLock()
	defer f.RUnlock()
	return f.last == nil
}

func (f *Fifo) Pop() *Element {
	f.Lock()
	defer f.Unlock()
	if f.last == nil {
		return nil
	}

	popped := f.last
	f.last = f.last.prev
	if f.last != nil {
		f.last.next = nil
	}
	f.size--
	return popped
}

func (f *Fifo) Len() int {
	f.RLock()
	defer f.RUnlock()
	return f.size
}
