package clog

import (
	"log/slog"
)

// ConvertToAttrs converts a map of custom fields to a slice of slog.Attr
func ConvertToAttrs(fields map[string]interface{}) []any {
	var attrs []any

	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	return attrs
}
