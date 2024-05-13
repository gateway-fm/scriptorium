package clog

import "log/slog"

type fieldKey string

type fields map[fieldKey]interface{}

// convertToAttrs converts a map of custom fields to a slice of slog.Attr
func convertToAttrs(fields fields) []any {
	var attrs []any
	for k, v := range fields {
		attrs = append(attrs, slog.Any(string(k), v))
	}
	return attrs
}
