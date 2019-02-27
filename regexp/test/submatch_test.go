package main

import (
	"fmt"
	"regexp"
	"regexp/submatch"
	"testing"
	"time"
)

type TestStructure struct {
	M1 string    `regexp:"m1"`
	M2 int8      `regexp:"m2"`
	M3 int16     `regexp:"m3"`
	M4 time.Time `regexp:"m4"`
}

func TestGetByte(t *testing.T) {
	rex := regexp.MustCompile("(?P<test>.*)")
	line := "shouldmatcheverything"
	sh := submatch.NewHelper(rex)
	sh.Prepare([]byte(line))
	val, err := sh.GetBytes("test")
	t.Logf("Retrieved value: %s", val)
	if string(val) != line || err != nil {
		if err != nil {
			t.Errorf("Failed to retrieve field: %s", err)
		}
		t.Errorf("Retrieved field value not expected")
	}
}

func TestUnmarshal(t *testing.T) {
	rex := regexp.MustCompile("((?P<m1>.*?),(?P<m2>.*?),(?P<m3>.*?),(?P<m4>.*),)")
	line := fmt.Sprintf("thisisastring,4,42,%s,", time.Now().Format(time.RFC1123Z))
	sh := submatch.NewHelper(rex)
	ts := TestStructure{}
	sh.SetTimeLayout(time.RFC1123Z)
	sh.Prepare([]byte(line))
	err := sh.Unmarshal(&ts)
	if err != nil {
		t.Errorf("Failed to unmarshal: %s", err)
	}
	t.Logf("%v", ts)
}
