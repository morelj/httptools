package request

import (
	"context"
	"net/http"

	"github.com/morelj/httptools/body"
	"github.com/morelj/httptools/header"
)

// NewWithContext returns a new http.Request with the given body and options.
// The body may be nil. When non-nil, the body will also be able to define the Content-Type header.
// Options are processed in order after the request is created. Using options it is possible to override all aspects
// of the request.
func NewWithContext(ctx context.Context, method, url string, b body.Body, options ...Option) (*http.Request, error) {
	r, err := body.ProvideBody(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return nil, err
	}

	if b != nil {
		if contentType := b.ContentType(); contentType != "" {
			req.Header.Set(header.ContentType, contentType)
		}
	}

	return req, Apply(req, options...)
}

func New(method, url string, b body.Body, options ...Option) (*http.Request, error) {
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

// Do executes the request on the given client. In case of success, the action is executed in turn.
// The action may be nil.
// If the action returns an error, Do will return it along with the response.
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
