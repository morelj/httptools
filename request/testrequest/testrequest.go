package testrequest

import (
	"net/http"
	"net/http/httptest"

	"github.com/morelj/httptools/body"
	"github.com/morelj/httptools/header"
	"github.com/morelj/httptools/request"
)

func NewRequest(method, target string, requestBody body.Body, options ...request.Option) *http.Request {
	r, err := body.ProvideBody(requestBody)
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(method, target, r)
	if contentType := requestBody.ContentType(); contentType != "" {
		req.Header.Set(header.ContentType, contentType)
	}
	request.Apply(req, options...)

	return req
}
