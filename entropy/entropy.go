package entropy

import (
	"bufio"
	"bytes"
	"io"
	"math"
)

func readerEntropy(r io.Reader) float64 {
	var b byte
	var err error
	var sum, entropy float64
	counter := [256]float64{}
	buff := bufio.NewReader(r)
	for {
		b, err = buff.ReadByte()
		if err == io.EOF {
			break
		}
		counter[b]++
		sum++
	}
	if sum > 0 {
		for _, count := range counter {
			if count > 0 {
				p := count / sum
				entropy += p * math.Log2(p)
			}
		}
	}
	if entropy == 0 {
		return 0
	}
	return -entropy
}

// ReaderEntropy computes the entropy of a reader
func ReaderEntropy(r io.Reader) float64 {
	return readerEntropy(r)
}

// BytesEntropy computes the entropy of a buffer
func BytesEntropy(b []byte) float64 {
	return readerEntropy(bytes.NewReader(b))
}

// StringEntropy computes the entropy of a string
func StringEntropy(s string) float64 {
	return readerEntropy(bytes.NewReader([]byte(s)))
}
