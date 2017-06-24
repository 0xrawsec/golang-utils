package ngram

import (
	"errors"
	"hash/fnv"
	"io"
	"os"
)

type Ngram []byte

type FastGenerator struct {
	Init   bool
	Reader io.ReadSeeker
	Ngram  Ngram
}

// ErrBadNgramSize is raised when the ngram size is not in the correct range
var ErrBadNgramSize = errors.New("ErrBadNgramSize: ngram size must be in ]0;MAXINT]")

// New new ngram from buffer
func NewNgram(buf []byte) Ngram {
	ngram := make(Ngram, len(buf))
	copy(ngram, buf)
	return ngram
}

// Hash hashes a ngram
func (ngram *Ngram) Hash() uint64 {
	h := fnv.New64a()
	h.Write((*ngram)[:])
	return h.Sum64()
}

// Generator generates Ngrams of size sNgram from a file
func Generator(reader io.Reader, sNgram int) chan Ngram {
	if sNgram <= 0 {
		panic(ErrBadNgramSize)
	}
	yield := make(chan Ngram)

	go func() {
		feed(yield, reader, sNgram)
	}()

	return yield
}

func feed(generator chan Ngram, reader io.Reader, sNgram int) {
	defer close(generator)
	buf := make([]byte, 4096)
	ngram := make(Ngram, sNgram)
	read, err := reader.Read(ngram)
	if read < sNgram {
		return
	}
	switch err {
	case io.EOF, nil:
		generator <- NewNgram(ngram)
	default:
		panic(err)
	}
	for read, err := reader.Read(buf); err != io.EOF; read, err = reader.Read(buf) {
		for i := 0; i < read; i++ {
			copy(ngram, ngram[1:])
			ngram[len(ngram)-1] = buf[i]
			generator <- NewNgram(ngram)
		}
	}
}

func NewFastGenerator(reader io.ReadSeeker, sNgram int) (fg FastGenerator) {
	fg.Init = true
	fg.Reader = reader
	fg.Ngram = make(Ngram, sNgram)
	return fg
}

func (fg *FastGenerator) Next() (err error) {
	if !fg.Init {
		fg.Reader.Seek(int64(-(len(fg.Ngram) - 1)), os.SEEK_CUR)
	} else {
		fg.Init = false
	}
	read, err := fg.Reader.Read(fg.Ngram)
	if read < len(fg.Ngram) {
		return ErrBadNgramSize
	}
	if err != nil {
		return err
	}
	return nil
}
