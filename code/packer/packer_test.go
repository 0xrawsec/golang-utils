package packer

import (
	"os"
	"testing"
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
	p := Packer{}
	p.AddResourceReader("/bin/ls", file)
	p.Dump("resources", "resources.go")
}

func TestUnpacker(t *testing.T) {
	/*d, err := resources.Resources.GetResource(dataFile)
	if err != nil {
		panic(err)
	}
	if data.Md5(d) != file.Md5(dataFile) {
		t.Fail()
	}*/
}
