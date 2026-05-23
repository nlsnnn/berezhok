package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type orderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	CreateOrGetActiveOrder(ctx context.Context, order *domain.Order) (bool, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
	GetOrderDetailsByID(ctx context.Context, orderID uuid.UUID) (*domain.OrderDetails, error)
	GetOrderProjection(ctx context.Context, orderID uuid.UUID) (*domain.OrderProjection, error)
	GetPartnerOrderByPickupCode(ctx context.Context, pickupCode string, partnerID uuid.UUID) (*domain.PartnerOrderByCode, error)
	GetLocationOrderByPickupCode(ctx context.Context, pickupCode string, locationID uuid.UUID) (*domain.PartnerOrderByCode, error)
	ListOrdersByPartnerID(ctx context.Context, partnerID uuid.UUID, status string, limit, offset int) ([]domain.PartnerOrderListItem, int, error)
	ListActiveOrdersByLocationID(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]domain.PartnerOrderListItem, int, error)
	MarkOrderPickedUp(ctx context.Context, orderID, partnerID, employeeID uuid.UUID) error
	MarkLocationOrderPickedUp(ctx context.Context, orderID, locationID, employeeID uuid.UUID) error
	ListOrdersFiltered(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) ([]domain.OrderListItem, int, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error
	ReserveBox(ctx context.Context, boxID uuid.UUID) (bool, error)
}

type orderProjectionPublisher interface {
	PublishOrderCreated(ctx context.Context, projection domain.OrderProjection) error
	PublishOrderStatusChanged(ctx context.Context, projection domain.OrderProjection) error
}

type paymentProvider interface {
	EnsurePaymentLink(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error)
}

type boxProvider interface {
	GetBoxForOrder(ctx context.Context, boxID uuid.UUID) (*BoxForOrder, error)
}

type BoxForOrder struct {
	LocationID uuid.UUID
	Amount     decimal.Decimal
	PickupTime sharedDomain.PickupTime
	Available  bool
}
