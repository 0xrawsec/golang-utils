package logfile

import (
	"compress/gzip"
	"fmt"
	"fsutil"
	"os"
	"sync"
	"time"
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

// LogFile structure definition
// A LogFile is a GZIP compressed file which rotates automatically
type LogFile struct {
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
func OpenFile(path string, perm os.FileMode, size int64) (*LogFile, error) {
	l := LogFile{}
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

// Path returns the path of the LogFile
func (l *LogFile) Path() string {
	if l.idx == 0 {
		return l.path
	}
	return fmt.Sprintf("%s.%d", l.path, l.idx)
}

// helper function to retrieve the first available file for writing
func (l *LogFile) searchFirstAvPath() error {
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

// helper function which rotate the logfile when needed
func (l *LogFile) rotateRoutine() {
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

// SetRefreshRate sets the rate at which the LogFile should check for rotating
// check DefaultRotationRate for default value
func (l *LogFile) SetRefreshRate(rate time.Duration) {
	l.rate = rate
}

// Rotate rotates the current LogFile
func (l *LogFile) Rotate() error {
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

// Close closes the LogFile properly
func (l *LogFile) Close() error {
	l.writer.Flush()
	l.writer.Close()
	return l.file.Close()
}

// WriteString writes a string into the LogFile
func (l *LogFile) WriteString(s string) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.writer.Write([]byte(s))
}

// Write writes bytes into the LogFile
func (l *LogFile) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.writer.Write(b)
}

// Flush method
func (l *LogFile) Flush() error {
	return l.writer.Flush()
}
