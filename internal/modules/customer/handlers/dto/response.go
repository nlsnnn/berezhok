package dto

import "time"

// ProfileResponse represents customer profile
type ProfileResponse struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// LocationSearchResponse represents location in search results
type LocationSearchResponse struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Category         CategoryResponse    `json:"category"`
	Address          string              `json:"address"`
	Distance         *float64            `json:"distance,omitempty"` // in meters, optional for now
	Coordinates      CoordinatesResponse `json:"coordinates"`
	Rating           *RatingResponse     `json:"rating,omitempty"`
	LogoURL          string              `json:"logo_url,omitempty"`
	ActiveBoxesCount int                 `json:"active_boxes_count"`
}

// LocationDetailsResponse represents full location details
type LocationDetailsResponse struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Category      CategoryResponse    `json:"category"`
	Address       string              `json:"address"`
	Coordinates   CoordinatesResponse `json:"coordinates"`
	Phone         string              `json:"phone,omitempty"`
	WorkingHours  map[string]string   `json:"working_hours,omitempty"`
	LogoURL       string              `json:"logo_url,omitempty"`
	CoverImageURL string              `json:"cover_image_url,omitempty"`
	Gallery       []string            `json:"gallery,omitempty"`
	Rating        *RatingResponse     `json:"rating,omitempty"`
	ActiveBoxes   []BoxResponse       `json:"active_boxes"`
}

// CategoryResponse represents location category
type CategoryResponse struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	IconURL string `json:"icon_url,omitempty"`
}

// CoordinatesResponse represents geographic coordinates
type CoordinatesResponse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// RatingResponse represents location rating (stub for now)
type RatingResponse struct {
	Average      float64             `json:"average"`
	TotalReviews int                 `json:"total_reviews"`
	Distribution *RatingDistribution `json:"distribution,omitempty"`
}

// RatingDistribution represents rating breakdown
type RatingDistribution struct {
	Five  int `json:"5"`
	Four  int `json:"4"`
	Three int `json:"3"`
	Two   int `json:"2"`
	One   int `json:"1"`
}

// BoxResponse represents a surprise box
type BoxResponse struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	OriginalPrice     float64            `json:"original_price"`
	DiscountPrice     float64            `json:"discount_price"`
	QuantityAvailable int                `json:"quantity_available"`
	PickupTime        PickupTimeResponse `json:"pickup_time"`
	ImageURL          string             `json:"image_url,omitempty"`
}

// PickupTimeResponse represents pickup time window
type PickupTimeResponse struct {
	Start string `json:"start"` // HH:MM format
	End   string `json:"end"`   // HH:MM format
}

// OrderResponse represents order details (stub)
type OrderResponse struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	PickupCode string    `json:"pickup_code,omitempty"`
	QRCodeURL  string    `json:"qr_code_url,omitempty"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateOrderResponse represents response after creating order
type CreateOrderResponse struct {
	OrderID    string    `json:"order_id"`
	PaymentURL string    `json:"payment_url"`
	Amount     float64   `json:"amount"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// ReviewResponse represents a review
type ReviewResponse struct {
	ID        string    `json:"id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

// LocationSearchResultResponse wraps search results with pagination
type LocationSearchResultResponse struct {
	Items      []LocationSearchResponse `json:"items"`
	Pagination PaginationResponse       `json:"pagination"`
}

// OrderListResponse wraps order list with pagination
type OrderListResponse struct {
	Items      []OrderResponse    `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

// ReviewListResponse wraps review list with pagination
type ReviewListResponse struct {
	Items      []ReviewResponse   `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}
