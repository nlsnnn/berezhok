package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
	yoowebhook "github.com/rvinnie/yookassa-sdk-go/yookassa/webhook"

	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type paymentService interface {
	ProccessEvent(ctx context.Context, orderID uuid.UUID, eventType string, payload interface{}) error
}

type webhookHandler struct {
	service   paymentService
	log       *slog.Logger
	validator *validator.Validator
}

func NewWebhookHandler(service paymentService, log *slog.Logger, v *validator.Validator) *webhookHandler {
	return &webhookHandler{
		service:   service,
		log:       log,
		validator: v,
	}
}

type YookassaWebhookRequest struct {
	Event   string `json:"event"`
	Payload struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"payload"`
}

func (h *webhookHandler) Yookassa(w http.ResponseWriter, r *http.Request) {
	const op = "payment.handlers.Yookassa"
	log := h.log.With(slog.String("op", op))

	var webhookEvent yoowebhook.WebhookEvent[yoopayment.Payment]

	err := json.NewDecoder(r.Body).Decode(&webhookEvent)
	if err != nil {
		log.Warn("invalid webhook data", slog.String("error", err.Error()))
		response.BadRequest(w, "Invalid webhook data")
		return
	}

	metadata, ok := webhookEvent.Object.Metadata.(map[string]interface{})
	if !ok {
		log.Warn("invalid metadata format")
		response.BadRequest(w, "Invalid metadata format")
		return
	}

	orderID, ok := metadata["order_id"].(string)
	if !ok {
		log.Warn("missing or invalid order_id in metadata")
		response.BadRequest(w, "Missing order_id")
		return
	}

	orderUID, err := uuid.Parse(orderID)
	if err != nil {
		log.Warn("invalid order_id format")
		response.BadRequest(w, "Invalid order_id format")
		return
	}

	var eventType string
	switch webhookEvent.Event {
	case "payment.succeeded":
		eventType = "succeeded"
	case "payment.failed":
		eventType = "failed"
	default:
		log.Warn("unhandled webhook event type", slog.String("event_type", string(webhookEvent.Type)))
		response.BadRequest(w, "Unhandled webhook event type")
		return
	}

	err = h.service.ProccessEvent(r.Context(), orderUID, eventType, webhookEvent.Object)
	if err != nil {
		log.Error("failed to process webhook event", slog.String("error", err.Error()), slog.String("order_id", orderID), slog.String("payment_id", webhookEvent.Object.ID), slog.String("event_type", string(webhookEvent.Event)))
		response.InternalError(w, nil)
		return
	}

	log.Info("webhook processed", slog.String("type", string(webhookEvent.Type)), slog.String("event", string(webhookEvent.Event)), slog.String("payment_id", webhookEvent.Object.ID), slog.String("status", string(webhookEvent.Object.Status)))

	response.Success(w, "Webhook processed successfully")
}
