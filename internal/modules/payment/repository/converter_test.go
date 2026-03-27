package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
)

// Unit tests for the toDomain converter function

// Tests that the toDomain function correctly maps SQL payment status to domain payment status
func TestToDomainMapsStatus(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	orderID := uuid.New()
	now := time.Now().UTC()

	repo := &PaymentRepo{}
	payment := repo.toDomain(sqlc.Payment{
		ID:      id,
		OrderID: orderID,
		Status:  sqlc.PaymentStatusSucceeded,
		PaidAt: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	if payment.Status != domain.PaymentStatusSucceeded {
		t.Fatalf("expected status %q, got %q", domain.PaymentStatusSucceeded, payment.Status)
	}

	if payment.PaidAt == nil {
		t.Fatal("expected PaidAt to be set")
	}
}
