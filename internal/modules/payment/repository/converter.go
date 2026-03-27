package repository

import (
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
)

func (r *PaymentRepo) toDomain(sqlPayment sqlc.Payment) *domain.Payment {
	payment := &domain.Payment{
		ID:        sqlPayment.ID,
		OrderID:   sqlPayment.OrderID,
		Amount:    pgconverter.NumericToDecimalOrZero(sqlPayment.Amount),
		Status:    domain.PaymentStatus(sqlPayment.Status),
		CreatedAt: sqlPayment.CreatedAt,
		UpdatedAt: sqlPayment.UpdatedAt,
	}

	if sqlPayment.Method.Valid {
		payment.Method = string(sqlPayment.Method.PaymentMethod)
	}

	if sqlPayment.Provider.Valid {
		payment.Provider.ProviderName = domain.ProviderName(sqlPayment.Provider.PaymentProvider)
	}

	payment.Provider.PaymentID = pgconverter.TextToString(sqlPayment.ProviderPaymentID)
	payment.Provider.PaymentLink = pgconverter.TextToString(sqlPayment.PaymentUrl)

	if sqlPayment.PaidAt.Valid {
		paidAt := sqlPayment.PaidAt.Time
		payment.PaidAt = &paidAt
	}

	return payment
}

// func (r *PaymentRepo) toDomainList(sqlPayments []sqlc.Payment) []*domain.Payment {
// 	payments := make([]*domain.Payment, len(sqlPayments))
// 	for i, sqlPayment := range sqlPayments {
// 		payments[i] = r.toDomain(sqlPayment)
// 	}
// 	return payments
// }
