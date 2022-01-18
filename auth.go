package mux

import (
	"log"
	"net/http"
)

const cookieName = "id"

type Validator func(username, password string) bool

type AuthOptions struct {
	Validator Validator
	FailURL   string
}

func Auth(options AuthOptions) Middleware {
	if options.FailURL == "" || options.Validator == nil {
		log.Fatal("invalid options for auth middleware")
	}
	sessions := &sessions{list: map[string]*Session{}}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				var session *Session
				cookie, err := r.Cookie(cookieName)
				cookieIsValid := false

				if err == nil {
					// cookie exists
					id := cookie.Value
					session, err = sessions.get(id)
					if err == nil {
						// found session for cookie
						cookieIsValid = true
					}
				}

				if !cookieIsValid {
					// cookie is missing or expired
					// try basic auth first
					username, password, ok := r.BasicAuth()
					if !ok {
						// fall back to form data
						username = r.PostFormValue("username")
						password = r.PostFormValue("password")
					}
					if !options.Validator(username, password) {
						// not authorized
						http.Redirect(w, r, options.FailURL, http.StatusFound)
						return
					}
					// create a session
					session = &Session{Username: username}
					id := sessions.add(session)
					// create a cookie
					cookie = &http.Cookie{
						Name:  cookieName,
						Value: id,
					}
				}

				// add session to context
				ctx, err := getContext(r)
				if err == nil {
					ctx.session = session
				}
				http.SetCookie(w, cookie)
				next.ServeHTTP(w, r)
			})
	}
}
