package check

import (
	"testing"
)

var (
	md5    = "68b329da9893e34099c7d8ad5cb9c940"
	sha1   = "adc83b19e793491b1c6ea0fd8b46cd9f32e592fc"
	sha256 = "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b"
	sha512 = "be688838ca8686e5c90689bf2ab585cef1137c999b48c70b92f67a5c34dc15697b5d11c982ed6d71be1e1e7f7b4e0733884aa97c3f7a339a8ed03577cf74be09"
)

func TestHashVerification(t *testing.T) {
	switch {
	case !IsValidHash(md5):
		t.Log("MD5 not valid")
		t.Fail()
	case !IsValidHash(sha1):
		t.Log("SHA1 not valid")
		t.Fail()
	case !IsValidHash(sha256):
		t.Log("SHA256 not valid")
		t.Fail()
	case !IsValidHash(sha512):
		t.Log("SHA512 not valid")
		t.Fail()
	}
}
