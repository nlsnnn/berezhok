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
	UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error
}

type paymentProvider interface {
	Create(ctx context.Context, amount string, description string, returnURL string, metadata map[string]string) (domain.ProviderPaymentResult, error)
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
	description := fmt.Sprintf("Оплата заказа #%s", orderID.String())
	returnURL := "berezhok://orders/" + orderID.String() // Deep link для мобильного приложения
	metadata := map[string]string{
		"order_id": orderID.String(),
	}
	providerResult, err := s.provider.Create(ctx, amount.StringFixed(2), description, returnURL, metadata)
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
	}

	err = s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return "", err
	}

	return providerResult.PaymentLink, nil
}

func (s *paymentService) ProccessEvent(ctx context.Context, orderID string, providerPaymentID string, eventType string) error {

	return nil
}

func (s *paymentService) handleSuccess(ctx context.Context, orderID, providerPaymentID string) error {

	return nil
}

// func (s *paymentService) GetPayment(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
// 	return s.repo.GetPaymentByID(ctx, paymentID)
// }
