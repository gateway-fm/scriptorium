package tracing

// Code is an 32-bit representation of a status state.
// It's design purpose is to remove the OpenTelemetry's codes.Code dependency from our interfaces
type Code uint32

// WARNING: please do not modify the order of constants or remove the existing ones, unless it is necessary to comply
// with OpenTelemetry's Go SDK. Otherwise, it will break the compatibility with OpenTelemetry's status codes.
const (
	Unset Code = iota

	Error

	Ok
)
