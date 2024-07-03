package clog_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/gateway-fm/scriptorium/clog"
)

func BenchmarkCustomLogger(b *testing.B) {
	var buf bytes.Buffer
	logger := clog.NewCustomLogger(&buf, clog.LevelDebug, true)

	ctx := logger.AddKeysValuesToCtx(context.Background(), map[string]interface{}{
		"userID":    12345,
		"userName":  "testuser",
		"timestamp": time.Now(),
		"data":      []int{1, 2, 3},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		logger.InfoCtx(ctx, "User %d logged in", 12345)
	}
}
