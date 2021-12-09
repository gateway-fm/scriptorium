package error_source

import (
	"fmt"
	"strings"
)

// TODO: TEMPORARY!!!

// ErrorSource represent available services that can cause error
type ErrorSource int

const (
	errorSourceUnsupported ErrorSource = iota // TODO: REPLACING/REFACTORING IS REQUIRED!!!
	None
)

// errorSources is slice of error source string representations
var errorSources = [...]string{}

// String return ErrorSource enum as a string
func (s ErrorSource) String() string {
	return errorSources[s]
}

// ErrorSourceFromString return new ErrorSource enum from given string
func ErrorSourceFromString(s string) ErrorSource {
	for i, r := range errorSources {
		if strings.ToLower(s) == r {
			return ErrorSource(i)
		}
	}
	return errorSourceUnsupported
}

// ErrorSourceFromStringE return new ErrorSource enum from given string or return an error
func ErrorSourceFromStringE(s string) (ErrorSource, error) {
	for i, r := range errorSources {
		if strings.ToLower(s) == r {
			return ErrorSource(i), nil
		}
	}
	return errorSourceUnsupported, fmt.Errorf("invalid error source value %q", s)
}
