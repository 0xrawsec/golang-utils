package shred

import (
	"crypto/rand"
	"os"
)

// Shred a file
func Shred(fpath string) error {
	stat, err := os.Stat(fpath)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fpath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	b := make([]byte, stat.Size())
	rand.Read(b)
	_, err = file.Write(b)

	err = os.Remove(fpath)
	if err != nil {
		return err
	}

	return err
}
