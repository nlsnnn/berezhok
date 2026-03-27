package middleware

import (
	"testing"
)

// Unit tests for the webhookMiddleware

// Tests that IsIPAllowed correctly allows single IPs and CIDR ranges
func TestIsIPAllowedSupportsSingleIPAndCIDR(t *testing.T) {
	t.Parallel()

	allowed := []string{"127.0.0.1/32", "77.75.156.11"}

	if !IsIPAllowed("127.0.0.1", allowed) {
		t.Fatal("expected CIDR IP to be allowed")
	}

	if !IsIPAllowed("77.75.156.11", allowed) {
		t.Fatal("expected single IP to be allowed")
	}

	if IsIPAllowed("10.10.10.10", allowed) {
		t.Fatal("expected unrelated IP to be blocked")
	}
}

// Tests that the middleware blocks requests from disallowed IPs
// func TestIPFilterMiddlewareIgnoresSpoofedXRealIP(t *testing.T) {
// 	t.Parallel()

// 	wm := NewWebhookMiddleware([]string{"127.0.0.1/32"}, slog.New(slog.NewTextHandler(io.Discard, nil)))

// 	h := wm.IPFilterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	}))

// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/yookassa", nil)
// 	req.RemoteAddr = "127.0.0.1:12345"
// 	req.Header.Set("X-Real-IP", "8.8.8.8")

// 	rr := httptest.NewRecorder()
// 	h.ServeHTTP(rr, req)

// 	if rr.Code != http.StatusForbidden {
// 		t.Fatalf("expected status 403, got %d", rr.Code)
// 	}
// }
