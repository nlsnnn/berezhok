package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	RoutingKeyEmail = "notification.email"
	RoutingKeySMS   = "notification.sms"
	RoutingKeyPush  = "notification.push"
)

type NotificationMessage struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Recipient string          `json:"recipient"`
	Template  string          `json:"template"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

type NotificationPublisher struct {
	client *Client
}

func NewNotificationPublisher(client *Client) *NotificationPublisher {
	return &NotificationPublisher{client: client}
}

func (p *NotificationPublisher) SendEmail(ctx context.Context, to, template string, payload any) error {
	return p.send(ctx, RoutingKeyEmail, "email", to, template, payload)
}

func (p *NotificationPublisher) SendSMS(ctx context.Context, phone, template string, payload any) error {
	return p.send(ctx, RoutingKeySMS, "sms", phone, template, payload)
}

func (p *NotificationPublisher) SendPush(ctx context.Context, deviceToken, template string, payload any) error {
	return p.send(ctx, RoutingKeyPush, "push", deviceToken, template, payload)
}

func (p *NotificationPublisher) send(ctx context.Context, routingKey, msgType, recipient, template string, payload any) error {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := NotificationMessage{
		ID:        uuid.New().String(),
		Type:      msgType,
		Recipient: recipient,
		Template:  template,
		Payload:   rawPayload,
		CreatedAt: time.Now(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	return p.client.Publish(ctx, routingKey, body)
}
