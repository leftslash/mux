package main

import (
	"fmt"
	"net/http"

	"github.com/leftslash/mux"
)

func Extra(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			fmt.Fprintf(w, "Extra after\n")
		})
}

type App struct{ db string }

func (a App) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, %s!\n", a.db)
}

func (a App) Goodbye(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goodbye, %s.\n", a.db)
}

// Closure approach for injecting db connection into handler
func QuePasa(db string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Que Pasa, %s?\n", db)
	}
}

func User(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "user: %s\n", mux.URLParam(r, "user"))
}

func UserGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "user: %s, group: %s\n", mux.URLParam(r, "user"), mux.URLParam(r, "group"))
}

func main() {

	app := App{db: "World"}

	router := mux.NewRouter()
	router.Use(mux.Logger(nil))

	router.HandleFunc(http.MethodGet, "/h", app.Hello)
	router.HandleFunc(http.MethodGet, "/q", QuePasa("Amigo"))
	router.Handle(http.MethodGet, "/g", Extra(http.HandlerFunc(app.Goodbye)))

	router.HandleFunc(http.MethodGet, "/users/{user}/", User)
	router.HandleFunc(http.MethodGet, "/users/{user}/group/{group}", UserGroup)
	router.HandleFunc(http.MethodGet, "/user/{user}/*", User)

	router.Run("127.0.0.1:8080")
}
