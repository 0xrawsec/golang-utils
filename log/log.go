package log

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

const (
	// LDebug log level
	LDebug = 1
	// LInfo log level
	LInfo = 1 << 1
	// LError log level
	LError = 1 << 2
	// LCritical log level
	LCritical       = 1 << 3
	defaultFileMode = 0640
)

var (
	gLogLevel       = LInfo
	gLogLevelBackup = gLogLevel

	MockAbort = false
)

func init() {
	//gLogger.Set
	InitLogger(LInfo)
}

// InitLogger Initialize the global logger
func InitLogger(logLevel int) {
	SetLogLevel(logLevel)
	if logLevel <= LDebug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

// SetLogfile sets output file to put logging messages
func SetLogfile(logfilePath string, opts ...os.FileMode) {
	var err error
	mode := os.FileMode(defaultFileMode)
	// We open the file in append mode
	if len(opts) > 0 {
		mode = opts[0]
	}
	gLogFile, err := os.OpenFile(logfilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, mode)
	if err != nil {
		panic(err)
	}
	if _, err := gLogFile.Seek(0, os.SEEK_END); err != nil {
		panic(err)
	}
	log.SetOutput(gLogFile)
}

// SetLogLevel backup gLoglevel and set gLogLevel to logLevel
func SetLogLevel(logLevel int) {
	gLogLevelBackup = gLogLevel
	switch logLevel {
	case LInfo:
		gLogLevel = logLevel
	case LDebug:
		gLogLevel = logLevel
	case LCritical:
		gLogLevel = logLevel
	case LError:
		gLogLevel = logLevel
	default:
		gLogLevel = LInfo
	}
}

// RestoreLogLevel restore gLogLevel to gLogLevelBackup
func RestoreLogLevel() {
	gLogLevel = gLogLevelBackup
}

func logMessage(prefix string, i ...interface{}) {
	format := fmt.Sprintf("%s%s", prefix, strings.Repeat("%v ", len(i)))
	msg := fmt.Sprintf(format, i...)
	log.Output(3, msg)
}

// Info log message if gLogLevel <= LInfo
func Info(i ...interface{}) {
	if gLogLevel <= LInfo {
		logMessage("INFO - ", i...)
	}
}

// Infof log message with format if gLogLevel <= LInfo
func Infof(format string, i ...interface{}) {
	if gLogLevel <= LInfo {
		logMessage("INFO - ", fmt.Sprintf(format, i...))
	}
}

// Warning log message if gLogLevel <= LInfo
func Warn(i ...interface{}) {
	if gLogLevel <= LInfo {
		logMessage("WARNING - ", i...)
	}
}

// Warnf log message with format if gLogLevel <= LInfo
func Warnf(format string, i ...interface{}) {
	if gLogLevel <= LInfo {
		logMessage("WARNING - ", fmt.Sprintf(format, i...))
	}
}

// Debug log message if gLogLevel <= LDebug
func Debug(i ...interface{}) {
	if gLogLevel <= LDebug {
		logMessage("DEBUG - ", i...)
	}
}

// Debugf log message with format if gLogLevel <= LDebug
func Debugf(format string, i ...interface{}) {
	if gLogLevel <= LDebug {
		logMessage("DEBUG - ", fmt.Sprintf(format, i...))
	}
}

// Error log message if gLogLevel <= LError
func Error(i ...interface{}) {
	if gLogLevel <= LError {
		logMessage("ERROR - ", i...)
	}
}

// Errorf log message with format if gLogLevel <= LError
func Errorf(format string, i ...interface{}) {
	if gLogLevel <= LError {
		logMessage("ERROR - ", fmt.Sprintf(format, i...))
	}
}

// Abort logs an error and exit with return code
func Abort(rc int, i ...interface{}) {
	if gLogLevel <= LError {
		logMessage("ABORT - ", i...)
	}
	if !MockAbort {
		os.Exit(rc)
	}
}

// Critical log message if gLogLevel <= LCritical
func Critical(i ...interface{}) {
	if gLogLevel <= LCritical {
		logMessage("CRITICAL - ", i...)
	}
}

// Criticalf log message with format if gLogLevel <= LCritical
func Criticalf(format string, i ...interface{}) {
	if gLogLevel <= LCritical {
		logMessage("CRITICAL - ", fmt.Sprintf(format, i...))
	}
}

// DontPanic only prints panic information but don't panic
func DontPanic(i interface{}) {
	msg := fmt.Sprintf("%v\n %s", i, debug.Stack())
	logMessage("PANIC - ", msg)
}

// DebugDontPanic only prints panic information but don't panic
func DebugDontPanic(i interface{}) {
	if gLogLevel <= LDebug {
		msg := fmt.Sprintf("%v\n %s", i, debug.Stack())
		logMessage("PANIC - ", msg)
	}
}

// DontPanicf only prints panic information but don't panic
func DontPanicf(format string, i ...interface{}) {
	msg := fmt.Sprintf("%v\n %s", fmt.Sprintf(format, i...), debug.Stack())
	logMessage("PANIC - ", msg)
}

// DebugDontPanicf only prints panic information but don't panic
func DebugDontPanicf(format string, i ...interface{}) {
	if gLogLevel <= LDebug {
		msg := fmt.Sprintf("%v\n %s", fmt.Sprintf(format, i...), debug.Stack())
		logMessage("PANIC - ", msg)
	}
}

// Panic prints panic information and call panic
func Panic(i interface{}) {
	msg := fmt.Sprintf("%v\n %s", i, debug.Stack())
	logMessage("PANIC - ", msg)
	panic(i)
}
