package readers

import (
	"io"
	"os"
	"strings"
	"testing"

	"text/scanner"
)

const (
	text = `This is a text containing
	lines that should be
	printed in reversed
	order.`

	testfile        = "../LICENSE"
	testfileCharCnt = 35141
	testfileLineCnt = 674
)

func TestReverseReaderBasic(t *testing.T) {
	var c [4]byte
	r := strings.NewReader(text)
	rr := NewReverseReader(r)
	for _, err := rr.Read(c[:]); err != io.EOF; _, err = rr.Read(c[:]) {
		t.Logf("%q", c)
	}
}

func TestReverseReaderReadRune(t *testing.T) {
	r := strings.NewReader(text)
	rr := NewReverseReader(r)
	for ru, _, err := rr.ReadRune(); err != io.EOF; ru, _, err = rr.ReadRune() {
		t.Logf("%q", ru)
	}
}

func TestReversedReadline(t *testing.T) {
	r := strings.NewReader(text)
	for line := range ReversedReadlines(r) {
		t.Log(string(line))
	}
}

func TestReversedReadlineOnFile(t *testing.T) {
	var lCnt, cCnt int
	fd, err := os.Open(testfile)
	if err != nil {
		t.Logf("Failed to open test file: %s", testfile)
		t.FailNow()
	}
	defer fd.Close()
	for line := range ReversedReadlines(fd) {
		// +1 for \n removed
		cCnt += len(line) + 1
		lCnt++
	}
	t.Logf("#char: %d #lines: %d", cCnt, lCnt)
	// we add one to the original counters because last line of the file is counted
	// as a line while it is not really since there is no \n
	if lCnt != testfileLineCnt+1 || cCnt != testfileCharCnt+1 {
		t.Log("Bad number of lines or characters")
		t.FailNow()
	}
}

func TestReverseReaderAdvanced(t *testing.T) {
	r := strings.NewReader(text)
	rr := NewReverseReader(r)
	s := scanner.Scanner{}
	s.Init(rr)
	s.Whitespace = 0
	s.Whitespace ^= 0x1 << '\n'
	for ru := s.Scan(); ru != scanner.EOF; ru = s.Scan() {
		switch ru {
		case '\n':
			t.Logf("%q", s.TokenText())
		default:
			t.Logf("%q", s.TokenText())
		}
	}
}
