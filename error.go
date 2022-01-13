package mux

import (
	"net/http"
)

type ErrorHandler struct {
	err        error
	statusCode int
}

func (e ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, e.err.Error(), e.statusCode)
}
