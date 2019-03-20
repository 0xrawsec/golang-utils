package logfile

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

var (
	path = filepath.Join("test", "output", "logfile.log")
)

func TestLogfile(t *testing.T) {
	size := int64(MB * 10)
	lf, err := OpenFile(path, 0600, size)
	if err != nil {
		t.Fail()
		t.Logf("Cannot create logfile: %s", err)
		return
	}
	lf.SetRefreshRate(time.Nanosecond * 5)
	defer lf.Close()
	buff := make([]byte, 10)
	lwritten := 0
	for i := int64(0); i < size/5; i++ {
		rand.Read(buff)
		lf.WriteString(fmt.Sprintf("%q\n", buff))
		lwritten++
	}
	t.Logf("Written %d lines", lwritten)
}
