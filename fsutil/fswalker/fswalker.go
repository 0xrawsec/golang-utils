package fswalker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/0xrawsec/golang-utils/log"
)

const (
	chanBuffSize = 4096
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func perror(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// WalkItem returned by the Walk method
type WalkItem struct {
	Dirpath string
	Dirs    []os.FileInfo
	Files   []os.FileInfo
	Err     error
}

// NormalizePath normalizes a given path
func NormalizePath(path string) string {
	pointer, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	abs, err := filepath.Abs(pointer)
	if err != nil {
		return pointer
	}
	return abs
}

// Walk : walks recursively through the FS
func Walk(root string) <-chan WalkItem {
	// probably more efficient since wait only when chan is full
	iterChannel := make(chan WalkItem, chanBuffSize)
	dirsAlreadyProcessed := make(map[string]bool)
	dirsToProcess := []string{root}
	go func() {
		for len(dirsToProcess) > 0 {
			dirs, files := []os.FileInfo{}, []os.FileInfo{}
			dirpath := NormalizePath(dirsToProcess[len(dirsToProcess)-1])
			dirsToProcess = dirsToProcess[:len(dirsToProcess)-1]
			if _, ok := dirsAlreadyProcessed[dirpath]; !ok {
				dirsAlreadyProcessed[dirpath] = true
				filesInfo, err := ioutil.ReadDir(NormalizePath(dirpath))
				if err != nil {
					log.Errorf("Error reading directory (%s): %s\n", err.Error(), dirpath)
				} else {
					for _, fileInfo := range filesInfo {
						switch {
						case fileInfo.Mode().IsDir():
							dirs = append(dirs, fileInfo)
							dirsToProcess = append(dirsToProcess, filepath.Join(dirpath, fileInfo.Name()))
						case fileInfo.Mode().IsRegular():
							files = append(files, fileInfo)
						case fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink:
							var pointerFI os.FileInfo
							sympath := NormalizePath(filepath.Join(dirpath, fileInfo.Name()))
							pointerFI, err = os.Stat(sympath)
							if err != nil {
								log.Errorf("Error reading symlink (%s): %s\n", err.Error(), sympath)
							} else {
								switch {
								case pointerFI.Mode().IsDir():
									dirs = append(dirs, fileInfo)
									dirsToProcess = append(dirsToProcess, sympath)
								case pointerFI.Mode().IsRegular():
									files = append(files, fileInfo)
								}
							}
						}
					}
				}
				iterChannel <- WalkItem{dirpath, dirs, files, err}
			}
		}
		close(iterChannel)
	}()
	return iterChannel
}
