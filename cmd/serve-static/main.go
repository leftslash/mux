package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/leftslash/mux"
)

func main() {

	addr := flag.String("addr", "127.0.0.1", "ip address on which to listen")
	port := flag.Int("port", 8080, "tcp port on which to listen")
	dir := flag.String("dir", ".", "directory of content to serve")
	flag.Parse()

	r := mux.NewRouter()
	r.Use(mux.Logger(nil))
	r.Handle("GET", "/*", http.FileServer(http.Dir(*dir)))
	r.Run(fmt.Sprintf("%s:%d", *addr, *port))
}
