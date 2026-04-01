package dto

import "time"

type ApplicationResponse struct {
	ID              string     `json:"id"`
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
	Partner   PartnerResponse    `json:"partner"`
	Employee  EmployeeResponse   `json:"employee"`
	Location  *LocationResponse  `json:"location,omitempty"` // Employee's assigned location (backwards compatibility)
	Locations []LocationResponse `json:"locations"`          // All partner locations
}

type PartnerResponse struct {
	ID             string     `json:"id"`
	BrandName      string     `json:"brand_name"`
	Status         string     `json:"status"`
	CommissionRate float64    `json:"commission_rate"`
	PromoUntil     *time.Time `json:"promo_until,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type EmployeeResponse struct {
	ID                 string    `json:"id"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Role               string    `json:"role"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
}

type LocationResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
