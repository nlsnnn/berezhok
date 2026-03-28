package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	catalogDomain "github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
	orderErrors "github.com/nlsnnn/berezhok/internal/modules/order/errors"
	orderRepos "github.com/nlsnnn/berezhok/internal/modules/order/repository"
	"github.com/shopspring/decimal"
)

type orderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
	ListOrdersFiltered(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) ([]orderRepos.OrderListItem, int, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error
	ReserveBox(ctx context.Context, boxID uuid.UUID) (bool, error)
}

type paymentProvider interface {
	Create(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error)
}

type boxProvider interface {
	GetBoxByID(ctx context.Context, id string) (*catalogDomain.SurpriseBox, error)
}

type orderService struct {
	repo            orderRepository
	paymentProvider paymentProvider
	boxProvider     boxProvider
	log             *slog.Logger
}

type CreateOrderResult struct {
	PaymentLink string
	OrderID     uuid.UUID
}

// ListOrdersResult contains paginated order list data
type ListOrdersResult struct {
	Items  []orderRepos.OrderListItem
	Total  int
	Limit  int
	Offset int
}

func NewOrderService(repo orderRepository, boxProvider boxProvider, paymentProvider paymentProvider, log *slog.Logger) *orderService {
	return &orderService{
		repo:            repo,
		boxProvider:     boxProvider,
		paymentProvider: paymentProvider,
		log:             log,
	}
}

// CreateOrder creates a new order with box reservation
func (s *orderService) CreateOrder(ctx context.Context, boxID uuid.UUID, customerID uuid.UUID) (*CreateOrderResult, error) {
	const op = "order.service.CreateOrder"
	log := s.log.With(slog.String("op", op))

	// 1. Validate box exists and is active
	box, err := s.boxProvider.GetBoxByID(ctx, boxID.String())
	if err != nil {
		if errors.Is(err, catalogErrors.ErrBoxNotFound) {
			return nil, orderErrors.ErrBoxNotAvailable
		}
		return nil, fmt.Errorf("%s: failed to get box: %w", op, err)
	}

	// 2. Validate box status
	if !box.IsAvailable() {
		log.Warn("box is not available", slog.String("box_id", boxID.String()), slog.String("status", string(box.Status)))
		return nil, orderErrors.ErrBoxNotAvailable
	}

	// 3. Atomically reserve the box
	// TODO: Move reservation logic to a catalog service
	reserved, err := s.repo.ReserveBox(ctx, boxID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to reserve box: %w", op, err)
	}

	if !reserved {
		log.Warn("failed to reserve box - no rows affected", slog.String("box_id", boxID.String()))
		return nil, orderErrors.ErrBoxNotAvailable
	}

	// 5. Create the order
	order := domain.NewOrder(customerID, boxID, box.LocationID, box.PickupTime, box.Price.Discount)

	err = s.repo.CreateOrder(ctx, order)
	if err != nil {
		// TODO: Consider implementing rollback logic to unreserve the box
		return nil, fmt.Errorf("%s: failed to create order: %w", op, err)
	}

	log.Info("order created successfully",
		slog.String("order_id", order.ID.String()),
		slog.String("customer_id", customerID.String()),
		slog.String("box_id", boxID.String()),
	)

	// 6. Create payment link
	paymentLink, err := s.paymentProvider.Create(ctx, box.Price.Discount, order.ID)
	if err != nil {
		log.Error("failed to create payment link", slog.String("order_id", order.ID.String()), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, orderErrors.ErrPaymentFailed)
	}

	return &CreateOrderResult{
		PaymentLink: paymentLink,
		OrderID:     order.ID,
	}, nil
}

// GetOrderByID retrieves an order by its ID
func (s *orderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	const op = "order.service.GetOrderByID"

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

// ListOrdersByCustomerID retrieves filtered, paginated orders for a customer
func (s *orderService) ListOrdersByCustomerID(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) (*ListOrdersResult, error) {
	const op = "order.service.ListOrdersByCustomerID"

	items, total, err := s.repo.ListOrdersFiltered(ctx, customerID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ListOrdersResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// UpdateOrderStatus updates the status of an order (called by payment webhook handler)
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	const op = "order.service.UpdateOrderStatus"
	log := s.log.With(slog.String("op", op))

	// TODO: Add validation for valid status transitions
	// For example: pending -> paid, paid -> confirmed, etc.

	err := s.repo.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("order status updated",
		slog.String("order_id", orderID.String()),
		slog.String("new_status", string(status)),
	)

	return nil
}
