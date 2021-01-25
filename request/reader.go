package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Reader wraps an http.Request and provides helper functions to read from it.
type Reader struct {
	r *http.Request
}

// NewReader returns a new Reader initialized with the given request
func NewReader(r *http.Request) Reader {
	return Reader{r: r}
}

// Bytes returns the request's body bytes
func (r Reader) Bytes() ([]byte, error) {
	defer r.r.Body.Close()
	return ioutil.ReadAll(r.r.Body)
}

// MustBytes returns the request's body bytes, or panics in case of error
func (r Reader) MustBytes() []byte {
	data, err := r.Bytes()
	if err != nil {
		panic(err)
	}
	return data
}

// String returns the request's body as a string
func (r Reader) String() (string, error) {
	data, err := r.Bytes()
	return string(data), err
}

// MustString returns the request's body as a string, or panics in case of error
func (r Reader) MustString() string {
	return string(r.MustBytes())
}

// JSON parses the request's body as JSON into v
func (r Reader) JSON(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// JSON parses the request's body as JSON into v, or panics in case of error
func (r Reader) MustJSON(v interface{}) {
	if err := r.JSON(v); err != nil {
		panic(err)
	}
}
