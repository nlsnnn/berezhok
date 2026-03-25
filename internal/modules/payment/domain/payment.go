package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Payment struct {
	ID      uuid.UUID
	OrderID uuid.UUID

	// Provider details
	Provider Provider
	Method   string
	Amount   decimal.Decimal

	// Timestamps
	PaidAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Provider struct {
	PaymentID    string
	PaymentLink  string
	ProviderName string
}

type ProviderPaymentResult struct {
	PaymentLink       string
	ProviderPaymentID string
}
