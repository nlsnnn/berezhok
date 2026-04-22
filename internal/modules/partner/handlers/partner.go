package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
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

func (h *partnerHandler) Stats(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.Stats"
	log := h.log.With(slog.String("op", op))

	userID, err := contextx.UserID(r)
	if err != nil {
		log.Error("user_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	filter, err := statsFilterFromRequest(r)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	stats, err := h.partService.Stats(r.Context(), userID.String(), filter)
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrInvalidStatsPeriod),
			errors.Is(err, partnerErrors.ErrInvalidStatsDateRange),
			errors.Is(err, partnerErrors.ErrInvalidStatsSort):
			response.BadRequest(w, err.Error())
		default:
			log.Error("failed to get stats", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Success(w, toPartnerStatsResponse(stats))
}

func (h *partnerHandler) AddLegalInfo(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.AddLegalInfo"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	var req dto.AddLegalInfoRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	err = h.partService.AddLegalInfo(r.Context(), req.ToInput(partnerID.String()))
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrPartnerStatusInvalid):
			log.Warn("partner status is invalid", sl.Err(err))
			response.BadRequest(w, err.Error())
		case errors.Is(err, partnerErrors.ErrInvalidINN),
			errors.Is(err, partnerErrors.ErrInvalidOGRN),
			errors.Is(err, partnerErrors.ErrInvalidKPP):
			log.Warn("invalid legal info", sl.Err(err))
			response.BadRequest(w, err.Error())
		default:
			log.Error("failed to save legal info", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Success(w, map[string]string{"message": "legal info saved successfully"})
}

func statsFilterFromRequest(r *http.Request) (domain.StatsFilter, error) {
	query := r.URL.Query()
	filter := domain.StatsFilter{
		Period:           query.Get("period"),
		LocationID:       query.Get("location_id"),
		Status:           query.Get("status"),
		TopLocationsSort: query.Get("top_locations_sort"),
		TopBoxesSort:     query.Get("top_boxes_sort"),
		OrdersSort:       query.Get("orders_sort"),
		Limit:            20,
	}

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return domain.StatsFilter{}, partnerErrors.ErrInvalidStatsDateRange
		}
		filter.Limit = parsed
	}

	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return domain.StatsFilter{}, partnerErrors.ErrInvalidStatsDateRange
		}
		filter.Offset = parsed
	}

	if value := query.Get("date_from"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			return domain.StatsFilter{}, partnerErrors.ErrInvalidStatsDateRange
		}
		filter.DateFrom = parsed
	}

	if value := query.Get("date_to"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			return domain.StatsFilter{}, partnerErrors.ErrInvalidStatsDateRange
		}
		filter.DateTo = parsed
	}

	return filter, nil
}

func toPartnerStatsResponse(stats domain.PartnerStats) dto.PartnerStatsResponse {
	timeline := make([]dto.PartnerStatsTimelinePointResponse, len(stats.Timeline))
	for i, item := range stats.Timeline {
		timeline[i] = dto.PartnerStatsTimelinePointResponse{
			Date:            item.Date,
			OrdersTotal:     item.OrdersTotal,
			OrdersCompleted: item.OrdersCompleted,
			GrossRevenue:    item.GrossRevenue,
			NetRevenue:      item.NetRevenue,
		}
	}

	statusBreakdown := make([]dto.PartnerStatsStatusBreakdownItemResponse, len(stats.StatusBreakdown))
	for i, item := range stats.StatusBreakdown {
		statusBreakdown[i] = dto.PartnerStatsStatusBreakdownItemResponse{
			Status: item.Status,
			Count:  item.Count,
			Share:  item.Share,
		}
	}

	topLocations := make([]dto.PartnerStatsLocationResponse, len(stats.TopLocations))
	for i, item := range stats.TopLocations {
		topLocations[i] = dto.PartnerStatsLocationResponse{
			LocationID:      item.LocationID,
			Name:            item.Name,
			Address:         item.Address,
			OrdersTotal:     item.OrdersTotal,
			OrdersCompleted: item.OrdersCompleted,
			GrossRevenue:    item.GrossRevenue,
			NetRevenue:      item.NetRevenue,
			AvgRating:       item.AvgRating,
		}
	}

	topBoxes := make([]dto.PartnerStatsBoxResponse, len(stats.TopBoxes))
	for i, item := range stats.TopBoxes {
		topBoxes[i] = dto.PartnerStatsBoxResponse{
			BoxID:           item.BoxID,
			Name:            item.Name,
			ImageURL:        item.ImageURL,
			LocationName:    item.LocationName,
			OrdersTotal:     item.OrdersTotal,
			OrdersCompleted: item.OrdersCompleted,
			GrossRevenue:    item.GrossRevenue,
			NetRevenue:      item.NetRevenue,
		}
	}

	orders := make([]dto.PartnerStatsOrderResponse, len(stats.Orders))
	for i, item := range stats.Orders {
		orders[i] = dto.PartnerStatsOrderResponse{
			ID:              item.ID,
			Status:          item.Status,
			PickupCode:      item.PickupCode,
			Amount:          item.Amount,
			BoxName:         item.BoxName,
			BoxImageURL:     item.BoxImageURL,
			CustomerPhone:   item.CustomerPhone,
			CustomerName:    item.CustomerName,
			LocationID:      item.LocationID,
			LocationName:    item.LocationName,
			LocationAddress: item.LocationAddress,
			PickupTimeStart: item.PickupTimeStart,
			PickupTimeEnd:   item.PickupTimeEnd,
			CreatedAt:       item.CreatedAt,
			CanPickup:       item.CanPickup,
		}
	}

	return dto.PartnerStatsResponse{
		Summary: dto.PartnerStatsSummaryResponse{
			OrdersTotal:               stats.Summary.OrdersTotal,
			OrdersCompleted:           stats.Summary.OrdersCompleted,
			OrdersCancelled:           stats.Summary.OrdersCancelled,
			OrdersPendingConfirmation: stats.Summary.OrdersPendingConfirmation,
			GrossRevenue:              stats.Summary.GrossRevenue,
			NetRevenue:                stats.Summary.NetRevenue,
			AvgOrderValue:             stats.Summary.AvgOrderValue,
			AvgRating:                 stats.Summary.AvgRating,
			ReviewsCount:              stats.Summary.ReviewsCount,
		},
		Timeline:        timeline,
		StatusBreakdown: statusBreakdown,
		TopLocations:    topLocations,
		TopBoxes:        topBoxes,
		Orders:          orders,
		Meta: dto.PartnerStatsMetaResponse{
			Period:           stats.Meta.Period,
			DateFrom:         stats.Meta.DateFrom,
			DateTo:           stats.Meta.DateTo,
			LocationID:       stats.Meta.LocationID,
			Status:           stats.Meta.Status,
			TopLocationsSort: stats.Meta.TopLocationsSort,
			TopBoxesSort:     stats.Meta.TopBoxesSort,
			OrdersSort:       stats.Meta.OrdersSort,
			Pagination: dto.PartnerStatsPaginationResponse{
				Total:   stats.Meta.Pagination.Total,
				Limit:   stats.Meta.Pagination.Limit,
				Offset:  stats.Meta.Pagination.Offset,
				HasMore: stats.Meta.Pagination.HasMore,
			},
		},
	}
}
