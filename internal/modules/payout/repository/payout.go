package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
	payoutErrors "github.com/nlsnnn/berezhok/internal/modules/payout/errors"
	"github.com/nlsnnn/berezhok/internal/modules/payout/service"
)

type PayoutRepo struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewPayoutRepo(q *sqlc.Queries, pool *pgxpool.Pool) *PayoutRepo {
	return &PayoutRepo{q: q, pool: pool}
}

func (r *PayoutRepo) GetDestination(ctx context.Context, partnerID uuid.UUID) (domain.PayoutDestination, error) {
	row, err := r.q.GetPayoutDestination(ctx, partnerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PayoutDestination{}, payoutErrors.ErrDestinationNotFound
		}

		return domain.PayoutDestination{}, err
	}

	return destinationToDomain(row), nil
}

func (r *PayoutRepo) UpsertDestination(ctx context.Context, d domain.PayoutDestination) (domain.PayoutDestination, error) {
	row, err := r.q.UpsertPayoutDestination(ctx, sqlc.UpsertPayoutDestinationParams{
		PartnerID:     d.PartnerID,
		Type:          string(d.Type),
		SbpPhone:      pgconverter.StringToText(d.SBPPhone),
		SbpBankID:     pgconverter.StringToText(d.SBPBankID),
		RecipientName: pgconverter.StringToText(d.RecipientName),
	})
	if err != nil {
		return domain.PayoutDestination{}, err
	}

	return destinationToDomain(row), nil
}

func (r *PayoutRepo) ListActivePartnersForPayout(ctx context.Context) ([]sqlc.ListActivePartnersForPayoutRow, error) {
	return r.q.ListActivePartnersForPayout(ctx)
}

func (r *PayoutRepo) ListUnsettledCompletedOrders(ctx context.Context, partnerID uuid.UUID, from, to time.Time) ([]domain.OrderForPayout, error) {
	rows, err := r.q.ListUnsettledCompletedOrders(ctx, sqlc.ListUnsettledCompletedOrdersParams{
		PartnerID:   partnerID,
		UpdatedAt:   from,
		UpdatedAt_2: to,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]domain.OrderForPayout, len(rows))
	for i, row := range rows {
		orders[i] = domain.OrderForPayout{
			ID:     row.ID,
			Amount: pgconverter.NumericToDecimalOrZero(row.Amount),
		}
	}

	return orders, nil
}

func (r *PayoutRepo) CreatePayoutWithOrders(ctx context.Context, p *domain.Payout, orders []domain.PayoutOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.q.WithTx(tx)

	_, err = qtx.CreatePayout(ctx, sqlc.CreatePayoutParams{
		ID:                    p.ID,
		PartnerID:             p.PartnerID,
		PeriodStart:           p.PeriodStart,
		PeriodEnd:             p.PeriodEnd,
		GrossAmount:           pgconverter.DecimalToNumeric(p.Gross, true),
		CommissionAmount:      pgconverter.DecimalToNumeric(p.Commission, true),
		CommissionRateApplied: pgconverter.DecimalToNumeric(p.CommissionRateApplied, true),
		NetAmount:             pgconverter.DecimalToNumeric(p.Net, true),
		Status:                string(p.Status),
		Provider:              p.Provider,
		IdempotencyKey:        pgconverter.StringToText(p.IdempotencyKey),
	})
	if err != nil {
		return err
	}

	for _, o := range orders {
		if err := qtx.AddPayoutOrder(ctx, sqlc.AddPayoutOrderParams{
			PayoutID:       o.PayoutID,
			OrderID:        o.OrderID,
			OrderAmount:    pgconverter.DecimalToNumeric(o.OrderAmount, true),
			CommissionPart: pgconverter.DecimalToNumeric(o.CommissionPart, true),
		}); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PayoutRepo) LockPendingForDispatch(ctx context.Context, limit int32) ([]uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.q.WithTx(tx)

	ids, err := qtx.LockPendingPayoutsForDispatch(ctx, limit)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		if err := qtx.MarkPayoutProcessing(ctx, id); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *PayoutRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Payout, error) {
	row, err := r.q.GetPayoutByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Payout{}, payoutErrors.ErrPayoutNotFound
		}

		return domain.Payout{}, err
	}

	return payoutToDomain(row), nil
}

func (r *PayoutRepo) MarkCompleted(ctx context.Context, id uuid.UUID, providerPayoutID string) error {
	return r.q.MarkPayoutCompleted(ctx, sqlc.MarkPayoutCompletedParams{
		ID:               id,
		ProviderPayoutID: pgconverter.StringToText(providerPayoutID),
	})
}

func (r *PayoutRepo) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return r.q.MarkPayoutFailed(ctx, sqlc.MarkPayoutFailedParams{
		ID:           id,
		ErrorMessage: pgconverter.StringToText(errMsg),
	})
}

func (r *PayoutRepo) SetProviderPayoutID(ctx context.Context, id uuid.UUID, providerPayoutID string) error {
	return r.q.SetProviderPayoutID(ctx, sqlc.SetProviderPayoutIDParams{
		ID:               id,
		ProviderPayoutID: pgconverter.StringToText(providerPayoutID),
	})
}

func (r *PayoutRepo) ListProcessingPayouts(ctx context.Context, limit int32) ([]service.ProcessingPayout, error) {
	rows, err := r.q.ListProcessingPayouts(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]service.ProcessingPayout, len(rows))
	for i, row := range rows {
		result[i] = service.ProcessingPayout{
			ID:               row.ID,
			ProviderPayoutID: pgconverter.TextToString(row.ProviderPayoutID),
		}
	}

	return result, nil
}

func (r *PayoutRepo) ListByPartner(ctx context.Context, partnerID uuid.UUID, limit, offset int32) ([]domain.Payout, int64, error) {
	rows, err := r.q.ListPayoutsByPartner(ctx, sqlc.ListPayoutsByPartnerParams{
		PartnerID: partnerID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountPayoutsByPartner(ctx, partnerID)
	if err != nil {
		return nil, 0, err
	}

	payouts := make([]domain.Payout, len(rows))
	for i, row := range rows {
		payouts[i] = payoutToDomain(row)
	}

	return payouts, total, nil
}

func (r *PayoutRepo) ListOrdersByPayoutID(ctx context.Context, payoutID uuid.UUID) ([]domain.PayoutOrder, error) {
	rows, err := r.q.ListPayoutOrders(ctx, payoutID)
	if err != nil {
		return nil, err
	}

	orders := make([]domain.PayoutOrder, len(rows))
	for i, row := range rows {
		orders[i] = domain.PayoutOrder{
			PayoutID:       row.PayoutID,
			OrderID:        row.OrderID,
			OrderAmount:    pgconverter.NumericToDecimalOrZero(row.OrderAmount),
			CommissionPart: pgconverter.NumericToDecimalOrZero(row.CommissionPart),
		}
	}

	return orders, nil
}

func payoutToDomain(p sqlc.PartnerPayout) domain.Payout {
	d := domain.Payout{
		ID:                    p.ID,
		PartnerID:             p.PartnerID,
		PeriodStart:           p.PeriodStart,
		PeriodEnd:             p.PeriodEnd,
		Gross:                 pgconverter.NumericToDecimalOrZero(p.GrossAmount),
		Commission:            pgconverter.NumericToDecimalOrZero(p.CommissionAmount),
		CommissionRateApplied: pgconverter.NumericToDecimalOrZero(p.CommissionRateApplied),
		Net:                   pgconverter.NumericToDecimalOrZero(p.NetAmount),
		Status:                domain.PayoutStatus(p.Status),
		Provider:              p.Provider,
		IdempotencyKey:        pgconverter.TextToString(p.IdempotencyKey),
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}

	if p.ProviderPayoutID.Valid {
		s := p.ProviderPayoutID.String
		d.ProviderPayoutID = &s
	}

	if p.ErrorMessage.Valid {
		s := p.ErrorMessage.String
		d.ErrorMessage = &s
	}

	if p.ProcessedAt.Valid {
		t := p.ProcessedAt.Time
		d.ProcessedAt = &t
	}

	return d
}

func destinationToDomain(d sqlc.PartnerPayoutDestination) domain.PayoutDestination {
	return domain.PayoutDestination{
		PartnerID:     d.PartnerID,
		Type:          domain.DestinationType(d.Type),
		SBPPhone:      pgconverter.TextToString(d.SbpPhone),
		SBPBankID:     pgconverter.TextToString(d.SbpBankID),
		RecipientName: pgconverter.TextToString(d.RecipientName),
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
