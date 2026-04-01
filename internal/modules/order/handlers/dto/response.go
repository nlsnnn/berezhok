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

// OrderDetailResponse represents detailed order information
type OrderDetailResponse struct {
	ID          string                  `json:"id"`
	Status      string                  `json:"status"`
	PickupCode  string                  `json:"pickup_code,omitempty"`
	QRCodeURL   string                  `json:"qr_code_url,omitempty"`
	Amount      float64                 `json:"amount"`
	Box         OrderBoxResponse        `json:"box"`
	Location    OrderLocationResponse   `json:"location"`
	PickupTime  OrderPickupTimeResponse `json:"pickup_time"`
	CreatedAt   time.Time               `json:"created_at"`
	ConfirmedAt *time.Time              `json:"confirmed_at,omitempty"`
}

type OrderBoxResponse struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type OrderCoordinatesResponse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type OrderLocationResponse struct {
	Name        string                   `json:"name"`
	Address     string                   `json:"address"`
	Phone       string                   `json:"phone"`
	Coordinates OrderCoordinatesResponse `json:"coordinates"`
}

type OrderPickupTimeResponse struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type PartnerOrderByCodeResponse struct {
	ID         string                             `json:"id"`
	PickupCode string                             `json:"pickup_code"`
	Status     string                             `json:"status"`
	Box        PartnerOrderByCodeBoxResponse      `json:"box"`
	Customer   PartnerOrderByCodeCustomerResponse `json:"customer"`
	PickupTime OrderPickupTimeResponse            `json:"pickup_time"`
	CreatedAt  time.Time                          `json:"created_at"`
}

type PartnerOrderByCodeBoxResponse struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type PartnerOrderByCodeCustomerResponse struct {
	Phone string `json:"phone"`
	Name  string `json:"name,omitempty"`
}

type PartnerPickupResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
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

// ToOrderDetailResponse converts order details to API response contract for GET /customer/orders/{order_id}
func ToOrderDetailResponse(order *domain.OrderDetails) OrderDetailResponse {
	return OrderDetailResponse{
		ID:         order.ID.String(),
		Status:     string(order.Status),
		PickupCode: order.PickupCode,
		QRCodeURL:  order.QRCodeURL,
		Amount:     order.Amount,
		Box: OrderBoxResponse{
			Name:     order.BoxName,
			ImageURL: order.BoxImageURL,
		},
		Location: OrderLocationResponse{
			Name:    order.LocationName,
			Address: order.LocationAddress,
			Phone:   order.LocationPhone,
			Coordinates: OrderCoordinatesResponse{
				Lat: order.LocationLat,
				Lng: order.LocationLng,
			},
		},
		PickupTime: OrderPickupTimeResponse{
			Start: order.PickupTimeStart,
			End:   order.PickupTimeEnd,
		},
		CreatedAt:   order.CreatedAt,
		ConfirmedAt: order.ConfirmedAt,
	}
}

func ToPartnerOrderByCodeResponse(order *domain.PartnerOrderByCode) PartnerOrderByCodeResponse {
	return PartnerOrderByCodeResponse{
		ID:         order.ID.String(),
		PickupCode: order.PickupCode,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt,
		Box: PartnerOrderByCodeBoxResponse{
			Name:     order.BoxName,
			ImageURL: order.BoxImageURL,
		},
		Customer: PartnerOrderByCodeCustomerResponse{
			Phone: order.CustomerPhone,
			Name:  order.CustomerName,
		},
		PickupTime: OrderPickupTimeResponse{
			Start: order.PickupTimeStart,
			End:   order.PickupTimeEnd,
		},
	}
}

func ToPartnerPickupResponse(orderID, status, message string) PartnerPickupResponse {
	return PartnerPickupResponse{
		OrderID: orderID,
		Status:  status,
		Message: message,
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
func ToCreateOrderResponse(orderID, paymentURL string, amount float64, expiresAt time.Time) CreateOrderResponse {
	return CreateOrderResponse{
		OrderID:    orderID,
		PaymentURL: paymentURL,
		Amount:     amount,
		ExpiresAt:  expiresAt,
	}
}
