package main

import (
	"net/http"

	"github.com/leftslash/mux"
)

func main() {
	addr := "127.0.0.1:8080"
	r := mux.NewRouter()
	r.Use(mux.Logger(nil))
	r.Handle("GET", "/*", http.FileServer(http.Dir(".")))
	r.Run(addr)
}
