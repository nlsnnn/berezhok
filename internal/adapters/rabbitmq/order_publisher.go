package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	ExchangeOrders         = "orders"
	RoutingKeyOrderCreated = "order.created"
	RoutingKeyOrderUpdated = "order.status_changed"
)

type OrderProjectionMessage struct {
	EventID    uuid.UUID `json:"event_id"`
	EventType  string    `json:"event_type"`
	OrderID    uuid.UUID `json:"order_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	PartnerID  uuid.UUID `json:"partner_id"`
	LocationID uuid.UUID `json:"location_id"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderPublisher struct {
	client *Client
}

func NewOrderPublisher(client *Client) *OrderPublisher {
	return &OrderPublisher{client: client}
}

func (p *OrderPublisher) PublishProjection(ctx context.Context, routingKey, eventType string, msg OrderProjectionMessage) error {
	msg.EventID = uuid.New()
	msg.EventType = eventType

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal order projection: %w", err)
	}
	return p.client.PublishToExchange(ctx, ExchangeOrders, routingKey, body)
}
