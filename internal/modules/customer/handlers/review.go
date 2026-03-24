package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type reviewHandler struct {
	log *slog.Logger
}

func NewReviewHandler(log *slog.Logger) *reviewHandler {
	return &reviewHandler{log: log}
}

// CreateReview handles POST /customer/reviews
func (h *reviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Review module not implemented yet", http.StatusNotImplemented)
}

// ListLocationReviews handles GET /customer/locations/{location_id}/reviews
func (h *reviewHandler) ListLocationReviews(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Review module not implemented yet", http.StatusNotImplemented)
}
