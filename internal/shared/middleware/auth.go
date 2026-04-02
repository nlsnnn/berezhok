package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/auth"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type authMiddleware struct {
	tokenService tokenSvc
}

type tokenSvc interface {
	Validate(token string) (*auth.TokenClaims, error)
}

func NewAuthMiddleware(tokenService tokenSvc) *authMiddleware {
	return &authMiddleware{tokenService: tokenService}
}

func (a *authMiddleware) RequireAuth(allowedTypes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w, "invalid authorization header")
				return
			}

			tokenString := parts[1]

			claims, err := a.tokenService.Validate(tokenString)
			if err != nil {
				response.Unauthorized(w, "invalid token")
				return
			}

			if len(allowedTypes) > 0 {
				allowed := false
				for _, userType := range allowedTypes {
					if claims.UserType == string(userType) {
						allowed = true
						break
					}
				}
				if !allowed {
					response.Forbidden(w, "access denied for this user")
					return
				}
			}

			ctx := context.WithValue(r.Context(), contextx.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, contextx.UserTypeKey, claims.UserType)

			if claims.UserType == "partner" {
				if pid, ok := claims.UserData.(uuid.UUID); ok {
					ctx = context.WithValue(ctx, contextx.PartnerIDKey, pid)
				} else if pidStr, ok := claims.UserData.(string); ok {
					if pid, err := uuid.Parse(pidStr); err == nil {
						ctx = context.WithValue(ctx, contextx.PartnerIDKey, pid)
					}
				}
			}

			if claims.UserType == "customer" {
				ctx = context.WithValue(ctx, contextx.CustomerIDKey, claims.UserID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
