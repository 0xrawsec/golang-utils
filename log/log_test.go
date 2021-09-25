package log

import (
	"errors"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	InitLogger(LDebug)
	Debug("We are debugging", "this part", "of", "code")
	Info("I log into console", "this", "and", "that")
	Infof("%s %s", "we print", "a formated string")
	Warn("This", "is", "dangerous")
	Warnf("%s %s %s", "This", "is", "dangerous formated string")
	Error(fmt.Errorf("error encountered in program"), "but also this strange number:", 42)
	Errorf("%s %d", "encountered error", 666)
	Critical("Dammit", "we are in a bad", errors.New("situation"))
	Criticalf("%s %s %s", "Dammit", "we are in a bad", errors.New("situation"))
	DontPanic(errors.New("no stress"))
	DebugDontPanic(errors.New("no stress"))
	DontPanicf("%s %s", "manual say", "we should not panic")
	DebugDontPanicf("%s %s", "manual say", "we should not panic")
	MockAbort = true
	Abort(0, "Aborting because of", fmt.Errorf("error raised by some function"))
}
