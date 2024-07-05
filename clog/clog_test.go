package clog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gateway-fm/scriptorium/clog"
)

const msgKey = "msg"

func TestCustomLogger(t *testing.T) {
	var buf bytes.Buffer

	logger := clog.NewCustomLogger(&buf, clog.LevelDebug, true)

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

	logger := clog.NewCustomLogger(&buf, clog.LevelInfo, true)

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

func TestAddKeysValuesToCtxConcurrentAccess(t *testing.T) {
	var buffer bytes.Buffer
	ctx := context.Background()

	logger := clog.NewCustomLogger(&buffer, clog.LevelInfo, true)

	var wg sync.WaitGroup

	repeat := 100
	wg.Add(repeat)

	for i := 0; i < repeat; i++ {
		go func(ctx context.Context) {
			defer wg.Done()

			ctx = logger.AddKeysValuesToCtx(ctx, map[string]interface{}{
				"timestamp": time.Now(),
			})

			logger.InfoCtx(ctx, "sample log message")

		}(ctx)
	}

	wg.Wait()
}

type testStruct struct {
	Field1 string
	Field2 int
}

func TestCustomLoggerWithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := clog.NewCustomLogger(&buf, clog.LevelInfo, true)

	testCh := make(chan int, 1)
	testCh <- 0

	s := testStruct{
		Field1: "value1",
		Field2: 100,
	}

	ctx := logger.AddKeysValuesToCtx(context.Background(), map[string]interface{}{
		"userID":     12345,
		"userName":   "testuser",
		"time":       time.Now(),
		"data":       []int{1, 2, 3},
		"testCh":     testCh,
		"testStruct": s,
	})

	logger.InfoCtx(ctx, "User %d logged in", 12345)
	require.Contains(t, buf.String(), "User 12345 logged in")
	require.Contains(t, buf.String(), "userID")
	require.Contains(t, buf.String(), "userName")
	require.Contains(t, buf.String(), "time")
	require.Contains(t, buf.String(), "data")
	require.Contains(t, buf.String(), "testCh")
	require.Contains(t, buf.String(), "testStruct")

	buf.Reset()

	err := errors.New("something went wrong")
	logger.ErrorCtx(ctx, err, "Failed to process user %d", 12345)
	require.Contains(t, buf.String(), "Failed to process user 12345")
	require.Contains(t, buf.String(), "something went wrong")
	require.Contains(t, buf.String(), "userID")
	require.Contains(t, buf.String(), "userName")
	require.Contains(t, buf.String(), "time")
	require.Contains(t, buf.String(), "data")
	require.Contains(t, buf.String(), "testCh")
	require.Contains(t, buf.String(), "testStruct")
}
