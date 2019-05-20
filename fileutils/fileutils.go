package fileutils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// GzipFile compresses a file to gzip and deletes the original file
func GzipFile(path string) (err error) {
	var buf [4096]byte
	f, err := os.Open(path)
	if err != nil {
		return
	}

	fname := fmt.Sprintf("%s.gz", path)
	partname := fmt.Sprintf("%s.part", fname)

	// to keep permission of compressed file
	stat, err := os.Stat(path)
	if err != nil {
		return
	}

	of, err := os.OpenFile(partname, os.O_CREATE|os.O_WRONLY, stat.Mode())
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
