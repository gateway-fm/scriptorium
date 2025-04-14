package pointer

func SafeDeref[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

func Ref[T any](val T) *T {
	return &val
}
