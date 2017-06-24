package builder

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/tools/imports"
)

/////////////////////////////// Utility function ///////////////////////////////

func joinFold(ss []string, sep string, foldWidth int) string {
	tmpWidth := 0
	var out []string
	for _, s := range ss {
		tmpWidth += len(s)
		if tmpWidth >= foldWidth {
			out = append(out, fmt.Sprintf("%s%s\n", s, sep))
			tmpWidth = 0
			continue
		}
		out = append(out, fmt.Sprintf("%s%s ", s, sep))
	}
	return strings.Join(out, "")
}

func repr(i interface{}) string {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Struct:
		structStr := make([]string, 0, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			structStr = append(structStr, repr(v.Field(i).Interface()))
		}
		return fmt.Sprintf("%T{%s}", v.Interface(), strings.Join(structStr, ","))
	case reflect.String:
		return fmt.Sprintf("\"%s\"", v)
	case reflect.Slice:
		elements := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			elements = append(elements, repr(v.Index(i).Interface()))
		}
		return fmt.Sprintf("%s{%v}", v.Type(), joinFold(elements, ", ", 80))
	case reflect.Map:
		elements := make([]string, 0, v.Len())
		for _, k := range v.MapKeys() {
			elements = append(elements, fmt.Sprintf("%s: %s", repr(k.Interface()), repr(v.MapIndex(k).Interface())))
		}
		return fmt.Sprintf("%s{\n%s}", v.Type(), strings.Join(elements, ",\n"))
	case reflect.Interface:
		return repr(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type CodeBuilder struct {
	bytes.Buffer
}

func (cb *CodeBuilder) Package(packageName string) {
	cb.WriteString(fmt.Sprintf("package %s\n", packageName))
}

func (cb *CodeBuilder) DefVariable(name string, value interface{}) {
	cb.WriteString(fmt.Sprintf("var %s = %s\n", name, repr(value)))
}

func (cb *CodeBuilder) ResolveImports() {
	out, err := imports.Process("", cb.Bytes(), nil)
	if err != nil {
		panic(err)
	}
	cb.Reset()
	cb.Write(out)
}
