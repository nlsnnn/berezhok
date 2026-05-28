package middleware

import (
	"net/http"

	"github.com/nlsnnn/berezhok/internal/shared/response"
)

func InternalAuthMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" || r.Header.Get("X-Internal-Token") != token {
				response.Unauthorized(w, "invalid internal token")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
