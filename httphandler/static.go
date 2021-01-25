package httphandler

import (
	"net/http"

	"github.com/morelj/httptools/response"
)

type StaticHandler struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func (h StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := response.NewBuilder()
	if h.StatusCode != 0 {
		b.WithStatus(h.StatusCode)
	}
	if len(h.Header) > 0 {
		b.WithHeaders(h.Header)
	}
	if len(h.Body) > 0 {
		b.WithBody(h.Body)
	}
	b.MustWrite(w)
}
