package httperror

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/morelj/httptools/header"
	"github.com/morelj/httptools/response"
)

// An ErrorResponseWriterFunc is a function which writes an Error into a ResponseWriter
type ErrorResponseWriterFunc func(err Error, w http.ResponseWriter) error

// NewMiddleware returns a middleware which will recover when subsequent handlers panics.
// The panic value is used to produce an error response using the ErrorResponseWriterFunc and write it to the
// ResponseWriter.
// If the panic value is an Error, it is used as is. Otherwise, the error is wrapped into an Error with the error
// code 500.
func NewMiddleware(ew ErrorResponseWriterFunc) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovering from panic: %v\n", r)
					debug.PrintStack()

					var err error

					switch r := r.(type) {
					case Error:
						err = ew(r, w)

					case error:
						err = ew(httpError{
							Message: r.Error(),
							Code:    http.StatusInternalServerError,
						}, w)

					default:
						err = ew(httpError{
							Message: fmt.Sprintf("Panic: %v", r),
							Code:    http.StatusInternalServerError,
						}, w)
					}

					if err != nil {
						log.Printf("Error writing error: %v\n", err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	})
}

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
//     {
//         "message": "error message",
//         "code": 404
//     }
var WriteDefaultJSONErrorResponse = NewJSONErrorResponseWriter(func(err Error) interface{} {
	return err
})
