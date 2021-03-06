package progress

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/0xrawsec/golang-utils/fsutil/fswalker"
)

func TestProgress(t *testing.T) {
	progress := New(80)
	progress.SetPre("This is a test")
	for wi := range fswalker.Walk("../../../") {
		if wi.Err != nil {
			//panic(wi.Err)
		}
		for _, fileInfo := range wi.Files {
			progress.Update(filepath.Join(wi.Dirpath, fileInfo.Name()))
			progress.Print()
		}
	}
	fmt.Print("\n")
}
