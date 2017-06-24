package packer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/0xrawsec/golang-utils/code/builder"
)

type ErrResourceNotFound struct {
	Name string
}

func (e ErrResourceNotFound) Error() string {
	return fmt.Sprintf("Resource %s not found", e.Name)
}

type Packer map[string][]byte

func PackReader(reader io.Reader) []byte {
	buf := new(bytes.Buffer)
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	packer := gzip.NewWriter(buf)
	packer.Write(all)
	packer.Close()
	return buf.Bytes()
}

func UnpackReader(reader io.Reader) []byte {
	unpacker, err := gzip.NewReader(reader)
	if err != nil {
		panic(err)
	}
	all, err := ioutil.ReadAll(unpacker)
	if err != nil {
		panic(err)
	}
	return all
}

func (p *Packer) AddResource(name string, data []byte) {
	buf := bytes.NewBuffer(data)
	(*p)[name] = PackReader(buf)
}

func (p *Packer) AddResourceReader(name string, reader io.Reader) {
	(*p)[name] = PackReader(reader)

}

func (p *Packer) GetResource(name string) ([]byte, error) {
	if data, ok := (*p)[name]; ok {
		buf := bytes.NewBuffer(data)
		return UnpackReader(buf), nil
	}
	return []byte{}, ErrResourceNotFound{name}
}

func (p *Packer) Dumps(packageName string) []byte {
	b := builder.CodeBuilder{}
	b.Package(packageName)
	b.DefVariable("Resources", *p)
	b.ResolveImports()
	return b.Bytes()
}

func (p *Packer) Dump(packageName, outfile string) {
	err := os.Mkdir(packageName, 0700)
	if !os.IsExist(err) && err != nil {
		panic(err)
	}
	out, err := os.Create(filepath.Join(packageName, outfile))
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	defer out.Close()
	out.Write(p.Dumps(packageName))
}
