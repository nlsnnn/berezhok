package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	orderErrors "github.com/nlsnnn/berezhok/internal/modules/order/errors"
	reviewErrors "github.com/nlsnnn/berezhok/internal/modules/review/errors"
	"github.com/nlsnnn/berezhok/internal/modules/review/handlers/dto"
	reviewService "github.com/nlsnnn/berezhok/internal/modules/review/service"
	sharedErrors "github.com/nlsnnn/berezhok/internal/shared/errors"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type reviewSvc interface {
	CreateReview(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error)
	ListLocationReviews(ctx context.Context, locationID uuid.UUID, limit, offset int) (*reviewService.ListLocationReviewsResult, error)
}

type reviewHandler struct {
	service reviewSvc
	log     *slog.Logger
	v       *validator.Validator
}

func NewReviewHandler(service reviewSvc, log *slog.Logger, v *validator.Validator) *reviewHandler {
	return &reviewHandler{
		service: service,
		log:     log,
		v:       v,
	}
}

func (h *reviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	const op = "review.handler.CreateReview"
	log := h.log.With(slog.String("op", op))

	userID, err := h.getCustomerIDFromContext(r)
	if err != nil {
		log.Error("failed to get customer_id from context", sl.Err(err))
		response.Unauthorized(w, "authentication required")
		return
	}

	var req dto.CreateReviewRequest
	if errs := h.v.DecodeAndValidate(r, &req); errs != nil {
		log.Warn("invalid request", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		response.BadRequest(w, "invalid order_id format")
		return
	}

	created, err := h.service.CreateReview(r.Context(), reviewService.CreateReviewInput{
		UserID:  userID,
		OrderID: orderID,
		Rating:  req.Rating,
		Comment: req.Comment,
	})
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrOrderNotFound):
			response.NotFound(w, "order not found")
		case errors.Is(err, reviewErrors.ErrOrderOwnershipMismatch):
			response.BadRequest(w, "order does not belong to customer")
		case errors.Is(err, reviewErrors.ErrOrderNotCompleted):
			response.BadRequest(w, "order is not completed")
		case errors.Is(err, reviewErrors.ErrReviewAlreadyExists):
			response.Error(w, "review already exists", http.StatusConflict)
		default:
			log.Error("failed to create review", sl.Err(err))
			response.InternalError(w, nil)
		}

		return
	}

	response.Created(w, dto.ReviewResponse{
		ReviewID:  created.ID.String(),
		OrderID:   created.OrderID.String(),
		Rating:    created.Rating,
		Comment:   created.Comment,
		CreatedAt: created.CreatedAt,
	})
}

func (h *reviewHandler) ListLocationReviews(w http.ResponseWriter, r *http.Request) {
	const op = "review.handler.ListLocationReviews"
	log := h.log.With(slog.String("op", op))

	locationIDStr := chi.URLParam(r, "location_id")
	locationID, err := uuid.Parse(locationIDStr)
	if err != nil {
		response.BadRequest(w, "invalid location_id format")
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, parseErr := strconv.Atoi(limitStr); parseErr == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, parseErr := strconv.Atoi(offsetStr); parseErr == nil && o >= 0 {
			offset = o
		}
	}

	result, err := h.service.ListLocationReviews(r.Context(), locationID, limit, offset)
	if err != nil {
		log.Error("failed to list location reviews", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	items := make([]dto.ReviewResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = dto.ReviewResponse{
			ID:        item.ID.String(),
			Rating:    item.Rating,
			Comment:   item.Comment,
			UserName:  item.UserName,
			CreatedAt: item.CreatedAt,
		}
	}

	response.Success(w, dto.ReviewListResponse{
		Items: items,
		Pagination: dto.PaginationResponse{
			Total:   result.Total,
			Limit:   result.Limit,
			Offset:  result.Offset,
			HasMore: result.HasMore,
		},
	})
}

func (h *reviewHandler) getCustomerIDFromContext(r *http.Request) (uuid.UUID, error) {
	customerIDRaw := r.Context().Value("customer_id")
	if customerIDRaw == nil {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}

	customerID, ok := customerIDRaw.(uuid.UUID)
	if !ok {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}

	return customerID, nil
}
