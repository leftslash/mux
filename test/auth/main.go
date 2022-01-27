package main

import (
	"fmt"
	"net/http"

	"github.com/leftslash/mux"
)

func isValid(username, password string) (ok bool) {
	ok = username == "todd" && password == "hi"
	return
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mux.GetSession(r)
	if err != nil {
		fmt.Fprintln(w, "ok")
	} else {
		fmt.Fprintln(w, session.Username)
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(mux.NewAuth("/login", "/logout", isValid).Handle)
	r.HandleFunc(http.MethodGet, "/", okHandler)
	r.HandleFunc(http.MethodPost, "/", okHandler)
	r.Run("127.0.0.1:8080")
}
