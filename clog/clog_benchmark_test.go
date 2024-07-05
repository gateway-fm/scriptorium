package clog_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gateway-fm/scriptorium/clog"
)

func BenchmarkCustomLogger(b *testing.B) {
	logger := clog.NewCustomLogger(os.Stdout, clog.LevelDebug, true)

	ctx := logger.AddKeysValuesToCtx(context.Background(), map[string]interface{}{
		"userID":    12345,
		"userName":  "testuser",
		"timestamp": time.Now(),
		"data":      []int{1, 2, 3},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoCtx(ctx, "Some test message")
	}
}
