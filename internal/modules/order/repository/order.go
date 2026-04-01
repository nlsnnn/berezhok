package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

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

// GetOrderDetailsByID retrieves enriched order details by order ID
func (r *OrderRepo) GetOrderDetailsByID(ctx context.Context, orderID uuid.UUID) (*domain.OrderDetails, error) {
	row, err := r.q.GetOrderDetailsByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orderErrors.ErrOrderNotFound
		}
		return nil, err
	}

	lat, _ := row.LocationLat.(float64)
	lng, _ := row.LocationLng.(float64)

	var confirmedAt *time.Time
	if row.PartnerConfirmedAt.Valid {
		value := row.PartnerConfirmedAt.Time
		confirmedAt = &value
	}

	return &domain.OrderDetails{
		ID:              row.ID,
		CustomerID:      row.UserID,
		Status:          domain.OrderStatus(row.Status),
		PickupCode:      row.PickupCode,
		QRCodeURL:       row.QrCodeUrl,
		Amount:          pgconverter.NumericToDecimalOrZero(row.Amount).InexactFloat64(),
		BoxName:         row.BoxName,
		BoxImageURL:     row.BoxImageUrl,
		LocationName:    row.LocationName,
		LocationAddress: row.LocationAddress,
		LocationPhone:   row.LocationPhone,
		LocationLat:     lat,
		LocationLng:     lng,
		PickupTimeStart: row.PickupTimeStart,
		PickupTimeEnd:   row.PickupTimeEnd,
		CreatedAt:       row.CreatedAt,
		ConfirmedAt:     confirmedAt,
	}, nil
}

// GetPartnerOrderByPickupCode retrieves partner-scoped order details by pickup code.
func (r *OrderRepo) GetPartnerOrderByPickupCode(ctx context.Context, pickupCode string, partnerID uuid.UUID) (*domain.PartnerOrderByCode, error) {
	row, err := r.q.GetPartnerOrderByPickupCode(ctx, sqlc.GetPartnerOrderByPickupCodeParams{
		PickupCode: pickupCode,
		PartnerID:  partnerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orderErrors.ErrOrderNotFound
		}

		return nil, err
	}

	return &domain.PartnerOrderByCode{
		ID:              row.ID,
		PickupCode:      row.PickupCode,
		Status:          domain.OrderStatus(row.Status),
		BoxName:         row.BoxName,
		BoxImageURL:     row.BoxImageUrl,
		CustomerPhone:   row.CustomerPhone,
		CustomerName:    row.CustomerName,
		PickupTimeStart: row.PickupTimeStart,
		PickupTimeEnd:   row.PickupTimeEnd,
		CreatedAt:       row.CreatedAt,
	}, nil
}

// MarkOrderPickedUp marks partner-owned order as picked up by employee.
func (r *OrderRepo) MarkOrderPickedUp(ctx context.Context, orderID, partnerID, employeeID uuid.UUID) error {
	partnerOrder, err := r.q.GetPartnerOrderByID(ctx, sqlc.GetPartnerOrderByIDParams{
		ID:        orderID,
		PartnerID: partnerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orderErrors.ErrOrderNotFound
		}

		return err
	}

	if domain.OrderStatus(partnerOrder.Status) != domain.OrderStatusConfirmed {
		return orderErrors.ErrOrderNotReady
	}

	rowsAffected, err := r.q.MarkOrderPickedUp(ctx, sqlc.MarkOrderPickedUpParams{
		ID: orderID,
		PickedUpConfirmedBy: pgtype.UUID{
			Bytes: employeeID,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return orderErrors.ErrOrderNotReady
	}

	return nil
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

// ListOrdersFiltered returns paginated, optionally filtered orders with box/location names
func (r *OrderRepo) ListOrdersFiltered(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) ([]domain.OrderListItem, int, error) {
	sqlItems, err := r.q.ListOrdersByCustomerIDFiltered(ctx, sqlc.ListOrdersByCustomerIDFilteredParams{
		UserID:  customerID,
		Column2: status,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.OrderListItem, len(sqlItems))
	for i, row := range sqlItems {
		items[i] = domain.OrderListItem{
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
