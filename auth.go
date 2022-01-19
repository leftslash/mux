package mux

import (
	"net/http"
)

const cookieName = "id"

type Authenticator func(username, password string) bool

type Auth struct {
	loginURL      string
	logoutURL     string
	authenticator Authenticator
	sessions      *sessions
}

func NewAuth(login, logout string, auth Authenticator) *Auth {
	return &Auth{
		loginURL:      login,
		logoutURL:     logout,
		authenticator: auth,
		sessions:      &sessions{list: map[string]*Session{}},
	}
}

func (a *Auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var session *Session
			cookie, err := r.Cookie(cookieName)
			cookieIsValid := false
			if err == nil {
				// cookie exists
				id := cookie.Value
				session, err = a.sessions.get(id)
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
				if !a.authenticator(username, password) {
					// not authorized
					http.Redirect(w, r, a.loginURL, http.StatusFound)
					return
				}
				// create a session
				session = &Session{Username: username}
				id := a.sessions.add(session)
				// create a cookie
				cookie = &http.Cookie{Name: cookieName, Value: id}
			}
			// add session to context
			ctx, err := getContext(r)
			if err == nil {
				ctx.session = session
			}
			// add secure attributes to cookie
			cookie.Path = "/"
			cookie.Secure = true
			cookie.HttpOnly = true
			// set cookie and serve content
			http.SetCookie(w, cookie)
			next.ServeHTTP(w, r)
		})
}

func (a *Auth) HandleFunc(h http.HandlerFunc) http.Handler {
	return a.Handle(http.HandlerFunc(h))
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		a.sessions.remove(cookie.Value)
	}
	http.Redirect(w, r, a.logoutURL, http.StatusFound)
}
