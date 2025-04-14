package pointer_test

import (
	"testing"

	"github.com/gateway-fm/scriptorium/pointer"
)

func TestSafeDeref(t *testing.T) {
	var nilPtr *int
	if got := pointer.SafeDeref(nilPtr); got != 0 {
		t.Errorf("SafeDeref(nilPtr) = %v, want %v", got, 0)
	}

	value := 123
	ptr := &value
	if got := pointer.SafeDeref(ptr); got != 123 {
		t.Errorf("SafeDeref(ptr) = %v, want %v", got, 123)
	}
}

func TestRef(t *testing.T) {
	value := 456
	ptr := pointer.Ref(value)
	if *ptr != 456 {
		t.Errorf("*Ref(value) = %v, want %v", *ptr, 456)
	}

	value = 789 // nolint:ineffassign
	if *ptr == 789 {
		t.Errorf("pointer should not be affected by changes to original variable")
	}
}
