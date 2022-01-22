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
		Errorf(err, 0, "no context").Log()
		return
	}
	value, _ = ctx.parms[key]
	return
}

func GetSession(r *http.Request) (s *Session, e Error) {
	var ctx *muxContext
	ctx, err := getContext(r)
	if err != nil {
		e = Errorf(err, 0, "no session")
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

func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) (e Error) {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		e = Errorf(err, 0, "decoding JSON")
		http.Error(w, e.Error(), http.StatusBadRequest)
	}
	return
}

func WriteJSON(w http.ResponseWriter, v interface{}) (e Error) {
	var b []byte
	b, err := json.Marshal(v)
	if err != nil {
		e = Errorf(err, 0, "encoding JSON")
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
	return
}
