package httperror

import "fmt"

type Error interface {
	error
	StatusCode() int
}

type httpError struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func (e httpError) Error() string {
	return e.Message
}

func (e httpError) StatusCode() int {
	return e.Code
}

func New(statusCode int, message string) Error {
	return httpError{
		Message: message,
		Code:    statusCode,
	}
}

func Newf(statusCode int, format string, a ...interface{}) Error {
	return httpError{
		Message: fmt.Sprintf(format, a...),
		Code:    statusCode,
	}
}
