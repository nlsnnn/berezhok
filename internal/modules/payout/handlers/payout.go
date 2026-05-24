package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
	payoutErrors "github.com/nlsnnn/berezhok/internal/modules/payout/errors"
	"github.com/nlsnnn/berezhok/internal/modules/payout/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type QueryServiceI interface {
	ListByPartner(ctx context.Context, partnerID uuid.UUID, limit, offset int32) ([]domain.Payout, int64, error)
	GetWithOrders(ctx context.Context, id uuid.UUID) (domain.Payout, []domain.PayoutOrder, error)
}

type PayoutHandler struct {
	svc QueryServiceI
	log *slog.Logger
}

func NewPayoutHandler(svc QueryServiceI, log *slog.Logger) *PayoutHandler {
	return &PayoutHandler{svc: svc, log: log}
}

func (h *PayoutHandler) List(w http.ResponseWriter, r *http.Request) {
	actor, err := contextx.PartnerActor(r)
	if err != nil {
		response.Unauthorized(w, "unauthorized")

		return
	}

	limit := int32(20)
	offset := int32(0)

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = int32(v)
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = int32(v)
		}
	}

	payouts, total, err := h.svc.ListByPartner(r.Context(), actor.PartnerID, limit, offset)
	if err != nil {
		h.log.Error("list payouts", slog.String("err", err.Error()))
		response.InternalError(w, err)

		return
	}

	items := make([]dto.PayoutResponse, len(payouts))
	for i, p := range payouts {
		items[i] = dto.ToPayoutResponse(p)
	}

	response.Success(w, dto.PayoutsListResponse{
		Payouts: items,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset)+int64(len(items)) < total,
	})
}

func (h *PayoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	_, err := contextx.PartnerActor(r)
	if err != nil {
		response.Unauthorized(w, "unauthorized")

		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "invalid payout id")

		return
	}

	payout, orders, err := h.svc.GetWithOrders(r.Context(), id)
	if err != nil {
		if errors.Is(err, payoutErrors.ErrPayoutNotFound) {
			response.NotFound(w, "payout not found")

			return
		}

		h.log.Error("get payout", slog.String("err", err.Error()))
		response.InternalError(w, err)

		return
	}

	response.Success(w, dto.ToPayoutDetailResponse(payout, orders))
}
