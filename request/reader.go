package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Reader struct {
	r *http.Request
}

func NewReader(r *http.Request) Reader {
	return Reader{r: r}
}

func (r Reader) Bytes() ([]byte, error) {
	defer r.r.Body.Close()
	return ioutil.ReadAll(r.r.Body)
}

func (r Reader) MustBytes() []byte {
	data, err := r.Bytes()
	if err != nil {
		panic(err)
	}
	return data
}

func (r Reader) String() (string, error) {
	data, err := r.Bytes()
	return string(data), err
}

func (r Reader) MustString() string {
	return string(r.MustBytes())
}

func (r Reader) UnmarshalJSON(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (r Reader) MustUnmarshalJSON(v interface{}) {
	if err := r.UnmarshalJSON(v); err != nil {
		panic(err)
	}
}
