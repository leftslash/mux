package mux

import (
	"net/http"
)

func MaxReader(n int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				r2 := *r
				r2.Body = http.MaxBytesReader(w, r.Body, n)
				next.ServeHTTP(w, &r2)
			})
	}
}
