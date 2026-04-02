package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/validator"
	orderErrors "github.com/nlsnnn/berezhok/internal/modules/order/errors"
	reviewDomain "github.com/nlsnnn/berezhok/internal/modules/review/domain"
	reviewErrors "github.com/nlsnnn/berezhok/internal/modules/review/errors"
	reviewService "github.com/nlsnnn/berezhok/internal/modules/review/service"
)

type testReviewService struct {
	createReviewFn        func(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error)
	listLocationReviewsFn func(ctx context.Context, locationID uuid.UUID, limit, offset int) (*reviewService.ListLocationReviewsResult, error)
}

func (s *testReviewService) CreateReview(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error) {
	if s.createReviewFn != nil {
		return s.createReviewFn(ctx, input)
	}

	return nil, nil
}

func (s *testReviewService) ListLocationReviews(ctx context.Context, locationID uuid.UUID, limit, offset int) (*reviewService.ListLocationReviewsResult, error) {
	if s.listLocationReviewsFn != nil {
		return s.listLocationReviewsFn(ctx, locationID, limit, offset)
	}

	return nil, nil
}

func TestCreateReviewReturnsConflictWhenAlreadyExists(t *testing.T) {
	t.Parallel()

	customerID := uuid.New()

	h := NewReviewHandler(&testReviewService{
		createReviewFn: func(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error) {
			return nil, reviewErrors.ErrReviewAlreadyExists
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body, _ := json.Marshal(map[string]any{
		"order_id": uuid.New().String(),
		"rating":   5,
		"comment":  "great",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/reviews", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), "customer_id", customerID))
	rr := httptest.NewRecorder()

	h.CreateReview(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rr.Code)
	}
}

func TestCreateReviewReturnsNotFoundWhenOrderMissing(t *testing.T) {
	t.Parallel()

	customerID := uuid.New()

	h := NewReviewHandler(&testReviewService{
		createReviewFn: func(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error) {
			return nil, orderErrors.ErrOrderNotFound
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body, _ := json.Marshal(map[string]any{
		"order_id": uuid.New().String(),
		"rating":   5,
		"comment":  "great",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/reviews", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), "customer_id", customerID))
	rr := httptest.NewRecorder()

	h.CreateReview(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestCreateReviewSuccess(t *testing.T) {
	t.Parallel()

	customerID := uuid.New()
	createdID := uuid.New()
	now := time.Now().UTC()

	h := NewReviewHandler(&testReviewService{
		createReviewFn: func(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error) {
			return &reviewService.CreateReviewResult{
				ID:        createdID,
				OrderID:   input.OrderID,
				Rating:    input.Rating,
				Comment:   input.Comment,
				CreatedAt: now,
			}, nil
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body, _ := json.Marshal(map[string]any{
		"order_id": uuid.New().String(),
		"rating":   4,
		"comment":  "nice",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/reviews", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), "customer_id", customerID))
	rr := httptest.NewRecorder()

	h.CreateReview(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}
}

func TestListLocationReviewsSuccess(t *testing.T) {
	t.Parallel()

	locationID := uuid.New()
	now := time.Now().UTC()

	h := NewReviewHandler(&testReviewService{
		listLocationReviewsFn: func(ctx context.Context, id uuid.UUID, limit, offset int) (*reviewService.ListLocationReviewsResult, error) {
			if id != locationID {
				t.Fatalf("expected location id %s, got %s", locationID, id)
			}

			if limit != 20 || offset != 0 {
				t.Fatalf("expected limit/offset 20/0, got %d/%d", limit, offset)
			}

			return &reviewService.ListLocationReviewsResult{
				Items: []reviewDomain.ReviewWithUser{{
					ID:        uuid.New(),
					Rating:    5,
					Comment:   "great",
					UserName:  "Egor",
					CreatedAt: now,
				}},
				Total:   1,
				Limit:   limit,
				Offset:  offset,
				HasMore: false,
			}, nil
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/customer/locations/"+locationID.String()+"/reviews", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("location_id", locationID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	h.ListLocationReviews(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestCreateReviewReturnsBadRequestWhenOrderNotCompleted(t *testing.T) {
	t.Parallel()

	customerID := uuid.New()

	h := NewReviewHandler(&testReviewService{
		createReviewFn: func(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error) {
			return nil, reviewErrors.ErrOrderNotCompleted
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body, _ := json.Marshal(map[string]any{
		"order_id": uuid.New().String(),
		"rating":   5,
		"comment":  "great",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/reviews", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), "customer_id", customerID))
	rr := httptest.NewRecorder()

	h.CreateReview(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateReviewReturnsBadRequestWhenOrderBelongsToAnotherCustomer(t *testing.T) {
	t.Parallel()

	customerID := uuid.New()

	h := NewReviewHandler(&testReviewService{
		createReviewFn: func(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error) {
			return nil, reviewErrors.ErrOrderOwnershipMismatch
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body, _ := json.Marshal(map[string]any{
		"order_id": uuid.New().String(),
		"rating":   5,
		"comment":  "great",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/reviews", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), "customer_id", customerID))
	rr := httptest.NewRecorder()

	h.CreateReview(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateReviewReturnsInternalErrorOnUnknownServiceError(t *testing.T) {
	t.Parallel()

	customerID := uuid.New()

	h := NewReviewHandler(&testReviewService{
		createReviewFn: func(ctx context.Context, input reviewService.CreateReviewInput) (*reviewService.CreateReviewResult, error) {
			return nil, errors.New("db exploded")
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body, _ := json.Marshal(map[string]any{
		"order_id": uuid.New().String(),
		"rating":   5,
		"comment":  "great",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/customer/reviews", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), "customer_id", customerID))
	rr := httptest.NewRecorder()

	h.CreateReview(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}
