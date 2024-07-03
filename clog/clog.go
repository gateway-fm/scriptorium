package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

func NewCustomLogger(dest io.Writer, level slog.Level, addSource bool) *CustomLogger {
	return &CustomLogger{
		Logger: slog.New(slog.NewJSONHandler(
			dest,
			&slog.HandlerOptions{
				AddSource: addSource,
				Level:     level,
			})),
		ctxKeys: []fieldKey{},
	}
}

type CustomLogger struct {
	*slog.Logger

	mu      sync.RWMutex
	ctxKeys []fieldKey
}

// ErrorCtx logs an error message with fmt.SprintF()
func (l *CustomLogger) ErrorCtx(ctx context.Context, err error, msg string, args ...any) {
	l.With(ConvertToAttrs(l.fromCtx(ctx))...).With(slog.String("error", err.Error())).ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

// InfoCtx logs an informational message with fmt.SprintF()
func (l *CustomLogger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.With(ConvertToAttrs(l.fromCtx(ctx))...).InfoContext(ctx, fmt.Sprintf(msg, args...))
}

// DebugCtx logs a debug message with fmt.SprintF()
func (l *CustomLogger) DebugCtx(ctx context.Context, msg string, args ...any) {
	l.With(ConvertToAttrs(l.fromCtx(ctx))...).DebugContext(ctx, fmt.Sprintf(msg, args...))
}

// WarnCtx logs a debug message with fmt.SprintF()
func (l *CustomLogger) WarnCtx(ctx context.Context, msg string, args ...any) {
	l.With(ConvertToAttrs(l.fromCtx(ctx))...).WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *CustomLogger) AddKeysValuesToCtx(ctx context.Context, kv map[string]interface{}) context.Context {
	l.mu.Lock()
	defer l.mu.Unlock()

	for k, v := range kv {
		ctx = context.WithValue(ctx, fieldKey(k), v)
		l.ctxKeys = append(l.ctxKeys, fieldKey(k))
	}

	return ctx
}

func (l *CustomLogger) fromCtx(ctx context.Context) Fields {
	l.mu.Lock()
	defer l.mu.Unlock()

	f := make(Fields)
	for _, ctxKey := range l.ctxKeys {
		if ctx.Value(ctxKey) != nil {
			f[ctxKey] = ctx.Value(ctxKey)
		}
	}

	return f
}
