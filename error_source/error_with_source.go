package error_source

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IErrorWithSource is an interface which has additional Source() method;
// This interface is used for propagating errors through the application
// without loss of an error source
type IErrorWithSource interface {
	Source() ErrorSource
	Error() string
	Unwrap() error
}

// errorWithSource is a basic implementation of IErrorWithSource interface
type errorWithSource struct {
	source ErrorSource
	err    error
}

// Source returns the error source
func (e *errorWithSource) Source() ErrorSource {
	if e == nil {
		return None
	}
	return e.source
}

// Error returns the error message
func (e errorWithSource) Error() string {
	return e.err.Error()
}

// Unwrap returns the error which is inside of errorWithSource;
// This method is implemented to make support for usage of errors.Is function
func (e *errorWithSource) Unwrap() error {
	return e.err
}

// ErrorWithSource returns a new instance of IErrorWithSource by
// a given error source and an error
func ErrorWithSource(source ErrorSource, err error) IErrorWithSource {
	if err == nil {
		return nil
	}
	return &errorWithSource{source: source, err: err}
}

// FromError returns the ErrorSource of a given error;
// Returns error_source.None if a given error is not an instance of IErrorWithSource
// or if it doesn't have wrapped instance of IErrorWithSource inside;
// Returns error_source.None if no source is wrapped inside a given error
func FromError(err error) ErrorSource {
	if err == nil {
		return None
	}

	if errorWithSource, ok := err.(IErrorWithSource); ok {
		return errorWithSource.Source()
	}

	// recursively unwrap error until there is nothing to unwrap or the source is
	// successfully got from the error
	for errors.Unwrap(err) != nil {
		unwrappedErr := errors.Unwrap(err)
		if errorWithSource, ok := unwrappedErr.(IErrorWithSource); ok {
			return errorWithSource.Source()
		}
		err = unwrappedErr
	}
	return None
}

// IsNone returns bool indicating whether error_source.None is wrapped inside a given error;
// Returns true if there is no source wrapped inside an error
func IsNone(err error) bool {
	if err == nil {
		return true
	}
	return FromError(err) == None
}

// GrpcErrorWithSource returns a new grpc error with an error source wrapped inside;
// The result of this function should be used only as a returning error for grpc endpoint;
// If you need to locally propagate an error through the application please use ErrorWithSource function;
// Returns nil if provided error is nil;
// Returns the provided error without changes if any errors are encountered during wrapping
func GrpcErrorWithSource(source ErrorSource, err error) error {
	if err == nil {
		return nil
	}

	st := status.New(codes.Unknown, err.Error())
	stDetails := errdetails.ErrorInfo{Reason: source.String()}

	st, e := st.WithDetails(&stDetails)
	if e != nil {
		return err
	}
	return st.Err()
}

// FromGrpcError returns an errors source wrapped inside a given grpc error;
// Returns error_source.None if no source is wrapped inside;
// Returns error_source.None if provided error is nil
func FromGrpcError(err error) ErrorSource {
	if err == nil {
		return None
	}

	st := grpcStatusFromError(err)
	if st == nil {
		return None
	}

	for _, detail := range st.Details() {
		if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			return ErrorSourceFromString(errorInfo.Reason)
		}
	}
	return None
}

// grpcStatusFromError returns grpc status wrapped inside a given error;
// Returns nil if there is no status wrapped
func grpcStatusFromError(err error) *status.Status {
	if err == nil {
		return nil
	}

	// try to get the status from original error
	if st, ok := status.FromError(err); ok {
		return st
	}

	// recursively unwrap error until there is nothing to unwrap or the status is
	// successfully got from the error
	for errors.Unwrap(err) != nil {
		unwrappedErr := errors.Unwrap(err)
		if st, ok := status.FromError(unwrappedErr); ok {
			return st
		}
		err = unwrappedErr
	}
	return nil
}
