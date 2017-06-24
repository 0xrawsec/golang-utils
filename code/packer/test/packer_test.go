package main

import (
	"os"
	"testing"

	"github.com/0xrawsec/golang-utils/code/packer"
	"github.com/0xrawsec/golang-utils/crypto/data"
	"github.com/0xrawsec/golang-utils/crypto/file"
)

var (
	dataFile = "/bin/ls"
)

func TestPacker(t *testing.T) {
	file, err := os.Open(dataFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := packer.Packer{}
	p.AddResourceReader("/bin/ls", file)
	p.Dump("resources", "resources.go")
}

func TestUnpacker(t *testing.T) {
	d, err := resources.Resources.GetResource(dataFile)
	if err != nil {
		panic(err)
	}
	if data.Md5(d) != file.Md5(dataFile) {
		t.Fail()
	}
}
