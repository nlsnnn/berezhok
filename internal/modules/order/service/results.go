package service

import (
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
)

type CreateOrderResult struct {
	PaymentLink string
	OrderID     uuid.UUID
}

// ListOrdersResult contains paginated order list data.
type ListOrdersResult struct {
	Items  []domain.OrderListItem
	Total  int
	Limit  int
	Offset int
}
