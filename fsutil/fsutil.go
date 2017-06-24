package fsutil

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	// ErrSrcNotRegularFile : src not a regular file
	ErrSrcNotRegularFile = errors.New("Source file is not a regular file")
	// ErrDstNotRegularFile : dst not a regular file
	ErrDstNotRegularFile = errors.New("Destination file is not a regular file")
)

// CopyFile : copies src file to dst file
func CopyFile(src, dst string) (err error) {
	srcStats, err := os.Stat(src)
	if err != nil {
		return
	}
	if !srcStats.Mode().IsRegular() {
		return ErrSrcNotRegularFile
	}
	dstStats, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		// The file already exists
		if !dstStats.Mode().IsRegular() {
			return ErrDstNotRegularFile
		}
	}
	return copyFileContents(src, dst)
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// AbsFromRelativeToBin : function that returns the absolute path to a file/directory
// computed with the directory of the binary as the root
// Example : if bin is /opt/program this function will return an absolute path computed as relative to /opt/
// @relPath : the parts of the path you want to Join in your final path
// return (string, error) : the absolute path and an error if necessary
func AbsFromRelativeToBin(relPath ...string) (string, error) {
	rootDirname := filepath.Dir(os.Args[0])
	absRootDirname, err := filepath.Abs(rootDirname)
	if err != nil {
		return "", err
	}
	return filepath.Join(absRootDirname, filepath.Join(relPath...)), nil
}
