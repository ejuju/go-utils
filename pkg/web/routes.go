package web

import (
	"errors"
	"net/http"
	"strings"
)

// Routes represents a list of routes.
type Routes []*Route

func (rhs Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rhs.Match(r).ServeHTTP(w, r)
}

func (rhs Routes) Match(r *http.Request) http.Handler {
	for _, rh := range rhs {
		if rh.IsMatch(r) {
			return rh.handler
		}
	}
	panic(errors.New("no handler found"))
}

func (rhs *Routes) Handle(h http.Handler, matchers ...RequestMatcher) {
	*rhs = append(*rhs, &Route{handler: h, matchers: matchers})
}

// Route represents a HTTP handler with request matchers.
type Route struct {
	handler  http.Handler
	matchers []RequestMatcher
}

// All matchers must be true for this function to return true.
// If no matchers is provided, true is also returned.
func (rh *Route) IsMatch(r *http.Request) bool {
	for _, matcher := range rh.matchers {
		if !matcher(r) {
			return false
		}
	}
	return true
}

// RequestMatcher represents a function that reports whether a HTTP request matches a certain criteria.
//
// For example: MatchPath and MatchMethod allow you to match a HTTP request
// with a certain path and method request to a handler.
type RequestMatcher func(r *http.Request) bool

// Checks if the request URL path is the same as the provided one.
func MatchPath(path string) RequestMatcher {
	return func(r *http.Request) bool { return r.URL.Path == path }
}

// Checks if the request URL path starts with the provided prefix.
func MatchPathPrefix(prefix string) RequestMatcher {
	return func(r *http.Request) bool { return strings.HasPrefix(r.URL.Path, prefix) }
}

// Checks if the request method is one of the provided ones.
func MatchMethod(methods ...string) RequestMatcher {
	return func(r *http.Request) bool {
		for _, m := range methods {
			if r.Method == m {
				return true
			}
		}
		return false
	}
}

// Defined for readability purposes, to make it explicit that the handlers intents to catch all requests.
// Used for handling 404s.
// Equivalent to not using any RequestMatcher in routes.Handle
func CatchAll(_ *http.Request) bool { return true }

// Utility request method matchers defined for convenience and conciseness.
func MatchMethodGET(r *http.Request) bool  { return r.Method == http.MethodGet }
func MatchMethodPOST(r *http.Request) bool { return r.Method == http.MethodPost }
