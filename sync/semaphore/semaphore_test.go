package semaphore

import (
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"

	"github.com/0xrawsec/golang-utils/fsutil/fswalker"
)

func TestSemaphore(t *testing.T) {
	var wg sync.WaitGroup
	sem := New(1024)
	for wi := range fswalker.Walk("../../..") {
		if wi.Err != nil {
			panic(wi.Err)
		}
		for _, fileInfo := range wi.Files {
			fileInfo := fileInfo
			wg.Add(1)
			go func() {
				sem.Acquire()
				defer sem.Release()
				t.Logf("Reading: %s", filepath.Join(wi.Dirpath, fileInfo.Name()))
				ioutil.ReadFile(filepath.Join(wi.Dirpath, fileInfo.Name()))
				defer wg.Done()
			}()
		}
	}
	wg.Wait()
}
