package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/customer/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/modules/customer/service"
	"github.com/nlsnnn/berezhok/internal/shared/logger/sl"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type locationHandler struct {
	service locationSvc
	log     *slog.Logger
}

type locationSvc interface {
	SearchLocations(ctx context.Context, input service.SearchLocationsInput) (service.SearchLocationsResult, error)
	GetLocationDetails(ctx context.Context, locationID uuid.UUID) (service.LocationDetails, error)
}

func NewLocationHandler(service locationSvc, log *slog.Logger) *locationHandler {
	return &locationHandler{
		service: service,
		log:     log,
	}
}

// SearchLocations handles GET /customer/locations
func (h *locationHandler) SearchLocations(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	categoryCode := r.URL.Query().Get("category")
	var categoryPtr *string
	if categoryCode != "" {
		categoryPtr = &categoryCode
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	result, err := h.service.SearchLocations(r.Context(), service.SearchLocationsInput{
		CategoryCode: categoryPtr,
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		h.log.Error("failed to search locations", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	// Convert to response DTO
	items := make([]dto.LocationSearchResponse, len(result.Locations))
	for i, loc := range result.Locations {
		items[i] = dto.LocationSearchResponse{
			ID:   loc.ID.String(),
			Name: loc.Name,
			Category: dto.CategoryResponse{
				Code:    loc.Category.Code,
				Name:    loc.Category.Name,
				IconURL: loc.Category.IconURL,
			},
			Address: loc.Address,
			Coordinates: dto.CoordinatesResponse{
				Lat: loc.Coords.Latitude,
				Lng: loc.Coords.Longitude,
			},
			LogoURL: loc.LogoURL,
			Rating:  nil,
			ActiveBoxesCount: 0,
		}
	}

	resp := dto.LocationSearchResultResponse{
		Items: items,
		Pagination: dto.PaginationResponse{
			Total:   int(result.Total),
			Limit:   result.Limit,
			Offset:  result.Offset,
			HasMore: result.HasMore,
		},
	}

	response.Success(w, resp)
}

// GetLocationDetails handles GET /customer/locations/{location_id}
func (h *locationHandler) GetLocationDetails(w http.ResponseWriter, r *http.Request) {
	locationIDStr := chi.URLParam(r, "location_id")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		response.BadRequest(w, "invalid location ID")
		return
	}

	details, err := h.service.GetLocationDetails(r.Context(), locationID)
	if err != nil {
		h.log.Error("failed to get location details", sl.Err(err))
		response.NotFound(w, "location not found")
		return
	}

	// Convert active boxes
	activeBoxes := make([]dto.BoxResponse, len(details.ActiveBoxes))
	for i, box := range details.ActiveBoxes {
		activeBoxes[i] = dto.BoxResponse{
			ID:                box.ID.String(),
			Name:              box.Name,
			Description:       box.Description,
			OriginalPrice:     box.Price.Original.InexactFloat64(),
			DiscountPrice:     box.Price.Discount.InexactFloat64(),
			QuantityAvailable: box.Quantity,
			PickupTime: dto.PickupTimeResponse{
				Start: box.PickupTime.Start.Format("15:04"),
				End:   box.PickupTime.End.Format("15:04"),
			},
			ImageURL: box.Image,
		}
	}

	resp := dto.LocationDetailsResponse{
		ID:   details.Location.ID.String(),
		Name: details.Location.Name,
		Category: dto.CategoryResponse{
			Code:    details.Location.Category.Code,
			Name:    details.Location.Category.Name,
			IconURL: details.Location.Category.IconURL,
		},
		Address: details.Location.Address,
		Coordinates: dto.CoordinatesResponse{
			Lat: details.Location.Coords.Latitude,
			Lng: details.Location.Coords.Longitude,
		},
		Phone:         details.Location.Phone,
		WorkingHours:  details.Location.WorkingHours,
		LogoURL:       details.Location.LogoURL,
		CoverImageURL: details.Location.CoverImageURL,
		Gallery:       details.Location.GalleryURLs,
		Rating:        nil, // Stub for now
		ActiveBoxes:   activeBoxes,
	}

	response.Success(w, resp)
}
