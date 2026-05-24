package dto

import "time"

type SBPBankResponse struct {
	BankID string `json:"bank_id"`
	Name   string `json:"name"`
	BIC    string `json:"bic"`
}

type DestinationResponse struct {
	Type          string    `json:"type"`
	SBPPhone      string    `json:"sbp_phone"`
	SBPBankID     string    `json:"sbp_bank_id"`
	RecipientName string    `json:"recipient_name"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PayoutResponse struct {
	ID                    string     `json:"id"`
	PeriodStart           time.Time  `json:"period_start"`
	PeriodEnd             time.Time  `json:"period_end"`
	GrossAmount           string     `json:"gross_amount"`
	CommissionAmount      string     `json:"commission_amount"`
	CommissionRateApplied string     `json:"commission_rate_applied"`
	NetAmount             string     `json:"net_amount"`
	Status                string     `json:"status"`
	Provider              string     `json:"provider"`
	ProviderPayoutID      *string    `json:"provider_payout_id,omitempty"`
	ErrorMessage          *string    `json:"error_message,omitempty"`
	ProcessedAt           *time.Time `json:"processed_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}

type PayoutOrderResponse struct {
	OrderID        string `json:"order_id"`
	OrderAmount    string `json:"order_amount"`
	CommissionPart string `json:"commission_part"`
}

type PayoutDetailResponse struct {
	PayoutResponse
	Orders []PayoutOrderResponse `json:"orders"`
}

type PayoutsListResponse struct {
	Payouts []PayoutResponse `json:"payouts"`
	Total   int64            `json:"total"`
	Limit   int32            `json:"limit"`
	Offset  int32            `json:"offset"`
	HasMore bool             `json:"has_more"`
}
