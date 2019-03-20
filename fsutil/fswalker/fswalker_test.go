package fswalker

import (
	"path/filepath"
	"testing"
)

var (
	loopDir = "./test/test_dir/"
)

func TestWalk(t *testing.T) {
	for wi := range Walk(loopDir) {
		t.Log("Directories")
		for _, di := range wi.Dirs {
			t.Logf("\t%s", filepath.Join(wi.Dirpath, di.Name()))
		}

		t.Log("Files")
		for _, fi := range wi.Files {
			t.Logf("\t%s", filepath.Join(wi.Dirpath, fi.Name()))
		}
	}
}
