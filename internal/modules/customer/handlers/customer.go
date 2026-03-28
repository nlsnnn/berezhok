package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/customer/domain"
	"github.com/nlsnnn/berezhok/internal/modules/customer/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type customerHandler struct {
	service customerSvc
	log     *slog.Logger
	v       *validator.Validator
}

type customerSvc interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, name string) (domain.User, error)
}

func NewCustomerHandler(service customerSvc, log *slog.Logger, v *validator.Validator) *customerHandler {
	return &customerHandler{
		service: service,
		log:     log,
		v:       v,
	}
}

// GetProfile handles GET /customer/profile
func (h *customerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.log.Error("user_id not found in context")
		response.InternalError(w, nil)
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get profile", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	resp := dto.ProfileResponse{
		ID:        user.ID.String(),
		Phone:     user.Phone.Number,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response.Success(w, resp)
}

// UpdateProfile handles PATCH /customer/profile
func (h *customerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.log.Error("user_id not found in context")
		response.InternalError(w, nil)
		return
	}

	var req dto.UpdateProfileRequest
	if errs := h.v.DecodeAndValidate(r, &req); errs != nil {
		h.log.Error("validation failed", slog.Any("errors", errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	user, err := h.service.UpdateProfile(r.Context(), userID, req.Name)
	if err != nil {
		h.log.Error("failed to update profile", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	resp := dto.ProfileResponse{
		ID:        user.ID.String(),
		Phone:     user.Phone.Number,
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}

	response.Success(w, resp)
}
