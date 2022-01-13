package mux

import (
	"fmt"
	"sort"
	"testing"
)

var (
	leadingSlashErr = fmt.Errorf("leading slash required")
	whitespaceErr   = fmt.Errorf("pattern must not contain whitespace")
	emptyNestedErr  = fmt.Errorf("parameters must not be empty or nested")
	slashesErr      = fmt.Errorf("parameters must not contain slashes")
	wildcardErr     = fmt.Errorf("wildcard must follow slash and be last character")
)

func TestValidate(t *testing.T) {
	tests := []struct {
		pattern string
		want    error
	}{
		{"a", leadingSlashErr},
		{"/", nil},
		{"/*", nil},
		{"/a*", wildcardErr},
		{"/a/*", nil},
		{"/*a", wildcardErr},
		{"/*/", wildcardErr},
		{"/a*/", wildcardErr},
		{"/a**/", wildcardErr},
		{"/ a", whitespaceErr},
		{"/ a/", whitespaceErr},
		{"/a ", whitespaceErr},
		{"/a /", whitespaceErr},
		{"/{ }", whitespaceErr},
		{"/{id }", whitespaceErr},
		{"/{ id}", whitespaceErr},
		{"/{i d}", whitespaceErr},
		{"{id}", leadingSlashErr},
		{"/{id}", nil},
		{"/{id}*", wildcardErr},
		{"/{id}/", nil},
		{"/{id}/*", nil},
		{"/a{id}", nil},
		{"/{id}a", nil},
		{"/a{id}a", nil},
		{"/{}", emptyNestedErr},
		{"/{{}}", emptyNestedErr},
		{"/{{id}}", emptyNestedErr},
		{"/{a{id}", emptyNestedErr},
		{"/{id}a}", emptyNestedErr},
		{"/{a{id}b}", emptyNestedErr},
		{"/{/}", slashesErr},
		{"/{/id}", slashesErr},
		{"/{id/}", slashesErr},
		{"/{i/d}", slashesErr},
		{"//", nil},
		{"//a", nil},
		{"/a//", nil},
		{"//a//", nil},
	}
	for i, test := range tests {
		name := fmt.Sprintf("%d", i)
		t.Run(name, func(t *testing.T) {
			got := validate(test.pattern)
			if got == nil {
				if test.want == nil {
					return
				}
				t.Errorf("%s: got: %v; want: %v", test.pattern, got, test.want)
				return
			}
			if test.want == nil {
				t.Errorf("%s: got: %v; want: %v", test.pattern, got, test.want)
				return
			}
			if got.Error() != test.want.Error() {
				t.Errorf("%s: got: %v; want: %v", test.pattern, got, test.want)
				return
			}
		})
	}
}

func TestParse(t *testing.T) {
	type want struct {
		regex  string
		length int
	}
	tests := []struct {
		pattern string
		want    want
	}{
		{"/", want{"/", 0}},
		{"//", want{"/", 0}},
		{"///", want{"/", 0}},
		{"//a", want{"/a", 1}},
		{"//a//", want{"/a/", 1}},
		{"/*", want{"/*", 1}},
		{"/a/b", want{"/a/b", 2}},
		{"/a/b/", want{"/a/b/", 2}},
	}
	for i, test := range tests {
		name := fmt.Sprintf("%d", i)
		t.Run(name, func(t *testing.T) {
			got_regex, got_length := normalize(test.pattern)
			if got_length != test.want.length || got_regex != test.want.regex {
				t.Errorf("%s: got: %s, %d; want: %s, %d", test.pattern,
					got_regex, got_length, test.want.regex, test.want.length)
			}
		})
	}
}

func TestCompile(t *testing.T) {
	tests := []struct {
		pattern string
		url     string
		want    bool
	}{
		{"/", "/", true},
		{"/", "//", true},
		{"/", "/a", false},
		{"/a", "/a", true},
		{"/a", "/a/", false},
		{"/a/", "/a", true},
		{"/a/", "/a/", true},
		{"/a/*", "/a", false},
		{"/a/*", "/a/", true},
		{"/*", "/a", true},
		{"/*", "/a/", true},
		{"/*", "/a/b", true},
		{"/*", "/a/b/", true},
		{"/{id}", "/42", true},
		{"/{id}", "/42/", false},
		{"/{id}/", "/42", true},
		{"/{id}/", "/42/", true},
		{"/{id}/*", "/42", false},
		{"/{id}/*", "/42/", true},
		{"/{id}/*", "/42/43", true},
		{"/a{id}", "/a42", true},
		{"/a{id}", "/a42/", false},
		{"/{id}a", "/42a", true},
		{"/{id}a", "/42a/", false},
		{"/a{id}a", "/a42a", true},
		{"/a{id}a", "/a42a/", false},
	}
	for i, test := range tests {
		name := fmt.Sprintf("%d", i)
		t.Run(name, func(t *testing.T) {
			regex := compile(test.pattern)
			got := regex.MatchString(test.url)
			if got != test.want {
				t.Errorf("%s: url: %s; got: %t; want: %t", test.pattern, test.url, got, test.want)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	type want struct {
		ok    bool
		parms [][]string
	}
	tests := []struct {
		pattern string
		url     string
		want    want
	}{
		{"/u", "/x", want{false, [][]string{}}},
		{"/u", "/u", want{true, [][]string{}}},
		{"/u/{u}/g/{u}", "/u/joe/g/admin", want{true, [][]string{{"u", "admin"}}}},
		{"/u/{u}/g/{g}", "/u/joe/g/admin", want{true, [][]string{{"g", "admin"}, {"u", "joe"}}}},
		{"/u/{u}-{g}", "/u/joe-admin", want{true, [][]string{{"g", "admin"}, {"u", "joe"}}}},
		{"/u/{g}/g/{u}", "/u/joe/g/admin", want{true, [][]string{{"g", "joe"}, {"u", "admin"}}}},
	}
	for i, test := range tests {
		name := fmt.Sprintf("%d", i)
		t.Run(name, func(t *testing.T) {
			got_ok, got_parms := match(compile(test.pattern), test.url)
			if got_ok != test.want.ok {
				t.Errorf("%s: url: %s, got: %t; want: %t", test.pattern, test.url, got_ok, test.want.ok)
				return
			}
			got := [][]string{}
			for k, v := range got_parms {
				got = append(got, []string{k, v})
			}
			sort.Slice(got, func(i, j int) bool { return got[i][0] < got[j][0] })
			if len(got) != len(test.want.parms) {
				t.Errorf("%s: url: %s, got: %v; want: %v",
					test.pattern, test.url, got, test.want.parms)
				return
			}
			for i := range got {
				if len(got[i]) != len(test.want.parms[i]) {
					t.Errorf("%s: url: %s, got: %v %d; want: %v, %d",
						test.pattern, test.url, got, len(got[i]), test.want.parms, len(test.want.parms[i]))
					return
				}
				if got[i][0] != test.want.parms[i][0] {
					t.Errorf("%s: url: %s, got: %v; want: %v",
						test.pattern, test.url, got[i][0], test.want.parms[i][0])
					return
				}
				if got[i][1] != test.want.parms[i][1] {
					t.Errorf("%s: url: %s, got: %v; want: %v",
						test.pattern, test.url, got[i][1], test.want.parms[i][1])
					return
				}
				continue
			}
		})
	}
}
