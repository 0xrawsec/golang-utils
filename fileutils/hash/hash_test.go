package hash

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/0xrawsec/golang-utils/fsutil/fswalker"
	"github.com/0xrawsec/golang-utils/log"
)

func TestHash(t *testing.T) {
	log.InitLogger(log.LDebug)
	wg := sync.WaitGroup{}
	for wi := range fswalker.Walk("../../") {
		if wi.Err != nil {
			//panic(wi.Err)
		}
		for _, fileInfo := range wi.Files {
			path := filepath.Join(wi.Dirpath, fileInfo.Name())
			wg.Add(1)
			go func() {
				h := New()
				log.Info(path)
				h.Update(path)
				log.Debug(h)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
