package logfile

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	dir  = filepath.Join("test", "output")
	path = filepath.Join(dir, "logfile.log")
)

func init() {
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0777)
}

func TestLogfile(t *testing.T) {
	var lf LogFile
	size := int64(MB * 10)
	lf, err := OpenSizeRotateLogFile(path, 0600, size)

	if err != nil {
		t.Fail()
		t.Logf("Cannot create logfile: %s", err)
		return
	}
	//lf.(*SizeRotateLogFile).SetRefreshRate(time.Nanosecond * 5)
	defer lf.Close()
	buff := make([]byte, 10)
	lwritten := 0
	for i := int64(0); i < size/5; i++ {
		rand.Read(buff)
		lf.(*SizeRotateLogFile).WriteString(fmt.Sprintf("%q\n", buff))
		lwritten++
	}
	t.Logf("Written %d lines", lwritten)
}

func TestTimeRotateLFBasic(t *testing.T) {
	var lf LogFile
	lf, err := OpenTimeRotateLogFile(path, 0600, 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to create logfile")
		t.FailNow()
	}
	//defer lf.Close()
	buff := make([]byte, 10)
	lwritten := 0
	for i := int64(0); i < 1000; i++ {
		if i%500 == 0 {
			time.Sleep(1 * time.Second)
		}
		rand.Read(buff)
		if _, err := lf.(*TimeRotateLogFile).Write([]byte(fmt.Sprintf("%q\n", buff))); err != nil {
			t.Logf("Error writting: %s", err)
		}
		lwritten++
	}
	t.Logf("Written %d lines", lwritten)
	time.Sleep(3 * time.Second)
	lf.Close()
}
