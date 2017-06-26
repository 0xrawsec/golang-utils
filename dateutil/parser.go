package dateutil

import (
	"fmt"
	"regexp"
	"time"
)

var (
	layouts = [...]string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
		time.RFC850}
)

const (
	// ANSICRe regexp for ANSIC format
	// "Mon Jan _2 15:04:05 2006"
	ANSICRe = "[A-Z][a-z]{2} [A-Z][a-z]{2} [0-1]{0,1}[0-9] [0-2][0-9]:[0-6][0-9]:[0-6][0-9] [0-9]{4}"
	// UnixDateRe regexp for UnixDate format
	// "Mon Jan _2 15:04:05 MST 2006"
	UnixDateRe = "[A-Z][a-z]{2} [A-Z][a-z]{2} [0-1]{0,1}[0-9] [0-2][0-9]:[0-6][0-9]:[0-6][0-9] [A-Z]+ [0-9]{4}"
	// RubyDateRe regexp for RubyDate format
	// "Mon Jan 02 15:04:05 -0700 2006"
	RubyDateRe = "[A-Z][a-z]{2} [A-Z][a-z]{2} [0-1][0-9] [0-2][0-9]:[0-6][0-9]:[0-6][0-9] [+-][0-9]{4} [0-9]{4}"
	//RFC822Re RFC822 regexp
	RFC822Re = "[0-1][0-9] [A-Z][a-z]{2} [0-9]{2} [0-2][0-9]:[0-6][0-9] [A-Z]+"
	//RFC822ZRe RFC822 regexp
	RFC822ZRe = "[0-1][0-9] [A-Z][a-z]{2} [0-9]{2} [0-2][0-9]:[0-6][0-9] [+-][0-9]{4}"
	// RFC850Re RFC850 regexp
	// "Monday, 02-Jan-06 15:04:05 MST"
	RFC850Re = "[A-Z][a-z]*, [0-1][0-9]-[A-Z][a-z]*-[0-9]{2} [0-2][0-9]:[0-6][0-9]:[0-6][0-9] [A-Z]+"
	// RFC1123Re RFC1123 regexp
	RFC1123Re = "[A-Z][a-z]{2}, [0-1][0-9] [A-Z][a-z]* [0-9]{4} [0-2][0-9]:[0-6][0-9]:[0-6][0-9] [A-Z]+"
	// RFC1123ZRe RFC1123Z regexp
	RFC1123ZRe = "[A-Z][a-z]{2}, [0-1][0-9] [A-Z][a-z]* [0-9]{4} [0-2][0-9]:[0-6][0-9]:[0-6][0-9] [+-][0-9]{4}"
	// RFC3339Re RFC3339 regexp
	RFC3339Re = "[0-9]{4}-[0-1][0-9]-[0-3][0-9]T[0-2][0-9]:[0-6][0-9]:[0-6][0-9](Z|[+-][0-1][0-9]:[0-6][0-9])"
	// RFC3339NanoRe RFC3339Nano regexp
	RFC3339NanoRe = `[0-9]{4}-[0-1][0-9]-[0-3][0-9]T[0-2][0-9]:[0-6][0-9]:[0-6][0-9]\.[0-9]{9}(Z|[+-][0-1][0-9]:[0-6][0-9])`
)

var (
	// ANSIC DateString
	ANSIC = NewDateString(ANSICRe, time.ANSIC)
	// UnixDate DateString
	UnixDate = NewDateString(UnixDateRe, time.UnixDate)
	// RubyDate DateString
	RubyDate = NewDateString(RubyDateRe, time.RubyDate)
	// RFC822 DateString
	RFC822 = NewDateString(RFC822Re, time.RFC822)
	// RFC822Z DateString
	RFC822Z = NewDateString(RFC822ZRe, time.RFC822Z)
	// RFC850 DateString
	RFC850 = NewDateString(RFC850Re, time.RFC850)
	// RFC1123 DateString
	RFC1123 = NewDateString(RFC1123Re, time.RFC1123)
	// RFC1123Z DateString
	RFC1123Z = NewDateString(RFC1123ZRe, time.RFC1123Z)
	// RFC3339 DateString
	RFC3339 = NewDateString(RFC3339Re, time.RFC3339)
	// RFC3339Nano DateString
	RFC3339Nano = NewDateString(RFC3339NanoRe, time.RFC3339Nano)

	allDateStrings = []*DateString{
		&ANSIC,
		&UnixDate,
		&RubyDate,
		&RFC822,
		&RFC822Z,
		&RFC850,
		&RFC1123,
		&RFC1123Z,
		&RFC3339,
		&RFC3339Nano}
)

// DateString structure
type DateString struct {
	Regexp *regexp.Regexp
	Layout string
}

// UnknownDateFormatError error
type UnknownDateFormatError struct {
	DateStr string
}

// Error error implementation
func (u *UnknownDateFormatError) Error() string {
	return fmt.Sprintf("Unknown date format: %s", u.DateStr)
}

// AddDateString Adds a NewDateString to the list of default DateStrings
func AddDateString(ds DateString) {
	allDateStrings = append(allDateStrings, &ds)
}

// Parse attempts to parse a time string with all the knowns DateStrings
func Parse(value string) (time.Time, error) {
	for _, ds := range allDateStrings {
		if ds.Match(value) {
			return ds.Parse(value)
		}
	}
	return time.Time{}, &UnknownDateFormatError{value}
}

// NewDateString creates a DateString structure
func NewDateString(dateRe, layout string) DateString {
	return DateString{regexp.MustCompile(dateRe), layout}
}

// Match returns true if the DateString Regexp matches b
func (d *DateString) Match(value string) bool {
	return d.Regexp.Match([]byte(value))
}

// Parse parses value and returns the corresponding time.Time
func (d *DateString) Parse(value string) (time.Time, error) {
	return time.Parse(d.Layout, value)
}
