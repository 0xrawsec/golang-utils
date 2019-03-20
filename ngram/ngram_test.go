package ngram

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var (
	SizeNgram = 3
	TestFile  = "./test/1M.bin"
)

func BenchmarkNgramOnFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, err := os.Open(TestFile)
		if err != nil {
			panic(err)
		}
		ng := Generator(file, SizeNgram)
		for _ = range ng {
			//fmt.Println(string(ngram))
		}
	}
}

func BenchmarkFastOnFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, err := os.Open(TestFile)
		if err != nil {
			panic(err)
		}
		fng := NewFastGenerator(file, SizeNgram)
		for err := fng.Next(); err == nil; err = fng.Next() {
			//fmt.Println(string(ngram))
			//fmt.Println(fng.Ngram)
		}
	}
}

func BenchmarkNgramOnData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dat, err := ioutil.ReadFile(TestFile)
		if err != nil {
			panic(err)
		}
		ng := Generator(bytes.NewReader(dat), SizeNgram)
		for _ = range ng {
			//fmt.Println(string(ngram))
		}
	}
}
