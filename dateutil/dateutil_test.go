package dateutil

import (
	"testing"
)

func TestDateStringANSIC(t *testing.T) {
	shouldMatch := [...]string{
		"Mon Jan 2 15:04:05 2006",
		"Mon Jan 12 15:04:05 2006"}
	ds := ANSIC
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestDateStringUnix(t *testing.T) {
	shouldMatch := [...]string{
		"Mon Jan 2 15:04:05 MST 2006",
		"Mon Jan 12 15:04:05 MST 2006"}
	ds := UnixDate
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestDateRuby(t *testing.T) {
	shouldMatch := [...]string{
		"Mon Jan 02 15:04:05 -0700 2006",
		"Mon Jan 12 15:04:05 +0000 2006"}
	ds := RubyDate
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestRFC822(t *testing.T) {
	shouldMatch := [...]string{
		"02 Jan 06 15:04 MST",
		"12 Jan 06 15:04 MST"}
	ds := RFC822
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestRFC822Z(t *testing.T) {
	shouldMatch := [...]string{
		"02 Jan 06 15:04 -0700",
		"12 Jan 06 15:04 +0700"}
	ds := RFC822Z
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestRFC850(t *testing.T) {
	shouldMatch := [...]string{
		"Monday, 02-Jan-06 15:04:05 MST",
		"Friday, 02-Jan-99 15:04:05 MST"}
	ds := RFC850
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestRFC1123(t *testing.T) {
	shouldMatch := [...]string{
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Fri, 08 Jan 2006 15:04:05 CET"}
	ds := RFC1123
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestRFC1123Z(t *testing.T) {
	shouldMatch := [...]string{
		"Mon, 02 Jan 2006 15:04:05 +1200",
		"Fri, 08 Jan 2006 15:04:05 -0700"}
	ds := RFC1123Z
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestRFC3339(t *testing.T) {
	shouldMatch := [...]string{
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:00Z",
		"2006-01-02T15:04:05+07:00"}
	ds := RFC3339
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestRFC3339Nano(t *testing.T) {
	shouldMatch := [...]string{
		"2006-01-02T15:04:05.999999999Z",
		"2006-01-02T15:04:05.999999999-07:00",
		"2006-01-02T15:04:05.999999999+07:00"}
	ds := RFC3339Nano
	for _, s := range shouldMatch {
		if ok := ds.Match(s); !ok {
			t.Logf("Does not match: %s", s)
			t.Fail()
		}
		if _, err := ds.Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		}
		t.Logf("Valid : %s", s)
	}
}

func TestParse(t *testing.T) {
	shouldMatch := [...]string{
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05+07:00",
		"Mon, 02 Jan 2006 15:04:05 +1200",
		"Fri, 08 Jan 2006 15:04:05 -0700",
		"Mon Jan 2 15:04:05 2006",
		"Mon Jan 12 15:04:05 2006",
		"2006-01-02T15:04:05.999999999-07:00",
		"2006-01-02T15:04:05.999999999+07:00"}
	for _, s := range shouldMatch {
		if _, err := Parse(s); err != nil {
			t.Logf("Cannot parse %s: %s", s, err)
			t.Fail()
		} else {
			t.Logf("Valid : %s", s)
		}
	}
}
