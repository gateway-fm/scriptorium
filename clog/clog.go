package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

func NewCustomLogger(writer io.Writer, level slog.Level, addSource bool) *CustomLogger {
	return &CustomLogger{
		Logger: slog.New(
			slog.NewJSONHandler(
				writer,
				&slog.HandlerOptions{
					AddSource: addSource,
					Level:     level,
				},
			),
		),
		ctxKeys: []fieldKey{},
	}
}

type CustomLogger struct {
	*slog.Logger

	mu      sync.RWMutex
	ctxKeys []fieldKey
}

// ErrorfCtx logs an error message with fmt.SprintF()
func (l *CustomLogger) ErrorfCtx(ctx context.Context, err error, msg string, args ...any) {
	l.With(convertToAttrs(l.fromCtx(ctx))...).With(slog.String("error", err.Error())).ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

// InfofCtx logs an informational message with fmt.SprintF()
func (l *CustomLogger) InfofCtx(ctx context.Context, msg string, args ...any) {
	l.With(convertToAttrs(l.fromCtx(ctx))...).InfoContext(ctx, fmt.Sprintf(msg, args...))
}

// DebugfCtx logs a debug message with fmt.SprintF()
func (l *CustomLogger) DebugfCtx(ctx context.Context, msg string, args ...any) {
	l.With(convertToAttrs(l.fromCtx(ctx))...).DebugContext(ctx, fmt.Sprintf(msg, args...))
}

// WarnfCtx logs a debug message with fmt.SprintF()
func (l *CustomLogger) WarnfCtx(ctx context.Context, msg string, args ...any) {
	l.With(convertToAttrs(l.fromCtx(ctx))...).WarnContext(ctx, fmt.Sprintf(msg, args...))
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

func (l *CustomLogger) fromCtx(ctx context.Context) fields {
	l.mu.Lock()
	defer l.mu.Unlock()

	f := make(fields)
	for _, ctxKey := range l.ctxKeys {
		f[ctxKey] = ctx.Value(ctxKey)
	}

	return f
}
