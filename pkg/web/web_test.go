package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPAddressFromRequest(t *testing.T) {
	t.Run("can get IP address from HTTP request", func(t *testing.T) {
		xForwardedFor := "127.0.0.2"
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", xForwardedFor)

		// Check X-Forwarded-For header
		ipAddr := IPAddressFromRequest(req, true)
		if ipAddr.String() != xForwardedFor {
			t.Fatalf("want %s but got %s", xForwardedFor, ipAddr)
		}

		// Use standard http.Request.RemoteAddr
		remoteAddr := "127.0.0.1"
		req.RemoteAddr = remoteAddr + ":58000"
		ipAddr = IPAddressFromRequest(req, false)
		if ipAddr.String() != remoteAddr {
			t.Fatalf("want %s but got %s", xForwardedFor, ipAddr)
		}
	})
}
