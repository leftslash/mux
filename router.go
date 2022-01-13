package mux

import (
	"log"
	"net/http"
	"sort"
)

type Router struct {
	routes     routelist
	middleware []Middleware
}

func NewRouter() *Router {
	return &Router{}
}

type Middleware func(http.Handler) http.Handler

func (r *Router) Use(m Middleware) {
	r.middleware = append(r.middleware, m)
}

func (r *Router) Handle(method, pattern string, handler http.Handler) {
	for _, middleware := range r.middleware {
		handler = middleware(handler)
	}
	r.routes.register(method, pattern, handler)
}

func (r *Router) HandleFunc(method, pattern string, h http.HandlerFunc) {
	r.Handle(method, pattern, http.HandlerFunc(h))
}

func (r *Router) Run(addr string) {
	sort.Sort(r.routes)
	log.Printf("listening on http://%s\n", addr)
	err := http.ListenAndServe(addr, r)
	log.Println(err)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler
	// find route and passed parms
	route, parms, err := r.routes.resolve(req.Method, req.URL.Path)
	if err != nil {
		// set error handler
		handler = http.HandlerFunc(http.NotFound)
		// queue up all the middleware
		for _, middleware := range r.middleware {
			handler = middleware(handler)
		}
		handler.ServeHTTP(w, req)
		return
	}
	// set url parms in context
	req = setContext(req, &muxContext{parms: parms})
	// execute handler
	route.handler.ServeHTTP(w, req)
}
