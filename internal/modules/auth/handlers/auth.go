package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	adminDomain "github.com/nlsnnn/berezhok/internal/modules/admin/domain"
	"github.com/nlsnnn/berezhok/internal/modules/auth"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type partAuth interface {
	Authenticate(ctx context.Context, email, password string) (*auth.TokenClaims, error)
}

type customerAuth interface {
	Authenticate(ctx context.Context, phone, code string) (*auth.TokenClaims, error)
	SendCode(ctx context.Context, phone string) error
}

type adminAuth interface {
	Authenticate(ctx context.Context, email, password string) (*auth.TokenClaims, error)
}

type authHandler struct {
	partnerAuthenticator  partAuth
	customerAuthenticator customerAuth
	adminAuthenticator    adminAuth
	validator             *validator.Validator
	log                   *slog.Logger
}

func NewAuthHandler(
	validator *validator.Validator,
	log *slog.Logger,
	partnerAuthenticator partAuth,
	customerAuthenticator customerAuth,
	adminAuthenticator adminAuth,
) *authHandler {
	return &authHandler{
		partnerAuthenticator:  partnerAuthenticator,
		customerAuthenticator: customerAuthenticator,
		adminAuthenticator:    adminAuthenticator,
		validator:             validator,
		log:                   log,
	}
}

func (h *authHandler) PartnerLogin(w http.ResponseWriter, r *http.Request) {
	const op = "auth.handlers.PartnerLogin"
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

	response.Success(w, LoginPartnerResponse{
		UserID:     claims.UserID.String(),
		EmployeeID: claims.UserID.String(),
		PartnerID:  claims.PartnerID.String(),
		LocationID: employeeLocationID(claims.UserData.(domain.Employee)),
		Role:       claims.Role,
		Token:      claims.Access,
		MustChange: claims.UserData.(domain.Employee).MustChangePassword,
	})
}

func (h *authHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	const op = "auth.handlers.AdminLogin"
	log := h.log.With(slog.String("operation", op))
	log.Info("admin login attempt")

	var req LoginEmailPasswordRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	claims, err := h.adminAuthenticator.Authenticate(r.Context(), req.Email, req.Password)
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

	admin := claims.UserData.(adminDomain.AdminUser)
	response.Success(w, LoginAdminResponse{
		UserID: claims.UserID.String(),
		Email:  admin.Email,
		Name:   admin.Name,
		Role:   claims.Role,
		Token:  claims.Access,
	})
}

func (h *authHandler) CustomerSendCode(w http.ResponseWriter, r *http.Request) {
	const op = "auth.handlers.CustomerSendCode"
	log := h.log.With(slog.String("operation", op))
	log.Info("customer send code attempt")

	var req SendCodeRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	if err := h.customerAuthenticator.SendCode(r.Context(), req.Phone); err != nil {
		log.Error("failed to send code", sl.Err(err))
		response.Error(w, "failed to send code", http.StatusInternalServerError)
		return
	}

	response.Success(w, map[string]string{"message": "code sent"})
}

func (h *authHandler) CustomerLogin(w http.ResponseWriter, r *http.Request) {
	const op = "auth.handlers.CustomerLogin"
	log := h.log.With(slog.String("operation", op))
	log.Info("customer login attempt")

	var req LoginPhoneRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	claims, err := h.customerAuthenticator.Authenticate(r.Context(), req.Phone, req.Code)
	if err != nil {
		log.Error("authentication failed", sl.Err(err))
		response.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	response.Success(w, LoginResponse{
		UserID: claims.UserID.String(),
		Token:  claims.Access,
	})
}

func employeeLocationID(employee domain.Employee) *string {
	if employee.LocationID == "" {
		return nil
	}

	locationID := employee.LocationID
	return &locationID
}
