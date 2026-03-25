package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MockPaymentProvider is a stub payment provider for development
type MockPaymentProvider struct{}

func NewMockPaymentProvider() *MockPaymentProvider {
	return &MockPaymentProvider{}
}

// Create returns a mock payment URL
func (m *MockPaymentProvider) Create(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error) {
	// Return a fake payment link
	return "https://mock-payment-provider.com/pay/" + orderID.String(), nil
}
