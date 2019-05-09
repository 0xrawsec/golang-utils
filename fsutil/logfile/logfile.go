package logfile

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/0xrawsec/golang-utils/fileutils"
	"github.com/0xrawsec/golang-utils/fsutil"
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
	// Write used to write to the LogFile
	Write([]byte) (int, error)
	// Close used to close the LogFile
	Close() error
}

// BaseLogFile structure definition
type BaseLogFile struct {
	sync.Mutex
	path   string
	prefix string
	perm   os.FileMode
	file   *os.File
	writer io.Writer
	done   chan bool
}

// Write implements LogFile interface
func (b *BaseLogFile) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.writer.Write(p)
}

// Path implements LogFile interface
func (b *BaseLogFile) Path() string {
	return fmt.Sprintf("%s.%s", b.path, b.prefix)
}

// Close implements LogFile interface
func (b *BaseLogFile) Close() error {
	b.Lock()
	defer b.Unlock()
	b.done <- true
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
	gzipDelay     time.Duration
	wg            sync.WaitGroup
}

// OpenTimeRotateLogFile opens a new TimeRotateLogFile drot controls
// the rotation delay and dgzip the time to wait before the latest file is GZIPed
func OpenTimeRotateLogFile(path string, perm os.FileMode, drot time.Duration, dgzip time.Duration) (l *TimeRotateLogFile, err error) {

	l = &TimeRotateLogFile{
		rotationDelay: drot,
		timer:         time.NewTimer(drot),
		gzipDelay:     dgzip,
		wg:            sync.WaitGroup{},
	}

	l.done = make(chan bool)
	l.path = path
	l.prefix = fmt.Sprintf("%d", time.Now().Unix())
	l.perm = perm

	l.file, err = os.OpenFile(l.Path(), os.O_CREATE|os.O_RDWR, l.perm)

	if err != nil {
		return
	}

	l.writer = l.file

	// Go routine responsible of log rotation
	l.wg.Add(1)
	go func() {
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
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	return
}

// Rotate implements LogFile interface
func (l *TimeRotateLogFile) Rotate() (err error) {
	l.Lock()
	defer l.Unlock()
	l.file.Close()
	// Go Routine to gzip out the previous file
	p := l.Path()
	l.wg.Add(1)
	go func() {
		time.Sleep(l.gzipDelay)
		if err := fileutils.GzipFile(p); err != nil {
			log.Errorf("Failed to gzip LogFile: %s", err)
		}
		l.wg.Done()
	}()

	l.prefix = fmt.Sprintf("%d", time.Now().Unix())
	l.file, err = os.OpenFile(l.Path(), os.O_CREATE|os.O_RDWR, l.perm)
	l.writer = l.file
	l.timer.Reset(l.rotationDelay)
	return err
}

// Close implements LogFile interface. Whenever the TimeRotateLogFile is closed
// the last file in use is GZIP compressed before Close returns
func (l *TimeRotateLogFile) Close() error {
	l.Lock()
	defer l.Unlock()
	l.done <- true
	// timer needs to be stopped not to try to Rotate while
	// some member have been uninitialized
	l.timer.Stop()
	l.wg.Wait()

	if err := fileutils.GzipFile(l.Path()); err != nil {
		log.Errorf("Failed to gzip LogFile: %s", err)
	}

	return nil
}

// SizeRotateLogFile structure definition
// A SizeRotateLogFile is a GZIP compressed file which rotates automatically
type SizeRotateLogFile struct {
	sync.Mutex
	idx    int
	path   string
	perm   os.FileMode
	file   *os.File
	writer *gzip.Writer
	size   int64
	rate   time.Duration
}

// OpenFile opens a new file for logging
func OpenFile(path string, perm os.FileMode, size int64) (*SizeRotateLogFile, error) {
	l := SizeRotateLogFile{}
	l.path = path
	l.perm = perm
	l.size = size
	l.rate = DefaultRotationRate

	// Search for the first available path
	err := l.searchFirstAvPath()
	if err != nil {
		return nil, err
	}
	// Open the file descriptor
	f, err := os.OpenFile(l.Path(), os.O_APPEND|os.O_CREATE|os.O_RDWR, l.perm)
	if err != nil {
		return nil, err
	}
	l.file = f
	l.writer = gzip.NewWriter(f)
	// We start the rotate routine
	l.rotateRoutine()
	return &l, nil
}

// Path returns the path of the LogFileSizeRotate
func (l *SizeRotateLogFile) Path() string {
	if l.idx == 0 {
		return l.path
	}
	return fmt.Sprintf("%s.%d", l.path, l.idx)
}

// helper function to retrieve the first available file for writing
func (l *SizeRotateLogFile) searchFirstAvPath() error {
	for {
		if fsutil.Exists(l.Path()) {
			stats, err := os.Stat(l.Path())
			if err != nil {
				return err
			}
			if stats.Size() >= l.size {
				l.idx++
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
}

// helper function which rotate the LogFileSizeRotate when needed
func (l *SizeRotateLogFile) rotateRoutine() {
	go func() {
		for {
			if stats, err := os.Stat(l.Path()); err == nil {
				if stats.Size() >= l.size {
					l.Rotate()
				}
			}
			time.Sleep(l.rate)
		}
	}()
}

// SetRefreshRate sets the rate at which the LogFileSizeRotate should check for rotating
// check DefaultRotationRate for default value
func (l *SizeRotateLogFile) SetRefreshRate(rate time.Duration) {
	l.rate = rate
}

// Rotate rotates the current LogFileSizeRotate
func (l *SizeRotateLogFile) Rotate() error {
	l.Lock()
	defer l.Unlock()
	// We close everything first
	l.Close()
	if err := l.searchFirstAvPath(); err != nil {
		return err
	}
	f, err := os.OpenFile(l.Path(), os.O_APPEND|os.O_CREATE|os.O_RDWR, l.perm)
	if err != nil {
		return err
	}
	l.file = f
	l.writer = gzip.NewWriter(l.file)
	return nil
}

// Close closes the LogFileSizeRotate properly
func (l *SizeRotateLogFile) Close() error {
	l.writer.Flush()
	l.writer.Close()
	return l.file.Close()
}

// WriteString writes a string into the LogFileSizeRotate
func (l *SizeRotateLogFile) WriteString(s string) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.writer.Write([]byte(s))
}

// Write writes bytes into the LogFileSizeRotate
func (l *SizeRotateLogFile) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.writer.Write(b)
}

// Flush method
func (l *SizeRotateLogFile) Flush() error {
	l.Lock()
	defer l.Unlock()
	return l.writer.Flush()
}
