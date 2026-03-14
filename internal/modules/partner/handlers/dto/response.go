package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
)

type ApplicationResponse struct {
	ID              uuid.UUID  `json:"id"`
	ContactName     string     `json:"contact_name"`
	ContactEmail    string     `json:"contact_email"`
	ContactPhone    string     `json:"contact_phone"`
	BusinessName    string     `json:"business_name"`
	CategoryCode    string     `json:"category_code,omitempty"`
	Address         string     `json:"address,omitempty"`
	Description     string     `json:"description,omitempty"`
	Status          string     `json:"status"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PartnerProfileResponse struct {
	Partner  PartnerResponse   `json:"partner"`
	Employee EmployeeResponse  `json:"employee"`
	Location *LocationResponse `json:"location,omitempty"`
}

type PartnerResponse struct {
	ID             uuid.UUID  `json:"id"`
	LegalName      string     `json:"legal_name"`
	BrandName      string     `json:"brand_name"`
	Status         string     `json:"status"`
	CommissionRate float64    `json:"commission_rate"`
	PromoUntil     *time.Time `json:"promo_until,omitempty"`
}

type EmployeeResponse struct {
	ID                 uuid.UUID `json:"id"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Role               string    `json:"role"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
}

type LocationResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func FromApplication(m sqlc.PartnerApplication) ApplicationResponse {
	var reviewedAt *time.Time
	if m.ReviewedAt.Valid {
		reviewedAt = &m.ReviewedAt.Time
	}

	return ApplicationResponse{
		ID:              m.ID,
		ContactName:     m.ContactName,
		ContactEmail:    m.ContactEmail,
		ContactPhone:    m.ContactPhone,
		BusinessName:    m.BusinessName,
		CategoryCode:    m.CategoryCode.String,
		Address:         m.Address.String,
		Description:     m.Description.String,
		Status:          m.Status,
		ReviewedAt:      reviewedAt,
		RejectionReason: m.RejectionReason.String,
		CreatedAt:       m.CreatedAt,
	}
}
