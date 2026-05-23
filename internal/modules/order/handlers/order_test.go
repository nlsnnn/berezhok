package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
	"github.com/nlsnnn/berezhok/internal/modules/order/service"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
)

type orderServiceStub struct {
	listPartnerOrdersByIDFn func(ctx context.Context, actor authz.PartnerActor, status string, limit, offset int) (*service.ListPartnerOrdersResult, error)
}

func (s *orderServiceStub) CreateOrder(ctx context.Context, boxID, customerID uuid.UUID) (*service.CreateOrderResult, error) {
	return nil, nil
}

func (s *orderServiceStub) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	return nil, nil
}

func (s *orderServiceStub) GetOrderDetailsByID(ctx context.Context, orderID uuid.UUID) (*domain.OrderDetails, error) {
	return nil, nil
}

func (s *orderServiceStub) GetOrderProjection(ctx context.Context, orderID uuid.UUID) (*domain.OrderProjection, error) {
	return nil, nil
}

func (s *orderServiceStub) GetPartnerOrderByPickupCode(ctx context.Context, actor authz.PartnerActor, pickupCode string) (*domain.PartnerOrderByCode, error) {
	return nil, nil
}

func (s *orderServiceStub) ListOrdersByPartnerID(ctx context.Context, actor authz.PartnerActor, status string, limit, offset int) (*service.ListPartnerOrdersResult, error) {
	if s.listPartnerOrdersByIDFn != nil {
		return s.listPartnerOrdersByIDFn(ctx, actor, status, limit, offset)
	}

	return &service.ListPartnerOrdersResult{}, nil
}

func (s *orderServiceStub) MarkOrderPickedUp(ctx context.Context, actor authz.PartnerActor, orderID uuid.UUID) error {
	return nil
}

func (s *orderServiceStub) ListOrdersByCustomerID(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) (*service.ListOrdersResult, error) {
	return &service.ListOrdersResult{}, nil
}

func TestListPartnerOrdersSuccess(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()
	pickupStart := time.Date(2026, 4, 20, 18, 0, 0, 0, time.UTC)
	pickupEnd := pickupStart.Add(2 * time.Hour)
	createdAt := pickupStart.Add(-4 * time.Hour)

	h := NewOrderHandler(&orderServiceStub{
		listPartnerOrdersByIDFn: func(ctx context.Context, actor authz.PartnerActor, status string, limit, offset int) (*service.ListPartnerOrdersResult, error) {
			if actor.PartnerID != partnerID {
				t.Fatalf("expected partner id %s, got %s", partnerID, actor.PartnerID)
			}
			if actor.Role != authz.RoleOwner {
				t.Fatalf("expected role owner, got %s", actor.Role)
			}
			if status != "confirmed" {
				t.Fatalf("expected status confirmed, got %s", status)
			}
			if limit != 10 {
				t.Fatalf("expected limit 10, got %d", limit)
			}
			if offset != 20 {
				t.Fatalf("expected offset 20, got %d", offset)
			}

			return &service.ListPartnerOrdersResult{
				Items: []domain.PartnerOrderListItem{{
					ID:              uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					Status:          domain.OrderStatusConfirmed,
					PickupCode:      "AB12CD34",
					BoxName:         "Сюрприз-бокс",
					BoxImageURL:     "https://example.com/box.jpg",
					CustomerPhone:   "+79990001122",
					CustomerName:    "Иван",
					LocationID:      uuid.MustParse("22222222-2222-2222-2222-222222222222"),
					LocationName:    "Кофейня",
					LocationAddress: "Ленина, 1",
					PickupTimeStart: pickupStart,
					PickupTimeEnd:   pickupEnd,
					CreatedAt:       createdAt,
					CanPickup:       true,
				}},
				Total:  1,
				Limit:  10,
				Offset: 20,
			}, nil
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/partner/orders?status=confirmed&limit=10&offset=20", nil)
	ctx := context.WithValue(req.Context(), contextx.PartnerIDKey, partnerID)
	ctx = context.WithValue(ctx, contextx.PartnerActorKey, authz.PartnerActor{
		PartnerID:  partnerID,
		EmployeeID: uuid.New(),
		Role:       authz.RoleOwner,
	})
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	h.ListPartnerOrders(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["success"] != true {
		t.Fatalf("expected success=true, got %v", body["success"])
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", body["data"])
	}

	items, ok := data["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("expected 1 item, got %v", data["items"])
	}

	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("expected item object, got %T", items[0])
	}

	if item["pickup_code"] != "AB12CD34" {
		t.Fatalf("expected pickup code AB12CD34, got %v", item["pickup_code"])
	}

	if item["can_pickup"] != true {
		t.Fatalf("expected can_pickup=true, got %v", item["can_pickup"])
	}
}

func TestListPartnerOrdersUnauthorized(t *testing.T) {
	t.Parallel()

	h := NewOrderHandler(&orderServiceStub{}, slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/partner/orders", nil)
	rr := httptest.NewRecorder()

	h.ListPartnerOrders(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}
}
