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
		if v != nil && !isZeroValue(v) {
			attrs = append(attrs, slog.Any(string(k), v))
		}
	}
	return attrs
}

func isZeroValue(v interface{}) bool {
	t := reflect.TypeOf(v)
	if !t.Comparable() {
		return false
	}
	return v == reflect.Zero(t).Interface()
}
