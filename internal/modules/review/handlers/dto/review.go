package dto

import "time"

type CreateReviewRequest struct {
	OrderID string `json:"order_id" validate:"required,uuid"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"max=500"`
}

type ReviewResponse struct {
	ReviewID  string    `json:"review_id,omitempty"`
	OrderID   string    `json:"order_id,omitempty"`
	ID        string    `json:"id,omitempty"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	UserName  string    `json:"user_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginationResponse struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

type ReviewListResponse struct {
	Items      []ReviewResponse   `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}
