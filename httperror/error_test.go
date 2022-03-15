package httperror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type customErr string

func (c customErr) Error() string {
	return string(c)
}

func TestNew(t *testing.T) {
	cases := []struct {
		err     Error
		as      interface{}
		asOK    bool
		status  int
		message string
	}{
		{
			err:     NewWithError(customErr("error"), http.StatusInternalServerError, "message"),
			as:      new(customErr),
			asOK:    true,
			status:  http.StatusInternalServerError,
			message: "message",
		},
		{
			err:     New(http.StatusInternalServerError, "message"),
			as:      new(customErr),
			asOK:    false,
			status:  http.StatusInternalServerError,
			message: "message",
		},
		{
			err:     NewWithErrorf(customErr("error"), http.StatusInternalServerError, "error: %d", 42),
			as:      new(customErr),
			asOK:    true,
			status:  http.StatusInternalServerError,
			message: "error: 42",
		},
		{
			err:     Newf(http.StatusInternalServerError, "error: %d", 42),
			as:      new(customErr),
			asOK:    false,
			status:  http.StatusInternalServerError,
			message: "error: 42",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert := assert.New(t)

			if c.as != nil {
				assert.Equal(c.asOK, errors.As(c.err, c.as))
			}
			assert.Equal(c.status, c.err.StatusCode())
			assert.Equal(c.message, c.err.Error())
		})
	}
}
