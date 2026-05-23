package service

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type orderRepoStub struct {
	listByPartnerCalled bool
	listByLocationFn    func(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]domain.PartnerOrderListItem, int, error)
	markPartnerCalled   bool
	markLocationFn      func(ctx context.Context, orderID, locationID, employeeID uuid.UUID) error
	createOrGetFn       func(ctx context.Context, order *domain.Order) (bool, error)
	createOrderFn       func(ctx context.Context, order *domain.Order) error
	reserveBoxFn        func(ctx context.Context, boxID uuid.UUID) (bool, error)
	getProjectionFn     func(ctx context.Context, orderID uuid.UUID) (*domain.OrderProjection, error)
}

func (r *orderRepoStub) CreateOrder(ctx context.Context, order *domain.Order) error {
	if r.createOrderFn != nil {
		return r.createOrderFn(ctx, order)
	}
	return nil
}

func (r *orderRepoStub) CreateOrGetActiveOrder(ctx context.Context, order *domain.Order) (bool, error) {
	if r.createOrGetFn != nil {
		return r.createOrGetFn(ctx, order)
	}
	return true, nil
}

func (r *orderRepoStub) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	return nil, nil
}

func (r *orderRepoStub) GetOrderDetailsByID(ctx context.Context, orderID uuid.UUID) (*domain.OrderDetails, error) {
	return nil, nil
}

func (r *orderRepoStub) GetOrderProjection(ctx context.Context, orderID uuid.UUID) (*domain.OrderProjection, error) {
	if r.getProjectionFn != nil {
		return r.getProjectionFn(ctx, orderID)
	}
	return nil, nil
}

func (r *orderRepoStub) GetPartnerOrderByPickupCode(ctx context.Context, pickupCode string, partnerID uuid.UUID) (*domain.PartnerOrderByCode, error) {
	return nil, nil
}

func (r *orderRepoStub) GetLocationOrderByPickupCode(ctx context.Context, pickupCode string, locationID uuid.UUID) (*domain.PartnerOrderByCode, error) {
	return nil, nil
}

func (r *orderRepoStub) ListOrdersByPartnerID(ctx context.Context, partnerID uuid.UUID, status string, limit, offset int) ([]domain.PartnerOrderListItem, int, error) {
	r.listByPartnerCalled = true
	return nil, 0, nil
}

func (r *orderRepoStub) ListActiveOrdersByLocationID(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]domain.PartnerOrderListItem, int, error) {
	if r.listByLocationFn != nil {
		return r.listByLocationFn(ctx, locationID, limit, offset)
	}
	return nil, 0, nil
}

func (r *orderRepoStub) MarkOrderPickedUp(ctx context.Context, orderID, partnerID, employeeID uuid.UUID) error {
	r.markPartnerCalled = true
	return nil
}

func (r *orderRepoStub) MarkLocationOrderPickedUp(ctx context.Context, orderID, locationID, employeeID uuid.UUID) error {
	if r.markLocationFn != nil {
		return r.markLocationFn(ctx, orderID, locationID, employeeID)
	}
	return nil
}

func (r *orderRepoStub) ListOrdersFiltered(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) ([]domain.OrderListItem, int, error) {
	return nil, 0, nil
}

func (r *orderRepoStub) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	return nil
}

func (r *orderRepoStub) ReserveBox(ctx context.Context, boxID uuid.UUID) (bool, error) {
	if r.reserveBoxFn != nil {
		return r.reserveBoxFn(ctx, boxID)
	}
	return true, nil
}

type boxProviderStub struct {
	getBoxForOrderFn func(ctx context.Context, boxID uuid.UUID) (*BoxForOrder, error)
}

func (p *boxProviderStub) GetBoxForOrder(ctx context.Context, boxID uuid.UUID) (*BoxForOrder, error) {
	if p.getBoxForOrderFn != nil {
		return p.getBoxForOrderFn(ctx, boxID)
	}
	return nil, nil
}

type paymentProviderStub struct {
	ensurePaymentLinkFn func(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error)
}

func (p *paymentProviderStub) EnsurePaymentLink(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error) {
	if p.ensurePaymentLinkFn != nil {
		return p.ensurePaymentLinkFn(ctx, amount, orderID)
	}
	return "", nil
}

type projectionPublisherStub struct {
	createdCount int
	changedCount int
}

func (p *projectionPublisherStub) PublishOrderCreated(ctx context.Context, projection domain.OrderProjection) error {
	p.createdCount++
	return nil
}

func (p *projectionPublisherStub) PublishOrderStatusChanged(ctx context.Context, projection domain.OrderProjection) error {
	p.changedCount++
	return nil
}

func testLogger(t *testing.T) *slog.Logger {
	t.Helper()
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestCreateOrderCreatesNewOrderAndPublishesEvent(t *testing.T) {
	t.Parallel()

	boxID := uuid.New()
	customerID := uuid.New()
	locationID := uuid.New()
	orderID := uuid.New()
	pickupStart := time.Date(2026, 5, 22, 18, 0, 0, 0, time.UTC)
	pickupEnd := pickupStart.Add(2 * time.Hour)
	publisher := &projectionPublisherStub{}

	repo := &orderRepoStub{
		createOrGetFn: func(ctx context.Context, order *domain.Order) (bool, error) {
			if order.CustomerID != customerID {
				t.Fatalf("expected customer %s, got %s", customerID, order.CustomerID)
			}
			if order.BoxID != boxID {
				t.Fatalf("expected box %s, got %s", boxID, order.BoxID)
			}
			order.ID = orderID
			return true, nil
		},
		reserveBoxFn: func(ctx context.Context, gotBoxID uuid.UUID) (bool, error) {
			t.Fatal("CreateOrder must reserve boxes through CreateOrGetActiveOrder")
			return false, nil
		},
		createOrderFn: func(ctx context.Context, order *domain.Order) error {
			t.Fatal("CreateOrder must create orders through CreateOrGetActiveOrder")
			return nil
		},
		getProjectionFn: func(ctx context.Context, gotOrderID uuid.UUID) (*domain.OrderProjection, error) {
			if gotOrderID != orderID {
				t.Fatalf("expected projection order %s, got %s", orderID, gotOrderID)
			}
			return &domain.OrderProjection{OrderID: gotOrderID}, nil
		},
	}
	payments := &paymentProviderStub{
		ensurePaymentLinkFn: func(ctx context.Context, amount decimal.Decimal, gotOrderID uuid.UUID) (string, error) {
			if gotOrderID != orderID {
				t.Fatalf("expected payment order %s, got %s", orderID, gotOrderID)
			}
			if !amount.Equal(decimal.RequireFromString("250")) {
				t.Fatalf("expected amount 250, got %s", amount)
			}
			return "https://pay.example/new", nil
		},
	}
	boxes := &boxProviderStub{
		getBoxForOrderFn: func(ctx context.Context, gotBoxID uuid.UUID) (*BoxForOrder, error) {
			if gotBoxID != boxID {
				t.Fatalf("expected box id %s, got %s", boxID, gotBoxID)
			}
			return &BoxForOrder{
				LocationID: locationID,
				Amount:     decimal.RequireFromString("250"),
				PickupTime: sharedDomain.PickupTime{Start: pickupStart, End: pickupEnd},
				Available:  true,
			}, nil
		},
	}

	svc := NewOrderService(repo, boxes, payments, testLogger(t), publisher)

	result, err := svc.CreateOrder(context.Background(), boxID, customerID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.OrderID != orderID {
		t.Fatalf("expected order %s, got %s", orderID, result.OrderID)
	}
	if result.PaymentLink != "https://pay.example/new" {
		t.Fatalf("expected new payment link, got %s", result.PaymentLink)
	}
	if publisher.createdCount != 1 {
		t.Fatalf("expected one created event, got %d", publisher.createdCount)
	}
}

func TestCreateOrderReturnsExistingActiveOrderWithoutRepublishing(t *testing.T) {
	t.Parallel()

	boxID := uuid.New()
	customerID := uuid.New()
	locationID := uuid.New()
	existingOrderID := uuid.New()
	pickupStart := time.Date(2026, 5, 22, 18, 0, 0, 0, time.UTC)
	pickupEnd := pickupStart.Add(2 * time.Hour)
	publisher := &projectionPublisherStub{}

	repo := &orderRepoStub{
		createOrGetFn: func(ctx context.Context, order *domain.Order) (bool, error) {
			order.ID = existingOrderID
			return false, nil
		},
		reserveBoxFn: func(ctx context.Context, gotBoxID uuid.UUID) (bool, error) {
			t.Fatal("duplicate order must not reserve another box")
			return false, nil
		},
		createOrderFn: func(ctx context.Context, order *domain.Order) error {
			t.Fatal("duplicate order must not insert another order")
			return nil
		},
		getProjectionFn: func(ctx context.Context, orderID uuid.UUID) (*domain.OrderProjection, error) {
			t.Fatal("duplicate order must not publish order.created")
			return nil, nil
		},
	}
	payments := &paymentProviderStub{
		ensurePaymentLinkFn: func(ctx context.Context, amount decimal.Decimal, orderID uuid.UUID) (string, error) {
			if orderID != existingOrderID {
				t.Fatalf("expected existing order %s, got %s", existingOrderID, orderID)
			}
			return "https://pay.example/existing", nil
		},
	}
	boxes := &boxProviderStub{
		getBoxForOrderFn: func(ctx context.Context, gotBoxID uuid.UUID) (*BoxForOrder, error) {
			if gotBoxID != boxID {
				t.Fatalf("expected box id %s, got %s", boxID, gotBoxID)
			}
			return &BoxForOrder{
				LocationID: locationID,
				Amount:     decimal.RequireFromString("250"),
				PickupTime: sharedDomain.PickupTime{Start: pickupStart, End: pickupEnd},
				Available:  true,
			}, nil
		},
	}

	svc := NewOrderService(repo, boxes, payments, testLogger(t), publisher)

	result, err := svc.CreateOrder(context.Background(), boxID, customerID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result.OrderID != existingOrderID {
		t.Fatalf("expected existing order %s, got %s", existingOrderID, result.OrderID)
	}
	if result.PaymentLink != "https://pay.example/existing" {
		t.Fatalf("expected existing payment link, got %s", result.PaymentLink)
	}
	if publisher.createdCount != 0 {
		t.Fatalf("expected no created event for duplicate order, got %d", publisher.createdCount)
	}
}

func TestListOrdersScopesManagerToAssignedLocation(t *testing.T) {
	t.Parallel()

	locationID := uuid.New()
	repo := &orderRepoStub{
		listByLocationFn: func(ctx context.Context, gotLocationID uuid.UUID, limit, offset int) ([]domain.PartnerOrderListItem, int, error) {
			if gotLocationID != locationID {
				t.Fatalf("expected location %s, got %s", locationID, gotLocationID)
			}
			return nil, 0, nil
		},
	}
	svc := NewOrderService(repo, nil, nil, nil)

	_, err := svc.ListOrdersByPartnerID(context.Background(), authz.PartnerActor{
		PartnerID:  uuid.New(),
		EmployeeID: uuid.New(),
		Role:       authz.RoleManager,
		LocationID: &locationID,
	}, "confirmed", 20, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if repo.listByPartnerCalled {
		t.Fatalf("expected manager list to use location scoped query")
	}
}

func TestMarkOrderPickedUpScopesManagerToAssignedLocation(t *testing.T) {
	t.Parallel()

	locationID := uuid.New()
	orderID := uuid.New()
	employeeID := uuid.New()
	repo := &orderRepoStub{
		markLocationFn: func(ctx context.Context, gotOrderID, gotLocationID, gotEmployeeID uuid.UUID) error {
			if gotOrderID != orderID {
				t.Fatalf("expected order %s, got %s", orderID, gotOrderID)
			}
			if gotLocationID != locationID {
				t.Fatalf("expected location %s, got %s", locationID, gotLocationID)
			}
			if gotEmployeeID != employeeID {
				t.Fatalf("expected employee %s, got %s", employeeID, gotEmployeeID)
			}
			return nil
		},
	}
	svc := NewOrderService(repo, nil, nil, nil)

	err := svc.MarkOrderPickedUp(context.Background(), authz.PartnerActor{
		PartnerID:  uuid.New(),
		EmployeeID: employeeID,
		Role:       authz.RoleManager,
		LocationID: &locationID,
	}, orderID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if repo.markPartnerCalled {
		t.Fatalf("expected manager pickup to use location scoped query")
	}
}
