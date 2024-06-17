package httperror

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/morelj/httptools/header"
	"github.com/morelj/httptools/response"
	"github.com/morelj/httptools/stack"
	"github.com/morelj/log"
)

// An ErrorResponseWriterFunc is a function which writes an Error into a ResponseWriter
type ErrorResponseWriterFunc func(err Error, w http.ResponseWriter) error

// A LoggerFunc purpose is to log an error, after it has been wrapped.
// The original HTTP request and the call stack of the current goroutine are also provided.
type LoggerFunc func(r *http.Request, err Error, stack stack.Stack)

// A WrapperFunc must wrap a panic (r) into an Error.
type WrapperFunc func(r interface{}, stack stack.Stack) Error

// A WrapperContextFunc must wrap a panic (r) into an Error.
// ctx is the original http.Request's context.
type WrapperContextFunc func(ctx context.Context, r any, stack stack.Stack) Error

// NewMiddleware returns a middleware which will recover when subsequent handlers panics.
// The panic value is used to produce an error response using the ErrorResponseWriterFunc and write it to the
// ResponseWriter.
// If the panic value is an Error, it is used as is. Otherwise, the error is wrapped into an Error with the error
// code 500.
// Calling NewMiddleware(ew) is equivalent to calling NewCustomMiddleware(ew, Wrap, Log)
func NewMiddleware(ew ErrorResponseWriterFunc) mux.MiddlewareFunc {
	return NewCustomMiddleware(ew, Wrap, Log)
}

// NewCustomMiddleware returns a middleware which will recover when subsequent handlers panics.
//
// In case of panic:
// - logger is called to log the error
// - then wrap is called to obtain an Error from the value returned by recover
// - finally the error is serialized using ew
//
// This function is similar to NewCustomContextMiddleware but uses a WrapperFunc insteads of a WrapperContextFunc
func NewCustomMiddleware(ew ErrorResponseWriterFunc, wrap WrapperFunc, logger LoggerFunc) mux.MiddlewareFunc {
	return NewCustomContextMiddleware(ew, func(ctx context.Context, r any, stack stack.Stack) Error {
		return wrap(r, stack)
	}, logger)
}

// NewCustomContextMiddleware returns a middleware which will recover when subsequent handlers panics.
//
// In case of panic:
// - logger is called to log the error
// - then wrap is called to obtain an Error from the value returned by recover
// - finally the error is serialized using ew
func NewCustomContextMiddleware(ew ErrorResponseWriterFunc, wrap WrapperContextFunc, logger LoggerFunc) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack, err := stack.Parse(debug.Stack())
					if err != nil {
						log.Errorf("Failed to parse stack: %v", err)
					}

					// Wrap the error
					wrappedErr := wrap(r.Context(), rec, stack)

					// Log the error
					logger(r, wrappedErr, stack)

					if err := ew(wrappedErr, w); err != nil {
						log.Errorf("Error writing error: %v\n", err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	})
}

// Wrap is the default WrapperFunc.
// - If r is an Error, it is returned as is
// - If it is any other error type, it is wrapped into an Error with a 500 status code
// - If it is any other value, it returns a 500 Error with an error message
func Wrap(r interface{}, stack stack.Stack) Error {
	switch r := r.(type) {
	case Error:
		return r

	case error:
		return httpError{
			Message: r.Error(),
			Code:    http.StatusInternalServerError,
		}

	default:
		return httpError{
			Message: fmt.Sprintf("Panic: %v", r),
			Code:    http.StatusInternalServerError,
		}
	}
}

// Log is the default LoggerFunc.
// It logs the value of r and the raw stack.
func Log(r *http.Request, err Error, stack stack.Stack) {
	log.Errorf("Error (%d): %s\n", err.StatusCode(), err.Error())
	log.Errorf(string(stack.Raw))
}

// NoOpLog is a no-operation LoggerFunc.
// It does nothing at all.
func NoOpLog(*http.Request, Error, stack.Stack) {}

// WriteTextErrorResponse is an error response writer to be used with NewMiddleware.
// It serializes the error message in plain text.
func WriteTextErrorResponse(err Error, w http.ResponseWriter) error {
	return response.NewBuilder().
		WithStatus(err.StatusCode()).
		WithHeader(header.ContentType, "text/plain").
		WithBody(err.Error()).
		Write(w)
}

// NewJSONErrorResponseWriter returns an ErrorResponseWriterFunc which serializes errors into JSON.
// The passed newValue function must returns a new value which will be used as the serialization target.
func NewJSONErrorResponseWriter(newValue func(err Error) interface{}) ErrorResponseWriterFunc {
	return ErrorResponseWriterFunc(func(err Error, w http.ResponseWriter) error {
		return response.NewBuilder().
			WithStatus(err.StatusCode()).
			WithJSONBody(newValue(err)).
			Write(w)
	})
}

// WriteDefaultJSONErrorResponse is ErrorResponseWriterFunc which serializes errors into JSON using the default format.
//
// Results will look like:
//
//	{
//	    "message": "error message",
//	    "code": 404
//	}
var WriteDefaultJSONErrorResponse = NewJSONErrorResponseWriter(func(err Error) interface{} {
	return err
})
