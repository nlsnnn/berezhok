package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/nlsnnn/berezhok/internal/modules/auth"
	"github.com/nlsnnn/berezhok/internal/shared/logger/sl"
	"github.com/nlsnnn/berezhok/internal/shared/response"
	"github.com/nlsnnn/berezhok/internal/shared/validator"
)

type partAuth interface {
	Authenticate(ctx context.Context, email, password string) (*auth.TokenClaims, error)
}

type authHandler struct {
	partnerAuthenticator partAuth
	validator            *validator.Validator
	log                  *slog.Logger
}

func NewAuthHandler(
	validator *validator.Validator,
	log *slog.Logger,
	partnerAuthenticator partAuth,
) *authHandler {
	return &authHandler{
		partnerAuthenticator: partnerAuthenticator,
		validator:            validator,
		log:                  log,
	}
}

func (h *authHandler) PartnerLogin(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.PartnerLogin"
	log := h.log.With(slog.String("operation", op))
	log.Info("partner login attempt")

	var req LoginEmailPasswordRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	claims, err := h.partnerAuthenticator.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			log.Warn("invalid credentials", sl.Err(err))
			response.Unauthorized(w, "invalid credentials")
			return
		}

		log.Error("authentication failed", sl.Err(err))
		response.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	response.Success(w, LoginResponse{UserID: claims.UserID.String(), Token: claims.Access})
}
