package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	"github.com/nlsnnn/berezhok/internal/modules/catalog/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/modules/catalog/service"
	"github.com/nlsnnn/berezhok/internal/shared/logger/sl"
	"github.com/nlsnnn/berezhok/internal/shared/response"
	"github.com/nlsnnn/berezhok/internal/shared/validator"
)

type boxHandler struct {
	boxService BoxService
	log        *slog.Logger
	validator  *validator.Validator
}

type BoxService interface {
	CreateBox(ctx context.Context, partnerID uuid.UUID, input service.CreateBoxInput) (domain.SurpriseBox, error)
	GetBoxByID(ctx context.Context, id string) (*domain.SurpriseBox, error)
	UpdateBox(ctx context.Context, input service.UpdateBoxInput) (domain.SurpriseBox, error)
	DeleteBox(ctx context.Context, id string) error
	GetBoxesByLocationID(ctx context.Context, locationID uuid.UUID) ([]domain.SurpriseBox, error)
	GetBoxesByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]domain.SurpriseBox, error)
}

func NewBoxHandler(boxService BoxService, log *slog.Logger, validator *validator.Validator) *boxHandler {
	return &boxHandler{
		boxService: boxService,
		log:        log,
		validator:  validator,
	}
}

func (h *boxHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.handler.box.Create"
	log := h.log.With(slog.String("op", op))

	partnerID, ok := r.Context().Value("partner_id").(string)
	if !ok {
		log.Error("partner_id not found in context")
		response.InternalError(w, nil)
		return
	}

	var req dto.CreateBoxRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	box, err := h.boxService.CreateBox(r.Context(), uuid.MustParse(partnerID), req.ToInput())
	if err != nil {
		switch {
		case errors.Is(err, catalogErrors.ErrInvalidPickupTimeFormat):
			log.Warn("invalid pickup time format", sl.Err(err))
			response.BadRequest(w, "invalid pickup time format, expected HH:MM")
		case errors.Is(err, catalogErrors.ErrInvalidPickupTimeRange):
			log.Warn("invalid pickup time range", sl.Err(err))
			response.BadRequest(w, "pickup_time_end must be after pickup_time_start")
		case errors.Is(err, catalogErrors.ErrLocationNotFound):
			log.Warn("location not found", sl.Err(err))
			response.NotFound(w, "location not found")
		case errors.Is(err, catalogErrors.ErrUnauthorizedLocation):
			log.Warn("unauthorized location", sl.Err(err))
			response.Forbidden(w, "partner does not own the specified location")
		default:
			log.Error("failed to create box", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Created(w, dto.BoxToResponse(box))
}

func (h *boxHandler) GetAllByLocationID(w http.ResponseWriter, r *http.Request) {
	locationID := chi.URLParam(r, "location_id")

	boxes, err := h.boxService.GetBoxesByLocationID(r.Context(), uuid.MustParse(locationID))
	if err != nil {
		h.log.Error("failed to get boxes by location id", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, dto.BoxesToResponses(boxes))
}

func (h *boxHandler) GetAllByPartnerID(w http.ResponseWriter, r *http.Request) {
	partnerID, ok := r.Context().Value("partner_id").(string)
	if !ok {
		h.log.Error("partner_id not found in context")
		response.InternalError(w, nil)
		return
	}

	boxes, err := h.boxService.GetBoxesByPartnerID(r.Context(), uuid.MustParse(partnerID))
	if err != nil {
		h.log.Error("failed to get boxes by partner id", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, dto.BoxesToResponses(boxes))
}

func (h *boxHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	box, err := h.boxService.GetBoxByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, catalogErrors.ErrInvalidBoxID):
			response.BadRequest(w, "invalid box id")
		case errors.Is(err, catalogErrors.ErrBoxNotFound):
			response.NotFound(w, "box not found")
		default:
			h.log.Error("failed to get box by id", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Success(w, dto.BoxToResponse(*box))
}

func (h *boxHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateBoxRequest
	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		h.log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	box, err := h.boxService.UpdateBox(r.Context(), req.ToInput(id))
	if err != nil {
		switch {
		case errors.Is(err, catalogErrors.ErrInvalidBoxID):
			response.BadRequest(w, "invalid box id")
		case errors.Is(err, catalogErrors.ErrBoxNotFound):
			response.NotFound(w, "box not found")
		case errors.Is(err, catalogErrors.ErrInvalidPickupTimeFormat):
			response.BadRequest(w, "invalid pickup time format, expected HH:MM")
		case errors.Is(err, catalogErrors.ErrInvalidPickupTimeRange):
			response.BadRequest(w, "pickup_time_end must be after pickup_time_start")
		default:
			h.log.Error("failed to update box", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Success(w, dto.BoxToResponse(box))
}

func (h *boxHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.boxService.DeleteBox(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, catalogErrors.ErrInvalidBoxID):
			response.BadRequest(w, "invalid box id")
		case errors.Is(err, catalogErrors.ErrBoxNotFound):
			response.NotFound(w, "box not found")
		default:
			h.log.Error("failed to delete box", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.NoContent(w)
}
