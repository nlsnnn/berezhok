package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/modules/eco/domain"
	"github.com/nlsnnn/berezhok/internal/modules/eco/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type ecoSvc interface {
	GetForUser(ctx context.Context, userID uuid.UUID) (domain.EcoStats, error)
}

type ecoHandler struct {
	service ecoSvc
	log     *slog.Logger
}

func NewEcoHandler(service ecoSvc, log *slog.Logger) *ecoHandler {
	return &ecoHandler{service: service, log: log}
}

// GetEcoStats handles GET /customer/eco-stats.
func (h *ecoHandler) GetEcoStats(w http.ResponseWriter, r *http.Request) {
	userID, err := contextx.UserID(r)
	if err != nil {
		h.log.Error("user_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	stats, err := h.service.GetForUser(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get eco stats", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	resp := dto.EcoStatsResponse{
		BoxesPickedUp:   stats.BoxesPickedUp,
		TotalKg:         stats.TotalKg,
		CO2SavedKg:      stats.CO2SavedKg,
		SavingsRub:      stats.SavingsRub.Round(0).IntPart(),
		MealsEquivalent: stats.MealsEquivalent,
		Tier:            string(stats.Tier),
		TierProgress:    stats.TierProgress,
		KgToNextTier:    stats.KgToNextTier,
	}
	if stats.NextTier != nil {
		nt := string(*stats.NextTier)
		resp.NextTier = &nt
	}

	response.Success(w, resp)
}
