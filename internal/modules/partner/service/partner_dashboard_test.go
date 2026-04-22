package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
)

type dashboardRepoStub struct {
	dashboardFn func(ctx context.Context, employeeID string) (domain.PartnerDashboard, error)
	statsFn     func(ctx context.Context, employeeID string, filter domain.StatsFilter) (domain.PartnerStats, error)
}

func (r *dashboardRepoStub) FindByID(ctx context.Context, id string) (domain.Partner, error) {
	return domain.Partner{}, nil
}

func (r *dashboardRepoStub) List(ctx context.Context) ([]domain.Partner, error) {
	return nil, nil
}

func (r *dashboardRepoStub) Create(ctx context.Context, name string) (domain.Partner, error) {
	return domain.Partner{}, nil
}

func (r *dashboardRepoStub) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (r *dashboardRepoStub) GetProfile(ctx context.Context, employeeID string) (domain.PartnerProfile, error) {
	return domain.PartnerProfile{}, nil
}

func (r *dashboardRepoStub) GetDashboard(ctx context.Context, employeeID string) (domain.PartnerDashboard, error) {
	if r.dashboardFn != nil {
		return r.dashboardFn(ctx, employeeID)
	}

	return domain.PartnerDashboard{}, nil
}

func (r *dashboardRepoStub) GetStats(ctx context.Context, employeeID string, filter domain.StatsFilter) (domain.PartnerStats, error) {
	if r.statsFn != nil {
		return r.statsFn(ctx, employeeID, filter)
	}

	return domain.PartnerStats{}, nil
}

func (r *dashboardRepoStub) UpdateEmployeePassword(ctx context.Context, employeeID, newHash string) error {
	return nil
}

func (r *dashboardRepoStub) UpsertLegalInfo(ctx context.Context, info domain.LegalInfo) error {
	return nil
}

func (r *dashboardRepoStub) UpdateStatus(ctx context.Context, partnerID string, status domain.PartnerStatus) error {
	return nil
}

type dashboardEmployeeRepoStub struct{}

func (r *dashboardEmployeeRepoStub) FindByID(ctx context.Context, id string) (domain.Employee, error) {
	return domain.Employee{}, nil
}

func TestPartnerServiceDashboard(t *testing.T) {
	t.Parallel()

	expected := domain.PartnerDashboard{
		Partner: domain.Partner{ID: "partner-id", BrandName: "Mama pechet"},
	}

	repo := &dashboardRepoStub{
		dashboardFn: func(ctx context.Context, employeeID string) (domain.PartnerDashboard, error) {
			if employeeID != "employee-id" {
				t.Fatalf("unexpected employee id: %s", employeeID)
			}

			return expected, nil
		},
	}

	svc := NewPartnerService(repo, &dashboardEmployeeRepoStub{})

	actual, err := svc.Dashboard(context.Background(), "employee-id")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if actual.Partner.ID != expected.Partner.ID {
		t.Fatalf("expected partner id %s, got %s", expected.Partner.ID, actual.Partner.ID)
	}
}

func TestPartnerServiceDashboardError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("db unavailable")

	repo := &dashboardRepoStub{
		dashboardFn: func(ctx context.Context, employeeID string) (domain.PartnerDashboard, error) {
			return domain.PartnerDashboard{}, expectedErr
		},
	}

	svc := NewPartnerService(repo, &dashboardEmployeeRepoStub{})

	_, err := svc.Dashboard(context.Background(), "employee-id")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestNormalizeStatsFilterPreset(t *testing.T) {
	t.Parallel()

	filter, err := NormalizeStatsFilter(domain.StatsFilter{Period: "last_7_days"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if filter.Period != "last_7_days" {
		t.Fatalf("expected period last_7_days, got %s", filter.Period)
	}
	if filter.DateFrom.After(filter.DateTo) {
		t.Fatalf("expected valid range, got %v > %v", filter.DateFrom, filter.DateTo)
	}
}

func TestNormalizeStatsFilterCustomRange(t *testing.T) {
	t.Parallel()

	from := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	to := time.Date(2026, 4, 10, 18, 0, 0, 0, time.UTC)

	filter, err := NormalizeStatsFilter(domain.StatsFilter{
		DateFrom: from,
		DateTo:   to,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if filter.Period != "custom" {
		t.Fatalf("expected custom period, got %s", filter.Period)
	}
	if filter.DateFrom.Hour() != 0 || filter.DateTo.Hour() != 0 {
		t.Fatalf("expected normalized dates, got %v and %v", filter.DateFrom, filter.DateTo)
	}
}
