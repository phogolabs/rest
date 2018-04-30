package rho

import (
	"reflect"
	"strings"

	"github.com/go-openapi/inflect"
)

func typeName(data interface{}) string {
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

	if typ.Kind() == reflect.Struct {
		name := typ.Name()
		name = inflect.Singularize(name)
		name = strings.ToLower(name)
		return name
	}

	return ""
}
