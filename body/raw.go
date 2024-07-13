package body

import (
	"bytes"
	"encoding/hex"
	"io"
	"unicode/utf8"
)

type Raw []byte

func (b Raw) Bytes() []byte {
	return []byte(b)
}

func (b *Raw) ReadBody(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err == nil {
		*b = data
	}
	return err
}

func (b *Raw) ContentType() string {
	return ""
}

func (b Raw) ProvideBody() (io.Reader, error) {
	if b != nil {
		return bytes.NewReader(b), nil
	}
	return nil, nil
}

func (b Raw) WriteBody(w io.Writer) (int, error) {
	if b != nil {
		return w.Write(b.Bytes())
	}
	return 0, nil
}

func (b Raw) String() string {
	if utf8.Valid(b) {
		return string(b)
	}
	return hex.EncodeToString(b)
}

var (
	_ Provider     = Raw{}
	_ Writer       = Raw{}
	_ Reader       = (*Raw)(nil)
	_ ContentTyper = (*Raw)(nil)
)
