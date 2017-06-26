package args

import (
	"fmt"
	"strconv"
	"strings"
)

// ListVar structure to deal with the flag module. Parse the argument to the
// flag as comma separated string
type ListVar []string

// String interface implementation
func (la *ListVar) String() string {
	return fmt.Sprintf("%s", *la)
}

// Set interface implementation
func (la *ListVar) Set(input string) error {
	*la = strings.Split(input, ",")
	return nil
}

type ListIntVar []int

// String interface implementation
func (lia *ListIntVar) String() string {
	return fmt.Sprintf("%v", *lia)
}

// Set interface implementation
func (lia *ListIntVar) Set(input string) error {
	lsa := strings.Split(input, ",")
	*lia = make([]int, len(lsa))
	for i, s := range lsa {
		iv, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		(*lia)[i] = iv
	}
	return nil
}
