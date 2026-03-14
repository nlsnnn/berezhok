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

type locationHandler struct {
	log       *slog.Logger
	validator *validator.Validator
	svc       locationSvc
	// partService partnerSvc
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
		// partService: partSvc,
	}
}

func (h *locationHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *locationHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.location.Create"
	log := h.log.With(slog.String("op", op))

	partnerID, ok := r.Context().Value("partner_id").(string)
	if !ok {
		log.Error("partner_id not found in context")
		response.InternalError(w, nil)
		return
	}
	partnerUUID, err := uuid.Parse(partnerID)
	if err != nil {
		log.Error("invalid partner_id in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	var req dto.CreateLocationRequest

	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	// location, err := h.svc.Create(r.Context(), sqlc.CreateLocationParams{
	// 	Name:          req.Name,
	// 	Address:       req.Address,
	// 	PartnerID:     partnerUUID,
	// 	Status:        "inactive",
	// 	CategoryCode:  req.CategoryCode,
	// 	StMakepoint:   req.Latitude,
	// 	StMakepoint_2: req.Longitude,
	// })

	location, err := h.svc.Create(r.Context(), partnerUUID,
		req.CategoryCode,
		req.Name,
		req.Address,
		req.Latitude,
		req.Longitude,
	)

	if err != nil {
		switch {
		case errors.Is(err, partner.ErrPartnerNotFound):
			log.Error("partner not found", sl.Err(err))
			response.BadRequest(w, "partner not found")
		case errors.Is(err, partner.ErrLocationCategoryNotFound):
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
