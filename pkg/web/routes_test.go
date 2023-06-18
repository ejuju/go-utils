package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	t.Run("can route requests to the desired handler", func(t *testing.T) {
		// Handler for a certain path and method
		handlerNumCalls := 0
		handler := func(w http.ResponseWriter, r *http.Request) { handlerNumCalls++ }
		handlerPath := "/abc1"
		handlerMethod := http.MethodDelete

		// Handler that is used only for a certain path prefix
		prefixedHandlerNumCalls := 0
		prefixedHandlerPath := "/with-prefix/"
		prefixedHandler := func(w http.ResponseWriter, r *http.Request) { prefixedHandlerNumCalls++ }

		// 404 Handler (should be the last route)
		matchAllHandlerNumCalls := 0
		matchAllHandler := func(w http.ResponseWriter, r *http.Request) { matchAllHandlerNumCalls++ }

		routes := Routes{}
		routes.Handle(http.HandlerFunc(handler), MatchPath(handlerPath), MatchMethod(handlerMethod))
		routes.Handle(http.HandlerFunc(prefixedHandler), MatchPathPrefix(prefixedHandlerPath))
		routes.Handle(http.HandlerFunc(matchAllHandler))

		// Send request for normal handler and check number of calls
		req := httptest.NewRequest(handlerMethod, handlerPath, nil)
		resrec := httptest.NewRecorder()
		routes.ServeHTTP(resrec, req)
		if handlerNumCalls != 1 {
			t.Fatalf("normal handler: want 1 call but got %d", handlerNumCalls)
		}

		// Send request for prefixed handler and check number of calls
		req = httptest.NewRequest(http.MethodGet, prefixedHandlerPath, nil)
		routes.ServeHTTP(resrec, req)
		if prefixedHandlerNumCalls != 1 {
			t.Fatalf("prefixed handler: want 1 call but got %d", prefixedHandlerNumCalls)
		}

		// Send unmatched request and check number of calls for 404 handler
		req = httptest.NewRequest(http.MethodGet, "/unhandled-path", nil)
		routes.ServeHTTP(resrec, req)
		if matchAllHandlerNumCalls != 1 {
			t.Fatalf("404 handler: want 1 call but got %d", matchAllHandlerNumCalls)
		}
	})
}

func TestRequestMatchers(t *testing.T) {
	t.Run("can match sub-domain", func(t *testing.T) {
		tests := []struct {
			subdomain  string
			requestURL string
			wantMatch  bool
		}{
			{
				subdomain:  "www",
				requestURL: "www.example.com",
				wantMatch:  true,
			},
			{
				subdomain:  "www",
				requestURL: "admin.example.com",
				wantMatch:  false,
			},
			{
				subdomain:  "example",
				requestURL: "example.com",
				wantMatch:  false, // should not match, SLDs are not treated as subdomains
			},
			{
				subdomain:  "example",
				requestURL: "example.co.uk",
				wantMatch:  true, // should match, treat ccSLD just as normal SLDs
			},
		}

		for _, test := range tests {
			req := httptest.NewRequest(http.MethodGet, "http://"+test.requestURL, nil)
			gotMatch := MatchSubdomain(test.subdomain)(req)
			if gotMatch != test.wantMatch {
				t.Fatalf("subdomain %q and URL %q: want %v but got %v", test.subdomain, test.requestURL, test.wantMatch, gotMatch)
			}
		}
	})

	t.Run("can match a domain name", func(t *testing.T) {
		tests := []struct {
			domain     string
			requestURL string
			wantMatch  bool
		}{
			{
				domain:     "example.com",
				requestURL: "example.com",
				wantMatch:  true,
			},
			{
				domain:     "example.com",
				requestURL: "www.example.com",
				wantMatch:  true,
			},
			{
				domain:     "example.com",
				requestURL: "subdomain1.subdomain2.example.com",
				wantMatch:  true,
			},
			{
				domain:     "example.com",
				requestURL: "another.com",
				wantMatch:  false,
			},
			{
				domain:     "localhost",
				requestURL: "localhost",
				wantMatch:  true,
			},
		}

		for _, test := range tests {
			req := httptest.NewRequest(http.MethodGet, "http://"+test.requestURL, nil)
			gotMatch := MatchDomainName(test.domain)(req)
			if gotMatch != test.wantMatch {
				t.Fatalf("domain %q and URL %q: want %v but got %v", test.domain, test.requestURL, test.wantMatch, gotMatch)
			}
		}
	})
}
