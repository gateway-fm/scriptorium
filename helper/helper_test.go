package helper

import (
	"context"
	"testing"
)

func TestGetRequestID(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKeyRequestID, "ReqId")

	tests := []struct {
		name string
		want string
	}{
		{"ctx", "ReqId"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRequestID(ctx)
			if got != tt.want {
				t.Errorf("GetRequestID() = %v, want %v", got, tt.want)
			}
		})
	}
}
