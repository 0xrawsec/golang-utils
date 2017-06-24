package main

import (
	"testing"

	"github.com/0xrawsec/golang-utils/log"
)

func TestLog(t *testing.T) {
	log.InitLogger(log.LInfo)
	log.Info("I log into console")
	log.SetLogfile("test.log")
	log.Infof("I log into %s", "file")
}
