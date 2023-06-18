package web

import (
	"net/http"
	"strings"
)

// Routes represents a list of routes.
type Routes []*Route

func (rhs Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rhs.RequestHandler(r).ServeHTTP(w, r)
}

func (rhs Routes) RequestHandler(r *http.Request) http.Handler {
	for _, rh := range rhs {
		if rh.Match(r) {
			return rh.handler
		}
	}
	panic("no match found")
}

func (rhs *Routes) Handle(h http.Handler, matchers ...RequestMatcher) {
	*rhs = append(*rhs, &Route{handler: h, matchers: matchers})
}

// Route represents a HTTP handler with request matchers.
type Route struct {
	handler  http.Handler
	matchers []RequestMatcher
}

// If all matchers yield true, this function returns true.
// If there are no matchers provided, true is also returned.
func (rh *Route) Match(r *http.Request) bool {
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

// Utility request method matchers defined for convenience and conciseness.
func MatchMethodGet(r *http.Request) bool     { return r.Method == http.MethodGet }
func MatchMethodHead(r *http.Request) bool    { return r.Method == http.MethodHead }
func MatchMethodPost(r *http.Request) bool    { return r.Method == http.MethodPost }
func MatchMethodPut(r *http.Request) bool     { return r.Method == http.MethodPut }
func MatchMethodPatch(r *http.Request) bool   { return r.Method == http.MethodPatch }
func MatchMethodDelete(r *http.Request) bool  { return r.Method == http.MethodDelete }
func MatchMethodConnect(r *http.Request) bool { return r.Method == http.MethodConnect }
func MatchMethodOptions(r *http.Request) bool { return r.Method == http.MethodOptions }
func MatchMethodTrace(r *http.Request) bool   { return r.Method == http.MethodTrace }

// Defined for readability purposes, to make it explicit that the handlers intents to catch all requests.
// Used for handling 404s.
// Equivalent to not using any RequestMatcher in routes.Handle, as this will match any request.
func CatchAll(_ *http.Request) bool { return true }

// Matches a certain host. Ignores port if present.
func MatchHostname(s string) RequestMatcher {
	return func(r *http.Request) bool { return r.URL.Hostname() == s }
}

// Matches a certain domain name (SLD + TLD).
func MatchDomainName(s string) RequestMatcher {
	return func(r *http.Request) bool {
		hostname := r.URL.Hostname()
		parts := strings.Split(hostname, ".")
		if len(parts) >= 3 {
			// if subdomain is included, remove subdomain(s) from parts
			hostname = strings.Join(parts[len(parts)-2:], ".")
		}
		return hostname == s
	}
}

// Matches a certain subdomain.
//
// NB: Treats ccSLD just as normal SLDs (cf. unit test)
func MatchSubdomain(s string) RequestMatcher {
	return func(r *http.Request) bool {
		if strings.Count(r.URL.Host, ".") <= 1 {
			return false
		}
		return strings.Split(r.URL.Host, ".")[0] == s
	}
}

// Matches a certain HTTP header key and value.
func MatchHeader(k, v string) RequestMatcher {
	return func(r *http.Request) bool { return r.Header.Get(k) == v }
}
