package data

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// Md5 returns the md5 sum of data
func Md5(data []byte) string {
	md5 := md5.New()
	md5.Write(data)
	return hex.EncodeToString(md5.Sum(nil))
}

// Sha1 returns the sha1 sum of data
func Sha1(data []byte) string {
	sha1 := sha1.New()
	sha1.Write(data)
	return hex.EncodeToString(sha1.Sum(nil))
}

// Sha256 returns the sha256 sum of data
func Sha256(data []byte) string {
	sha256 := sha256.New()
	sha256.Write(data)
	return hex.EncodeToString(sha256.Sum(nil))
}

// Sha512 returns the sha512 sum of data
func Sha512(data []byte) string {
	sha512 := sha512.New()
	sha512.Write(data)
	return hex.EncodeToString(sha512.Sum(nil))
}
