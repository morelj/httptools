package handlers

import (
	"io"
	"net/http"

	"github.com/morelj/httptools/body"
	"github.com/morelj/httptools/header"
	"github.com/morelj/httptools/response"
)

// Static is an implementation of http.Handler which always writes the same pre-defined response.
type Static struct {
	StatusCode int
	Body       body.Body
	Header     http.Header
}

func NewStatic(statusCode int, bdy body.Body, headers http.Header) (Static, error) {
	handler := Static{
		StatusCode: statusCode,
		Header:     http.Header{},
	}

	if bdy != nil {
		bodyReader, err := body.ProvideBody(bdy)
		if err != nil {
			return Static{}, err
		}
		data, err := io.ReadAll(bodyReader)
		if err != nil {
			return Static{}, err
		}
		handler.Body = body.Raw(data)
		if contentType := bdy.ContentType(); contentType != "" {
			handler.Header.Set(header.ContentType, contentType)
		}
	}

	for k, v := range headers {
		handler.Header[k] = v
	}

	return handler, nil
}

func (h Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := response.NewBuilder()
	if h.StatusCode != 0 {
		b.WithStatus(h.StatusCode)
	}
	if len(h.Header) > 0 {
		b.WithHeaders(h.Header)
	}
	if h.Body != nil {
		b.WithBody(h.Body)
	}
	b.MustWrite(w)
}
