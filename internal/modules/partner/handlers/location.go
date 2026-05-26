package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type locationHandler struct {
	log       *slog.Logger
	validator *validator.Validator
	svc       locationSvc
}

func NewLocationHandler(
	log *slog.Logger,
	validator *validator.Validator,
	svc locationSvc,
	partSvc partnerSvc,
) locationHandler {
	return locationHandler{
		log:       log,
		svc:       svc,
		validator: validator,
	}
}

func (h *locationHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *locationHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.location.Create"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	var req dto.CreateLocationRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	input := req.ToInput(partnerID.String())
	input.Status = domain.LocationStatusPendingReview

	location, err := h.svc.Create(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrPartnerNotFound):
			log.Error("partner not found", sl.Err(err))
			response.BadRequest(w, "partner not found")
		case errors.Is(err, partnerErrors.ErrLocationCategoryNotFound):
			log.Error("location category not found", sl.Err(err))
			response.BadRequest(w, "invalid category code")
		default:
			log.Error("failed to create location", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	log.Info("location created successfully", slog.String("location_id", location.ID))

	response.Success(w, dto.FromLocation(location))
}

func (h *locationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.location.GetByID"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	locationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid location id")
		return
	}

	location, err := h.svc.GetByID(r.Context(), partnerID, locationID)
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrLocationNotFound):
			response.NotFound(w, "location not found")
		case errors.Is(err, partnerErrors.ErrLocationNotOwnedByPartner):
			response.Forbidden(w, "access denied")
		default:
			log.Error("failed to get location", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Success(w, dto.FromLocation(location))
}

func (h *locationHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.location.Update"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	locationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid location id")
		return
	}

	var req dto.UpdateLocationRequest
	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	location, err := h.svc.Update(r.Context(), req.ToInput(locationID, partnerID))
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrLocationNotFound):
			response.NotFound(w, "location not found")
		case errors.Is(err, partnerErrors.ErrLocationNotOwnedByPartner):
			response.Forbidden(w, "access denied")
		case errors.Is(err, partnerErrors.ErrInvalidWorkingHours):
			response.BadRequest(w, "invalid working hours")
		default:
			log.Error("failed to update location", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	log.Info("location updated successfully", slog.String("location_id", location.ID))

	response.Success(w, dto.FromLocation(location))
}
