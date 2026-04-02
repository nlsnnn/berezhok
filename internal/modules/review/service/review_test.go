package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	orderDomain "github.com/nlsnnn/berezhok/internal/modules/order/domain"
	reviewDomain "github.com/nlsnnn/berezhok/internal/modules/review/domain"
	reviewErrors "github.com/nlsnnn/berezhok/internal/modules/review/errors"
)

type testReviewRepo struct {
	createFn            func(ctx context.Context, review *reviewDomain.Review) error
	listByLocationIDFn  func(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]reviewDomain.ReviewWithUser, error)
	countByLocationIDFn func(ctx context.Context, locationID uuid.UUID) (int, error)
}

func (r *testReviewRepo) Create(ctx context.Context, review *reviewDomain.Review) error {
	if r.createFn != nil {
		return r.createFn(ctx, review)
	}

	return nil
}

func (r *testReviewRepo) ListByLocationID(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]reviewDomain.ReviewWithUser, error) {
	if r.listByLocationIDFn != nil {
		return r.listByLocationIDFn(ctx, locationID, limit, offset)
	}

	return nil, nil
}

func (r *testReviewRepo) CountByLocationID(ctx context.Context, locationID uuid.UUID) (int, error) {
	if r.countByLocationIDFn != nil {
		return r.countByLocationIDFn(ctx, locationID)
	}

	return 0, nil
}

type testOrderReader struct {
	getOrderByIDFn func(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error)
}

func (o *testOrderReader) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error) {
	if o.getOrderByIDFn != nil {
		return o.getOrderByIDFn(ctx, orderID)
	}

	return nil, nil
}

func TestCreateReviewOrderMustBeCompleted(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	service := NewReviewService(&testReviewRepo{}, &testOrderReader{
		getOrderByIDFn: func(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error) {
			return &orderDomain.Order{
				ID:         orderID,
				CustomerID: userID,
				LocationID: uuid.New(),
				Status:     orderDomain.OrderStatusPickedUp,
			}, nil
		},
	})

	_, err := service.CreateReview(context.Background(), CreateReviewInput{
		UserID:  userID,
		OrderID: uuid.New(),
		Rating:  5,
		Comment: "great",
	})

	if !errors.Is(err, reviewErrors.ErrOrderNotCompleted) {
		t.Fatalf("expected ErrOrderNotCompleted, got %v", err)
	}
}

func TestCreateReviewOrderMustBelongToCustomer(t *testing.T) {
	t.Parallel()

	service := NewReviewService(&testReviewRepo{}, &testOrderReader{
		getOrderByIDFn: func(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error) {
			return &orderDomain.Order{
				ID:         orderID,
				CustomerID: uuid.New(),
				LocationID: uuid.New(),
				Status:     orderDomain.OrderStatusCompleted,
			}, nil
		},
	})

	_, err := service.CreateReview(context.Background(), CreateReviewInput{
		UserID:  uuid.New(),
		OrderID: uuid.New(),
		Rating:  5,
		Comment: "great",
	})

	if !errors.Is(err, reviewErrors.ErrOrderOwnershipMismatch) {
		t.Fatalf("expected ErrOrderOwnershipMismatch, got %v", err)
	}
}

func TestCreateReviewSuccess(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	orderID := uuid.New()
	locationID := uuid.New()
	called := false

	service := NewReviewService(&testReviewRepo{
		createFn: func(ctx context.Context, review *reviewDomain.Review) error {
			called = true

			if review.UserID != userID {
				t.Fatalf("expected user id %s, got %s", userID, review.UserID)
			}

			if review.OrderID != orderID {
				t.Fatalf("expected order id %s, got %s", orderID, review.OrderID)
			}

			if review.LocationID != locationID {
				t.Fatalf("expected location id %s, got %s", locationID, review.LocationID)
			}

			if review.Rating != 4 {
				t.Fatalf("expected rating 4, got %d", review.Rating)
			}

			if review.Comment != "nice" {
				t.Fatalf("expected comment nice, got %s", review.Comment)
			}

			return nil
		},
	}, &testOrderReader{
		getOrderByIDFn: func(ctx context.Context, id uuid.UUID) (*orderDomain.Order, error) {
			return &orderDomain.Order{
				ID:         id,
				CustomerID: userID,
				LocationID: locationID,
				Status:     orderDomain.OrderStatusCompleted,
			}, nil
		},
	})

	result, err := service.CreateReview(context.Background(), CreateReviewInput{
		UserID:  userID,
		OrderID: orderID,
		Rating:  4,
		Comment: "nice",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.OrderID != orderID {
		t.Fatalf("expected order id %s, got %s", orderID, result.OrderID)
	}

	if !called {
		t.Fatal("expected repository Create to be called")
	}
}

func TestCreateReviewDuplicate(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	service := NewReviewService(&testReviewRepo{
		createFn: func(ctx context.Context, review *reviewDomain.Review) error {
			return reviewErrors.ErrReviewAlreadyExists
		},
	}, &testOrderReader{
		getOrderByIDFn: func(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error) {
			return &orderDomain.Order{
				ID:         orderID,
				CustomerID: userID,
				LocationID: uuid.New(),
				Status:     orderDomain.OrderStatusCompleted,
			}, nil
		},
	})

	_, err := service.CreateReview(context.Background(), CreateReviewInput{
		UserID:  userID,
		OrderID: uuid.New(),
		Rating:  5,
		Comment: "great",
	})

	if !errors.Is(err, reviewErrors.ErrReviewAlreadyExists) {
		t.Fatalf("expected ErrReviewAlreadyExists, got %v", err)
	}
}

func TestListLocationReviewsSuccess(t *testing.T) {
	t.Parallel()

	locationID := uuid.New()
	now := time.Now().UTC()

	service := NewReviewService(&testReviewRepo{
		listByLocationIDFn: func(ctx context.Context, id uuid.UUID, limit, offset int) ([]reviewDomain.ReviewWithUser, error) {
			if id != locationID {
				t.Fatalf("expected location id %s, got %s", locationID, id)
			}

			if limit != 20 || offset != 0 {
				t.Fatalf("expected limit/offset 20/0, got %d/%d", limit, offset)
			}

			return []reviewDomain.ReviewWithUser{{
				ID:        uuid.New(),
				Rating:    5,
				Comment:   "great",
				UserName:  "Egor",
				CreatedAt: now,
			}}, nil
		},
		countByLocationIDFn: func(ctx context.Context, id uuid.UUID) (int, error) {
			if id != locationID {
				t.Fatalf("expected location id %s, got %s", locationID, id)
			}

			return 21, nil
		},
	}, &testOrderReader{})

	result, err := service.ListLocationReviews(context.Background(), locationID, 20, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}

	if !result.HasMore {
		t.Fatal("expected has_more to be true")
	}
}
