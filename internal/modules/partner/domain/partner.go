package domain

import "time"

type PartnerStatus string

const (
	PartnerStatusActive           PartnerStatus = "active"
	PartnerStatusPendingDocuments PartnerStatus = "pending_documents"
	PartnerStatusBlocked          PartnerStatus = "blocked"
)

type Partner struct {
	ID                   string
	LegalName            string
	BrandName            string
	LogoURL              string
	CommissionRate       float64
	PromoCommissionRate  float64
	PromoCommissionUntil *time.Time
	Status               PartnerStatus
	CreatedAt            time.Time
}
