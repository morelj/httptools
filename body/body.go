package body

import (
	"bytes"
	"io"
)

type Body interface {
	ContentType() string
}

type Provider interface {
	ProvideBody() (io.Reader, error)
}

type Reader interface {
	ReadBody(io.Reader) error
}

type Writer interface {
	WriteBody(io.Writer) error
}

func ProvideBody(body Body) (io.Reader, error) {
	if body != nil {
		if body, ok := body.(Provider); ok {
			return body.ProvideBody()
		}
		if body, ok := body.(Writer); ok {
			var buf bytes.Buffer
			err := body.WriteBody(&buf)
			if err != nil {
				return nil, err
			}
			return &buf, nil
		}
		panic("body is neither a Provider nor a Writer")
	}
	return nil, nil
}

func WriteBody(body Body, w io.Writer) error {
	if body != nil {
		if body, ok := body.(Writer); ok {
			return body.WriteBody(w)
		}
		if body, ok := body.(Provider); ok {
			r, err := body.ProvideBody()
			if err != nil {
				return err
			}
			_, err = io.Copy(w, r)
			return err
		}
		panic("body is neither a Provider nor a Writer")
	}
	return nil
}
