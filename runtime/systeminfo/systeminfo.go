package systeminfo

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/0xrawsec/golang-utils/log"
	"github.com/0xrawsec/golang-utils/readers"
)

type SystemInfo struct {
	SysLocale string
	OSName    string
	OSVersion string
}

func (si SystemInfo) String() string {
	return fmt.Sprintf("locale: %s; osname: %s; osversion: %s", si.SysLocale, si.OSName, si.OSVersion)
}

type SystemInfoGetter interface {
	Get() (SystemInfo, error)
}

// Utility

func trimString(str string) string {
	return strings.Trim(str, `\t\r\n `)
}

// EmptySystemInfoGetter
type EmptySystemInfoGetter struct{}

func (esig EmptySystemInfoGetter) Get() (si SystemInfo, err error) {
	si.SysLocale = "unk"
	si.OSName = "unk"
	si.OSVersion = "unk"
	return
}

// WindowsSystemInfoGetter
var (
	winOsNameRegexp    = regexp.MustCompile(`^OS Name:\s+(?P<osname>.*)`)
	winOsVersionRegexp = regexp.MustCompile(`^OS Version:\s+(?P<osversion>.*)`)
	winSysLocale       = regexp.MustCompile(`^System Locale:\s+(?P<syslocale>.*)`)
)

type WindowsSystemInfoGetter struct{}

func (wsig WindowsSystemInfoGetter) Get() (si SystemInfo, err error) {
	cmd := exec.Command("cmd", "/c", "systeminfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	for line := range readers.Readlines(bytes.NewReader(output)) {
		//log.Debug(string(line))
		switch {
		case winOsNameRegexp.Match(line):
			log.Debug(string(line))
			si.OSVersion = trimString(string(winOsNameRegexp.FindSubmatch(line)[1]))
		case winOsVersionRegexp.Match(line):
			si.OSName = trimString(string(winOsVersionRegexp.FindSubmatch(line)[1]))
		case winSysLocale.Match(line):
			si.SysLocale = trimString(string(winSysLocale.FindSubmatch(line)[1]))
		}
	}
	return
}

func New() SystemInfoGetter {
	switch runtime.GOOS {
	case "windows":
		return WindowsSystemInfoGetter{}
	}
	return EmptySystemInfoGetter{}
}
