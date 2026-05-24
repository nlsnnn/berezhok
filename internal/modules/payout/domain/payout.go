package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
)

type Payout struct {
	ID                    uuid.UUID
	PartnerID             uuid.UUID
	PeriodStart           time.Time
	PeriodEnd             time.Time
	Gross                 decimal.Decimal
	Commission            decimal.Decimal
	CommissionRateApplied decimal.Decimal
	Net                   decimal.Decimal
	Status                PayoutStatus
	Provider              string
	ProviderPayoutID      *string
	IdempotencyKey        string
	ErrorMessage          *string
	ProcessedAt           *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type PayoutOrder struct {
	PayoutID       uuid.UUID
	OrderID        uuid.UUID
	OrderAmount    decimal.Decimal
	CommissionPart decimal.Decimal
}

func NewPayout(
	partnerID uuid.UUID,
	periodStart, periodEnd time.Time,
	orders []OrderForPayout,
	commissionRate decimal.Decimal,
) (*Payout, []PayoutOrder) {
	id := uuid.New()

	var gross decimal.Decimal
	for _, o := range orders {
		gross = gross.Add(o.Amount)
	}

	commission := gross.Mul(commissionRate).RoundBank(2)
	net := gross.Sub(commission)

	payout := &Payout{
		ID:                    id,
		PartnerID:             partnerID,
		PeriodStart:           periodStart,
		PeriodEnd:             periodEnd,
		Gross:                 gross,
		Commission:            commission,
		CommissionRateApplied: commissionRate,
		Net:                   net,
		Status:                PayoutStatusPending,
		Provider:              "yookassa",
		IdempotencyKey:        id.String(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	payoutOrders := make([]PayoutOrder, len(orders))
	for i, o := range orders {
		commPart := o.Amount.Mul(commissionRate).RoundBank(2)
		payoutOrders[i] = PayoutOrder{
			PayoutID:       id,
			OrderID:        o.ID,
			OrderAmount:    o.Amount,
			CommissionPart: commPart,
		}
	}

	return payout, payoutOrders
}

type OrderForPayout struct {
	ID     uuid.UUID
	Amount decimal.Decimal
}
