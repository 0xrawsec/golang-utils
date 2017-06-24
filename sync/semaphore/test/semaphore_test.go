package main

import (
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"

	"github.com/0xrawsec/golang-utils/fsutil/fswalker"
	"github.com/0xrawsec/golang-utils/sync/semaphore"
)

func TestSemaphore(t *testing.T) {
	var wg sync.WaitGroup
	sem := semaphore.New(1024)
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
