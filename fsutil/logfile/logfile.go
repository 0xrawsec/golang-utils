package logfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0xrawsec/golang-utils/fileutils"
	"github.com/0xrawsec/golang-utils/fsutil"
	"github.com/0xrawsec/golang-utils/fsutil/fswalker"
	"github.com/0xrawsec/golang-utils/log"
)

const (
	_ = iota // ignore first value by assigning to blank identifier
	//KB Kilobytes
	KB = 1 << (10 * iota)
	//MB Megabytes
	MB
	//GB Gigabytes
	GB
	//TB Terabytes
	TB

	// DefaultRotationRate defines the default value for rotation refresh rate
	DefaultRotationRate = time.Millisecond * 250
)

// LogFile interface
type LogFile interface {
	// Path returns the path of the current LogFile
	Path() string
	// Rotate ensures file rotation
	Rotate() error
	// Rotation routine
	RotRoutine()
	// Write used to write to the LogFile
	Write([]byte) (int, error)
	// Close used to close the LogFile
	Close() error
}

// BaseLogFile structure definition
type BaseLogFile struct {
	sync.Mutex
	base string
	dir  string
	path string
	//prefix string
	perm   os.FileMode
	file   *os.File
	writer io.Writer
	done   chan bool
	wg     sync.WaitGroup
}

// Rotate implements LogFile interface
func (b *BaseLogFile) Rotate() (err error) {
	b.Lock()
	defer b.Unlock()
	b.file.Close()

	// First rename all the gzip files
	// First find max file index
	maxIdx := uint64(0)
	for wi := range fswalker.Walk(filepath.Dir(b.path)) {
		for _, fi := range wi.Files {
			if strings.HasPrefix(fi.Name(), b.base) && strings.HasSuffix(fi.Name(), ".gz") {
				ext := strings.TrimLeft(fi.Name(), fmt.Sprintf("%s.", b.base))
				sp := strings.SplitN(ext, ".", 2)
				if len(sp) == 2 {
					id, err := strconv.ParseUint(sp[0], 0, 64)
					if err != nil {
						log.Info(fi.Name())
						log.Errorf("Cannot parse logfile id: %s", err)
					}
					if id > maxIdx {
						maxIdx = id
					}
				}
			}
		}
	}

	// Actually renaming
	for i := maxIdx; i > 0; i-- {
		// renaming the zip file
		oldf := fmt.Sprintf("%s.%d.gz", b.path, i)
		newf := fmt.Sprintf("%s.%d.gz", b.path, i+1)
		// because there we do not guarantee that oldf exists due to previous loop
		if fsutil.IsFile(oldf) {
			if err := os.Rename(oldf, newf); err != nil {
				log.Errorf("Failed to rename old logfile: %s", err)
			}
		}
	}

	// Rename basename.1 to basename.2
	dot1 := fmt.Sprintf("%s.1", b.path)
	dot2 := fmt.Sprintf("%s.2", b.path)
	// path to part file to control that we are not already compressing
	dot2Part := fmt.Sprintf("%s.2.gz.part", b.path)

	// Should not happen but that's a precaution step not to overwrite dot2
	// without knowing it
	if fsutil.IsFile(dot2) && !fsutil.IsFile(dot2Part) {
		if err := fileutils.GzipFile(dot2); err != nil {
			log.Errorf("Failed to gzip LogFile: %s", err)
		}
	}

	if fsutil.IsFile(dot1) {
		if err := os.Rename(dot1, dot2); err != nil {
			log.Errorf("Failed to rename old file: %s", err)
		} else {
			// Start a routine to gzip dot2
			b.wg.Add(1)
			go func() {
				defer b.wg.Done()
				if fsutil.IsFile(dot2) && !fsutil.IsFile(dot2Part) {
					if err := fileutils.GzipFile(dot2); err != nil {
						log.Errorf("Failed to gzip LogFile: %s", err)
					}
				}
			}()
		}
	}

	// Move current to basename.1
	if err := os.Rename(b.path, dot1); err != nil {
		log.Errorf("Failed to rename old file: %s", err)
	}

	b.file, err = os.OpenFile(b.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, b.perm)
	b.writer = b.file
	//l.timer.Reset(l.rotationDelay)
	return err
}

// Write implements LogFile interface
func (b *BaseLogFile) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.writer.Write(p)
}

// WriteString implements LogFile interface
func (b *BaseLogFile) WriteString(s string) (int, error) {
	return b.Write([]byte(s))
}

// Path implements LogFile interface
func (b *BaseLogFile) Path() string {
	return b.path
}

// Close implements LogFile interface
func (b *BaseLogFile) Close() error {
	b.Lock()
	defer b.Unlock()
	b.done <- true
	b.wg.Wait()
	return nil
}

// TimeRotateLogFile structure definition.
// A TimeRotateLogFile rotates at whenever rotation delay expires.
// The current file being used is in plain-text. Whenever the rotation
// happens, the file is GZIP-ed to save space on disk. A delay can be
// specified in order to wait before the file is compressed.
type TimeRotateLogFile struct {
	BaseLogFile
	rotationDelay time.Duration
	timer         *time.Timer
}

// OpenTimeRotateLogFile opens a new TimeRotateLogFile drot controls
// the rotation delay and dgzip the time to wait before the latest file is GZIPed
func OpenTimeRotateLogFile(path string, perm os.FileMode, drot time.Duration) (l *TimeRotateLogFile, err error) {

	l = &TimeRotateLogFile{}
	// BaseLogfile fields
	l.base = filepath.Base(path)
	l.dir = filepath.Dir(path)
	l.path = path
	l.perm = perm
	l.wg = sync.WaitGroup{}
	l.done = make(chan bool)
	// TimeRotateLogFile fields
	l.timer = time.NewTimer(drot)
	l.rotationDelay = drot

	l.file, err = os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_RDWR, l.perm)

	if err != nil {
		return
	}

	l.writer = l.file

	// Go routine responsible for log rotation
	l.wg.Add(1)
	go l.RotRoutine()

	return
}

// RotRoutine implements LogFile
func (l *TimeRotateLogFile) RotRoutine() {
	defer l.wg.Done()
	for {
		select {
		case <-l.done:
			l.file.Close()
			return
		case <-l.timer.C:
			if err := l.Rotate(); err != nil {
				log.Errorf("Failed LogFile rotation: %s", err)
			}
			l.timer.Reset(l.rotationDelay)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// Close implements LogFile interface
func (l *TimeRotateLogFile) Close() error {
	l.Lock()
	defer l.Unlock()
	l.done <- true
	// timer needs to be stopped not to try to Rotate while
	// some member have been uninitialized
	l.timer.Stop()
	l.wg.Wait()

	return nil
}

// SizeRotateLogFile structure definition
// A SizeRotateLogFile is a GZIP compressed file which rotates automatically
type SizeRotateLogFile struct {
	BaseLogFile
	size int64
}

// OpenSizeRotateLogFile opens a new log file for logging rotating
// according to its own size
func OpenSizeRotateLogFile(path string, perm os.FileMode, size int64) (*SizeRotateLogFile, error) {
	l := SizeRotateLogFile{}
	l.base = filepath.Base(path)
	l.dir = filepath.Dir(path)
	l.path = path
	l.perm = perm
	l.wg = sync.WaitGroup{}
	l.done = make(chan bool)
	// fields specific to SizeRotateLogFile
	l.size = size

	// Open the file descriptor
	f, err := os.OpenFile(l.Path(), os.O_APPEND|os.O_CREATE|os.O_RDWR, l.perm)
	if err != nil {
		return nil, err
	}

	l.file = f
	l.writer = l.file
	// We start the rotate routine

	l.wg.Add(1)
	go l.RotRoutine()
	return &l, nil
}

// RotRoutine implements LogFile
func (l *SizeRotateLogFile) RotRoutine() {
	defer l.wg.Done()
	for {
		select {
		case <-l.done:
			l.file.Close()
			return
		default:
			if stats, err := os.Stat(l.path); err == nil {
				if stats.Size() >= l.size {
					l.Rotate()
				}
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}
