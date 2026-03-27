package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
	paymentErrors "github.com/nlsnnn/berezhok/internal/modules/payment/errors"
)

type PaymentRepo struct {
	q *sqlc.Queries
}

func NewPaymentRepo(q *sqlc.Queries) *PaymentRepo {
	return &PaymentRepo{q: q}
}

func (r *PaymentRepo) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	sqlPayment, err := r.q.CreatePayment(ctx, sqlc.CreatePaymentParams{
		OrderID:           payment.OrderID,
		ProviderPaymentID: pgconverter.StringToText(payment.Provider.PaymentID),
		PaymentUrl:        pgconverter.StringToText(payment.Provider.PaymentLink),
		Method: sqlc.NullPaymentMethod{
			PaymentMethod: sqlc.PaymentMethod(payment.Method),
			Valid:         payment.Method != "",
		},
		Provider: sqlc.NullPaymentProvider{
			PaymentProvider: sqlc.PaymentProvider(payment.Provider.ProviderName),
			Valid:           payment.Provider.ProviderName != "",
		},
		Amount: pgconverter.DecimalToNumeric(payment.Amount, true),
		Status: sqlc.PaymentStatusPending,
	})
	if err != nil {
		return err
	}

	payment.ID = sqlPayment.ID
	payment.CreatedAt = sqlPayment.CreatedAt
	payment.UpdatedAt = sqlPayment.UpdatedAt

	return nil
}

func (r *PaymentRepo) GetPaymentByID(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	sqlPayment, err := r.q.GetPaymentByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, paymentErrors.ErrPaymentNotFound
		}
		return nil, err
	}

	return r.toDomain(sqlPayment), nil
}

func (r *PaymentRepo) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	sqlPayment, err := r.q.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, paymentErrors.ErrPaymentNotFound
		}
		return nil, err
	}

	return r.toDomain(sqlPayment), nil
}

func (r *PaymentRepo) UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error {
	var paidAt pgtype.Timestamptz
	if status == domain.PaymentStatusSucceeded {
		now := time.Now()
		paidAt = pgtype.Timestamptz{Time: now, Valid: true}
	}

	_, err := r.q.UpdatePaymentStatus(ctx, sqlc.UpdatePaymentStatusParams{
		ID:     paymentID,
		Status: sqlc.PaymentStatus(status),
		PaidAt: paidAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return paymentErrors.ErrPaymentNotFound
		}
		return err
	}

	return nil
}

func (r *PaymentRepo) CreateEvent(ctx context.Context, paymentID uuid.UUID, eventType string, payload interface{}) error {
	payloadJSONB, err := pgconverter.InterfaceToJSONB(payload)
	if err != nil {
		return err
	}

	_, err = r.q.CreateEvent(ctx, sqlc.CreateEventParams{
		PaymentID: paymentID,
		EventType: eventType,
		Payload:   payloadJSONB,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return paymentErrors.ErrPaymentNotFound
		}
		return err
	}

	return nil
}
