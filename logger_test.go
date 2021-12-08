package scriptorium

import (
	"github.com/gateway-fm/scriptorium/logger"
	"testing"
)

func TestLogger_InitLogger(t *testing.T) {
	testlogger := logger.Log()
	goodenv := "production"
	type fields struct {
		logger *logger.Zaplog
		Env    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success logger init",
			fields: fields{
				logger: testlogger,
				Env:    "production",
			},
			wantErr: false,
		},
		{
			name: "Failed logger init",
			fields: fields{
				logger: testlogger,
				Env:    "nurgl",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Logger{
				logger: tt.fields.logger,
				Env:    tt.fields.Env,
			}
			if err := c.InitLogger(); (err != nil && tt.fields.Env != goodenv) != tt.wantErr {
				t.Errorf("InitLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
