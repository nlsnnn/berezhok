package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type pinsHandler struct {
	log       *slog.Logger
	validator *validator.Validator
	svc       pinsSvc
}

func NewPinsHandler(log *slog.Logger, v *validator.Validator, svc pinsSvc) pinsHandler {
	return pinsHandler{log: log, validator: v, svc: svc}
}

// ListAvailable handles GET /partner/pins
func (h *pinsHandler) ListAvailable(w http.ResponseWriter, r *http.Request) {
	pins, err := h.svc.ListAvailable(r.Context())
	if err != nil {
		h.log.Error("failed to list pins", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	result := make([]dto.LocationPinResponse, len(pins))
	for i, p := range pins {
		result[i] = dto.LocationPinResponse{Code: p.Code, NameRu: p.NameRu}
	}
	response.Success(w, result)
}

// UpdateLocationPins handles PUT /partner/locations/{id}/pins
func (h *pinsHandler) UpdateLocationPins(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.pins.UpdateLocationPins"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	locationIDStr := chi.URLParam(r, "id")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		response.BadRequest(w, "invalid location id")
		return
	}

	var req dto.UpdateLocationPinsRequest
	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		response.ValidationError(w, "validation failed", errs)
		return
	}

	if err := h.svc.UpdateLocationPins(r.Context(), partnerID.String(), locationID, req.PinCodes); err != nil {
		log.Error("failed to update location pins", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, map[string]string{"message": "pins updated"})
}
