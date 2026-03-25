package dto

import (
	"time"

	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
)

// CreateOrderResponse represents response after creating order
type CreateOrderResponse struct {
	OrderID    string    `json:"order_id"`
	PaymentURL string    `json:"payment_url"`
	Amount     float64   `json:"amount"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// OrderResponse represents order details
type OrderResponse struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	PickupCode string    `json:"pickup_code,omitempty"`
	QRCodeURL  string    `json:"qr_code_url,omitempty"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// OrderListResponse wraps order list
type OrderListResponse struct {
	Items []OrderResponse `json:"items"`
	Total int             `json:"total"`
}

// ToOrderResponse converts domain.Order to OrderResponse
func ToOrderResponse(order *domain.Order) OrderResponse {
	return OrderResponse{
		ID:         order.ID.String(),
		Status:     string(order.Status),
		PickupCode: order.PickupCode,
		QRCodeURL:  order.QRCode,
		Amount:     order.Amount().InexactFloat64(),
		CreatedAt:  order.CreatedAt,
	}
}

// ToOrderListResponse converts a slice of domain.Order to OrderListResponse
func ToOrderListResponse(orders []domain.Order) OrderListResponse {
	items := make([]OrderResponse, len(orders))
	for i, order := range orders {
		items[i] = ToOrderResponse(&order)
	}

	return OrderListResponse{
		Items: items,
		Total: len(orders),
	}
}

// ToCreateOrderResponse creates response for order creation
func ToCreateOrderResponse(orderID string, paymentURL string, amount float64, expiresAt time.Time) CreateOrderResponse {
	return CreateOrderResponse{
		OrderID:    orderID,
		PaymentURL: paymentURL,
		Amount:     amount,
		ExpiresAt:  expiresAt,
	}
}
