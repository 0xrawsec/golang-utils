package datastructs

import (
	"reflect"
)

//ToInterfaceSlice converts any slice object to an array of interface{}
//it can be usefull to initialize some datastructs
func ToInterfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	is := make([]interface{}, 0)
	if v.Kind() == reflect.Slice {
		is = make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			is = append(is, v.Index(i).Interface())
		}
	}
	return is
}
