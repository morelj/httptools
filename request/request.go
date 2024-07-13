package request

import (
	"context"
	"net/http"

	"github.com/morelj/httptools/body"
	"github.com/morelj/httptools/header"
)

// A Request wraps an *http.Request and an error.
// Request is built using NewRequestContext or NewRequest.
type Request struct {
	*http.Request
	Err error
}

// NewRequest calls NewRequestContext with the Background context.
func NewRequest(method, url string, bdy any, options ...Option) *Request {
	return NewRequestContext(context.Background(), method, url, bdy, options...)
}

// NewRequestContext is similar to NewWithContext, the difference being it returns
// a Request which wraps the error returned by NewWithContext, if any.
// The returned Request can directly be used with DoRequest or DoRequestDefault,
// allowing to delay the error handling.
func NewRequestContext(ctx context.Context, method, url string, bdy any, options ...Option) *Request {
	req, err := NewWithContext(ctx, method, url, bdy, options...)
	return &Request{
		Request: req,
		Err:     err,
	}
}

// DoRequest executes req using client.
// If req.Err is not nil, no request is made and the response's
// Err field will be set to req.Err.
func DoRequest(client *http.Client, req *Request) *Response {
	if req.Err != nil {
		return &Response{
			Err: req.Err,
		}
	}
	return Do(client, req.Request)
}

// DoRequestDefault calls DoRequest with the default HTTP client.
func DoRequestDefault(req *Request) *Response {
	return DoRequest(http.DefaultClient, req)
}

// NewWithContext returns a new http.Request with the given body and options.
// The body may be nil. When non-nil, the body will also be able to define the Content-Type header.
// Options are processed in order after the request is created. Using options it is possible to override all aspects
// of the request.
func NewWithContext(ctx context.Context, method, url string, b any, options ...Option) (*http.Request, error) {
	r, err := body.ReaderFor(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return nil, err
	}

	if b, ok := b.(body.ContentTyper); ok {
		if contentType := b.ContentType(); contentType != "" {
			req.Header.Set(header.ContentType, contentType)
		}
	}

	return req, Apply(req, options...)
}

func New(method, url string, b any, options ...Option) (*http.Request, error) {
	return NewWithContext(context.Background(), method, url, b, options...)
}

// Apply applies each option in order on the request.
func Apply(req *http.Request, options ...Option) error {
	if len(options) > 0 {
		for _, option := range options {
			if err := option(req); err != nil {
				return err
			}
		}
	}
	return nil
}

// Do executes the request on the given client.
func Do(client *http.Client, req *http.Request) *Response {
	res, err := client.Do(req)
	return &Response{
		Response: res,
		Err:      err,
	}
}

func DoDefault(req *http.Request) *Response {
	return Do(http.DefaultClient, req)
}
