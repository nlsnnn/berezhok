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

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
)

type partnerSvcDashboardStub struct {
	dashboardFn    func(ctx context.Context, userID string) (domain.PartnerDashboard, error)
	addLegalInfoFn func(ctx context.Context, input service.AddLegalInfoInput) error
}

func (s *partnerSvcDashboardStub) ChangePassword(ctx context.Context, input service.ChangePasswordInput) error {
	return nil
}

func (s *partnerSvcDashboardStub) Profile(ctx context.Context, userID string) (domain.PartnerProfile, error) {
	return domain.PartnerProfile{}, nil
}

func (s *partnerSvcDashboardStub) Dashboard(ctx context.Context, userID string) (domain.PartnerDashboard, error) {
	if s.dashboardFn != nil {
		return s.dashboardFn(ctx, userID)
	}

	return domain.PartnerDashboard{}, nil
}

func (s *partnerSvcDashboardStub) AddLegalInfo(ctx context.Context, input service.AddLegalInfoInput) error {
	if s.addLegalInfoFn != nil {
		return s.addLegalInfoFn(ctx, input)
	}

	return nil
}

func (s *partnerSvcDashboardStub) CanActivateBoxes(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	return true, nil
}

func TestPartnerDashboardSuccess(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	h := NewPartnerHandler(&partnerSvcDashboardStub{
		dashboardFn: func(ctx context.Context, gotUserID string) (domain.PartnerDashboard, error) {
			if gotUserID != userID.String() {
				t.Fatalf("expected user id %s, got %s", userID, gotUserID)
			}

			nextPayoutDate := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)

			return domain.PartnerDashboard{
				Partner:  domain.Partner{ID: "partner-id", BrandName: "Mama pechet", Status: domain.PartnerStatusActive},
				Employee: domain.Employee{ID: "employee-id", Name: "Egor", Email: "owner@example.com", Role: domain.EmployeeRoleOwner},
				Locations: []domain.DashboardLocation{{
					ID:               "location-id",
					Name:             "Location",
					Address:          "Lenina 1",
					Status:           domain.LocationStatusActive,
					ActiveBoxesCount: 3,
				}},
				Today: domain.DashboardTodayStats{
					PendingConfirmation: 2,
					Confirmed:           5,
					PickedUp:            1,
					Completed:           8,
				},
				Week: domain.DashboardWeekStats{
					OrdersCompleted: 42,
					GrossRevenue:    8550,
					NetRevenue:      7695,
					AvgRating:       4.5,
				},
				Finance: domain.DashboardFinance{
					BalancePending: 7695,
					NextPayoutDate: &nextPayoutDate,
				},
			}, nil
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/partner/dashboard", nil)
	req = req.WithContext(context.WithValue(req.Context(), contextx.UserIDKey, userID))
	rr := httptest.NewRecorder()

	h.Dashboard(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["success"] != true {
		t.Fatalf("expected success=true, got %v", body["success"])
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", body["data"])
	}

	today, ok := data["today"].(map[string]any)
	if !ok {
		t.Fatalf("expected today object, got %T", data["today"])
	}

	if today["pending_confirmation"] != float64(2) {
		t.Fatalf("expected pending_confirmation=2, got %v", today["pending_confirmation"])
	}

	week, ok := data["week"].(map[string]any)
	if !ok {
		t.Fatalf("expected week object, got %T", data["week"])
	}

	if week["net_revenue"] != float64(7695) {
		t.Fatalf("expected net_revenue=7695, got %v", week["net_revenue"])
	}

	finance, ok := data["finance"].(map[string]any)
	if !ok {
		t.Fatalf("expected finance object, got %T", data["finance"])
	}

	if finance["balance_pending"] != float64(7695) {
		t.Fatalf("expected balance_pending=7695, got %v", finance["balance_pending"])
	}
}

func TestPartnerDashboardInternalError(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	h := NewPartnerHandler(&partnerSvcDashboardStub{
		dashboardFn: func(ctx context.Context, gotUserID string) (domain.PartnerDashboard, error) {
			return domain.PartnerDashboard{}, errors.New("db down")
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/partner/dashboard", nil)
	req = req.WithContext(context.WithValue(req.Context(), contextx.UserIDKey, userID))
	rr := httptest.NewRecorder()

	h.Dashboard(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}
}

func TestPartnerAddLegalInfoSuccess(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()

	h := NewPartnerHandler(&partnerSvcDashboardStub{
		addLegalInfoFn: func(ctx context.Context, input service.AddLegalInfoInput) error {
			if input.PartnerID != partnerID.String() {
				t.Fatalf("expected partner id %s, got %s", partnerID, input.PartnerID)
			}

			if input.Inn != "1234567890" {
				t.Fatalf("expected inn 1234567890, got %s", input.Inn)
			}

			return nil
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))

	body := []byte(`{"inn":"1234567890","ogrn":"1234567890123","kpp":"123456789","legal_address":"Moscow, Lenina 1","legal_name":"OOO Berezhok"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/partner/legal-info", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), contextx.PartnerIDKey, partnerID))
	rr := httptest.NewRecorder()

	h.AddLegalInfo(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestPartnerAddLegalInfoValidationError(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()
	h := NewPartnerHandler(&partnerSvcDashboardStub{}, slog.New(slog.NewTextHandler(io.Discard, nil)))

	body := []byte(`{"inn":"abc","legal_address":"x","legal_name":"x"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/partner/legal-info", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), contextx.PartnerIDKey, partnerID))
	rr := httptest.NewRecorder()

	h.AddLegalInfo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}
