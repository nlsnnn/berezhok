package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

// Profile is a read model for customer profile endpoint.
type Profile struct {
	ID           uuid.UUID
	Phone        sharedDomain.Phone
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	OrdersCount  int
	ReviewsCount int
	SavedAmount  decimal.Decimal
}
