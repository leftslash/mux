package mux

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type route struct {
	method  string
	pattern string
	handler http.Handler
	length  int
	regex   *regexp.Regexp
}

type routelist []route

// methods necessary for sorting
func (r routelist) Len() int           { return len(r) }
func (r routelist) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r routelist) Less(i, j int) bool { return r[i].length > r[j].length }

func (r *routelist) register(method, pattern string, handler http.Handler) {

	route := route{method: method, pattern: pattern, handler: handler}

	err := validate(pattern)
	if err != nil {
		log.Fatalf("error registering pattern: %s; %s", pattern, err)
	}

	var regexStr string
	regexStr, route.length = normalize(pattern)
	route.regex = compile(regexStr)

	*r = append(*r, route)
}

func (r *routelist) resolve(method, urlpath string) (route route, parms map[string]string, err error) {
	var ok bool
	for _, route = range *r {
		if method != route.method {
			continue
		}
		if ok, parms = match(route.regex, urlpath); ok {
			return
		}
	}
	err = fmt.Errorf("cannot resolve url %q", urlpath)
	return
}
