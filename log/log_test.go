package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	InitLogger(LInfo)
	Info("I log into console")
	SetLogfile("test/test.log")
	Infof("I log into %s", "file")
}
