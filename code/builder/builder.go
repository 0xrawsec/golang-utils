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
	case reflect.Ptr:
		if v.IsZero() {
			return "nil"
		}
		return repr(v.Elem().Interface())
	case reflect.Struct:
		t := v.Type()
		structStr := make([]string, 0, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			fieldName := t.Field(i).Name
			if string(fieldName[0]) == strings.ToUpper(string(fieldName[0])) {
				structStr = append(structStr, repr(v.Field(i).Interface()))
			}
		}
		return fmt.Sprintf("%T{%s}", v.Interface(), strings.Join(structStr, ","))
	case reflect.String:
		s := fmt.Sprintf("%s", v)
		if strings.Contains(s, `"`) {
			s = strings.Replace(s, "`", "'", -1)
			return fmt.Sprintf("`%s`", s)
		}
		return fmt.Sprintf("\"%s\"", s)
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
