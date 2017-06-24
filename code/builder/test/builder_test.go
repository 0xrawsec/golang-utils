package main

import (
	"testing"

	"github.com/0xrawsec/golang-utils/code/builder"
	"github.com/0xrawsec/golang-utils/log"
)

type Matcher struct {
	Level  int
	Offset int
	Type   byte
	Range  int
	Mask   []byte
	Value  []byte
	Flags  int
}

func init() {
	log.InitLogger(log.LDebug)
}

func NewMatcher() (m Matcher) {
	m.Mask = []byte("foo")
	m.Value = []byte("bar")
	return
}

func TestMapBuilder1(t *testing.T) {
	m := map[string]bool{
		"foo":  true,
		"blop": false}

	cb := builder.CodeBuilder{}
	cb.Package("bar")
	cb.DefVariable("test", m)
	cb.ResolveImports()
	t.Log(cb.String())
}

func TestMapBuilder2(t *testing.T) {
	m := map[string][]byte{
		"foo": []byte("bar")}
	cb := builder.CodeBuilder{}
	cb.Package("blop")
	cb.DefVariable("test", m)
	cb.ResolveImports()
	t.Log(cb.String())
}

func TestStructBuilder(t *testing.T) {
	m := NewMatcher()
	cb := builder.CodeBuilder{}
	cb.Package("foo")
	cb.DefVariable("bar", m)
	t.Log(cb.String())
	//cb.ResolveImports()
}
