package dto

// CreateOrderRequest is the request to create an order
type CreateOrderRequest struct {
	BoxID string `json:"box_id" validate:"required,uuid"`
}
