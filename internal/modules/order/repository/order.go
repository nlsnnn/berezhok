package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
	orderErrors "github.com/nlsnnn/berezhok/internal/modules/order/errors"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type OrderRepo struct {
	q *sqlc.Queries
}

func NewOrderRepo(q *sqlc.Queries) *OrderRepo {
	return &OrderRepo{q: q}
}

// CreateOrder creates a new order in the database
func (r *OrderRepo) CreateOrder(ctx context.Context, order *domain.Order) error {
	sqlOrder, err := r.q.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID:                      order.CustomerID,
		BoxID:                       order.BoxID,
		LocationID:                  order.LocationID,
		PickupCode:                  order.PickupCode,
		QrCodeUrl:                   pgconverter.StringToText(order.QRCode),
		Amount:                      pgconverter.DecimalToNumeric(order.Amount(), false),
		PickupTimeStart:             order.PickupTime.Start,
		PickupTimeEnd:               order.PickupTime.End,
		Status:                      sqlc.OrderStatus(order.Status),
		PartnerConfirmationDeadline: order.ConfirmationDeadline,
	})
	if err != nil {
		return err
	}

	// Update domain object with generated ID
	order.ID = sqlOrder.ID
	order.CreatedAt = sqlOrder.CreatedAt
	order.UpdatedAt = sqlOrder.UpdatedAt

	return nil
}

// GetOrderByID retrieves an order by ID
func (r *OrderRepo) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	sqlOrder, err := r.q.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orderErrors.ErrOrderNotFound
		}
		return nil, err
	}

	return r.toDomain(sqlOrder), nil
}

// ListOrdersByCustomerID retrieves all orders for a customer
func (r *OrderRepo) ListOrdersByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Order, error) {
	sqlOrders, err := r.q.ListOrdersByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	orders := make([]domain.Order, len(sqlOrders))
	for i, sqlOrder := range sqlOrders {
		orders[i] = *r.toDomain(sqlOrder)
	}

	return orders, nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	_, err := r.q.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     orderID,
		Status: sqlc.OrderStatus(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orderErrors.ErrOrderNotFound
		}
		return err
	}

	return nil
}

// ReserveBox atomically reserves a box by decrementing quantity_available
func (r *OrderRepo) ReserveBox(ctx context.Context, boxID uuid.UUID) (bool, error) {
	rowsAffected, err := r.q.ReserveBox(ctx, boxID)
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// toDomain converts sqlc.Order to domain.Order
func (r *OrderRepo) toDomain(sqlOrder sqlc.Order) *domain.Order {
	order := &domain.Order{
		ID:                   sqlOrder.ID,
		CustomerID:           sqlOrder.UserID,
		BoxID:                sqlOrder.BoxID,
		LocationID:           sqlOrder.LocationID,
		PickupCode:           sqlOrder.PickupCode,
		QRCode:               pgconverter.TextToString(sqlOrder.QrCodeUrl),
		Status:               domain.OrderStatus(sqlOrder.Status),
		ConfirmationDeadline: sqlOrder.PartnerConfirmationDeadline,
		CreatedAt:            sqlOrder.CreatedAt,
		UpdatedAt:            sqlOrder.UpdatedAt,
	}

	// Set pickup time
	order.PickupTime = sharedDomain.PickupTime{
		Start: sqlOrder.PickupTimeStart,
		End:   sqlOrder.PickupTimeEnd,
	}

	// Set amount (private field through reflection or constructor)
	// Since amount is private, we need to use SetAmount or similar
	// For now, we'll leave it and handle in the constructor
	order.SetAmount(pgconverter.NumericToDecimalOrZero(sqlOrder.Amount))

	// Handle nullable fields
	if sqlOrder.PartnerConfirmedAt.Valid {
		confirmedAt := sqlOrder.PartnerConfirmedAt.Time
		order.ConfirmedAt = &confirmedAt
	}

	if sqlOrder.PartnerConfirmedBy.Valid {
		confirmedByUUID := uuid.UUID(sqlOrder.PartnerConfirmedBy.Bytes)
		order.ConfirmedBy = &confirmedByUUID
	}

	if sqlOrder.CancellationReason.Valid {
		reason := sqlOrder.CancellationReason.String
		order.CancellationReason = &reason
	}

	if sqlOrder.CancelledAt.Valid {
		cancelledAt := sqlOrder.CancelledAt.Time
		order.CancelledAt = &cancelledAt
	}

	if sqlOrder.PickedUpAt.Valid {
		pickedUpAt := sqlOrder.PickedUpAt.Time
		order.PickedUpAt = &pickedUpAt
	}

	if sqlOrder.PickedUpConfirmedBy.Valid {
		pickedUpConfirmedByUUID := uuid.UUID(sqlOrder.PickedUpConfirmedBy.Bytes)
		order.PickedUpConfirmedBy = &pickedUpConfirmedByUUID
	}

	if sqlOrder.UserConfirmedAt.Valid {
		userConfirmedAt := sqlOrder.UserConfirmedAt.Time
		order.UserConfirmedPickupAt = &userConfirmedAt
	}

	if sqlOrder.AutoCompletedAt.Valid {
		autoCompletedAt := sqlOrder.AutoCompletedAt.Time
		order.AutoCompletedAt = &autoCompletedAt
	}

	return order
}

// OrderListItem represents enriched order data for list view
type OrderListItem struct {
	ID              uuid.UUID
	Status          domain.OrderStatus
	PickupCode      string
	Amount          float64
	BoxName         string
	LocationName    string
	PickupTimeStart time.Time
	CreatedAt       time.Time
	HasReview       bool
}

// ListOrdersFiltered returns paginated, optionally filtered orders with box/location names
func (r *OrderRepo) ListOrdersFiltered(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) ([]OrderListItem, int, error) {
	sqlItems, err := r.q.ListOrdersByCustomerIDFiltered(ctx, sqlc.ListOrdersByCustomerIDFilteredParams{
		UserID:  customerID,
		Column2: status,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]OrderListItem, len(sqlItems))
	for i, row := range sqlItems {
		items[i] = OrderListItem{
			ID:              row.ID,
			Status:          domain.OrderStatus(row.Status),
			PickupCode:      row.PickupCode,
			Amount:          pgconverter.NumericToDecimalOrZero(row.Amount).InexactFloat64(),
			BoxName:         row.BoxName,
			LocationName:    row.LocationName,
			PickupTimeStart: row.PickupTimeStart,
			CreatedAt:       row.CreatedAt,
			HasReview:       row.HasReview,
		}
	}

	total, err := r.q.CountOrdersByCustomerID(ctx, sqlc.CountOrdersByCustomerIDParams{
		UserID:  customerID,
		Column2: status,
	})
	if err != nil {
		return nil, 0, err
	}

	return items, int(total), nil
}
