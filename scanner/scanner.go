package scanner

import (
	"bufio"
	"datastructs"
	"fmt"
	"io"
	"os"
)

const (
	// EOF scanner
	EOF = -(iota + 1)
	// MaxTokenLen length of a token
	MaxTokenLen = 1 << 20 // 1 Mega Token
)

var (
	// Whitespace default whitespaces characters taken by scanner
	Whitespace = "\t\n\r "
)

// Scanner structure definition
type Scanner struct {
	Offset     int64
	LastOffset int64
	Whitespace *datastructs.BitSet
	Error      func(error)
	r          *bufio.Reader
	tokenIdx   int
	token      []rune
}

// New creates a new scanner from reader
func New(r io.Reader) (s *Scanner) {
	s = &Scanner{}
	s.r = bufio.NewReader(r)
	s.Whitespace = datastructs.NewBitSet(256)
	s.Error = func(err error) {
		fmt.Fprintf(os.Stderr, "Scanner error: %s\n", err)
	}
	// initialize default with Whitespace variable
	s.InitWhitespace(Whitespace)
	s.token = make([]rune, MaxTokenLen)
	s.tokenIdx = 0
	return
}

// InitWhitespace initialised an new set of whitespaces for the scanner
func (s *Scanner) InitWhitespace(w string) {
	for _, c := range w {
		s.Whitespace.Set(int(c))
	}
}

// Scan scans until we reach a whitespace token
func (s *Scanner) Scan() (r rune) {
	s.tokenIdx = 0
	r, _, err := s.r.ReadRune()
	prevRune := r

	if s.Whitespace.Get(int(r)) {
		return r
	}

	for ; !s.Whitespace.Get(int(r)); r, _, err = s.r.ReadRune() {
		switch err {
		case nil:
			break
		case io.EOF:
			return EOF
		default:
			s.Error(err)
			return EOF
		}
		s.token[s.tokenIdx] = r
		s.tokenIdx++
		prevRune = r
	}

	// We have to UnreadRune because we went too far of one rune
	err = s.r.UnreadRune()
	if err != nil {
		s.Error(err)
	}

	return prevRune
}

// TokenText returns the string containing characters until the token was found
func (s *Scanner) TokenText() string {
	return string(s.token[:s.tokenIdx])
}
