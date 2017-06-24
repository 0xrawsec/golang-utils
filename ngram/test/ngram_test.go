package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/0xrawsec/golang-utils/ngram"
)

var (
	SizeNgram = 3
	TestFile  = "1M.bin"
)

func BenchmarkNgramOnFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, err := os.Open(TestFile)
		if err != nil {
			panic(err)
		}
		ng := ngram.Generator(file, SizeNgram)
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
		fng := ngram.NewFastGenerator(file, SizeNgram)
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
		ng := ngram.Generator(bytes.NewReader(dat), SizeNgram)
		for _ = range ng {
			//fmt.Println(string(ngram))
		}
	}
}
