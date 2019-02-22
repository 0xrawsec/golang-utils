package submatch

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)

// Helper structure definition
type Helper struct {
	IndexMap   map[string]int
	timeLayout string
	regex      *regexp.Regexp
	submatch   [][]byte
}

var (
	ErrNoSuchKey                 = errors.New("No such key")
	ErrIndexOutOfRange           = errors.New("Index out of range")
	ErrUnparsableDestinationType = errors.New("Unknown destination type")
	ErrNotValidPtr               = errors.New("Not valid pointer")
	TimeType                     = reflect.ValueOf(time.Time{}).Type()
)

type FieldNotSetError struct {
	Field string
}

func (fnse FieldNotSetError) Error() string {
	return fmt.Sprintf("Cannot set field: %s", fnse.Field)
}

// NewHelper : creates a new submatch helper from a regexp struct
func NewHelper(r *regexp.Regexp) (sm Helper) {
	sm.IndexMap = make(map[string]int)
	for i, name := range r.SubexpNames() {
		sm.IndexMap[name] = i
	}
	sm.timeLayout = time.RFC3339
	sm.regex = r
	return
}

// Prepare : this method must be called on any []byte/string you
// want the helper to work on. It basically apply regex.Regexp.FindSubmatch
// on b and initializes internal helper member for further processing.
func (sh *Helper) Prepare(b []byte) {
	sh.submatch = sh.regex.FindSubmatch(b)
}

// SetTimeLayout : setter for timeLayout field of SubmatchHelper
// to properly parse timestamps
func (sh *Helper) SetTimeLayout(layout string) {
	sh.timeLayout = layout
}

func strParse(s *string, k reflect.Kind) (interface{}, error) {
	switch k {
	// String
	case reflect.String:
		// return a copy of the string
		return string(*s), nil
	// Uints
	case reflect.Uint8:
		conv, err := strconv.ParseUint(*s, 10, 8)
		return uint8(conv), err
	case reflect.Uint16:
		conv, err := strconv.ParseUint(*s, 10, 16)
		return uint16(conv), err
	case reflect.Uint32:
		conv, err := strconv.ParseUint(*s, 10, 32)
		return uint32(conv), err
	case reflect.Uint64:
		return strconv.ParseUint(*s, 10, 64)
	case reflect.Uint:
		conv, err := strconv.ParseUint(*s, 10, 8*int(unsafe.Sizeof(uint(0))))
		return uint(conv), err
	// Ints
	case reflect.Int8:
		conv, err := strconv.ParseInt(*s, 10, 8)
		return int8(conv), err
	case reflect.Int16:
		conv, err := strconv.ParseInt(*s, 10, 16)
		return int16(conv), err
	case reflect.Int32:
		conv, err := strconv.ParseInt(*s, 10, 32)
		return int32(conv), err
	case reflect.Int64:
		return strconv.ParseInt(*s, 10, 64)
	case reflect.Int:
		conv, err := strconv.ParseInt(*s, 10, 8*int(unsafe.Sizeof(int(0))))
		return int(conv), err
	// Floats
	case reflect.Float32:
		conv, err := strconv.ParseFloat(*s, 32)
		return float32(conv), err
	case reflect.Float64:
		return strconv.ParseFloat(*s, 64)
	// Bool
	case reflect.Bool:
		return strconv.ParseBool(*s)
	}
	return "", ErrUnparsableDestinationType
}

// Unmarshal : unmarshal the data found by the Helper's regexp into v.
// Helper needs to be prepared first through the Prepare function.
func (sh *Helper) Unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrNotValidPtr
	}
	s := rv.Elem()
	t := s.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Unmarshal recursively if field of struct is a struct
		if s.Field(i).Kind() == reflect.Struct && s.Field(i).Type().Name() != TimeType.Name() {
			if err := sh.Unmarshal(s.Field(i).Addr().Interface()); err != nil {
				return err
			}
		} else {
			// We get the value in the tag named regexp
			key := field.Tag.Get("regexp")
			if key == "" {
				// If tag does not exist, we take field name
				key = field.Name
			}
			// Get the matched value and update the interface if necessary
			b, err := sh.GetBytes(key)
			switch err {
			case nil:
				var cast interface{}
				str := string(b)
				switch {
				case s.Field(i).Type().Name() == TimeType.Name() && s.Field(i).Kind() == reflect.Struct:
					cast, err = time.Parse(sh.timeLayout, str)
				default:
					cast, err = strParse(&str, s.Field(i).Kind())
				}
				s.Field(i).Set(reflect.ValueOf(cast))
				if err != nil {
					return err
				}
			case ErrNoSuchKey:
				// It means we cannot set the field
				//return FieldNotSetError{field.Name}
			default:
				return err
			}
		}
	}
	return nil
}

// GetBytes : Get the matching []byte (if any) extracted from
// the data matched by the Helper's regexp. Helper needs to be
// prepared using the Prepare function to work properly.
func (sh *Helper) GetBytes(key string) ([]byte, error) {
	if i, ok := sh.IndexMap[key]; ok {
		if i <= len(sh.submatch) {
			return sh.submatch[i], nil
		}
		return []byte{}, ErrIndexOutOfRange
	}
	return []byte{}, ErrNoSuchKey
}
