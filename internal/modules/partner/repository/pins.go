package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
)

type PinsRepo struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewPinsRepo(q *sqlc.Queries, pool *pgxpool.Pool) *PinsRepo {
	return &PinsRepo{q: q, pool: pool}
}

func (r *PinsRepo) ListAvailable(ctx context.Context) ([]domain.LocationPin, error) {
	rows, err := r.q.ListLocationPins(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.LocationPin, len(rows))
	for i, row := range rows {
		result[i] = domain.LocationPin{
			Code:   row.Code,
			NameRu: row.NameRu,
			Sort:   int(row.SortOrder),
		}
	}
	return result, nil
}

func (r *PinsRepo) GetForLocation(ctx context.Context, locationID uuid.UUID) ([]domain.LocationPin, error) {
	rows, err := r.q.GetLocationSelectedPins(ctx, locationID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.LocationPin, len(rows))
	for i, row := range rows {
		result[i] = domain.LocationPin{
			Code:   row.Code,
			NameRu: row.NameRu,
			Sort:   int(row.SortOrder),
		}
	}
	return result, nil
}

func (r *PinsRepo) SetForLocation(ctx context.Context, locationID uuid.UUID, pinCodes []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	qtx := r.q.WithTx(tx)

	if err := qtx.DeleteLocationSelectedPins(ctx, locationID); err != nil {
		return err
	}

	for _, code := range pinCodes {
		if err := qtx.InsertLocationSelectedPin(ctx, sqlc.InsertLocationSelectedPinParams{
			LocationID: locationID,
			PinCode:    code,
		}); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
