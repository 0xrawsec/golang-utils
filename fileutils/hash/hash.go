package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func openFile(path string, pan bool) *os.File {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		switch pan {
		case true:
			panic(err)
		default:
			fmt.Println(err)
		}
	}
	return file
}

// Hashes structure definition
type Hashes struct {
	Path   string `json:"path"`
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
	Sha512 string `json:"sha512"`
}

// New Hashes structure
func New() Hashes {
	return Hashes{}
}

// Update the current Hashes structure
func (h *Hashes) Update(path string) (err error) {
	var buffer [4096]byte
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	sha512 := sha512.New()

	for read, err := file.Read(buffer[:]); err != io.EOF && read != 0; read, err = file.Read(buffer[:]) {
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
		md5.Write(buffer[:read])
		sha1.Write(buffer[:read])
		sha256.Write(buffer[:read])
		sha512.Write(buffer[:read])
	}

	h.Path = path
	h.Md5 = hex.EncodeToString(md5.Sum(nil))
	h.Sha1 = hex.EncodeToString(sha1.Sum(nil))
	h.Sha256 = hex.EncodeToString(sha256.Sum(nil))
	h.Sha512 = hex.EncodeToString(sha512.Sum(nil))
	return nil
}
