package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderDetails represents enriched order data for detailed customer view.
type OrderDetails struct {
	ID              uuid.UUID
	CustomerID      uuid.UUID
	Status          OrderStatus
	PickupCode      string
	QRCodeURL       string
	Amount          float64
	BoxName         string
	BoxImageURL     string
	LocationName    string
	LocationAddress string
	LocationPhone   string
	LocationLat     float64
	LocationLng     float64
	PickupTimeStart time.Time
	PickupTimeEnd   time.Time
	CreatedAt       time.Time
	ConfirmedAt     *time.Time
}

// PartnerOrderByCode represents enriched order data for partner lookup by pickup code.
type PartnerOrderByCode struct {
	ID              uuid.UUID
	PickupCode      string
	Status          OrderStatus
	BoxName         string
	BoxImageURL     string
	CustomerPhone   string
	CustomerName    string
	PickupTimeStart time.Time
	PickupTimeEnd   time.Time
	CreatedAt       time.Time
}

// OrderListItem represents enriched order data for list view.
type OrderListItem struct {
	ID              uuid.UUID
	Status          OrderStatus
	PickupCode      string
	Amount          float64
	BoxName         string
	LocationName    string
	PickupTimeStart time.Time
	CreatedAt       time.Time
	HasReview       bool
}
