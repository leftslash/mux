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
	if val == nil {
		err = fmt.Errorf("no context for key %q", contextKey.String)
		return
	}
	var ok bool
	ctx, ok = val.(*muxContext)
	if !ok {
		err = fmt.Errorf("context type not *muxContext")
		return
	}
	return
}

func setContext(r *http.Request, c *muxContext) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contextKey, c))
}
