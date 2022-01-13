package mux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func URLParam(r *http.Request, key string) (value string) {
	if r == nil {
		return
	}
	ctx, err := getContext(r)
	if err != nil {
		return
	}
	value, _ = ctx.parms[key]
	return
}

func GetSession(r *http.Request) (s *Session, err error) {
	var ctx *muxContext
	ctx, err = getContext(r)
	if err != nil {
		return
	}
	s = ctx.session
	if s == nil {
		err = fmt.Errorf("no session")
		return
	}
	return
}

func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) (err error) {
	err = json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	return
}

func WriteJSON(w http.ResponseWriter, v interface{}) (err error) {
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	return
}
