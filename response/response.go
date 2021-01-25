package response

import (
	"encoding/json"
	"net/http"

	"github.com/morelj/httptools/header"
)

// A SerializerFunc serializes a body into bytes.
type SerializerFunc func(v interface{}) ([]byte, error)

// Builder provides an API to build an HTTP response before writing it to an http.Writer.
type Builder struct {
	headers    http.Header
	statusCode int
	body       interface{}
	serializer SerializerFunc
}

// NewBuilder returns a new, ready to use Builder.
// The buildser is initialized with a 200 response, no headers and a default serializer.
func NewBuilder() *Builder {
	return &Builder{
		headers:    http.Header{},
		statusCode: http.StatusOK,
		serializer: DefaultSerializer,
	}
}

// WithStatus sets the status code of the response
func (b *Builder) WithStatus(statusCode int) *Builder {
	b.statusCode = statusCode
	return b
}

// WithHeader sets an header on the response
func (b *Builder) WithHeader(key, value string) *Builder {
	b.headers.Set(key, value)
	return b
}

// WithHeaders sets all the headers of h to the response
func (b *Builder) WithHeaders(h http.Header) *Builder {
	for k, v := range h {
		b.headers[k] = v
	}
	return b
}

// WithBody sets the body of the response which must be either a []byte or a string.
// The body will be serialized with the default serializer.
func (b *Builder) WithBody(body interface{}) *Builder {
	b.body = body
	b.serializer = DefaultSerializer
	return b
}

// WithCustomBody sets the body of the response along with a serializer.
func (b *Builder) WithCustomBody(body interface{}, serializer SerializerFunc) *Builder {
	b.body = body
	b.serializer = serializer
	return b
}

// WithJSONBody sets the body of the response with a JSON serializer.
func (b *Builder) WithJSONBody(body interface{}) *Builder {
	return b.WithCustomJSONBody(body, false)
}

// WithCustomJSONBody sets the body of the response with a JSON serializer which can optionally be indented.
func (b *Builder) WithCustomJSONBody(body interface{}, indent bool) *Builder {
	b.body = body
	if indent {
		b.serializer = jsonMarshalIndent
	} else {
		b.serializer = json.Marshal
	}
	return b.WithHeader(header.ContentType, "application/json")
}

// Write writes the response to w.
// If set, the body is serialized using the serializer.
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

	return nil
}

// MustWrite writes the response to w or panics in case of error.
func (b *Builder) MustWrite(w http.ResponseWriter) {
	if err := b.Write(w); err != nil {
		panic(err)
	}
}

func jsonMarshalIndent(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
