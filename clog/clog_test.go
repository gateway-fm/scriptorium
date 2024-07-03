package clog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"

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

func TestIsZeroValue(t *testing.T) {
	tcs := []struct {
		name     string
		value    any
		expected bool
	}{
		{
			"non-empty string",
			"abc",
			false,
		},
		{
			"empty string",
			"",
			true,
		},
		{
			"empty slice",
			[]int{},
			true,
		},
		{
			"non-empty slice",
			[]int{1, 2, 3},
			false,
		},
		{
			"nil slice",
			([]int)(nil),
			true,
		},
		{
			"empty map",
			map[string]int{},
			true,
		},
		{
			"non-empty map",
			map[string]int{"key": 1},
			false,
		},
		{
			"nil map",
			(map[string]int)(nil),
			true,
		},
		{
			"zero int",
			0,
			true,
		},
		{
			"non-zero int",
			42,
			false,
		},
		{
			"zero float",
			0.0,
			true,
		},
		{
			"non-zero float",
			3.14,
			false,
		},
		{
			"empty struct",
			struct{}{},
			true,
		},
		{
			"non-zero struct",
			struct{ A int }{A: 1},
			false,
		},
		{
			"zero struct with fields",
			struct{ A int }{A: 0},
			true,
		},
		{
			"nil pointer",
			(*int)(nil),
			true,
		},
		{
			"non-nil pointer",
			func() *int { v := 42; return &v }(),
			false,
		},
		{
			"nil interface",
			(interface{})(nil),
			true,
		},
		{
			"non-nil interface",
			interface{}(42),
			false,
		},
		{
			"nil channel",
			(chan int)(nil),
			true,
		},
		{
			"non-nil channel",
			make(chan int),
			true,
		},
		{
			"nil function",
			(func())(nil),
			true,
		},
		{
			"non-nil function",
			func() {},
			false,
		},
		{
			"zero struct with multiple fields",
			struct {
				A int
				B string
			}{A: 0, B: ""},
			true,
		},
		{
			"non-zero struct with multiple fields",
			struct {
				A int
				B string
			}{A: 1, B: "non-zero"},
			false,
		},
		{
			"nested zero struct",
			struct {
				A int
				B struct{ C int }
			}{A: 0, B: struct{ C int }{C: 0}},
			true,
		},
		{
			"nested non-zero struct",
			struct {
				A int
				B struct{ C int }
			}{A: 0, B: struct{ C int }{C: 1}},
			false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			out := clog.IsZeroOfUnderlyingType(tc.value)
			require.Equal(t, tc.expected, out)
		})
	}
}

type testStruct struct {
	Field1 string
	Field2 int
}

func TestCustomLoggerWithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := clog.NewCustomLogger(&buf, slog.LevelInfo, true)

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
