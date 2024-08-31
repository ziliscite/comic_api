package helpers

import (
	"net/http"
)

type apiFunc func(http.ResponseWriter, *http.Request) (int, error)

func ServeHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := fn(w, r)
		if err != nil {
			RespondWithError(w, code, err.Error())
		}
	}
}
