package datastructs

import (
	"reflect"
)

//ToInterfaceSlice converts any slice object to an array of interface{}
//it can be usefull to initialize some datastructs
func ToInterfaceSlice(slice interface{}) (is []interface{}) {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Slice {
		is = make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			is = append(is, v.Index(i).Interface())
		}
	} else {
		panic("parameter must be a slice")
	}
	return is
}
