package httperror

import "net/http"

func Must(err error) {
	MustWithStatus(err, http.StatusInternalServerError)
}

func MustWithStatus(err error, status int) {
	if err != nil {
		switch err := err.(type) {
		case Error:
			panic(err)

		default:
			panic(New(status, err.Error()))
		}
	}
}
