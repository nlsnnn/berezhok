package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/shared/domain"
)

type BoxStatus string

const (
	BoxStatusActive   BoxStatus = "active"
	BoxStatusInactive BoxStatus = "inactive"
	BoxStatusSold     BoxStatus = "sold_out"
	BoxStatusDraft    BoxStatus = "draft"
)

type SurpriseBox struct {
	ID         uuid.UUID
	LocationID uuid.UUID

	Name        string
	Description string
	Price       Price
	PickupTime  domain.PickupTime
	Quantity    int
	Status      BoxStatus
	Image       string

	CreatedAt time.Time
}

type Price struct {
	Original decimal.Decimal
	Discount decimal.Decimal
}

func NewSurpriseBox(locationID uuid.UUID, name, description string, originalPrice, discountPrice decimal.Decimal, pickupTimeStart, pickupTimeEnd time.Time, quantity int, status BoxStatus, image string) (SurpriseBox, error) {
	if locationID == uuid.Nil {
		return SurpriseBox{}, fmt.Errorf("location ID is required")
	}

	return SurpriseBox{
		LocationID:  locationID,
		Name:        name,
		Description: description,
		Price: Price{
			Original: originalPrice,
			Discount: discountPrice,
		},
		PickupTime: domain.PickupTime{
			Start: pickupTimeStart,
			End:   pickupTimeEnd,
		},
		Quantity: quantity,
		Status:   status,
		Image:    image,
	}, nil
}

func (b *SurpriseBox) IsAvailable() bool {
	return b.Status == BoxStatusActive && b.Quantity > 0
}
