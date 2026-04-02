package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
)

type partnerSvcDashboardStub struct {
	dashboardFn func(ctx context.Context, userID string) (domain.PartnerDashboard, error)
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

func TestPartnerDashboardSuccess(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	h := NewPartnerHandler(&partnerSvcDashboardStub{
		dashboardFn: func(ctx context.Context, gotUserID string) (domain.PartnerDashboard, error) {
			if gotUserID != userID.String() {
				t.Fatalf("expected user id %s, got %s", userID, gotUserID)
			}

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
