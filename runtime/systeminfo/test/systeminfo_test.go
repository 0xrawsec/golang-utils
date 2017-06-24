package main

import (
	"testing"

	"github.com/0xrawsec/golang-utils/log"
	"github.com/0xrawsec/golang-utils/runtime/systeminfo"
)

func init() {
	log.InitLogger(log.LDebug)
}

func TestSystemInfo(t *testing.T) {
	sig := systeminfo.New()
	si, err := sig.Get()
	if err != nil {
		t.Error(err)
	}
	t.Log(si)
}
