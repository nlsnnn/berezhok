package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type BoxStatus string

const (
	BoxStatusActive   BoxStatus = "active"
	BoxStatusInactive BoxStatus = "inactive"
	BoxStatusSold     BoxStatus = "sold_out"
	BoxStatusDraft    BoxStatus = "draft"
)

type SurpriseBox struct {
	ID         string
	LocationID string

	Name        string
	Description string
	Price       Price
	PickupTime  PickupTime
	Quantity    int
	Status      BoxStatus
	Image       string
}

type PickupTime struct {
	Start time.Time
	End   time.Time
}

type Price struct {
	Original decimal.Decimal
	Discount decimal.Decimal
}
