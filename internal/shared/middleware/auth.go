package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/auth"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
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
			if claims.Role != "" {
				ctx = context.WithValue(ctx, contextx.UserRoleKey, claims.Role)
			}

			if claims.UserType == "partner" {
				ctx = context.WithValue(ctx, contextx.EmployeeIDKey, claims.UserID)
				actor := authz.PartnerActor{
					EmployeeID: claims.UserID,
					Role:       authz.Role(claims.Role),
				}
				if claims.LocationID != nil {
					ctx = context.WithValue(ctx, contextx.LocationIDKey, *claims.LocationID)
					actor.LocationID = claims.LocationID
				}
				if claims.PartnerID != nil {
					ctx = context.WithValue(ctx, contextx.PartnerIDKey, *claims.PartnerID)
					actor.PartnerID = *claims.PartnerID
				} else if pid, ok := claims.UserData.(uuid.UUID); ok {
					ctx = context.WithValue(ctx, contextx.PartnerIDKey, pid)
					actor.PartnerID = pid
				} else if pidStr, ok := claims.UserData.(string); ok {
					if pid, err := uuid.Parse(pidStr); err == nil {
						ctx = context.WithValue(ctx, contextx.PartnerIDKey, pid)
						actor.PartnerID = pid
					}
				}
				ctx = context.WithValue(ctx, contextx.PartnerActorKey, actor)
			}

			if claims.UserType == "customer" {
				ctx = context.WithValue(ctx, contextx.CustomerIDKey, claims.UserID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (a *authMiddleware) RequirePermission(permission authz.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userType, err := contextx.UserType(r)
			if err != nil || userType != "partner" {
				response.Forbidden(w, "access denied for this user")
				return
			}

			actor, err := contextx.PartnerActor(r)
			if err != nil {
				response.Forbidden(w, "access denied for this user")
				return
			}

			if !actor.Can(permission) {
				response.Forbidden(w, "access denied for this action")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
