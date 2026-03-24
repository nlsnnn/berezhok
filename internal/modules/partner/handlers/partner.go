package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type partnerHandler struct {
	partService partnerSvc
	log         *slog.Logger
	validator   *validator.Validator
}

func NewPartnerHandler(partService partnerSvc, log *slog.Logger) partnerHandler {
	return partnerHandler{
		partService: partService,
		log:         log,
		validator:   validator.New(),
	}
}

func (h *partnerHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.log.Error("user_id not found in context")
		response.InternalError(w, nil)
		return
	}

	var req dto.ChangePasswordRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		h.log.Error("validation failed", slog.Any("errors", errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	err := h.partService.ChangePassword(r.Context(), service.ChangePasswordInput{
		UserID:          userID.String(),
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		if errors.Is(err, partnerErrors.ErrInvalidCredentials) {
			h.log.Warn("invalid current password", sl.Err(err))
			response.BadRequest(w, "invalid current password")
			return
		}
		if errors.Is(err, partnerErrors.ErrPasswordUnchanged) {
			h.log.Warn("password unchanged", sl.Err(err))
			response.BadRequest(w, "password must be different from the current one")
			return
		}

		h.log.Error("failed to change password", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, map[string]string{"message": "password changed successfully"})
}

func (h *partnerHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.log.Error("user_id not found in context")
		response.InternalError(w, nil)
		return
	}

	profile, err := h.partService.Profile(r.Context(), userID.String())
	if err != nil {
		h.log.Error("failed to get profile", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	res := dto.PartnerProfileResponse{
		Partner: dto.PartnerResponse{
			ID:             profile.Partner.ID,
			BrandName:      profile.Partner.BrandName,
			Status:         string(profile.Partner.Status),
			CommissionRate: profile.Partner.Commission.Rate,
			PromoUntil:     profile.Partner.Commission.ValidUntil,
			CreatedAt:      profile.Partner.CreatedAt,
		},
		Employee: dto.EmployeeResponse{
			ID:                 profile.Employee.ID,
			Email:              profile.Employee.Email,
			Name:               profile.Employee.Name,
			Role:               string(profile.Employee.Role),
			MustChangePassword: profile.Employee.MustChangePassword,
			CreatedAt:          profile.Employee.CreatedAt,
		},
	}

	if profile.Location != nil {
		res.Location = &dto.LocationResponse{
			ID:        string(profile.Location.ID),
			Name:      profile.Location.Name,
			Address:   profile.Location.Address,
			CreatedAt: profile.Location.CreatedAt,
		}
	}

	// Map all partner locations
	res.Locations = make([]dto.LocationResponse, len(profile.Locations))
	for i, loc := range profile.Locations {
		res.Locations[i] = dto.LocationResponse{
			ID:        loc.ID,
			Name:      loc.Name,
			Address:   loc.Address,
			CreatedAt: loc.CreatedAt,
		}
	}

	response.Success(w, res)
}
