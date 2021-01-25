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
}

// Error returns the error's message
func (e httpError) Error() string {
	return e.Message
}

// StatusCode returns the error's status code
func (e httpError) StatusCode() int {
	return e.Code
}

// New returns a new Error with the given status code and message.
// The returned error can be serialized to JSON.
func New(statusCode int, message string) Error {
	return httpError{
		Message: message,
		Code:    statusCode,
	}
}

// Newf returns a new Error with the given status code and message, allowing it to be formatted using fmt.Sprintf.
// The returned error can be serialized to JSON.
func Newf(statusCode int, format string, a ...interface{}) Error {
	return httpError{
		Message: fmt.Sprintf(format, a...),
		Code:    statusCode,
	}
}
