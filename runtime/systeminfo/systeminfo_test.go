package systeminfo

import (
	"testing"

	"github.com/0xrawsec/golang-utils/log"
)

func init() {
	log.InitLogger(log.LDebug)
}

func TestSystemInfo(t *testing.T) {
	sig := New()
	si, err := sig.Get()
	if err != nil {
		t.Error(err)
	}
	t.Log(si)
}
