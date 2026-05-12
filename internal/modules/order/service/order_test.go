package service

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
)

type orderRepoStub struct {
	listByPartnerCalled bool
	listByLocationFn    func(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]domain.PartnerOrderListItem, int, error)
	markPartnerCalled   bool
	markLocationFn      func(ctx context.Context, orderID, locationID, employeeID uuid.UUID) error
}

func (r *orderRepoStub) CreateOrder(ctx context.Context, order *domain.Order) error { return nil }
func (r *orderRepoStub) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	return nil, nil
}

func (r *orderRepoStub) GetOrderDetailsByID(ctx context.Context, orderID uuid.UUID) (*domain.OrderDetails, error) {
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
	return true, nil
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
