package response

import (
	"net/http"

	"github.com/morelj/httptools/body"
	"github.com/morelj/httptools/header"
)

// Builder provides an API to build an HTTP response before writing it to an http.Writer.
type Builder struct {
	headers    http.Header
	statusCode int
	body       body.Body
}

// NewBuilder returns a new, ready to use Builder.
// The buildser is initialized with a 200 response, no headers and a default serializer.
func NewBuilder() *Builder {
	return &Builder{
		headers:    http.Header{},
		statusCode: http.StatusOK,
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

func (b *Builder) WithBody(bdy body.Body) *Builder {
	b.body = bdy
	return b
}

func (b *Builder) WithRawBody(data []byte) *Builder {
	b.body = body.Raw(data)
	return b
}

// WithJSONBody sets the body of the response with a JSON serializer.
func (b *Builder) WithJSONBody(bdy any) *Builder {
	return b.WithCustomJSONBody(bdy, false)
}

// WithCustomJSONBody sets the body of the response with a JSON serializer which can optionally be indented.
func (b *Builder) WithCustomJSONBody(bdy any, indent bool) *Builder {
	if indent {
		b.body = body.JSONIndent(bdy, "", "  ")
	} else {
		b.body = body.JSON(bdy)
	}
	return b
}

// Write writes the response to w.
// If set, the body is serialized using the serializer.
func (b *Builder) Write(w http.ResponseWriter) error {
	h := w.Header()

	if b.body != nil {
		if contentType := b.body.ContentType(); contentType != "" {
			h.Set(header.ContentType, contentType)
		}
	}

	for k, v := range b.headers {
		h[k] = v
	}
	w.WriteHeader(b.statusCode)

	if b.body != nil {
		return body.WriteBody(b.body, w)
	}

	return nil
}

// MustWrite writes the response to w or panics in case of error.
func (b *Builder) MustWrite(w http.ResponseWriter) {
	if err := b.Write(w); err != nil {
		panic(err)
	}
}
