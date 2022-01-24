package mux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Add("Content-Type", "text/plain")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func URLParam(r *http.Request, key string) (value string) {
	if r == nil {
		return
	}
	ctx, err := getContext(r)
	if err != nil {
		Errorf(err, 0x1a67, "no context").Log()
		return
	}
	value, _ = ctx.parms[key]
	return
}

func GetSession(r *http.Request) (s *Session, e Error) {
	var ctx *muxContext
	ctx, err := getContext(r)
	if err != nil {
		e = Errorf(err, 0x69d9, "no session")
		e.Log()
		return
	}
	s = ctx.session
	if s == nil {
		e = Errorf(fmt.Errorf("context does not contain session"), 0, "no session")
		e.Log()
		return
	}
	return
}

func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) (err error) {
	err = json.NewDecoder(r.Body).Decode(v)
	return
}

func WriteJSON(w http.ResponseWriter, v interface{}) (err error) {
	var b []byte
	b, err = json.Marshal(v)
	if err != nil {
		return
	}
	w.Write(b)
	return
}
