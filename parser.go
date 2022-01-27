// parser takes user inputted url patterns and in multiple steps
// validates them, parses them into their components and then
// compiles them into a regular expression that can be later used
// when matching the provided pattern against an actual URL.
// This process is very pedantic and perhaps could be written
// in a more concise manner, but given the complexity of regular
// expressions, the vagaries of how URL's can be represented and
// the number of times I tried to write this in a condensed manner
// lead to this more drawn-out and explicit approach.  Mea culpa.
package mux

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	spacesParmRegex       = regexp.MustCompile(" |\t|\r|\f|\n")
	emptyNestedParmRegex  = regexp.MustCompile("{[^}]*{|}[^{]*}|{}")
	slashParmRegex        = regexp.MustCompile("{[^}]*/")
	badWildcardParmRegex  = regexp.MustCompile("\\*.|[^/]\\*")
	dirRegex              = regexp.MustCompile("/+")
	trailingSlashRegex    = regexp.MustCompile("/$")
	nonTrailingSlashRegex = regexp.MustCompile("/(.)")
	parmRegex             = regexp.MustCompile("{([^}/]+)}")
)

// validate tests the pattern against a number of non-valid
// regular expressions to ensure that any errors are flagged.
// a nil error indicates that no errors are found and pattern is valid
func validate(pattern string) (err error) {

	if pattern == "" || pattern[0] != '/' {
		err = fmt.Errorf("leading slash required")
		return
	}

	if spacesParmRegex.MatchString(pattern) {
		err = fmt.Errorf("pattern must not contain whitespace")
		return
	}

	if emptyNestedParmRegex.MatchString(pattern) {
		err = fmt.Errorf("parameters must not be empty or nested")
		return
	}

	if slashParmRegex.MatchString(pattern) {
		err = fmt.Errorf("parameters must not contain slashes")
		return
	}

	if badWildcardParmRegex.MatchString(pattern) {
		err = fmt.Errorf("wildcard must follow slash and be last character")
		return
	}

	return
}

// normalize examines the pattern string and returns a simple
// regular expression string and an integer indicating the
// length of directories the pattern represents.  The regular
// expression string that is returned will be further refined
// by the compile() function.  Note that this function
// assumes pattern has been validated by validate()
func normalize(pattern string) (regexStr string, length int) {

	// split pattern into slash separated components
	// ignoring (via regex) duplicate slashes
	dirs := dirRegex.Split(pattern, -1)

	// recombine dirs into regexstr which
	// removes any duplicate slashes from input
	regexStr += strings.Join(dirs, "/")

	// last character in pattern
	last := pattern[len(pattern)-1]

	// calculate length accounting for the leading slash
	length = len(dirs) - 1

	// accounting for trailing slash
	if last == '/' {
		length -= 1
	}

	return
}

// compile further refines the regular expression string
// returned by normalize() by expanding slashes to patterns
// that can handle duplicate slashes, by adding leading
// and trailing matching patterns, by transforming the
// wildcard character into a pattern match and by
// transforming all parameters into matching patterns.
// It then compiles the string into a proper regular expression
func compile(regexStr string) *regexp.Regexp {

	// replace wildcard with regex match pattern
	regexStr = strings.ReplaceAll(regexStr, "*", ".*")

	// replace non-trailing slashes with dup friendly pattern
	regexStr = nonTrailingSlashRegex.ReplaceAllString(regexStr, "/+$1")

	// replace trailing slash with dup friendly pattern
	// regexStr = trailingSlashRegex.ReplaceAllString(regexStr, "/*")
	regexStr = trailingSlashRegex.ReplaceAllString(regexStr, "/+")

	// add beginning and ending match characters
	regexStr = "^" + regexStr + "$"

	// replace parameters with match patterns
	regexStr = parmRegex.ReplaceAllString(regexStr, "(?P<$1>[^/]+)")

	// compile into regexp
	return regexp.MustCompile(regexStr)

}

// match tests whether a regular expression matches a specific url.
// If it does, it returns a map of key/value pairs that
// represent the parameter names from the pattern with the
// actual values specified in the url.
func match(r *regexp.Regexp, url string) (ok bool, parms map[string]string) {

	// check whether the regex matches the url
	matches := r.FindAllStringSubmatch(url, -1)
	if ok = matches != nil; !ok {
		return
	}

	// create the parms map
	parms = map[string]string{}
	// retrieve the names from the regex
	names := r.SubexpNames()
	// walk through all the names
	for i := range names {
		// skip element 0 which is not a name or match
		if i == 0 {
			continue
		}
		// assign name/values from regex and url
		// note that matches is two dimensional, only
		// the first element is relevant to us
		parms[names[i]] = matches[0][i]
	}
	return
}
