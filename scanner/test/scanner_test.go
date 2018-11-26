package main

import (
	"scanner"
	"strings"
	"testing"
)

const (
	colontext   = `This:is:a:text:separated:by:colon:not:parseable:by:text/scanner:go:package`
	nlcolontext = "This:is\na:text:separated\nby:colon\nor:newline\nnot:parseable\nby:text/scanner:go\npackage"
)

func TestScannerBasic(t *testing.T) {
	s := scanner.New(strings.NewReader(colontext))
	s.InitWhitespace(":")
	for r := s.Scan(); r != scanner.EOF; r = s.Scan() {
		switch r {
		case ':':
			break
		default:
			t.Log(s.TokenText())
		}
	}
	t.Log(s.TokenText())
}

func TestScannerBasic2(t *testing.T) {
	s := scanner.New(strings.NewReader(nlcolontext))
	s.InitWhitespace(":\n")
	for r := s.Scan(); r != scanner.EOF; r = s.Scan() {
		switch r {
		case ':', '\n':
			break
		default:
			t.Log(s.TokenText())
		}
	}
	t.Log(s.TokenText())
}
