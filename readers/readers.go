package readers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

// ReverseReader reads bytes from the end of reader
type ReverseReader struct {
	rs     io.ReadSeeker
	size   int64
	offset int64
}

// NewReverseReader creates a new ReverseReader
func NewReverseReader(rs io.ReadSeeker) *ReverseReader {
	rr := ReverseReader{}
	rr.rs = rs
	s, err := rs.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	rr.size = s
	rr.offset = s
	return &rr
}

// Read implements io.Reader interface
// very likely bad performances due to the seeks
func (r *ReverseReader) Read(p []byte) (n int, err error) {
	//fmt.Printf("Offset %d\n", r.offset)
	//fmt.Printf("Reading %d bytes\n", len(p))
	switch {
	case r.offset <= 0:
		return 0, io.EOF
	case r.offset-int64(len(p)) <= 0:
		r.rs.Seek(0, os.SEEK_SET)
		n, err = r.rs.Read(p[:r.offset])
		r.rs.Seek(0, os.SEEK_SET)
		r.offset = 0
		return n, nil
	default:
		r.offset -= int64(len(p))
		r.rs.Seek(r.offset, os.SEEK_SET)
		n, err = r.rs.Read(p)
		r.rs.Seek(r.offset, os.SEEK_SET)
		return
	}
}

// ReadRune reads a rune backward
func (r *ReverseReader) ReadRune() (ru rune, size int, err error) {
	var rb [4]byte
	n, err := r.Read(rb[:])
	ru, size = utf8.DecodeLastRune(rb[:n])
	if err != nil {
		return
	}
	r.offset += int64(n - size)
	r.rs.Seek(r.offset, os.SEEK_SET)
	if ru == utf8.RuneError {
		return ru, size, fmt.Errorf("RuneError")
	}
	return
}

/***************************************************/
/************** Readlines functions ****************/
/***************************************************/

func min(i, k int) int {
	switch {
	case i < k:
		return i
	case k < i:
		return k
	default:
		return k
	}
}

func reversedCopy(dst, src []byte) {
	m := min(len(src), len(dst))
	for i, k := 0, m-1; i < m; i++ {
		dst[i] = src[k]
		k--
	}
}

// ReversedReadlines returns the lines found in a reader in reversed order
func ReversedReadlines(r io.ReadSeeker) (lines chan []byte) {
	lines = make(chan []byte)
	go func() {
		defer close(lines)

		var c [1]byte
		rr := NewReverseReader(r)
		line := make([]byte, 0, 4096)
		for n, err := rr.Read(c[:]); err != io.EOF && n != 0; n, err = rr.Read(c[:]) {
			if c[0] == '\n' {
				cpLine := make([]byte, len(line))
				reversedCopy(cpLine, line)
				lines <- cpLine
				line = make([]byte, 0, 4096)
			} else {
				// don't append newline
				line = append(line, c[0])
			}
		}
		// process the last line
		cpLine := make([]byte, len(line))
		reversedCopy(cpLine, line)
		lines <- cpLine
	}()
	return
}

// Readlines : returns a channel containing the lines of the reader
func Readlines(reader io.Reader) (generator chan []byte) {
	generator = make(chan []byte)
	go func() {
		defer close(generator)
		lreader := bufio.NewReader(reader)
		for line, isPrefix, err := lreader.ReadLine(); err != io.EOF; {
			fullLine := make([]byte, len(line))
			copy(fullLine, line)
			for isPrefix == true {
				line, isPrefix, err = lreader.ReadLine()
				fullLine = append(fullLine, line...)
			}
			generator <- fullLine
			line, isPrefix, err = lreader.ReadLine()
		}
	}()
	return generator
}

// ReadlinesUTF16 : returns a channel of []byte lines trimmed. To be used with a reader containing UTF16 encoded runes
func ReadlinesUTF16(reader io.Reader) (generator chan []byte) {
	generator = make(chan []byte)
	go func() {
		var utf16Rune [2]byte
		var r rune
		var u [1]uint16

		// closing generator
		defer close(generator)

		line := make([]byte, 0, 4096)
		i := 0
		for read, err := reader.Read(utf16Rune[:]); read == 2 && err == nil; {
			u[0] = *(*uint16)(unsafe.Pointer(&utf16Rune[0]))
			// Filter out UTF16 BOM
			if u[0] == 0xfeff && i == 0 {
				goto tail
			}

			r = utf16.Decode(u[:])[0]
			if r == '\n' {
				fullLine := make([]byte, len(line))
				copy(fullLine, line)
				fullLine = bytes.TrimRight(fullLine, "\r\n")
				generator <- fullLine
				line = make([]byte, 0, 4096)
			} else {
				line = append(line, byte(r))
			}

		tail:
			read, err = reader.Read(utf16Rune[:])
			i++

		}
	}()
	return generator
}
