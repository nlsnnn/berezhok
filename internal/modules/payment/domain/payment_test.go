package domain

import (
	"errors"
	"testing"

	paymentErrors "github.com/nlsnnn/berezhok/internal/modules/payment/errors"
)

// Unit tests for Payment domain logic

// TestPaymentSetSuccess tests the SetSuccess method of the Payment struct
func TestPaymentSetSuccess(t *testing.T) {
	t.Parallel()

	payment := &Payment{Status: PaymentStatusPending}

	err := payment.SetSuccess()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if payment.Status != PaymentStatusSucceeded {
		t.Fatalf("expected status %q, got %q", PaymentStatusSucceeded, payment.Status)
	}

	if payment.PaidAt == nil {
		t.Fatal("expected PaidAt to be set")
	}
}

// TestPaymentSetSuccessAlreadyHandled tests that SetSuccess returns an error if the payment is already handled
func TestPaymentSetSuccessAlreadyHandled(t *testing.T) {
	t.Parallel()

	payment := &Payment{Status: PaymentStatusSucceeded}

	err := payment.SetSuccess()
	if !errors.Is(err, paymentErrors.ErrPaymentAlreadyHandled) {
		t.Fatalf("expected ErrPaymentAlreadyHandled, got %v", err)
	}
}

// TestPaymentSetCanceled tests the SetCanceled method of the Payment struct
func TestPaymentSetCanceled(t *testing.T) {
	t.Parallel()

	payment := &Payment{Status: PaymentStatusPending}

	err := payment.SetCanceled()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if payment.Status != PaymentStatusCanceled {
		t.Fatalf("expected status %q, got %q", PaymentStatusCanceled, payment.Status)
	}

	if payment.PaidAt == nil {
		t.Fatal("expected PaidAt to be set")
	}
}

// TestPaymentSetCanceledAlreadyHandled tests that SetCanceled returns an error if the payment is already handled
func TestPaymentIsHandled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status PaymentStatus
		want   bool
	}{
		{name: "pending", status: PaymentStatusPending, want: false},
		{name: "succeeded", status: PaymentStatusSucceeded, want: true},
		{name: "canceled", status: PaymentStatusCanceled, want: true},
		{name: "failed", status: PaymentStatusFailed, want: true},
		{name: "refunded", status: PaymentStatusRefunded, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payment := &Payment{Status: tt.status}
			if got := payment.IsHandled(); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
