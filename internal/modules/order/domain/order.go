package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/shared/domain"
	"github.com/nlsnnn/berezhok/internal/shared/generator"
	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusPickedUp  OrderStatus = "picked_up"
	OrderStatusRefunded  OrderStatus = "refunded"
	OrderStatusDisputed  OrderStatus = "disputed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	BoxID      uuid.UUID
	LocationID uuid.UUID

	// Pickup details
	PickupCode          string
	QRCode              string
	PickedUpConfirmedBy *uuid.UUID

	// Financial details
	amount decimal.Decimal

	// Pickup time
	PickupTime domain.PickupTime

	Status OrderStatus

	// Confirmation details
	ConfirmationDeadline time.Time
	ConfirmedAt          *time.Time
	ConfirmedBy          *uuid.UUID

	// Cancellation details
	CancellationReason *string
	CancelledAt        *time.Time

	// Timestamps
	PickedUpAt            *time.Time
	UserConfirmedPickupAt *time.Time
	AutoCompletedAt       *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func NewOrder(
	customerID, boxID, locationID uuid.UUID,
	pickupTime domain.PickupTime,
	amount decimal.Decimal,
) *Order {
	if pickupTime.Start.IsZero() || pickupTime.End.IsZero() {
		panic("pickup time must be set")
	}

	code := generator.GeneratePickupCode()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	pickupTimeStart := today.Add(
		time.Duration(pickupTime.Start.Hour())*time.Hour +
			time.Duration(pickupTime.Start.Minute())*time.Minute,
	)
	pickupTimeEnd := today.Add(
		time.Duration(pickupTime.End.Hour())*time.Hour +
			time.Duration(pickupTime.End.Minute())*time.Minute,
	)

	if !pickupTimeEnd.After(pickupTimeStart) {
		pickupTimeEnd = pickupTimeEnd.Add(24 * time.Hour)
	}

	return &Order{
		ID:                   uuid.New(),
		CustomerID:           customerID,
		BoxID:                boxID,
		LocationID:           locationID,
		PickupCode:           code,
		PickupTime:           domain.PickupTime{Start: pickupTimeStart, End: pickupTimeEnd},
		amount:               amount,
		Status:               OrderStatusPending,
		ConfirmationDeadline: pickupTimeEnd,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// Amount returns the order amount
func (o *Order) Amount() decimal.Decimal {
	return o.amount
}

// SetAmount sets the order amount (used by repository layer)
func (o *Order) SetAmount(amount decimal.Decimal) {
	o.amount = amount
}
