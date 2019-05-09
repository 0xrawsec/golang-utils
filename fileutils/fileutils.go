package fileutils

import (
	"os"
	"fmt"
	"io"
	"compress/gzip"
)

// GzipFile compresses a file to gzip and deletes the original file
func GzipFile(path string) (err error) {
	var buf [4096]byte
	f, err := os.Open(path)
	if err != nil {
		return
	}
	//defer f.Close()
	fname := fmt.Sprintf("%s.gz", path)
	partname := fmt.Sprintf("%s.part", fname)
	of, err := os.Create(partname)
	if err != nil {
		return
	}

	w := gzip.NewWriter(of)
	for n, err := f.Read(buf[:]); err != io.EOF; {
		w.Write(buf[:n])
		n, err = f.Read(buf[:])
	}
	w.Flush()
	// gzip writer
	w.Close()
	// original file
	f.Close()
	// part file
	of.Close()
	if err = os.Remove(path); err != nil {
		return err
	}
	// rename the file to its final name
	return os.Rename(partname, fname)
}