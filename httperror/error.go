package httperror

import "fmt"

// Error represents an error which can be converted to an HTTP response
type Error interface {
	error

	// StatusCode returns the HTTP Status code corresponding to the error
	StatusCode() int
}

// httpError is the internal implementation of Error
type httpError struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	wrapped error  `json:"-"`
}

// Error returns the error's message
func (e httpError) Error() string {
	return e.Message
}

// StatusCode returns the error's status code
func (e httpError) StatusCode() int {
	return e.Code
}

func (e httpError) Unwrap() error {
	return e.wrapped
}

// New returns a new Error with the given status code and message.
// The returned error can be serialized to JSON.
func New(statusCode int, message string) Error {
	return NewWithError(nil, statusCode, message)
}

// New returns a new Error wrapping parent with the given status code and message.
// The returned error can be serialized to JSON.
func NewWithError(parent error, statusCode int, message string) Error {
	return httpError{
		Message: message,
		Code:    statusCode,
		wrapped: parent,
	}
}

// Newf returns a new Error with the given status code and message, allowing it to be formatted using fmt.Sprintf.
// The returned error can be serialized to JSON.
func Newf(statusCode int, format string, a ...interface{}) Error {
	return NewWithErrorf(nil, statusCode, format, a...)
}

// Newf returns a new Error wrapping parent with the given status code and message, allowing it to be formatted using fmt.Sprintf.
// The returned error can be serialized to JSON.
func NewWithErrorf(parent error, statusCode int, format string, a ...interface{}) Error {
	return httpError{
		wrapped: parent,
		Message: fmt.Sprintf(format, a...),
		Code:    statusCode,
	}
}
