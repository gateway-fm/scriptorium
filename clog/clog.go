package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

func NewCustomLogger(dest io.Writer, level Level, addSource bool) *CustomLogger {
	return &CustomLogger{
		Logger: slog.New(
			slog.NewJSONHandler(
				dest,
				&slog.HandlerOptions{
					AddSource: addSource,
					Level:     slog.Level(level),
				})),
	}
}

type CustomLogger struct {
	*slog.Logger
}

// ErrorCtx logs an error message with fmt.SprintF()
func (l *CustomLogger) ErrorCtx(ctx context.Context, err error, msg string, args ...any) {
	l.With(convertToAttrs(l.fieldsFromCtx(ctx))...).With(slog.String("error", err.Error())).ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

// InfoCtx logs an informational message with fmt.SprintF()
func (l *CustomLogger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.With(convertToAttrs(l.fieldsFromCtx(ctx))...).InfoContext(ctx, fmt.Sprintf(msg, args...))
}

// DebugCtx logs a debug message with fmt.SprintF()
func (l *CustomLogger) DebugCtx(ctx context.Context, msg string, args ...any) {
	l.With(convertToAttrs(l.fieldsFromCtx(ctx))...).DebugContext(ctx, fmt.Sprintf(msg, args...))
}

// WarnCtx logs a debug message with fmt.SprintF()
func (l *CustomLogger) WarnCtx(ctx context.Context, msg string, args ...any) {
	l.With(convertToAttrs(l.fieldsFromCtx(ctx))...).WarnContext(ctx, fmt.Sprintf(msg, args...))
}

// convertToAttrs converts a map of custom fields to a slice of slog.Attr
func convertToAttrs(fields map[string]interface{}) []any {
	var attrs []any

	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	return attrs
}
