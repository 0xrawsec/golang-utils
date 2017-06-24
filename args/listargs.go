package args

import (
	"fmt"
	"strings"
)

// ListArgs structure to deal with the flag module. Parse the argument to the
// flag as comma separated string
type ListArgs []string

// String interface implementation
func (la *ListArgs) String() string {
	return fmt.Sprintf("%s", *la)
}

// Set interface implementation
func (la *ListArgs) Set(input string) error {
	*la = strings.Split(input, ",")
	return nil
}
