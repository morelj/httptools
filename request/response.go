package request

import (
	"net/http"
	"net/url"

	"github.com/morelj/httptools/body"
	"github.com/morelj/httptools/httperror"
)

// Response is the response to an HTTP request.
// It is a thin wrapper around *http.Response.
type Response struct {
	*http.Response
	Err error
}

// Close closes r.Body, if any.
// If r.Body.Close() returned an error, this error is returned and r.Err is
// updated with this error (only if it was nil).
func (r *Response) Close() error {
	if r.Response != nil && r.Response != nil && r.Body != nil {
		err := r.Body.Close()
		if r.Err == nil {
			r.Err = err
		}
		return err
	}
	return nil
}

// HasSuccessStatus returns true if r.Err is nil and the response status code
// is between 200 (included) and 300 (excluded).
func (r *Response) HasSuccessStatus() bool {
	return r.Err == nil && r.StatusCode >= 200 && r.StatusCode < 300
}

// HasRedirectStatus returns true if r.Err is nil and the response status code
// is between 300 (included) and 400 (excluded).
func (r *Response) HasRedirectStatus() bool {
	return r.Err == nil && r.StatusCode >= 300 && r.StatusCode < 400
}

// HasClientErrorStatus returns true if r.Err is nil and the response status code
// is between 400 (included) and 500 (excluded).
func (r *Response) HasClientErrorStatus() bool {
	return r.Err == nil && r.StatusCode >= 400 && r.StatusCode < 500
}

// HasServerErrorStatus returns true if r.Err is nil and the response status code
// is greater or equals to 500.
func (r *Response) HasServerErrorStatus() bool {
	return r.Err == nil && r.StatusCode >= 500
}

// HasErrorStatus returns true if r.Err is nil and the response status code
// is greater or equals to 400.
func (r *Response) HasErrorStatus() bool {
	return r.Err == nil && r.StatusCode >= 400
}

// GlobalError returns:
// - r.Err if non nil
// - An httperror.Error if the HTTP status code is >= 400
// - nil otherwise
func (r *Response) GlobalError() error {
	if r.Err != nil {
		return r.Err
	}
	if r.HasErrorStatus() {
		return httperror.New(r.StatusCode, r.Status)
	}
	return nil
}

// ReadBody reads the response body using the given reader.
// ReadBody has no effect if r.Err is not nil.
// If the body.Reader returns an error, r.Err is set to this error.
func (r *Response) ReadBody(reader body.Reader) {
	if r.Err == nil && r.Body != nil {
		r.Err = reader.ReadBody(r.Body)
	}
}

// ReadJSONBody reads the response body as a JSON value into target.
// ReadJSONBody has no effect if r.Err is not nil.
// If there's an error reading or unmarshalling the body, r.Err is set to this error.
func (r *Response) ReadJSONBody(target any) {
	r.ReadBody(body.JSON(target))
}

// ReadRawBody reads the response body as a raw slice of bytes.
// ReadRawBody has no effect if r.Err is not nil.
// If there's an error reading the body, r.Err is set to this error.
func (r *Response) ReadRawBody() body.Raw {
	var raw body.Raw
	r.ReadBody(&raw)
	return raw
}

// ReadURLEncodedBody reads the response body as an URL-encoded string and returns it.
// ReadURLEncodedBody has no effect if r.Err is not nil.
// If there's an error reading or parsing the body, r.Err is set to this error.
func (r *Response) ReadURLEncodedBody() url.Values {
	var values body.URLEncodedFormBody
	r.ReadBody(&values)
	return url.Values(values)
}
