package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
	payoutErrors "github.com/nlsnnn/berezhok/internal/modules/payout/errors"
	"github.com/nlsnnn/berezhok/internal/modules/payout/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/modules/payout/service"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type DestinationServiceI interface {
	Get(ctx context.Context, partnerID uuid.UUID) (domain.PayoutDestination, error)
	Upsert(ctx context.Context, partnerID uuid.UUID, destType domain.DestinationType, phone, bankID, recipientName string) (domain.PayoutDestination, error)
	GetSBPBanks(ctx context.Context) ([]service.SBPBank, error)
}

type DestinationHandler struct {
	svc DestinationServiceI
	v   *validator.Validator
	log *slog.Logger
}

func NewDestinationHandler(svc DestinationServiceI, v *validator.Validator, log *slog.Logger) *DestinationHandler {
	return &DestinationHandler{svc: svc, v: v, log: log}
}

func (h *DestinationHandler) Get(w http.ResponseWriter, r *http.Request) {
	actor, err := contextx.PartnerActor(r)
	if err != nil {
		response.Unauthorized(w, "unauthorized")

		return
	}

	dest, err := h.svc.Get(r.Context(), actor.PartnerID)
	if err != nil {
		if errors.Is(err, payoutErrors.ErrDestinationNotFound) {
			response.NotFound(w, "payout destination not configured")

			return
		}

		h.log.Error("get payout destination", slog.String("err", err.Error()))
		response.InternalError(w, err)

		return
	}

	response.Success(w, dto.ToDestinationResponse(dest))
}

func (h *DestinationHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	actor, err := contextx.PartnerActor(r)
	if err != nil {
		response.Unauthorized(w, "unauthorized")

		return
	}

	var req dto.UpsertDestinationRequest
	if errs := h.v.DecodeAndValidate(r, &req); errs != nil {
		response.ValidationError(w, "validation failed", errs)

		return
	}

	dest, err := h.svc.Upsert(r.Context(), actor.PartnerID, domain.DestinationType(req.Type), req.SBPPhone, req.SBPBankID, req.RecipientName)
	if err != nil {
		h.log.Error("upsert payout destination", slog.String("err", err.Error()))
		response.InternalError(w, err)

		return
	}

	response.Success(w, dto.ToDestinationResponse(dest))
}

func (h *DestinationHandler) GetSBPBanks(w http.ResponseWriter, r *http.Request) {
	banks, err := h.svc.GetSBPBanks(r.Context())
	if err != nil {
		h.log.Error("get SBP banks", slog.String("err", err.Error()))
		response.InternalError(w, err)

		return
	}

	items := make([]dto.SBPBankResponse, len(banks))
	for i, b := range banks {
		items[i] = dto.SBPBankResponse{BankID: b.BankID, Name: b.Name, BIC: b.BIC}
	}

	response.Success(w, items)
}
