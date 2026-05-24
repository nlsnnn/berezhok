package dto

import (
	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
)

func ToDestinationResponse(d domain.PayoutDestination) DestinationResponse {
	return DestinationResponse{
		Type:          string(d.Type),
		SBPPhone:      d.SBPPhone,
		SBPBankID:     d.SBPBankID,
		RecipientName: d.RecipientName,
		UpdatedAt:     d.UpdatedAt,
	}
}

func ToPayoutResponse(p domain.Payout) PayoutResponse {
	return PayoutResponse{
		ID:                    p.ID.String(),
		PeriodStart:           p.PeriodStart,
		PeriodEnd:             p.PeriodEnd,
		GrossAmount:           p.Gross.StringFixed(2),
		CommissionAmount:      p.Commission.StringFixed(2),
		CommissionRateApplied: p.CommissionRateApplied.StringFixed(4),
		NetAmount:             p.Net.StringFixed(2),
		Status:                string(p.Status),
		Provider:              p.Provider,
		ProviderPayoutID:      p.ProviderPayoutID,
		ErrorMessage:          p.ErrorMessage,
		ProcessedAt:           p.ProcessedAt,
		CreatedAt:             p.CreatedAt,
	}
}

func ToPayoutDetailResponse(p domain.Payout, orders []domain.PayoutOrder) PayoutDetailResponse {
	orderResponses := make([]PayoutOrderResponse, len(orders))
	for i, o := range orders {
		orderResponses[i] = PayoutOrderResponse{
			OrderID:        o.OrderID.String(),
			OrderAmount:    o.OrderAmount.StringFixed(2),
			CommissionPart: o.CommissionPart.StringFixed(2),
		}
	}

	return PayoutDetailResponse{
		PayoutResponse: ToPayoutResponse(p),
		Orders:         orderResponses,
	}
}
