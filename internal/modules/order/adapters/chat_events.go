package adapters

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/adapters/rabbitmq"
	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
)

type ChatEventAdapter struct {
	publisher *rabbitmq.OrderPublisher
}

func NewChatEventAdapter(publisher *rabbitmq.OrderPublisher) *ChatEventAdapter {
	return &ChatEventAdapter{publisher: publisher}
}

func (a *ChatEventAdapter) PublishOrderCreated(ctx context.Context, projection domain.OrderChatProjection) error {
	return a.publisher.PublishProjection(ctx, rabbitmq.RoutingKeyOrderCreated, "order.created", toMessage(projection))
}

func (a *ChatEventAdapter) PublishOrderStatusChanged(ctx context.Context, projection domain.OrderChatProjection) error {
	return a.publisher.PublishProjection(ctx, rabbitmq.RoutingKeyOrderUpdated, "order.status_changed", toMessage(projection))
}

func toMessage(projection domain.OrderChatProjection) rabbitmq.OrderProjectionMessage {
	return rabbitmq.OrderProjectionMessage{
		OrderID:    projection.OrderID,
		CustomerID: projection.CustomerID,
		PartnerID:  projection.PartnerID,
		LocationID: projection.LocationID,
		Status:     string(projection.Status),
		UpdatedAt:  projection.UpdatedAt,
	}
}
