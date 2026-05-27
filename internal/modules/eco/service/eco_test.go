package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/eco/domain"
)

type stubRepo struct {
	rows []domain.CategoryAggregate
	err  error
}

func (s *stubRepo) GetAggregate(_ context.Context, _ uuid.UUID) ([]domain.CategoryAggregate, error) {
	return s.rows, s.err
}

type stubCache struct {
	hit      *domain.EcoStats
	getErr   error
	stored   *domain.EcoStats
	storeErr error
}

func (c *stubCache) Get(_ context.Context, _ uuid.UUID) (*domain.EcoStats, bool, error) {
	if c.getErr != nil {
		return nil, false, c.getErr
	}
	if c.hit != nil {
		return c.hit, true, nil
	}
	return nil, false, nil
}

func (c *stubCache) Set(_ context.Context, _ uuid.UUID, stats domain.EcoStats) error {
	c.stored = &stats
	return c.storeErr
}

func newDiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestGetForUser_CacheHit_SkipsRepo(t *testing.T) {
	cached := &domain.EcoStats{TotalKg: 42, Tier: domain.TierKeeper}
	repo := &stubRepo{err: errors.New("repo must not be called on cache hit")}
	cache := &stubCache{hit: cached}

	svc := NewEcoService(repo, cache, newDiscardLogger())
	got, err := svc.GetForUser(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TotalKg != 42 || got.Tier != domain.TierKeeper {
		t.Errorf("expected cached payload, got %+v", got)
	}
}

func TestGetForUser_CacheMiss_AggregatesAndStores(t *testing.T) {
	// 3 bakery + 2 grocery → 6.4 kg, helper, 16 kg CO₂
	repo := &stubRepo{rows: []domain.CategoryAggregate{
		{CategoryCode: "bakery", PickedCount: 3, Savings: decimal.NewFromInt(900)},
		{CategoryCode: "grocery", PickedCount: 2, Savings: decimal.NewFromInt(920)},
	}}
	cache := &stubCache{}

	svc := NewEcoService(repo, cache, newDiscardLogger())
	got, err := svc.GetForUser(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.BoxesPickedUp != 5 {
		t.Errorf("BoxesPickedUp = %d, want 5", got.BoxesPickedUp)
	}
	if got.TotalKg != 6.4 {
		t.Errorf("TotalKg = %v, want 6.4", got.TotalKg)
	}
	if got.CO2SavedKg != 16 {
		t.Errorf("CO2SavedKg = %v, want 16", got.CO2SavedKg)
	}
	if got.Tier != domain.TierHelper {
		t.Errorf("Tier = %q, want helper", got.Tier)
	}
	if !got.SavingsRub.Equal(decimal.NewFromInt(1820)) {
		t.Errorf("SavingsRub = %s, want 1820", got.SavingsRub.String())
	}
	if cache.stored == nil {
		t.Errorf("cache.Set was not called on miss")
	}
}

func TestGetForUser_NilCache_NoPanic(t *testing.T) {
	repo := &stubRepo{}
	svc := NewEcoService(repo, nil, newDiscardLogger())
	_, err := svc.GetForUser(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
