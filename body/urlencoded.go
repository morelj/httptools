package body

import (
	"io"
	"net/url"
	"strings"
)

// URLEncodedFormBody returns a Body which is the URL-encoded form of v with the application/x-www-form-urlencoded
// Content-Type.
type URLEncodedFormBody url.Values

func (b URLEncodedFormBody) ProvideBody() (io.Reader, error) {
	return strings.NewReader(url.Values(b).Encode()), nil
}

func (b URLEncodedFormBody) WriteBody(w io.Writer) (int, error) {
	return w.Write([]byte(url.Values(b).Encode()))
}

func (b *URLEncodedFormBody) ReadBody(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	res, err := url.ParseQuery(string(data))
	if err == nil {
		*b = URLEncodedFormBody(res)
	}
	return err
}

func (b *URLEncodedFormBody) ContentType() string {
	return "application/x-www-form-urlencoded"
}

var (
	_ Provider     = URLEncodedFormBody{}
	_ Writer       = URLEncodedFormBody{}
	_ Reader       = (*URLEncodedFormBody)(nil)
	_ ContentTyper = (*URLEncodedFormBody)(nil)
)
