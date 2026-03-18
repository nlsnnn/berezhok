package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BoxResponse struct {
	ID            uuid.UUID          `json:"id"`
	LocationID    uuid.UUID          `json:"location_id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	OriginalPrice decimal.Decimal    `json:"original_price"`
	DiscountPrice decimal.Decimal    `json:"discount_price"`
	PickupTime    PickupTimeResponse `json:"pickup_time"`
	Quantity      int                `json:"quantity"`
	Image         string             `json:"image_url"`
	Status        string             `json:"status"`
	CreatedAt     string             `json:"created_at"`
}

type PickupTimeResponse struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
