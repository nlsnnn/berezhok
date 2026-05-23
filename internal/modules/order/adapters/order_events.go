package adapters

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/adapters/rabbitmq"
	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
)

type OrderEventAdapter struct {
	publisher *rabbitmq.OrderPublisher
}

func NewOrderEventAdapter(publisher *rabbitmq.OrderPublisher) *OrderEventAdapter {
	return &OrderEventAdapter{publisher: publisher}
}

func (a *OrderEventAdapter) PublishOrderCreated(ctx context.Context, projection domain.OrderProjection) error {
	return a.publisher.PublishProjection(ctx, rabbitmq.RoutingKeyOrderCreated, "order.created", toMessage(projection))
}

func (a *OrderEventAdapter) PublishOrderStatusChanged(ctx context.Context, projection domain.OrderProjection) error {
	return a.publisher.PublishProjection(ctx, rabbitmq.RoutingKeyOrderUpdated, "order.status_changed", toMessage(projection))
}

func toMessage(projection domain.OrderProjection) rabbitmq.OrderProjectionMessage {
	return rabbitmq.OrderProjectionMessage{
		OrderID:    projection.OrderID,
		CustomerID: projection.CustomerID,
		PartnerID:  projection.PartnerID,
		LocationID: projection.LocationID,
		Status:     string(projection.Status),
		UpdatedAt:  projection.UpdatedAt,
	}
}
