package httperr

import (
	"reflect"
	"strings"
)

func pkgName(data interface{}) string {
	if data == nil {
		return ""
	}

	typ := reflect.TypeOf(data)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	parts := strings.SplitAfterN(typ.PkgPath(), "vendor/", 2)
	return parts[len(parts)-1]
}
