package request

import (
	"net/http"
	"net/url"

	"github.com/morelj/httptools/body"
)

type Response struct {
	*http.Response
	Err error
}

func (r *Response) Close() error {
	if r.Body != nil {
		err := r.Body.Close()
		if r.Err == nil {
			r.Err = err
		}
		return err
	}
	return nil
}

func (r *Response) ReadBody(reader body.Reader) {
	if r.Err == nil && r.Body != nil {
		r.Err = reader.ReadBody(r.Body)
	}
}

func (r *Response) ReadJSONBody(target any) {
	r.ReadBody(body.JSON(target))
}

func (r *Response) ReadRawBody() body.Raw {
	var raw body.Raw
	r.ReadBody(&raw)
	return raw
}

func (r *Response) ReadURLEncodedBody() url.Values {
	var values body.URLEncodedFormBody
	r.ReadBody(&values)
	return url.Values(values)
}
