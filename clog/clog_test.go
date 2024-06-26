package clog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gateway-fm/scriptorium/clog"
)

const msgKey = "msg"

func TestCustomLogger(t *testing.T) {
	var buf bytes.Buffer

	logger := clog.NewCustomLogger(&buf, slog.LevelDebug, true)

	ctx := context.Background()
	ctx = logger.AddKeysValuesToCtx(ctx, map[string]interface{}{"user": "testUser"})

	tests := []struct {
		name       string
		logFunc    func(ctx context.Context, msg string, args ...any)
		expected   map[string]interface{}
		errorInput error
	}{
		{
			name: "ErrorfCtx",
			logFunc: func(ctx context.Context, msg string, args ...any) {
				logger.ErrorCtx(ctx, fmt.Errorf("test error"), msg, args...)
			},
			expected:   map[string]interface{}{"level": "ERROR", "user": "testUser", "error": "test error", msgKey: "an error occurred"},
			errorInput: fmt.Errorf("test error"),
		},
		{
			name: "InfofCtx",
			logFunc: func(ctx context.Context, msg string, args ...any) {
				logger.InfoCtx(ctx, msg, args...)
			},
			expected: map[string]interface{}{"level": "INFO", "user": "testUser", msgKey: "informational message"},
		},
		{
			name: "DebugfCtx",
			logFunc: func(ctx context.Context, msg string, args ...any) {
				logger.DebugCtx(ctx, msg, args...)
			},
			expected: map[string]interface{}{"level": "DEBUG", "user": "testUser", msgKey: "debugging message"},
		},
		{
			name: "WarnfCtx",
			logFunc: func(ctx context.Context, msg string, args ...any) {
				logger.WarnCtx(ctx, msg, args...)
			},
			expected: map[string]interface{}{"level": "WARN", "user": "testUser", msgKey: "warning message"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			tc.logFunc(ctx, tc.expected[msgKey].(string))

			var actual map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &actual); err != nil {
				t.Fatalf("Failed to unmarshal log output: %v", err)
			}

			for key, expectedValue := range tc.expected {
				if actual[key] != expectedValue {
					t.Errorf("%s did not log correctly. Expected %v for %s, got %v", tc.name, expectedValue, key, actual[key])
				}
			}
		})
	}
}

func TestCustomLogger_Level(t *testing.T) {
	var buf bytes.Buffer

	logger := clog.NewCustomLogger(&buf, slog.LevelInfo, true)

	ctx := context.Background()
	ctx = logger.AddKeysValuesToCtx(ctx, map[string]interface{}{"user": "testUser"})

	tests := []struct {
		name       string
		logFunc    func(ctx context.Context, msg string, args ...any)
		expected   map[string]interface{}
		errorInput error
	}{
		{
			name: "DebugfCtx",
			logFunc: func(ctx context.Context, msg string, args ...any) {
				logger.DebugCtx(ctx, msg, args...)
			},
			expected: map[string]interface{}{"level": "DEBUG", "user": "testUser", msgKey: "debugging message"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			tc.logFunc(ctx, tc.expected[msgKey].(string))

			var actual map[string]interface{}
			require.Nil(t, actual)
		})
	}
}

func TestConvertToAttrsConcurrentAccess(t *testing.T) {
	testFields := clog.Fields{
		"user":    "testUser",
		"session": "xyz123",
		"role":    "admin",
	}

	var wg sync.WaitGroup

	repeat := 100
	wg.Add(repeat)

	for i := 0; i < repeat; i++ {
		go func() {
			defer wg.Done()
			_ = clog.ConvertToAttrs(testFields)
		}()
	}

	wg.Wait()
}
