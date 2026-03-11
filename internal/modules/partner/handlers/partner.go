package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/partner"
	"github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/logger/sl"
	"github.com/nlsnnn/berezhok/internal/shared/response"
	"github.com/nlsnnn/berezhok/internal/shared/validator"
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

	if err := h.partService.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, partner.ErrInvalidCredentials) {
			h.log.Warn("invalid current password", sl.Err(err))
			response.BadRequest(w, "invalid current password")
			return
		}
		if errors.Is(err, partner.ErrPasswordUnchanged) {
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

	profile, err := h.partService.Profile(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get profile", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	commissionRate, _ := profile.CommissionRate.Float64Value()

	res := dto.PartnerProfileResponse{
		Partner: dto.PartnerResponse{
			ID:             profile.PartnerID,
			LegalName:      profile.LegalName,
			BrandName:      profile.BrandName.String,
			Status:         profile.PartnerStatus,
			CommissionRate: commissionRate.Float64,
		},
		Employee: dto.EmployeeResponse{
			ID:    profile.EmployeeID,
			Email: profile.Email,
			Name:  profile.EmployeeName.String,
			Role:  profile.Role,
		},
	}

	if profile.PromoCommissionUntil.Valid {
		promoTime := profile.PromoCommissionUntil.Time
		res.Partner.PromoUntil = &promoTime
	}

	if profile.LocationID.Valid {
		res.Location = &dto.LocationResponse{
			ID:      profile.LocationID.Bytes,
			Name:    profile.LocationName.String,
			Address: profile.LocationAddress.String,
		}
	}

	response.Success(w, res)
}
