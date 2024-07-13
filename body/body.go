package body

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ContentTyper interface {
	ContentType() string
}

type Reader interface {
	ReadBody(r io.Reader) error
}

type Writer interface {
	WriteBody(w io.Writer) (int, error)
}

type Provider interface {
	ProvideBody() (io.Reader, error)
}

// ReaderFor returns an io.Reader which will read the bytes of body.
//
// body may be of one of these types:
// - nil. In this case, nil is returned.
// - a Provider: Returns the io.Reader returned by the ProvideBody() method.
// - a Writer: The WriteBody method is called.
// - an io.Reader: Returned as is
// - a string
// - a []byte
// - a fmt.Stringer. The string representation of body is used.
//
// If body does not match any of these types, ReaderFor panics.
func ReaderFor(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	if body, ok := body.(Provider); ok {
		return body.ProvideBody()
	}
	if body, ok := body.(Writer); ok {
		var buf bytes.Buffer
		if _, err := body.WriteBody(&buf); err != nil {
			return nil, err
		}
		return &buf, nil
	}
	if body, ok := body.(io.Reader); ok {
		return body, nil
	}
	if body, ok := body.(string); ok {
		return strings.NewReader(body), nil
	}
	if body, ok := body.([]byte); ok {
		return bytes.NewReader(body), nil
	}
	if body, ok := body.(fmt.Stringer); ok {
		return strings.NewReader(body.String()), nil
	}
	panic(fmt.Sprintf("unsupported body type %T", body))
}

func ContentTypeFor(body any) string {
	if contentTyper, ok := body.(ContentTyper); ok {
		return contentTyper.ContentType()
	}
	return ""
}

func Write(body any, w io.Writer) (int, error) {
	if body == nil {
		return 0, nil
	}
	if body, ok := body.(Writer); ok {
		return body.WriteBody(w)
	}
	if body, ok := body.(Provider); ok {
		r, err := body.ProvideBody()
		if err != nil {
			return 0, err
		}
		n, err := io.Copy(w, r)
		return int(n), err
	}
	if body, ok := body.(io.Reader); ok {
		n, err := io.Copy(w, body)
		return int(n), err
	}
	if body, ok := body.(string); ok {
		return w.Write([]byte(body))
	}
	if body, ok := body.([]byte); ok {
		return w.Write(body)
	}
	if body, ok := body.(fmt.Stringer); ok {
		return w.Write([]byte(body.String()))
	}
	panic(fmt.Sprintf("unsupported body type %T", body))
}
