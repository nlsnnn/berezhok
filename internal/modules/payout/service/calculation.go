package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
)

type CalculationRepository interface {
	ListActivePartnersForPayout(ctx context.Context) ([]sqlc.ListActivePartnersForPayoutRow, error)
	ListUnsettledCompletedOrders(ctx context.Context, partnerID uuid.UUID, from, to time.Time) ([]domain.OrderForPayout, error)
	CreatePayoutWithOrders(ctx context.Context, p *domain.Payout, orders []domain.PayoutOrder) error
}

type CalculationService struct {
	repo            CalculationRepository
	minPayoutAmount decimal.Decimal
	log             *slog.Logger
}

func NewCalculationService(repo CalculationRepository, minPayoutAmount decimal.Decimal, log *slog.Logger) *CalculationService {
	return &CalculationService{
		repo:            repo,
		minPayoutAmount: minPayoutAmount,
		log:             log,
	}
}

func (s *CalculationService) CalculateForPeriod(ctx context.Context, periodStart, periodEnd time.Time) (int, error) {
	partners, err := s.repo.ListActivePartnersForPayout(ctx)
	if err != nil {
		return 0, fmt.Errorf("list active partners: %w", err)
	}

	created := 0

	for _, partner := range partners {
		n, err := s.processPartner(ctx, partner, periodStart, periodEnd)
		if err != nil {
			s.log.Error("failed to process partner payout",
				slog.String("partner_id", partner.ID.String()),
				slog.String("err", err.Error()),
			)

			continue
		}

		created += n
	}

	return created, nil
}

func (s *CalculationService) processPartner(
	ctx context.Context,
	partner sqlc.ListActivePartnersForPayoutRow,
	periodStart, periodEnd time.Time,
) (int, error) {
	orders, err := s.repo.ListUnsettledCompletedOrders(ctx, partner.ID, periodStart, periodEnd)
	if err != nil {
		return 0, fmt.Errorf("list orders: %w", err)
	}

	if len(orders) == 0 {
		s.log.Info("no unsettled orders for partner", slog.String("partner_id", partner.ID.String()))

		return 0, nil
	}

	rate := s.effectiveRate(partner, periodEnd)

	payout, payoutOrders := domain.NewPayout(partner.ID, periodStart, periodEnd, orders, rate)

	if payout.Net.LessThan(s.minPayoutAmount) {
		s.log.Info("payout below minimum, skipping",
			slog.String("partner_id", partner.ID.String()),
			slog.String("net", payout.Net.String()),
			slog.String("min", s.minPayoutAmount.String()),
		)

		return 0, nil
	}

	if err := s.repo.CreatePayoutWithOrders(ctx, payout, payoutOrders); err != nil {
		return 0, fmt.Errorf("create payout: %w", err)
	}

	s.log.Info("payout created",
		slog.String("payout_id", payout.ID.String()),
		slog.String("partner_id", partner.ID.String()),
		slog.String("net", payout.Net.String()),
		slog.Int("orders", len(orders)),
	)

	return 1, nil
}

func (s *CalculationService) effectiveRate(partner sqlc.ListActivePartnersForPayoutRow, periodEnd time.Time) decimal.Decimal {
	if partner.PromoCommissionUntil.Valid && partner.PromoCommissionRate.Valid {
		// Compare UTC calendar dates to avoid timezone edge cases: a period ending
		// at 00:23 +03:00 is still the previous UTC day but should count as "today".
		periodEndDay := periodEnd.UTC().Truncate(24 * time.Hour)
		promoLastDay := time.Date(
			partner.PromoCommissionUntil.Time.Year(),
			partner.PromoCommissionUntil.Time.Month(),
			partner.PromoCommissionUntil.Time.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		if !periodEndDay.After(promoLastDay) {
			return numericToDecimal(partner.PromoCommissionRate)
		}
	}

	return numericToDecimal(partner.CommissionRate)
}

func numericToDecimal(n pgtype.Numeric) decimal.Decimal {
	if !n.Valid || n.Int == nil {
		return decimal.Zero
	}

	return decimal.NewFromBigInt(n.Int, n.Exp)
}
