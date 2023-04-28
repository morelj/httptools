package body

import (
	"io"
	"net/url"
	"strings"
)

// URLEncodedFormBody returns a Body which is the URL-encoded form of v with the application/x-www-form-urlencoded
// Content-Type.
type URLEncodedFormBody url.Values

func (b URLEncodedFormBody) ProvideBody() (io.Reader, string, error) {
	return strings.NewReader(url.Values(b).Encode()), "application/x-www-form-urlencoded", nil
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
