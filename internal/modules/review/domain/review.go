package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/review/errors"
)

type Review struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	LocationID uuid.UUID
	OrderID    uuid.UUID
	Rating     int
	Comment    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReviewWithUser struct {
	ID        uuid.UUID
	Rating    int
	Comment   string
	UserName  string
	CreatedAt time.Time
}

func NewReview(userID, locationID, orderID uuid.UUID, rating int, comment string) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.ErrRatingOutOfRange
	}

	return &Review{
		ID:         uuid.New(),
		UserID:     userID,
		LocationID: locationID,
		OrderID:    orderID,
		Rating:     rating,
		Comment:    comment,
	}, nil
}
