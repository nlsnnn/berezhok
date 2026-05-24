package service

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
)

// ----- mocks -----

type mockCalcRepo struct {
	partners []sqlc.ListActivePartnersForPayoutRow
	orders   []domain.OrderForPayout
	created  []*domain.Payout

	createErr error
}

func (m *mockCalcRepo) ListActivePartnersForPayout(_ context.Context) ([]sqlc.ListActivePartnersForPayoutRow, error) {
	return m.partners, nil
}

func (m *mockCalcRepo) ListUnsettledCompletedOrders(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.OrderForPayout, error) {
	return m.orders, nil
}

func (m *mockCalcRepo) CreatePayoutWithOrders(_ context.Context, p *domain.Payout, _ []domain.PayoutOrder) error {
	if m.createErr != nil {
		return m.createErr
	}

	m.created = append(m.created, p)

	return nil
}

// ----- helpers -----

func numericFromFloat(f float64) pgtype.Numeric {
	d := decimal.NewFromFloat(f)

	return pgtype.Numeric{Int: d.Coefficient(), Exp: d.Exponent(), Valid: true}
}

func makePartner(rate float64) sqlc.ListActivePartnersForPayoutRow {
	return sqlc.ListActivePartnersForPayoutRow{
		ID:             uuid.New(),
		CommissionRate: numericFromFloat(rate),
	}
}

func makeOrder(amount float64) domain.OrderForPayout {
	return domain.OrderForPayout{
		ID:     uuid.New(),
		Amount: decimal.NewFromFloat(amount),
	}
}

func newCalcSvc(repo *mockCalcRepo, minAmount float64) *CalculationService {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	return NewCalculationService(repo, decimal.NewFromFloat(minAmount), log)
}

// ----- tests -----

func TestCalculateForPeriod_CreatesPayoutForActivePartner(t *testing.T) {
	partner := makePartner(0.20)
	orders := []domain.OrderForPayout{makeOrder(1000), makeOrder(500)}

	repo := &mockCalcRepo{partners: []sqlc.ListActivePartnersForPayoutRow{partner}, orders: orders}
	svc := newCalcSvc(repo, 1)

	n, err := svc.CalculateForPeriod(context.Background(), time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)
	assert.Equal(t, 1, n)
	require.Len(t, repo.created, 1)

	p := repo.created[0]
	assert.True(t, decimal.NewFromFloat(1500).Equal(p.Gross), "gross mismatch: %s", p.Gross)
	assert.True(t, decimal.NewFromFloat(300).Equal(p.Commission), "commission mismatch: %s", p.Commission)
	assert.True(t, decimal.NewFromFloat(1200).Equal(p.Net), "net mismatch: %s", p.Net)
	assert.Equal(t, domain.PayoutStatusPending, p.Status)
}

func TestCalculateForPeriod_SkipsPartnerWithNoOrders(t *testing.T) {
	partner := makePartner(0.20)
	repo := &mockCalcRepo{
		partners: []sqlc.ListActivePartnersForPayoutRow{partner},
		orders:   []domain.OrderForPayout{},
	}
	svc := newCalcSvc(repo, 1)

	n, err := svc.CalculateForPeriod(context.Background(), time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, repo.created)
}

func TestCalculateForPeriod_SkipsBelowMinAmount(t *testing.T) {
	partner := makePartner(0.20)
	repo := &mockCalcRepo{
		partners: []sqlc.ListActivePartnersForPayoutRow{partner},
		orders:   []domain.OrderForPayout{makeOrder(10)}, // net = 8, below min 100
	}
	svc := newCalcSvc(repo, 100)

	n, err := svc.CalculateForPeriod(context.Background(), time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, repo.created)
}

func TestCalculateForPeriod_UsesPromoRateWhenActive(t *testing.T) {
	future := time.Now().Add(7 * 24 * time.Hour)
	partner := sqlc.ListActivePartnersForPayoutRow{
		ID:                   uuid.New(),
		CommissionRate:       numericFromFloat(0.20),
		PromoCommissionRate:  pgtype.Numeric{Int: decimal.NewFromFloat(0.10).Coefficient(), Exp: decimal.NewFromFloat(0.10).Exponent(), Valid: true},
		PromoCommissionUntil: pgtype.Date{Time: future, Valid: true},
	}
	orders := []domain.OrderForPayout{makeOrder(1000)}
	repo := &mockCalcRepo{partners: []sqlc.ListActivePartnersForPayoutRow{partner}, orders: orders}
	svc := newCalcSvc(repo, 1)

	n, err := svc.CalculateForPeriod(context.Background(), time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	p := repo.created[0]
	assert.True(t, decimal.NewFromFloat(100).Equal(p.Commission), "commission mismatch: %s", p.Commission)
	assert.True(t, decimal.NewFromFloat(900).Equal(p.Net), "net mismatch: %s", p.Net)
}

func TestCalculateForPeriod_UsesNormalRateWhenPromoExpired(t *testing.T) {
	past := time.Now().AddDate(0, 0, -7) // use a week ago so UTC day is clearly before periodEnd
	partner := sqlc.ListActivePartnersForPayoutRow{
		ID:                   uuid.New(),
		CommissionRate:       numericFromFloat(0.20),
		PromoCommissionRate:  pgtype.Numeric{Int: decimal.NewFromFloat(0.10).Coefficient(), Exp: decimal.NewFromFloat(0.10).Exponent(), Valid: true},
		PromoCommissionUntil: pgtype.Date{Time: past, Valid: true},
	}
	orders := []domain.OrderForPayout{makeOrder(1000)}
	repo := &mockCalcRepo{partners: []sqlc.ListActivePartnersForPayoutRow{partner}, orders: orders}
	svc := newCalcSvc(repo, 1)

	n, err := svc.CalculateForPeriod(context.Background(), time.Now().Add(-7*24*time.Hour), time.Now())
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	p := repo.created[0]
	assert.True(t, decimal.NewFromFloat(200).Equal(p.Commission), "commission mismatch: %s", p.Commission)
}

func TestCalculateForPeriod_IdempotencyKeyEqualsPayoutID(t *testing.T) {
	partner := makePartner(0.20)
	repo := &mockCalcRepo{
		partners: []sqlc.ListActivePartnersForPayoutRow{partner},
		orders:   []domain.OrderForPayout{makeOrder(500)},
	}
	svc := newCalcSvc(repo, 1)

	_, _ = svc.CalculateForPeriod(context.Background(), time.Now().Add(-7*24*time.Hour), time.Now())
	require.Len(t, repo.created, 1)
	assert.Equal(t, repo.created[0].ID.String(), repo.created[0].IdempotencyKey)
}
