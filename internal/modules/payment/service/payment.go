package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
)

type paymentRepo interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error)
	// UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error
}

type paymentProvider interface {
	Create(ctx context.Context, amount string, description string, method string, returnURL string) (domain.ProviderPaymentResult, error)
}

type paymentService struct {
	repo     paymentRepo
	provider paymentProvider
}

func NewPaymentService(repo paymentRepo, provider paymentProvider) *paymentService {
	return &paymentService{
		repo:     repo,
		provider: provider,
	}
}

// Create(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error)

func (s *paymentService) Create(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error) {
	// 1. Create payment in provider
	description := fmt.Sprintf("Payment for order %s", orderID.String())
	returnURL := "https://example.com/payment/return"
	method := "bank_card"
	providerResult, err := s.provider.Create(ctx, amount.String(), description, method, returnURL)
	if err != nil {
		return "", err
	}

	// 2. Save payment record in our database
	payment := &domain.Payment{
		OrderID: orderID,
		Provider: domain.Provider{
			PaymentID:    providerResult.ProviderPaymentID,
			PaymentLink:  providerResult.PaymentLink,
			ProviderName: "yookassa",
		},
		Amount: amount,
		Method: method,
	}

	err = s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return "", err
	}

	return providerResult.PaymentLink, nil
}
