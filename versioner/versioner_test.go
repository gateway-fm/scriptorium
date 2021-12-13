package version

import (
	"bytes"
	"testing"
	"time"
)

func TestVersion_GetName(t *testing.T) {
	type fields struct {
		Service string
		Tag     string
		Commit  string
		Branch  string
		URL     string
		Date    time.Time
		msg     bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct naming",
			fields: fields{Service: "scriptorium"},
			want:   "Scriptorium",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Service: tt.fields.Service,
				Tag:     tt.fields.Tag,
				Commit:  tt.fields.Commit,
				Branch:  tt.fields.Branch,
				URL:     tt.fields.URL,
				Date:    tt.fields.Date,
				msg:     tt.fields.msg,
			}
			if got := v.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_IfSpecified(t *testing.T) {
	type fields struct {
		Service string
		Tag     string
		Commit  string
		Branch  string
		URL     string
		Date    time.Time
		msg     bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "All Specified",
			fields: fields{
				Service: "GoMK",
				Tag:     "v0.0.1",
				Commit:  "5eadc39d3eba25d99273499aa3928bb70e6d9d35",
				Branch:  "main",
				URL:     "github.com/gateway-fm/scriptorium/versioner",
				Date:    time.Now(),
			},
			want: true,
		},
		{
			name: "Name not specified",
			fields: fields{
				Service: unspecified,
				Tag:     "v0.0.1",
				Commit:  "5eadc39d3eba25d99273499aa3928bb70e6d9d35",
				Branch:  "main",
				URL:     "github.com/gateway-fm/scriptorium/versioner",
				Date:    time.Now(),
			},
			want: false,
		},
		{
			name: "Tag not specified",
			fields: fields{
				Service: "Scriptorium",
				Tag:     unspecified,
				Commit:  "5c731d4c13d946bee43d9243b82537aaf6ef8e4b",
				Branch:  "main",
				URL:     "github.com/gateway-fm/scriptorium/versioner",
				Date:    time.Now(),
			},
			want: true,
		},
		{
			name: "Commit not specified",
			fields: fields{
				Service: "Scriptorium",
				Tag:     "v0.0.1",
				Commit:  unspecified,
				Branch:  "main",
				URL:     "github.com/gateway-fm/scriptorium/versioner",
				Date:    time.Now(),
			},
			want: false,
		},
		{
			name: "Branch not specified",
			fields: fields{
				Service: "Scriptorium",
				Tag:     "v0.0.1",
				Commit:  "5c731d4c13d946bee43d9243b82537aaf6ef8e4b",
				Branch:  unspecified,
				URL:     "github.com/gateway-fm/scriptorium/versioner",
				Date:    time.Now(),
			},
			want: false,
		},
		{
			name: "Url not specified",
			fields: fields{
				Service: "Scriptorium",
				Tag:     "v0.0.1",
				Commit:  "5c731d4c13d946bee43d9243b82537aaf6ef8e4b",
				Branch:  "main",
				URL:     unspecified,
				Date:    time.Now(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Service: tt.fields.Service,
				Tag:     tt.fields.Tag,
				Commit:  tt.fields.Commit,
				Branch:  tt.fields.Branch,
				URL:     tt.fields.URL,
				Date:    tt.fields.Date,
				msg:     tt.fields.msg,
			}
			if got := v.IfSpecified(); got != tt.want {
				t.Errorf("IfSpecified() = %v, want %v", got, tt.want)
			}
		})
	}
}
