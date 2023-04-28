package request

import (
	"net/http"
	"net/url"

	"github.com/morelj/httptools/body"
)

// Reader wraps an http.Request and provides helper functions to read from it.
type Reader struct {
	r *http.Request
}

// NewReader returns a new Reader initialized with the given request
func NewReader(r *http.Request) Reader {
	return Reader{r: r}
}

func (r Reader) ReadBody(bodyReader body.Reader) (err error) {
	defer func() {
		cErr := r.r.Body.Close()
		if err == nil {
			err = cErr
		}
	}()
	return bodyReader.ReadBody(r.r.Body)
}

func (r Reader) ReadRawBody() (body.Raw, error) {
	raw := body.Raw{}
	err := r.ReadBody(&raw)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (r Reader) MustReadRawBody() body.Raw {
	data, err := r.ReadRawBody()
	if err != nil {
		panic(err)
	}
	return data
}

func (r Reader) ReadJSONBody(v any) error {
	return r.ReadBody(body.JSON(v))
}

func (r Reader) MustReadJSONBody(v any) {
	if err := r.ReadJSONBody(v); err != nil {
		panic(err)
	}
}

func (r Reader) ReadURLEncodedBody() (url.Values, error) {
	var values body.URLEncodedFormBody
	err := r.ReadBody(&values)
	if err != nil {
		return nil, err
	}
	return url.Values(values), nil
}
