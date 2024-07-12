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
	Body       any
	Header     http.Header
	StatusCode int
}

func NewStatic(statusCode int, bdy any, headers http.Header) (Static, error) {
	handler := Static{
		StatusCode: statusCode,
		Header:     http.Header{},
	}

	if bdy != nil {
		bodyReader, err := body.ReaderFor(bdy)
		if err != nil {
			return Static{}, err
		}
		if closer, ok := bodyReader.(io.Closer); ok {
			defer closer.Close()
		}
		data, err := io.ReadAll(bodyReader)
		if err != nil {
			return Static{}, err
		}
		handler.Body = body.Raw(data)
		if contentType := body.ContentTypeFor(bdy); contentType != "" {
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
