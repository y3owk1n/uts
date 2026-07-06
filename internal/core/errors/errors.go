// Package derrors provides a domain-specific error type.
package derrors

import (
	"errors"
	"fmt"
)

// Code is a domain-specific error code.
type Code string

const (
	// CodeInvalidInput is returned when the input is invalid.
	CodeInvalidInput Code = "INVALID_INPUT"
	// CodeExecFailed is returned when an external command fails.
	CodeExecFailed Code = "EXEC_FAILED"
	// CodeToolNotFound is returned when a tool is not found.
	CodeToolNotFound Code = "TOOL_NOT_FOUND"
	// CodeFileNotFound is returned when a file is not found.
	CodeFileNotFound Code = "FILE_NOT_FOUND"
	// CodeUnsupportedFormat is returned when an unsupported format is used.
	CodeUnsupportedFormat Code = "UNSUPPORTED_FORMAT"
	// CodeCompressionFailed is returned when compression fails.
	CodeCompressionFailed Code = "COMPRESSION_FAILED"
	// CodeConversionFailed is returned when conversion fails.
	CodeConversionFailed Code = "CONVERSION_FAILED"
	// CodeArchiveFailed is returned when archiving fails.
	CodeArchiveFailed Code = "ARCHIVE_FAILED"
)

// Error is a domain-specific error.
type Error struct {
	code    Code
	message string
	cause   error
	context map[string]any
}

// New creates a new domain-specific error.
func New(code Code, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

// Newf creates a new domain-specific error with a formatted message.
func Newf(code Code, format string, args ...any) *Error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}

// Code returns the error code.
func (e *Error) Code() Code {
	return e.code
}

// Message returns the error message.
func (e *Error) Message() string {
	return e.message
}

// Cause returns the underlying error.
func (e *Error) Cause() error {
	return e.cause
}

// Context returns the error context.
func (e *Error) Context() map[string]any {
	return e.context
}

// Error returns the error message.
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.message, e.cause)
	}

	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.cause
}

// Is returns true if the error is the same type as the target error.
func (e *Error) Is(target error) bool {
	targetError, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.code == targetError.code
}

// WithContext returns a new error with the given context value.
func (e *Error) WithContext(key string, value any) *Error {
	if e.context == nil {
		e.context = make(map[string]any)
	}

	e.context[key] = value

	return e
}

// Wrap returns a new error with the given cause.
func Wrap(err error, code Code, message string) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		code:    code,
		message: message,
		cause:   err,
	}
}

// Wrapf returns a new error with the given cause and formatted message.
func Wrapf(err error, code Code, format string, args ...any) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
		cause:   err,
	}
}

// IsCode returns true if the error is the same type as the target error and.
func IsCode(err error, code Code) bool {
	if domainErr, ok := errors.AsType[*Error](err); ok {
		return domainErr.code == code
	}

	return false
}

// GetCode returns the error code.
func GetCode(err error) Code {
	if domainErr, ok := errors.AsType[*Error](err); ok {
		return domainErr.code
	}

	return CodeInvalidInput
}
