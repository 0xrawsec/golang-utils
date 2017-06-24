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

type SubmatchHelper struct {
	IndexMap   map[string]int
	timeLayout string
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

// NewSubmatchHelper : creates a submatch helper from a regexp struct
// @r: pointer to regexp struct
// return (SubmatchHelper)
func NewSubmatchHelper(r *regexp.Regexp) (sm SubmatchHelper) {
	sm.IndexMap = make(map[string]int)
	for i, name := range r.SubexpNames() {
		sm.IndexMap[name] = i
	}
	sm.timeLayout = time.RFC3339
	return
}

// SetTimeLayout : setter for timeLayout field of SubmatchHelper to properly parse
// timestamps
// @layout : layout to switch to
func (sh *SubmatchHelper) SetTimeLayout(layout string) {
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

// Unmarshal : unmarshal submatches resulting from regexp.FindSubmatch and fill
// structure accordingly
// @sm: pointer to result of regexp.FindSubmatch
// @v: pointer of struct to fill
// return (error)
func (sh *SubmatchHelper) Unmarshal(sm *[][]byte, v interface{}) error {
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
			if err := sh.Unmarshal(sm, s.Field(i).Addr().Interface()); err != nil {
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
			b, err := sh.GetBytes(key, sm)
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
				return FieldNotSetError{field.Name}
			default:
				return err
			}
		}
	}
	return nil
}

// GetString : Get the matching string (if any) from regexp.FindStringSubmatch,
// corresponding to a given key (named regexp)
// @key: name of the regex subexpression
// @sm: pointer to result of regexp.FindStringSubmatch
// return (string, error): the associated match and error
func (sh *SubmatchHelper) GetString(key string, sm *[]string) (string, error) {
	if i, ok := sh.IndexMap[key]; ok {
		if i <= len(*sm) {
			return (*sm)[i], nil
		}
		return "", ErrIndexOutOfRange
	}
	return "", ErrNoSuchKey
}

// GetString : Get the matching []byte (if any) from regexp.FindSubmatch,
// corresponding to a given key (named regexp)
// @key: name of the regex subexpression
// @sm: pointer to result of regexp.FindSubmatch
// return ([]byte, error): the associated match and error
func (sh *SubmatchHelper) GetBytes(key string, sm *[][]byte) ([]byte, error) {
	if i, ok := sh.IndexMap[key]; ok {
		if i <= len(*sm) {
			return (*sm)[i], nil
		}
		return []byte{}, ErrIndexOutOfRange
	}
	return []byte{}, ErrNoSuchKey
}
