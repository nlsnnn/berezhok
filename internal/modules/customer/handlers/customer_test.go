package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/customer/domain"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
)

type customerServiceStub struct {
	getProfileFn    func(ctx context.Context, userID uuid.UUID) (domain.Profile, error)
	updateProfileFn func(ctx context.Context, userID uuid.UUID, name string) (domain.User, error)
}

func (s *customerServiceStub) GetProfile(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	if s.getProfileFn != nil {
		return s.getProfileFn(ctx, userID)
	}

	return domain.Profile{}, nil
}

func (s *customerServiceStub) UpdateProfile(ctx context.Context, userID uuid.UUID, name string) (domain.User, error) {
	if s.updateProfileFn != nil {
		return s.updateProfileFn(ctx, userID, name)
	}

	return domain.User{}, nil
}

func TestGetProfileReturnsExtendedCustomerStats(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	createdAt := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	h := NewCustomerHandler(&customerServiceStub{
		getProfileFn: func(ctx context.Context, gotUserID uuid.UUID) (domain.Profile, error) {
			if gotUserID != userID {
				t.Fatalf("expected user id %s, got %s", userID, gotUserID)
			}

			return domain.Profile{
				ID:           userID,
				Name:         "егор",
				CreatedAt:    createdAt,
				OrdersCount:  7,
				ReviewsCount: 3,
				SavedAmount:  decimal.RequireFromString("1250.50"),
			}, nil
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/customer/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), contextx.UserIDKey, userID))
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", body["data"])
	}

	if data["name"] != "егор" {
		t.Fatalf("expected name егор, got %v", data["name"])
	}

	if data["orders_count"] != float64(7) {
		t.Fatalf("expected orders_count 7, got %v", data["orders_count"])
	}

	if data["reviews_count"] != float64(3) {
		t.Fatalf("expected reviews_count 3, got %v", data["reviews_count"])
	}

	if data["saved_amount"] != 1250.5 {
		t.Fatalf("expected saved_amount 1250.5, got %v", data["saved_amount"])
	}
}
