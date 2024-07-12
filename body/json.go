package body

import (
	"bytes"
	"encoding/json"
	"io"
)

type JSONBody struct {
	target any
	prefix string
	indent string
}

func JSON(v any) JSONBody {
	return JSONBody{target: v}
}

func JSONIndent(v any, prefix, indent string) JSONBody {
	return JSONBody{
		target: v,
		prefix: prefix,
		indent: indent,
	}
}

func (b JSONBody) ReadBody(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, b.target)
}

func (b JSONBody) marshal() ([]byte, error) {
	if b.indent != "" || b.prefix != "" {
		return json.MarshalIndent(b.target, b.prefix, b.indent)
	}
	return json.Marshal(b.target)
}

func (b JSONBody) ProvideBody() (io.Reader, error) {
	data, err := b.marshal()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func (b JSONBody) WriteBody(w io.Writer) (int, error) {
	data, err := b.marshal()
	if err != nil {
		return 0, err
	}
	return w.Write(data)
}

func (b JSONBody) ContentType() string {
	return "application/json"
}

var (
	_ Reader       = JSONBody{}
	_ Writer       = JSONBody{}
	_ Provider     = JSONBody{}
	_ ContentTyper = JSONBody{}
)
