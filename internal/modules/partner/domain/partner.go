package domain

import (
	"time"

	errs "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
)

type PartnerStatus string

const (
	PartnerStatusActive           PartnerStatus = "active"
	PartnerStatusPendingDocuments PartnerStatus = "pending_documents"
	PartnerStatusBlocked          PartnerStatus = "blocked"
)

type Partner struct {
	ID         string
	LegalName  string
	BrandName  string
	LogoURL    string
	Commission Commission
	Status     PartnerStatus
	CreatedAt  time.Time
}

type Commission struct {
	Rate       float64
	ValidUntil *time.Time
}

func NewCommission(rate float64, validUntil *time.Time) (Commission, error) {
	if rate < 0 || rate > 1 {
		return Commission{}, errs.ErrInvalidCommissionRate
	}

	return Commission{
		Rate:       rate,
		ValidUntil: validUntil,
	}, nil
}

func (c Commission) IsPromoActive() bool {
	if c.ValidUntil == nil {
		return false
	}
	return time.Now().Before(*c.ValidUntil)
}
