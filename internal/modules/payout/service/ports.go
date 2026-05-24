package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PayoutProvider interface {
	SendSBP(ctx context.Context, in SendSBPInput) (SendSBPResult, error)
	GetPayoutStatus(ctx context.Context, providerPayoutID string) (SendSBPResult, error)
}

type BanksProvider interface {
	GetSBPBanks(ctx context.Context) ([]SBPBank, error)
}

type SBPBank struct {
	BankID string
	Name   string
	BIC    string
}

// ProcessingPayout is a lightweight projection used by PollProcessing.
type ProcessingPayout struct {
	ID               uuid.UUID
	ProviderPayoutID string
}

type SendSBPInput struct {
	Amount         decimal.Decimal
	PhoneE164      string
	BankID         string
	RecipientName  string
	Description    string
	IdempotencyKey string
	Metadata       map[string]string
}

type SendSBPResult struct {
	ProviderPayoutID string
	Status           string
}
