package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

type Builder struct {
	r *http.Request
}

func NewTestBuilder(method, target string, body interface{}) Builder {
	var reader io.Reader
	if body != nil {
		switch body := body.(type) {
		case string:
			reader = bytes.NewReader(([]byte)(body))
		case []byte:
			reader = bytes.NewReader(body)
		case io.Reader:
			reader = body
		default:
			panic(fmt.Errorf("Unsupported body type: %T", body))
		}
	}
	return Builder{
		r: httptest.NewRequest(method, target, reader),
	}
}

func (b Builder) WithHeader(key, value string) Builder {
	b.r.Header.Set(key, value)
	return b
}

func (b Builder) Request() *http.Request {
	return b.r
}
