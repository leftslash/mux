package mux

import (
	"context"
	"fmt"
	"net/http"
)

var contextKey = &struct{ String string }{"context.mux.leftslash.com"}

type muxContext struct {
	parms   map[string]string
	session *Session
}

func getContext(r *http.Request) (ctx *muxContext, err error) {
	val := r.Context().Value(contextKey)
	var ok bool
	ctx, ok = val.(*muxContext)
	if !ok {
		err = fmt.Errorf("no context")
		return
	}
	return
}

func setContext(r *http.Request, c *muxContext) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contextKey, c))
}
