package readers

import (
	"bufio"
	"io"
)

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
