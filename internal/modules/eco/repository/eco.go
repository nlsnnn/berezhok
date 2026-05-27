package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/modules/eco/domain"
)

type EcoRepo struct {
	q *sqlc.Queries
}

func NewEcoRepo(q *sqlc.Queries) *EcoRepo {
	return &EcoRepo{q: q}
}

// GetAggregate returns picked-up order counts and savings grouped by location
// category for the given customer.
func (r *EcoRepo) GetAggregate(ctx context.Context, userID uuid.UUID) ([]domain.CategoryAggregate, error) {
	rows, err := r.q.GetEcoAggregateByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.CategoryAggregate, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.CategoryAggregate{
			CategoryCode: row.CategoryCode,
			PickedCount:  int(row.PickedCount),
			Savings:      pgconverter.NumericToDecimalOrZero(row.Savings),
		})
	}
	return result, nil
}
