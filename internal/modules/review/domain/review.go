package domain

import (
	"time"

	"github.com/google/uuid"
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
