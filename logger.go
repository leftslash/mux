package mux

import (
	"log"
	"net/http"
	"strings"
)

type LogOptions struct{}

func Logger(options *LogOptions) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w2 := &responseWriter{statusCode: 200, ResponseWriter: w}
				next.ServeHTTP(w2, r)
				statusCode := w2.statusCode
				log.Printf("%s %s %s %s %d %d %q\n",
					strings.Split(r.Host, ":")[0],
					r.Method,
					r.URL.Path,
					r.Proto,
					statusCode,
					w2.contentLength,
					r.Header.Get("User-Agent"),
				)
			})
	}
}

type responseWriter struct {
	statusCode    int
	contentLength int
	http.ResponseWriter
}

func (w *responseWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.contentLength += n
	return
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
