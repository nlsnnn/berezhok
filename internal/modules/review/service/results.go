package service

import (
	"time"

	"github.com/google/uuid"

	reviewDomain "github.com/nlsnnn/berezhok/internal/modules/review/domain"
)

type CreateReviewInput struct {
	UserID  uuid.UUID
	OrderID uuid.UUID
	Rating  int
	Comment string
}

type CreateReviewResult struct {
	ID        uuid.UUID
	Rating    int
	Comment   string
	OrderID   uuid.UUID
	CreatedAt time.Time
}

type ListLocationReviewsResult struct {
	Items   []reviewDomain.ReviewWithUser
	Total   int
	Limit   int
	Offset  int
	HasMore bool
}
