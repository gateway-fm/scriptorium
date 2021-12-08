package logger

import "testing"

func TestAppEnv_String(t *testing.T) {
	tests := []struct {
		name string
		s    AppEnv
		want string
	}{
		{
			name: "Production test",
			s:    Production,
			want: "prod",
		},
		{
			name: "Development test",
			s:    Development,
			want: "dev",
		},
		{
			name: "Local test",
			s:    Local,
			want: "local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvFromStr(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    AppEnv
		wantErr bool
	}{
		{
			name:    "Prod test",
			args:    args{s: "prod"},
			want:    Production,
			wantErr: false,
		},

		{
			name:    "Dev test",
			args:    args{s: "dev"},
			want:    Development,
			wantErr: false,
		},
		{
			name:    "Local test",
			args:    args{s: "local"},
			want:    Local,
			wantErr: false,
		},
		{
			name:    "Wrong test",
			args:    args{s: "nurgl"},
			want:    Wrong,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EnvFromStr(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvFromStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EnvFromStr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
