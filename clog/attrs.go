package clog

import (
	"log/slog"
	"reflect"
)

type fieldKey string

type Fields map[fieldKey]interface{}

// ConvertToAttrs converts a map of custom fields to a slice of slog.Attr
func ConvertToAttrs(fields Fields) []any {
	var attrs []any
	for k, v := range fields {
		if v != nil && !IsZeroOfUnderlyingType(v) {
			attrs = append(attrs, slog.Any(string(k), v))
		}
	}
	return attrs
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	if x == nil {
		return true
	}

	rv := reflect.ValueOf(x)
	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			if !field.CanInterface() {
				continue // Skip unexported fields
			}
			if !IsZeroOfUnderlyingType(field.Interface()) {
				return false
			}
		}
		return true
	}

	return reflect.DeepEqual(x, reflect.Zero(rv.Type()).Interface())
}
