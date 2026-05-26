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

type ManagedEmployeeResponse struct {
	ID                 string    `json:"id"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Role               string    `json:"role"`
	LocationID         string    `json:"location_id"`
	LocationName       string    `json:"location_name"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
}

type LocationResponse struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Address       string            `json:"address"`
	Phone         string            `json:"phone,omitempty"`
	CategoryCode  string            `json:"category_code,omitempty"`
	Status        string            `json:"status,omitempty"`
	LogoURL       string            `json:"logo_url,omitempty"`
	CoverImageURL string            `json:"cover_image_url,omitempty"`
	WorkingHours  map[string]string `json:"working_hours,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

type LocationPinResponse struct {
	Code   string `json:"code"`
	NameRu string `json:"name_ru"`
}

type PartnerDashboardResponse struct {
	Partner   DashboardPartnerResponse    `json:"partner"`
	Employee  DashboardEmployeeResponse   `json:"employee"`
	Locations []DashboardLocationResponse `json:"locations"`
	Today     DashboardTodayResponse      `json:"today"`
	Week      DashboardWeekResponse       `json:"week"`
	Finance   DashboardFinanceResponse    `json:"finance"`
}

type DashboardPartnerResponse struct {
	ID              string     `json:"id"`
	BrandName       string     `json:"brand_name"`
	Status          string     `json:"status"`
	CommissionRate  float64    `json:"commission_rate"`
	PromoUntil      *time.Time `json:"promo_until,omitempty"`
	HasLegalInfo    bool       `json:"has_legal_info"`
	LegalInfoStatus string     `json:"legal_info_status,omitempty"`
}

type DashboardEmployeeResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type DashboardLocationResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	Status           string `json:"status"`
	ActiveBoxesCount int    `json:"active_boxes_count"`
}

type DashboardTodayResponse struct {
	PendingConfirmation int `json:"pending_confirmation"`
	Confirmed           int `json:"confirmed"`
	PickedUp            int `json:"picked_up"`
	Completed           int `json:"completed"`
}

type DashboardWeekResponse struct {
	OrdersCompleted int     `json:"orders_completed"`
	GrossRevenue    int     `json:"gross_revenue"`
	NetRevenue      int     `json:"net_revenue"`
	AvgRating       float64 `json:"avg_rating"`
}

type DashboardFinanceResponse struct {
	BalancePending int        `json:"balance_pending"`
	NextPayoutDate *time.Time `json:"next_payout_date"`
}

type PartnerStatsResponse struct {
	Summary         PartnerStatsSummaryResponse               `json:"summary"`
	Timeline        []PartnerStatsTimelinePointResponse       `json:"timeline"`
	StatusBreakdown []PartnerStatsStatusBreakdownItemResponse `json:"status_breakdown"`
	TopLocations    []PartnerStatsLocationResponse            `json:"top_locations"`
	TopBoxes        []PartnerStatsBoxResponse                 `json:"top_boxes"`
	Orders          []PartnerStatsOrderResponse               `json:"orders"`
	Meta            PartnerStatsMetaResponse                  `json:"meta"`
}

type PartnerStatsSummaryResponse struct {
	OrdersTotal               int     `json:"orders_total"`
	OrdersCompleted           int     `json:"orders_completed"`
	OrdersCancelled           int     `json:"orders_cancelled"`
	OrdersPendingConfirmation int     `json:"orders_pending_confirmation"`
	GrossRevenue              int     `json:"gross_revenue"`
	NetRevenue                int     `json:"net_revenue"`
	AvgOrderValue             float64 `json:"avg_order_value"`
	AvgRating                 float64 `json:"avg_rating"`
	ReviewsCount              int     `json:"reviews_count"`
}

type PartnerStatsTimelinePointResponse struct {
	Date            time.Time `json:"date"`
	OrdersTotal     int       `json:"orders_total"`
	OrdersCompleted int       `json:"orders_completed"`
	GrossRevenue    int       `json:"gross_revenue"`
	NetRevenue      int       `json:"net_revenue"`
}

type PartnerStatsStatusBreakdownItemResponse struct {
	Status string  `json:"status"`
	Count  int     `json:"count"`
	Share  float64 `json:"share"`
}

type PartnerStatsLocationResponse struct {
	LocationID      string  `json:"location_id"`
	Name            string  `json:"name"`
	Address         string  `json:"address"`
	OrdersTotal     int     `json:"orders_total"`
	OrdersCompleted int     `json:"orders_completed"`
	GrossRevenue    int     `json:"gross_revenue"`
	NetRevenue      int     `json:"net_revenue"`
	AvgRating       float64 `json:"avg_rating"`
}

type PartnerStatsBoxResponse struct {
	BoxID           string `json:"box_id"`
	Name            string `json:"name"`
	ImageURL        string `json:"image_url"`
	LocationName    string `json:"location_name"`
	OrdersTotal     int    `json:"orders_total"`
	OrdersCompleted int    `json:"orders_completed"`
	GrossRevenue    int    `json:"gross_revenue"`
	NetRevenue      int    `json:"net_revenue"`
}

type PartnerStatsOrderResponse struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	PickupCode      string    `json:"pickup_code"`
	Amount          float64   `json:"amount"`
	BoxName         string    `json:"box_name"`
	BoxImageURL     string    `json:"box_image_url"`
	CustomerPhone   string    `json:"customer_phone"`
	CustomerName    string    `json:"customer_name"`
	LocationID      string    `json:"location_id"`
	LocationName    string    `json:"location_name"`
	LocationAddress string    `json:"location_address"`
	PickupTimeStart time.Time `json:"pickup_time_start"`
	PickupTimeEnd   time.Time `json:"pickup_time_end"`
	CreatedAt       time.Time `json:"created_at"`
	CanPickup       bool      `json:"can_pickup"`
}

type PartnerStatsMetaResponse struct {
	Period           string                         `json:"period"`
	DateFrom         time.Time                      `json:"date_from"`
	DateTo           time.Time                      `json:"date_to"`
	LocationID       string                         `json:"location_id,omitempty"`
	Status           string                         `json:"status,omitempty"`
	TopLocationsSort string                         `json:"top_locations_sort"`
	TopBoxesSort     string                         `json:"top_boxes_sort"`
	OrdersSort       string                         `json:"orders_sort"`
	Pagination       PartnerStatsPaginationResponse `json:"pagination"`
}

type PartnerStatsPaginationResponse struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}
