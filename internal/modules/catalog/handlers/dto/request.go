package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateBoxRequest struct {
	LocationID      uuid.UUID       `json:"location_id" validate:"required"`
	Name            string          `json:"name" validate:"required,min=2,max=100"`
	Description     string          `json:"description" validate:"required"`
	DiscountPrice   decimal.Decimal `json:"discount_price" validate:"required"`
	OriginalPrice   decimal.Decimal `json:"original_price" validate:"omitempty,gt=0"`
	PickupTimeStart string          `json:"pickup_time_start" validate:"required,datetime=15:04"`
	PickupTimeEnd   string          `json:"pickup_time_end"   validate:"required,datetime=15:04"`
	Quantity        int             `json:"quantity" validate:"required"`
	Image           string          `json:"image_url" validate:"omitempty"`
	Status          string          `json:"status" validate:"required,oneof=active inactive draft"`
}

type UpdateBoxRequest struct {
	Name            string          `json:"name" validate:"required,min=2,max=100"`
	Description     string          `json:"description" validate:"required"`
	DiscountPrice   decimal.Decimal `json:"discount_price" validate:"required"`
	OriginalPrice   decimal.Decimal `json:"original_price" validate:"omitempty,gt=0"`
	PickupTimeStart string          `json:"pickup_time_start" validate:"required,datetime=15:04"`
	PickupTimeEnd   string          `json:"pickup_time_end"   validate:"required,datetime=15:04"`
	Quantity        int             `json:"quantity" validate:"required"`
	Image           string          `json:"image_url" validate:"omitempty"`
	Status          string          `json:"status" validate:"required,oneof=active inactive draft sold_out"`
}
