package dto

type CreateApplicationRequest struct {
	ContactName  string `json:"contact_name" validate:"required,min=2,max=100"`
	ContactEmail string `json:"contact_email" validate:"required,email"`
	ContactPhone string `json:"contact_phone" validate:"required,e164"`
	BusinessName string `json:"business_name" validate:"required,min=2,max=200"`
	CategoryCode string `json:"category_code" validate:"omitempty"`
	Address      string `json:"address" validate:"omitempty"`
	Description  string `json:"description" validate:"omitempty"`
}

type RejectApplicationRequest struct {
	RejectionReason string `json:"rejection_reason" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type CreateLocationRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Address      string  `json:"address" validate:"required,min=5,max=200"`
	CategoryCode string  `json:"category_code" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude    float64 `json:"longitude" validate:"required,min=-180,max=180"`
	Phone        string  `json:"phone" validate:"omitempty,e164"`
}
