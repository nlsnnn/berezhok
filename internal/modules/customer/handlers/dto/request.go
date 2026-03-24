package dto

// UpdateProfileRequest is the request to update customer profile
type UpdateProfileRequest struct {
	Name string `json:"name" validate:"max=100"`
}

// CreateOrderRequest is the request to create an order
type CreateOrderRequest struct {
	BoxID string `json:"box_id" validate:"required,uuid"`
}

// CreateReviewRequest is the request to create a review
type CreateReviewRequest struct {
	OrderID string `json:"order_id" validate:"required,uuid"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"max=500"`
}

// CreateDisputeRequest is the request to open a dispute
type CreateDisputeRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}
