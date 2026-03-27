package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
	paymentErrors "github.com/nlsnnn/berezhok/internal/modules/payment/errors"
)

type paymentRepo interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	CreateEvent(ctx context.Context, paymentID uuid.UUID, eventType string, payload interface{}) error

	GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error
}

type paymentProvider interface {
	Create(ctx context.Context, amount string, description string, returnURL string, metadata map[string]string) (domain.ProviderPaymentResult, error)
}

type orderStatusUpdater interface {
	MarkOrderPaid(ctx context.Context, orderID uuid.UUID) error
	MarkOrderCanceled(ctx context.Context, orderID uuid.UUID) error
}

type paymentService struct {
	repo         paymentRepo
	provider     paymentProvider
	orderUpdater orderStatusUpdater
}

func NewPaymentService(repo paymentRepo, provider paymentProvider, orderUpdater orderStatusUpdater) *paymentService {
	return &paymentService{
		repo:         repo,
		provider:     provider,
		orderUpdater: orderUpdater,
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
			ProviderName: domain.ProviderYookassa,
		},
		Amount: amount,
	}

	err = s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return "", err
	}

	return providerResult.PaymentLink, nil
}

func (s *paymentService) ProccessEvent(ctx context.Context, orderID uuid.UUID, eventType string, payload interface{}) error {
	payment, err := s.repo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	// TODO: validate event type and payload structure based on provider's documentation
	err = s.repo.CreateEvent(ctx, payment.ID, eventType, payload)
	if err != nil {
		return err
	}

	switch eventType {
	case "succeeded":
		return s.handleSuccess(ctx, payment)
	case "failed", "canceled", "cancelled":
		return s.handleCancel(ctx, payment)
	}

	return nil
}

func (s *paymentService) handleSuccess(ctx context.Context, payment *domain.Payment) error {
	err := payment.SetSuccess()
	if err != nil {
		if errors.Is(err, paymentErrors.ErrPaymentAlreadyHandled) {
			return nil
		}
		return err
	}

	err = s.repo.UpdatePaymentStatus(ctx, payment.ID, payment.Status)
	if err != nil {
		return err
	}

	err = s.orderUpdater.MarkOrderPaid(ctx, payment.OrderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *paymentService) handleCancel(ctx context.Context, payment *domain.Payment) error {
	err := payment.SetCanceled()
	if err != nil {
		if errors.Is(err, paymentErrors.ErrPaymentAlreadyHandled) {
			return nil
		}
		return err
	}

	err = s.repo.UpdatePaymentStatus(ctx, payment.ID, payment.Status)
	if err != nil {
		return err
	}

	err = s.orderUpdater.MarkOrderCanceled(ctx, payment.OrderID)
	if err != nil {
		return err
	}

	return nil
}

// func (s *paymentService) GetPayment(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
// 	return s.repo.GetPaymentByID(ctx, paymentID)
// }
