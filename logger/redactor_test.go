package logger

import (
	"testing"
)

func Test_Redactor_RemovesIPAddresses(t *testing.T) {
	r := NewRedactor()

	m := map[string]string{
		"this contains a valid ip: 172.217.22.14":  "this contains a valid ip: xxx.xxx.xxx.xxx",
		"192.168.1.1 this started with a valid ip": "xxx.xxx.xxx.xxx this started with a valid ip",
		"this contains a non-ip 1234.12.15.2":      "this contains a non-ip 1234.12.15.2",
		"this has ip with port, 192.168.1.1:80":    "this has ip with port, xxx.xxx.xxx.xxx:80",
	}

	for msg, res := range m {
		if r.Redact(msg) != res {
			t.Errorf("Expected %s, got %s", res, r.Redact(msg))
		}
	}
}
