package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
	paymentErrors "github.com/nlsnnn/berezhok/internal/modules/payment/errors"
)

type testRepo struct {
	createPaymentFn       func(ctx context.Context, payment *domain.Payment) error
	createEventFn         func(ctx context.Context, paymentID uuid.UUID, eventType string, payload interface{}) error
	getPaymentByIDFn      func(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error)
	getPaymentByOrderIDFn func(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
	updatePaymentStatusFn func(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error
}

func (r *testRepo) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	if r.createPaymentFn != nil {
		return r.createPaymentFn(ctx, payment)
	}
	return nil
}

func (r *testRepo) CreateEvent(ctx context.Context, paymentID uuid.UUID, eventType string, payload interface{}) error {
	if r.createEventFn != nil {
		return r.createEventFn(ctx, paymentID, eventType, payload)
	}
	return nil
}

func (r *testRepo) GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	if r.getPaymentByIDFn != nil {
		return r.getPaymentByIDFn(ctx, paymentID)
	}
	return nil, nil
}

func (r *testRepo) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	if r.getPaymentByOrderIDFn != nil {
		return r.getPaymentByOrderIDFn(ctx, orderID)
	}
	return nil, nil
}

func (r *testRepo) UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error {
	if r.updatePaymentStatusFn != nil {
		return r.updatePaymentStatusFn(ctx, paymentID, status)
	}
	return nil
}

type testProvider struct {
	createFn func(ctx context.Context, amount, description, returnURL string, metadata map[string]string) (domain.ProviderPaymentResult, error)
}

func (p *testProvider) Create(ctx context.Context, amount, description, returnURL string, metadata map[string]string) (domain.ProviderPaymentResult, error) {
	if p.createFn != nil {
		return p.createFn(ctx, amount, description, returnURL, metadata)
	}
	return domain.ProviderPaymentResult{}, nil
}

type testOrderUpdater struct {
	markOrderPaidFn     func(ctx context.Context, orderID uuid.UUID) error
	markOrderCanceledFn func(ctx context.Context, orderID uuid.UUID) error
}

func (u *testOrderUpdater) MarkOrderPaid(ctx context.Context, orderID uuid.UUID) error {
	if u.markOrderPaidFn != nil {
		return u.markOrderPaidFn(ctx, orderID)
	}
	return nil
}

func (u *testOrderUpdater) MarkOrderCanceled(ctx context.Context, orderID uuid.UUID) error {
	if u.markOrderCanceledFn != nil {
		return u.markOrderCanceledFn(ctx, orderID)
	}
	return nil
}

func TestCreateSuccess(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()
	repoCalled := false

	repo := &testRepo{
		createPaymentFn: func(ctx context.Context, payment *domain.Payment) error {
			repoCalled = true
			if payment.OrderID != orderID {
				t.Fatalf("expected order id %s, got %s", orderID, payment.OrderID)
			}
			if payment.Provider.ProviderName != domain.ProviderYookassa {
				t.Fatalf("expected provider %q, got %q", domain.ProviderYookassa, payment.Provider.ProviderName)
			}
			return nil
		},
	}

	provider := &testProvider{
		createFn: func(ctx context.Context, amount, description, returnURL string, metadata map[string]string) (domain.ProviderPaymentResult, error) {
			if amount != "123.45" {
				t.Fatalf("expected amount 123.45, got %s", amount)
			}
			if metadata["order_id"] != orderID.String() {
				t.Fatalf("expected metadata order_id %s, got %s", orderID.String(), metadata["order_id"])
			}
			return domain.ProviderPaymentResult{
				PaymentLink:       "https://pay.example/123",
				ProviderPaymentID: "provider-123",
			}, nil
		},
	}

	svc := NewPaymentService(repo, provider, &testOrderUpdater{})

	link, err := svc.Create(context.Background(), decimal.RequireFromString("123.45"), orderID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if link != "https://pay.example/123" {
		t.Fatalf("unexpected payment link: %s", link)
	}

	if !repoCalled {
		t.Fatal("expected repository CreatePayment to be called")
	}
}

func TestCreateProviderError(t *testing.T) {
	t.Parallel()

	providerErr := errors.New("provider unavailable")
	repoCalled := false

	repo := &testRepo{
		createPaymentFn: func(ctx context.Context, payment *domain.Payment) error {
			repoCalled = true
			return nil
		},
	}

	provider := &testProvider{
		createFn: func(ctx context.Context, amount, description, returnURL string, metadata map[string]string) (domain.ProviderPaymentResult, error) {
			return domain.ProviderPaymentResult{}, providerErr
		},
	}

	svc := NewPaymentService(repo, provider, &testOrderUpdater{})

	_, err := svc.Create(context.Background(), decimal.RequireFromString("100"), uuid.New())
	if !errors.Is(err, providerErr) {
		t.Fatalf("expected provider error, got %v", err)
	}

	if repoCalled {
		t.Fatal("repository must not be called when provider create fails")
	}
}

func TestProcessEventSucceeded(t *testing.T) {
	t.Parallel()

	paymentID := uuid.New()
	orderID := uuid.New()

	createEventCalled := false
	updateStatusCalled := false
	markPaidCalled := false

	repo := &testRepo{
		getPaymentByOrderIDFn: func(ctx context.Context, reqOrderID uuid.UUID) (*domain.Payment, error) {
			if reqOrderID != orderID {
				t.Fatalf("expected order id %s, got %s", orderID, reqOrderID)
			}
			return &domain.Payment{
				ID:      paymentID,
				OrderID: orderID,
				Status:  domain.PaymentStatusPending,
			}, nil
		},
		createEventFn: func(ctx context.Context, id uuid.UUID, eventType string, payload interface{}) error {
			createEventCalled = true
			if id != paymentID {
				t.Fatalf("expected payment id %s, got %s", paymentID, id)
			}
			if eventType != "succeeded" {
				t.Fatalf("expected event type succeeded, got %s", eventType)
			}
			return nil
		},
		updatePaymentStatusFn: func(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
			updateStatusCalled = true
			if id != paymentID {
				t.Fatalf("expected payment id %s, got %s", paymentID, id)
			}
			if status != domain.PaymentStatusSucceeded {
				t.Fatalf("expected status %q, got %q", domain.PaymentStatusSucceeded, status)
			}
			return nil
		},
	}

	updater := &testOrderUpdater{
		markOrderPaidFn: func(ctx context.Context, id uuid.UUID) error {
			markPaidCalled = true
			if id != orderID {
				t.Fatalf("expected order id %s, got %s", orderID, id)
			}
			return nil
		},
	}

	svc := NewPaymentService(repo, &testProvider{}, updater)

	err := svc.ProccessEvent(context.Background(), orderID, "succeeded", map[string]string{"k": "v"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !createEventCalled {
		t.Fatal("expected CreateEvent to be called")
	}
	if !updateStatusCalled {
		t.Fatal("expected UpdatePaymentStatus to be called")
	}
	if !markPaidCalled {
		t.Fatal("expected MarkOrderPaid to be called")
	}
}

func TestProcessEventFailedMarksOrderCanceled(t *testing.T) {
	t.Parallel()

	paymentID := uuid.New()
	orderID := uuid.New()

	updateStatusCalled := false
	markCanceledCalled := false

	repo := &testRepo{
		getPaymentByOrderIDFn: func(ctx context.Context, reqOrderID uuid.UUID) (*domain.Payment, error) {
			return &domain.Payment{ID: paymentID, OrderID: reqOrderID, Status: domain.PaymentStatusPending}, nil
		},
		createEventFn: func(ctx context.Context, id uuid.UUID, eventType string, payload interface{}) error {
			if eventType != "failed" {
				t.Fatalf("expected event type failed, got %s", eventType)
			}
			return nil
		},
		updatePaymentStatusFn: func(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
			updateStatusCalled = true
			if status != domain.PaymentStatusCanceled {
				t.Fatalf("expected canceled status, got %s", status)
			}
			return nil
		},
	}

	updater := &testOrderUpdater{
		markOrderCanceledFn: func(ctx context.Context, id uuid.UUID) error {
			markCanceledCalled = true
			if id != orderID {
				t.Fatalf("expected order id %s, got %s", orderID, id)
			}
			return nil
		},
	}

	svc := NewPaymentService(repo, &testProvider{}, updater)

	err := svc.ProccessEvent(context.Background(), orderID, "failed", map[string]string{"k": "v"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !updateStatusCalled {
		t.Fatal("expected UpdatePaymentStatus to be called")
	}
	if !markCanceledCalled {
		t.Fatal("expected MarkOrderCanceled to be called")
	}
}

func TestProcessEventAlreadyHandledIsIdempotent(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()
	updateCalled := false
	markCalled := false

	repo := &testRepo{
		getPaymentByOrderIDFn: func(ctx context.Context, reqOrderID uuid.UUID) (*domain.Payment, error) {
			return &domain.Payment{
				ID:      uuid.New(),
				OrderID: reqOrderID,
				Status:  domain.PaymentStatusSucceeded,
			}, nil
		},
		createEventFn: func(ctx context.Context, paymentID uuid.UUID, eventType string, payload interface{}) error {
			return nil
		},
		updatePaymentStatusFn: func(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error {
			updateCalled = true
			return nil
		},
	}

	updater := &testOrderUpdater{
		markOrderPaidFn: func(ctx context.Context, orderID uuid.UUID) error {
			markCalled = true
			return nil
		},
	}

	svc := NewPaymentService(repo, &testProvider{}, updater)

	err := svc.ProccessEvent(context.Background(), orderID, "succeeded", map[string]string{"k": "v"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if updateCalled {
		t.Fatal("expected UpdatePaymentStatus not to be called for already handled payment")
	}
	if markCalled {
		t.Fatal("expected MarkOrderPaid not to be called for already handled payment")
	}
}

func TestProcessEventCreateEventError(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()
	expectedErr := errors.New("write event failed")

	repo := &testRepo{
		getPaymentByOrderIDFn: func(ctx context.Context, reqOrderID uuid.UUID) (*domain.Payment, error) {
			return &domain.Payment{ID: uuid.New(), OrderID: reqOrderID, Status: domain.PaymentStatusPending}, nil
		},
		createEventFn: func(ctx context.Context, paymentID uuid.UUID, eventType string, payload interface{}) error {
			return expectedErr
		},
	}

	svc := NewPaymentService(repo, &testProvider{}, &testOrderUpdater{})

	err := svc.ProccessEvent(context.Background(), orderID, "succeeded", map[string]string{"k": "v"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected create event error, got %v", err)
	}
}

func TestProcessEventPaymentNotFound(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()

	repo := &testRepo{
		getPaymentByOrderIDFn: func(ctx context.Context, reqOrderID uuid.UUID) (*domain.Payment, error) {
			return nil, paymentErrors.ErrPaymentNotFound
		},
	}

	svc := NewPaymentService(repo, &testProvider{}, &testOrderUpdater{})

	err := svc.ProccessEvent(context.Background(), orderID, "succeeded", map[string]string{"k": "v"})
	if !errors.Is(err, paymentErrors.ErrPaymentNotFound) {
		t.Fatalf("expected ErrPaymentNotFound, got %v", err)
	}
}
