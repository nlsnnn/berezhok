package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
)

type testPaymentService struct {
	processFn func(ctx context.Context, orderID uuid.UUID, eventType string, payload interface{}) error
}

func (s *testPaymentService) ProccessEvent(ctx context.Context, orderID uuid.UUID, eventType string, payload interface{}) error {
	if s.processFn != nil {
		return s.processFn(ctx, orderID, eventType, payload)
	}
	return nil
}

func TestWebhookYookassaSuccess(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()
	called := false

	svc := &testPaymentService{
		processFn: func(ctx context.Context, gotOrderID uuid.UUID, eventType string, payload interface{}) error {
			called = true
			if gotOrderID != orderID {
				t.Fatalf("expected order id %s, got %s", orderID, gotOrderID)
			}
			if eventType != "succeeded" {
				t.Fatalf("expected event type succeeded, got %s", eventType)
			}
			return nil
		},
	}

	h := NewWebhookHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body := map[string]interface{}{
		"type":  "notification",
		"event": "payment.succeeded",
		"object": map[string]interface{}{
			"id":     "pay_1",
			"status": "succeeded",
			"metadata": map[string]interface{}{
				"order_id": orderID.String(),
			},
		},
	}

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/yookassa", bytes.NewReader(raw))
	rr := httptest.NewRecorder()

	h.Yookassa(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if !called {
		t.Fatal("expected service ProccessEvent to be called")
	}
}

func TestWebhookYookassaFailedMappedToFailed(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()

	svc := &testPaymentService{
		processFn: func(ctx context.Context, gotOrderID uuid.UUID, eventType string, payload interface{}) error {
			if gotOrderID != orderID {
				t.Fatalf("expected order id %s, got %s", orderID, gotOrderID)
			}
			if eventType != "failed" {
				t.Fatalf("expected event type failed, got %s", eventType)
			}
			return nil
		},
	}

	h := NewWebhookHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body := map[string]interface{}{
		"type":  "notification",
		"event": "payment.failed",
		"object": map[string]interface{}{
			"id":     "pay_2",
			"status": "canceled",
			"metadata": map[string]interface{}{
				"order_id": orderID.String(),
			},
		},
	}

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/yookassa", bytes.NewReader(raw))
	rr := httptest.NewRecorder()

	h.Yookassa(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestWebhookYookassaInvalidJSON(t *testing.T) {
	t.Parallel()

	h := NewWebhookHandler(&testPaymentService{}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/yookassa", bytes.NewBufferString("{"))
	rr := httptest.NewRecorder()

	h.Yookassa(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestWebhookYookassaMissingOrderID(t *testing.T) {
	t.Parallel()

	h := NewWebhookHandler(&testPaymentService{}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body := map[string]interface{}{
		"type":  "notification",
		"event": "payment.succeeded",
		"object": map[string]interface{}{
			"id":       "pay_3",
			"status":   "succeeded",
			"metadata": map[string]interface{}{},
		},
	}

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/yookassa", bytes.NewReader(raw))
	rr := httptest.NewRecorder()

	h.Yookassa(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestWebhookYookassaUnhandledEvent(t *testing.T) {
	t.Parallel()

	h := NewWebhookHandler(&testPaymentService{}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body := map[string]interface{}{
		"type":  "notification",
		"event": "payment.waiting_for_capture",
		"object": map[string]interface{}{
			"id":     "pay_4",
			"status": "waiting_for_capture",
			"metadata": map[string]interface{}{
				"order_id": uuid.New().String(),
			},
		},
	}

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/yookassa", bytes.NewReader(raw))
	rr := httptest.NewRecorder()

	h.Yookassa(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestWebhookYookassaServiceErrorDoesNotLeakMessage(t *testing.T) {
	t.Parallel()

	orderID := uuid.New()

	h := NewWebhookHandler(&testPaymentService{
		processFn: func(ctx context.Context, orderID uuid.UUID, eventType string, payload interface{}) error {
			return errors.New("db connection refused")
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	body := map[string]interface{}{
		"type":  "notification",
		"event": "payment.succeeded",
		"object": map[string]interface{}{
			"id":     "pay_5",
			"status": "succeeded",
			"metadata": map[string]interface{}{
				"order_id": orderID.String(),
			},
		},
	}

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/yookassa", bytes.NewReader(raw))
	rr := httptest.NewRecorder()

	h.Yookassa(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}

	if bytes.Contains(rr.Body.Bytes(), []byte("db connection refused")) {
		t.Fatal("response must not leak internal error details")
	}
}
