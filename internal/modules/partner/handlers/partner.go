package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
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
	const op = "partner.handler.ChangePassword"
	log := h.log.With(slog.String("op", op))

	userID, err := contextx.UserID(r)
	if err != nil {
		log.Error("user_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	var req dto.ChangePasswordRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", slog.Any("errors", errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	err = h.partService.ChangePassword(r.Context(), service.ChangePasswordInput{
		UserID:          userID.String(),
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrInvalidCredentials), errors.Is(err, partnerErrors.ErrPasswordUnchanged):
			log.Warn("password change failed", sl.Err(err))
			response.BadRequest(w, err.Error())
		case errors.Is(err, partnerErrors.ErrPasswordUnchanged):
			log.Warn("password unchanged", sl.Err(err))
			response.BadRequest(w, "password must be different from the current one")
		default:
			log.Error("failed to change password", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Success(w, map[string]string{"message": "password changed successfully"})
}

func (h *partnerHandler) Profile(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.Profile"
	log := h.log.With(slog.String("op", op))

	userID, err := contextx.UserID(r)
	if err != nil {
		log.Error("user_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	profile, err := h.partService.Profile(r.Context(), userID.String())
	if err != nil {
		log.Error("failed to get profile", sl.Err(err))
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

func (h *partnerHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.Dashboard"
	log := h.log.With(slog.String("op", op))

	userID, err := contextx.UserID(r)
	if err != nil {
		log.Error("user_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	dashboard, err := h.partService.Dashboard(r.Context(), userID.String())
	if err != nil {
		log.Error("failed to get dashboard", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	locations := make([]dto.DashboardLocationResponse, len(dashboard.Locations))
	for i, loc := range dashboard.Locations {
		locations[i] = dto.DashboardLocationResponse{
			ID:               loc.ID,
			Name:             loc.Name,
			Address:          loc.Address,
			Status:           string(loc.Status),
			ActiveBoxesCount: loc.ActiveBoxesCount,
		}
	}

	res := dto.PartnerDashboardResponse{
		Partner: dto.DashboardPartnerResponse{
			ID:             dashboard.Partner.ID,
			BrandName:      dashboard.Partner.BrandName,
			Status:         string(dashboard.Partner.Status),
			CommissionRate: dashboard.Partner.Commission.Rate,
			PromoUntil:     dashboard.Partner.Commission.ValidUntil,
		},
		Employee: dto.DashboardEmployeeResponse{
			ID:    dashboard.Employee.ID,
			Name:  dashboard.Employee.Name,
			Email: dashboard.Employee.Email,
			Role:  string(dashboard.Employee.Role),
		},
		Locations: locations,
		Today: dto.DashboardTodayResponse{
			PendingConfirmation: dashboard.Today.PendingConfirmation,
			Confirmed:           dashboard.Today.Confirmed,
			PickedUp:            dashboard.Today.PickedUp,
			Completed:           dashboard.Today.Completed,
		},
		Week: dto.DashboardWeekResponse{
			OrdersCompleted: dashboard.Week.OrdersCompleted,
			GrossRevenue:    dashboard.Week.GrossRevenue,
			NetRevenue:      dashboard.Week.NetRevenue,
			AvgRating:       dashboard.Week.AvgRating,
		},
		Finance: dto.DashboardFinanceResponse{
			BalancePending: dashboard.Finance.BalancePending,
			NextPayoutDate: dashboard.Finance.NextPayoutDate,
		},
	}

	response.Success(w, res)
}
