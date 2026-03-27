package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/payment/errors"
	"github.com/shopspring/decimal"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusCanceled  PaymentStatus = "cancelled"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type ProviderName string

const (
	ProviderYookassa ProviderName = "yookassa"
)

type Payment struct {
	ID      uuid.UUID
	OrderID uuid.UUID

	// Provider details
	Provider Provider
	Method   string
	Amount   decimal.Decimal
	Status   PaymentStatus

	// Timestamps
	PaidAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Provider struct {
	PaymentID    string
	PaymentLink  string
	ProviderName ProviderName
}

type ProviderPaymentResult struct {
	PaymentLink       string
	ProviderPaymentID string
}

func (p *Payment) IsPaid() bool {
	return p.Status == PaymentStatusSucceeded
}

func (p *Payment) IsHandled() bool {
	return p.Status == PaymentStatusSucceeded || p.Status == PaymentStatusCanceled || p.Status == PaymentStatusFailed
}

func (p *Payment) SetSuccess() error {
	if p.IsHandled() {
		return errors.ErrPaymentAlreadyHandled
	}

	p.Status = PaymentStatusSucceeded
	now := time.Now()
	p.PaidAt = &now
	return nil
}

func (p *Payment) SetCanceled() error {
	if p.IsHandled() {
		return errors.ErrPaymentAlreadyHandled
	}
	p.Status = PaymentStatusCanceled
	now := time.Now()
	p.PaidAt = &now
	return nil
}
