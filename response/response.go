package response

import (
	"encoding/json"
	"net/http"
)

type SerializerFunc func(v interface{}) ([]byte, error)

type Builder struct {
	headers    http.Header
	statusCode int
	body       interface{}
	serializer SerializerFunc
}

func NewBuilder() *Builder {
	return &Builder{
		statusCode: http.StatusOK,
		serializer: DefaultSerializer,
	}
}

func (b *Builder) WithStatus(statusCode int) *Builder {
	b.statusCode = statusCode
	return b
}

func (b *Builder) WithHeader(key, value string) *Builder {
	b.headers.Set(key, value)
	return b
}

func (b *Builder) WithBody(body interface{}) *Builder {
	b.body = body
	return b
}

func (b *Builder) WithJSONBody(body interface{}) *Builder {
	b.body = body
	b.serializer = json.Marshal
	return b.WithHeader("content-type", "application/json")
}

func (b *Builder) Write(w http.ResponseWriter) error {
	h := w.Header()
	for k, v := range b.headers {
		h[k] = v
	}
	w.WriteHeader(b.statusCode)

	if b.body != nil {
		body, err := b.serializer(b.body)
		if err != nil {
			return err
		}
		_, err = w.Write(body)
		return err
	}

	_, err := w.Write(nil)
	return err
}

func (b *Builder) MustWrite(w http.ResponseWriter) {
	if err := b.Write(w); err != nil {
		panic(err)
	}
}
