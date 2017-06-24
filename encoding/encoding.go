package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

// Endianness interface definition specifies the endianness used to decode
type Endianness binary.ByteOrder

var (
	// ErrSeeking error definition
	ErrSeeking = errors.New("Error seeking")
	// ErrMultipleOffsets error definition
	ErrMultipleOffsets = errors.New("Only one offset argument is allowed")
	// ErrInvalidNilPointer
	ErrInvalidNilPointer = errors.New("Nil pointer is invalid")
	// No Pointer interface
	ErrNoPointerInterface = errors.New("Interface expect to be a pointer")
)

// Unpack data type from reader object. An optional offset can be specified.
func Unpack(reader io.ReadSeeker, endianness Endianness, data interface{}, offsets ...int64) error {

	switch {
	// No offset to deal with
	case len(offsets) == 0:
		if err := binary.Read(reader, endianness, data); err != nil {
			return err
		}
		// An offset to deal with
	case len(offsets) == 1:
		if soughtOffset, err := reader.Seek(offsets[0], os.SEEK_SET); soughtOffset != offsets[0] || err != nil {
			switch {
			case err != nil:
				return err
			case soughtOffset != offsets[0]:
				return ErrSeeking
			default:
				if err := binary.Read(reader, endianness, data); err != nil {
					return err
				}
			}
		}
		// Error if more than one offset
	default:
		return ErrMultipleOffsets
	}
	return nil
}

func marshalArray(data interface{}, endianness Endianness) ([]byte, error) {
	var out []byte
	val := reflect.ValueOf(data)
	if val.IsNil() {
		return out, ErrInvalidNilPointer
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Array {
		return out, fmt.Errorf("Not an Array structure")
	}
	for k := 0; k < elem.Len(); k++ {
		buff, err := Marshal(elem.Index(k).Addr().Interface(), endianness)
		if err != nil {
			return out, err
		}
		out = append(out, buff...)
	}
	return out, nil
}

func marshalSlice(data interface{}, endianness Endianness) ([]byte, error) {
	var out []byte
	val := reflect.ValueOf(data)
	if val.IsNil() {
		return out, ErrInvalidNilPointer
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Slice {
		return out, fmt.Errorf("Not a Slice object")
	}
	s := elem
	// We first serialize slice length as a int64
	sliceLen := int64(s.Len())
	buff, err := Marshal(&sliceLen, endianness)
	if err != nil {
		return out, err
	}
	out = append(out, buff...)
	for k := 0; k < s.Len(); k++ {
		buff, err := Marshal(s.Index(k).Addr().Interface(), endianness)
		if err != nil {
			return out, err
		}
		out = append(out, buff...)
	}
	return out, nil
}

func Marshal(data interface{}, endianness Endianness) ([]byte, error) {
	var out []byte
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr {
		return out, ErrNoPointerInterface
	}
	if val.IsNil() {
		return out, ErrInvalidNilPointer
	}
	elem := val.Elem()
	typ := elem.Type()
	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			tField := typ.Field(i)
			// Unmarshal recursively if field of struct is a struct
			switch tField.Type.Kind() {
			case reflect.Struct:
				buff, err := Marshal(elem.Field(i).Addr().Interface(), endianness)
				if err != nil {
					return out, err
				}
				out = append(out, buff...)
			case reflect.Array:
				buff, err := marshalArray(elem.Field(i).Addr().Interface(), endianness)
				if err != nil {
					return out, err
				}
				out = append(out, buff...)
			case reflect.Slice:
				buff, err := marshalSlice(elem.Field(i).Addr().Interface(), endianness)
				if err != nil {
					return out, err
				}
				out = append(out, buff...)
			default:
				buff, err := Marshal(elem.Field(i).Addr().Interface(), endianness)
				if err != nil {
					return out, err
				}
				out = append(out, buff...)
			}
		}
	case reflect.Array:
		buff, err := marshalArray(elem.Addr().Interface(), endianness)
		if err != nil {
			return out, err
		}
		out = append(out, buff...)

	case reflect.Slice:
		buff, err := marshalSlice(elem.Addr().Interface(), endianness)
		if err != nil {
			return out, err
		}
		out = append(out, buff...)

	default:
		writter := new(bytes.Buffer)
		if err := binary.Write(writter, endianness, elem.Interface()); err != nil {
			return out, err
		}
		out = append(out, writter.Bytes()...)
	}
	return out, nil
}

func UnmarshaInitSlice(reader io.Reader, data interface{}, endianness Endianness) error {
	val := reflect.ValueOf(data)
	if val.IsNil() {
		return ErrInvalidNilPointer
	}
	slice := val.Elem()
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("Not a slice object")
	}
	if slice.Len() == 0 {
		return fmt.Errorf("Not initialized slice")
	}
	for k := 0; k < slice.Len(); k++ {
		err := Unmarshal(reader, slice.Index(k).Addr().Interface(), endianness)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalArray(reader io.Reader, data interface{}, endianness Endianness) error {
	val := reflect.ValueOf(data)
	if val.IsNil() {
		return ErrInvalidNilPointer
	}
	array := val.Elem()
	if array.Kind() != reflect.Array {
		return fmt.Errorf("Not an Array structure")
	}
	for k := 0; k < array.Len(); k++ {
		err := Unmarshal(reader, array.Index(k).Addr().Interface(), endianness)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalSlice(reader io.Reader, data interface{}, endianness Endianness) error {
	var sliceLen int64
	val := reflect.ValueOf(data)
	if val.IsNil() {
		return ErrInvalidNilPointer
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Slice {
		return fmt.Errorf("Not a Slice object")
	}
	err := Unmarshal(reader, &sliceLen, endianness)
	if err != nil {
		return err
	}
	s := elem
	newS := reflect.MakeSlice(s.Type(), int(sliceLen), int(sliceLen))
	s.Set(newS)
	//return UnmarshaInitSlice(reader, newS.Interface(), endianness)

	for k := 0; k < s.Len(); k++ {
		err := Unmarshal(reader, s.Index(k).Addr().Interface(), endianness)
		if err != nil {
			return err
		}
	}
	return nil
}

func Unmarshal(reader io.Reader, data interface{}, endianness Endianness) error {
	val := reflect.ValueOf(data)
	if val.IsNil() {
		return ErrInvalidNilPointer
	}
	elem := val.Elem()
	typ := elem.Type()
	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			tField := typ.Field(i)
			// Unmarshal recursively if field of struct is a struct
			switch tField.Type.Kind() {
			case reflect.Struct:
				err := Unmarshal(reader, elem.Field(i).Addr().Interface(), endianness)
				if err != nil {
					return err
				}
			case reflect.Array:
				err := unmarshalArray(reader, elem.Field(i).Addr().Interface(), endianness)
				if err != nil {
					return err
				}
			case reflect.Slice:
				err := unmarshalSlice(reader, elem.Field(i).Addr().Interface(), endianness)
				if err != nil {
					return err
				}
			default:
				if err := Unmarshal(reader, elem.Field(i).Addr().Interface(), endianness); err != nil {
					return err
				}
			}
		}

	case reflect.Array:
		err := unmarshalArray(reader, elem.Addr().Interface(), endianness)
		if err != nil {
			return err
		}

	case reflect.Slice:
		err := unmarshalSlice(reader, elem.Addr().Interface(), endianness)
		if err != nil {
			return err
		}

	default:
		if err := binary.Read(reader, endianness, data); err != nil {
			return err
		}
	}
	return nil
}
