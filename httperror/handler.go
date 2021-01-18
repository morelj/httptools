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

type ErrorResponseWriterFunc func(err Error, w http.ResponseWriter) error

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

func WriteTextErrorResponse(err Error, w http.ResponseWriter) error {
	return response.NewBuilder().
		WithStatus(err.StatusCode()).
		WithHeader(header.ContentType, "text/plain").
		WithBody(err.Error()).
		Write(w)
}

func NewJSONErrorResponseWriter(newValue func(err Error) interface{}) ErrorResponseWriterFunc {
	return ErrorResponseWriterFunc(func(err Error, w http.ResponseWriter) error {
		return response.NewBuilder().
			WithStatus(err.StatusCode()).
			WithJSONBody(newValue(err)).
			Write(w)
	})
}

var WriteDefaultJSONErrorResponse = NewJSONErrorResponseWriter(func(err Error) interface{} {
	return err
})
