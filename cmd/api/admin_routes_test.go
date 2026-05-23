package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nlsnnn/berezhok/internal/shared/config"
)

func TestApplicationAdminRoutesAreProtected(t *testing.T) {
	t.Parallel()

	app := &application{
		cfg: &config.Config{},
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	handler := app.mount()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{name: "public application list removed", method: http.MethodGet, path: "/api/v1/applications", wantStatus: http.StatusMethodNotAllowed},
		{name: "admin application list requires auth", method: http.MethodGet, path: "/api/v1/admin/applications", wantStatus: http.StatusUnauthorized},
		{name: "admin application approve requires auth", method: http.MethodPost, path: "/api/v1/admin/applications/10000000-0000-0000-0000-000000000001/approve", wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}
		})
	}
}
