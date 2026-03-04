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
