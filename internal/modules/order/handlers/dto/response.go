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
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	PickupCode      string    `json:"pickup_code,omitempty"`
	QRCodeURL       string    `json:"qr_code_url,omitempty"`
	Amount          float64   `json:"amount"`
	BoxName         string    `json:"box_name,omitempty"`
	LocationName    string    `json:"location_name,omitempty"`
	PickupTimeStart time.Time `json:"pickup_time_start,omitempty"`
	HasReview       bool      `json:"has_review"`
	CreatedAt       time.Time `json:"created_at"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

// OrderListResponse wraps order list with pagination
type OrderListResponse struct {
	Items      []OrderListItem    `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
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

// OrderListItem represents enriched order data for list view
type OrderListItem struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	PickupCode      string    `json:"pickup_code"`
	Amount          float64   `json:"amount"`
	BoxName         string    `json:"box_name"`
	LocationName    string    `json:"location_name"`
	PickupTimeStart time.Time `json:"pickup_time_start"`
	CreatedAt       time.Time `json:"created_at"`
	HasReview       bool      `json:"has_review"`
}

// ToOrderListItem converts enriched order data to OrderListItem
func ToOrderListItem(id, status, pickupCode string, amount float64, boxName, locationName string, pickupTimeStart, createdAt time.Time, hasReview bool) OrderListItem {
	return OrderListItem{
		ID:              id,
		Status:          status,
		PickupCode:      pickupCode,
		Amount:          amount,
		BoxName:         boxName,
		LocationName:    locationName,
		PickupTimeStart: pickupTimeStart,
		CreatedAt:       createdAt,
		HasReview:       hasReview,
	}
}

// ToOrderListResponse converts items to paginated OrderListResponse
func ToOrderListResponse(items []OrderListItem, total, limit, offset int) OrderListResponse {
	return OrderListResponse{
		Items: items,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: offset+limit < total,
		},
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
