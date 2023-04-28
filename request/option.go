package request

import (
	"net/http"

	"github.com/morelj/httptools/header"
)

// An Option is a function which is able to alter a request
type Option func(req *http.Request) error

// WithHeader returns an Option which sets the given header on the request
func WithHeader(key, value string) Option {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

func WithContentType(value string) Option {
	return WithHeader(header.ContentType, value)
}

// WithBasicAuth returns an Option which enables Basic authentication on the request
func WithBasicAuth(username, password string) Option {
	return func(req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	}
}

// WithBearerAuth returns an Option which enables Bearer authentication on the request
func WithBearerAuth(token string) Option {
	return WithHeader(header.Authorization, "Bearer "+token)
}

func WithHost(host string) Option {
	return func(req *http.Request) error {
		req.Host = host
		return nil
	}
}

func WithClose(close bool) Option {
	return func(req *http.Request) error {
		req.Close = close
		return nil
	}
}

func WithContentLength(length int64) Option {
	return func(req *http.Request) error {
		req.ContentLength = length
		return nil
	}
}
