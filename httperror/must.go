package httperror

import "net/http"

// Must panics if err is not nil. It is equivalent to MustWithStatus(err, http.StatusInternalServerError)
func Must(err error) {
	MustWithStatus(err, http.StatusInternalServerError)
}

// MustWithStatus panics with the given status code if err is not nil.
func MustWithStatus(err error, status int) {
	if err != nil {
		switch err := err.(type) {
		case Error:
			panic(err)

		default:
			panic(NewWithError(err, status, err.Error()))
		}
	}
}
