package main

import (
	"path/filepath"
	"testing"

	"github.com/0xrawsec/golang-utils/fsutil/fswalker"
)

var (
	loopDir = "./test_dir/"
)

func TestWalk(t *testing.T) {
	for wi := range fswalker.Walk(loopDir) {
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
