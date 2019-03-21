package args

import (
	"time"

	"github.com/0xrawsec/golang-utils/dateutil"
)

// DateVar struct
type DateVar time.Time

// String argument implementation
func (da *DateVar) String() string {
	return time.Time(*da).String()
}

// Set argument implementation
func (da *DateVar) Set(input string) error {
	t, err := dateutil.Parse(input)
	(*da) = DateVar(t)
	return err
}

// DurationVar structure
type DurationVar time.Duration

// String argument implementation
func (da *DurationVar) String() string {
	return time.Duration(*da).String()
}

// Set argument implementation
func (da *DurationVar) Set(input string) error {
	tda, err := time.ParseDuration(input)
	if err == nil {
		*da = DurationVar(tda)
	}
	return err
}
