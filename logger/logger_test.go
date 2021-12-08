package logger

import (
	"reflect"
	"testing"
)

func TestLog(t *testing.T) {
	testlog := Log()

	tests := []struct {
		name string
		want *Zaplog
	}{
		{"name", testlog},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Log(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Log() = %v, want %v", got, tt.want)
			}
		})
	}
}
