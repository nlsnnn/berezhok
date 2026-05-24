package domain

import (
	"time"

	"github.com/google/uuid"
)

type DestinationType string

const (
	DestinationTypeSBP DestinationType = "sbp"
)

type PayoutDestination struct {
	PartnerID     uuid.UUID
	Type          DestinationType
	SBPPhone      string
	SBPBankID     string
	RecipientName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
