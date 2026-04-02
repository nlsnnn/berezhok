package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	reviewDomain "github.com/nlsnnn/berezhok/internal/modules/review/domain"
	reviewErrors "github.com/nlsnnn/berezhok/internal/modules/review/errors"
)

type ReviewRepo struct {
	q *sqlc.Queries
}

func NewReviewRepo(q *sqlc.Queries) *ReviewRepo {
	return &ReviewRepo{q: q}
}

func (r *ReviewRepo) Create(ctx context.Context, review *reviewDomain.Review) error {
	created, err := r.q.CreateReview(ctx, sqlc.CreateReviewParams{
		OrderID:    review.OrderID,
		UserID:     review.UserID,
		LocationID: review.LocationID,
		Rating:     int32(review.Rating),
		Comment:    pgconverter.StringToText(review.Comment),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return reviewErrors.ErrReviewAlreadyExists
		}

		return err
	}

	review.ID = created.ID
	review.CreatedAt = created.CreatedAt
	review.UpdatedAt = created.UpdatedAt

	return nil
}

func (r *ReviewRepo) ListByLocationID(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]reviewDomain.ReviewWithUser, error) {
	rows, err := r.q.ListLocationReviews(ctx, sqlc.ListLocationReviewsParams{
		LocationID: locationID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, err
	}

	items := make([]reviewDomain.ReviewWithUser, len(rows))
	for i, row := range rows {
		items[i] = reviewDomain.ReviewWithUser{
			ID:        row.ID,
			Rating:    int(row.Rating),
			Comment:   row.Comment,
			UserName:  row.UserName,
			CreatedAt: row.CreatedAt,
		}
	}

	return items, nil
}

func (r *ReviewRepo) CountByLocationID(ctx context.Context, locationID uuid.UUID) (int, error) {
	total, err := r.q.CountLocationReviews(ctx, locationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return int(total), nil
}
