package file

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
)

// Md5 return the md5 sum of a file
func Md5(path string) (string, error) {
	var buffer [4096]byte
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	md5 := md5.New()

	for read, err := file.Read(buffer[:]); err != io.EOF && read != 0; read, err = file.Read(buffer[:]) {
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return "", err
		}
		md5.Write(buffer[:read])
	}

	return hex.EncodeToString(md5.Sum(nil)), nil
}

// Sha1 return the sha1 sum of a file
func Sha1(path string) (string, error) {
	var buffer [4096]byte
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	sha1 := sha1.New()

	for read, err := file.Read(buffer[:]); err != io.EOF && read != 0; read, err = file.Read(buffer[:]) {
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return "", err
		}
		sha1.Write(buffer[:read])
	}

	return hex.EncodeToString(sha1.Sum(nil)), nil
}

// Sha256 return the sha256 sum of a file
func Sha256(path string) (string, error) {
	var buffer [4096]byte
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	sha256 := sha256.New()

	for read, err := file.Read(buffer[:]); err != io.EOF && read != 0; read, err = file.Read(buffer[:]) {
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return "", err
		}
		sha256.Write(buffer[:read])
	}

	return hex.EncodeToString(sha256.Sum(nil)), nil
}

// Sha512 return the sha512 sum of a file
func Sha512(path string) (string, error) {
	var buffer [4096]byte
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	sha512 := sha512.New()

	for read, err := file.Read(buffer[:]); err != io.EOF && read != 0; read, err = file.Read(buffer[:]) {
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return "", err
		}
		sha512.Write(buffer[:read])
	}

	return hex.EncodeToString(sha512.Sum(nil)), nil
}
