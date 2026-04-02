package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	orderDomain "github.com/nlsnnn/berezhok/internal/modules/order/domain"
	reviewDomain "github.com/nlsnnn/berezhok/internal/modules/review/domain"
	reviewErrors "github.com/nlsnnn/berezhok/internal/modules/review/errors"
)

type reviewRepo interface {
	Create(ctx context.Context, review *reviewDomain.Review) error
	ListByLocationID(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]reviewDomain.ReviewWithUser, error)
	CountByLocationID(ctx context.Context, locationID uuid.UUID) (int, error)
}

type orderReader interface {
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error)
}

type reviewService struct {
	repo    reviewRepo
	orders  orderReader
}

type CreateReviewInput struct {
	UserID  uuid.UUID
	OrderID uuid.UUID
	Rating  int
	Comment string
}

type CreateReviewResult struct {
	ID        uuid.UUID
	Rating    int
	Comment   string
	OrderID   uuid.UUID
	CreatedAt time.Time
}

type ListLocationReviewsResult struct {
	Items   []reviewDomain.ReviewWithUser
	Total   int
	Limit   int
	Offset  int
	HasMore bool
}

func NewReviewService(repo reviewRepo, orders orderReader) *reviewService {
	return &reviewService{repo: repo, orders: orders}
}

func (s *reviewService) CreateReview(ctx context.Context, input CreateReviewInput) (*CreateReviewResult, error) {
	const op = "review.service.CreateReview"

	order, err := s.orders.GetOrderByID(ctx, input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if order.CustomerID != input.UserID {
		return nil, reviewErrors.ErrOrderOwnershipMismatch
	}

	if order.Status != orderDomain.OrderStatusCompleted {
		return nil, reviewErrors.ErrOrderNotCompleted
	}

	review := &reviewDomain.Review{
		ID:         uuid.New(),
		UserID:     input.UserID,
		LocationID: order.LocationID,
		OrderID:    input.OrderID,
		Rating:     input.Rating,
		Comment:    input.Comment,
	}

	err = s.repo.Create(ctx, review)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &CreateReviewResult{
		ID:        review.ID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		OrderID:   review.OrderID,
		CreatedAt: review.CreatedAt,
	}, nil
}

func (s *reviewService) ListLocationReviews(ctx context.Context, locationID uuid.UUID, limit, offset int) (*ListLocationReviewsResult, error) {
	const op = "review.service.ListLocationReviews"

	items, err := s.repo.ListByLocationID(ctx, locationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := s.repo.CountByLocationID(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ListLocationReviewsResult{
		Items:   items,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+limit < total,
	}, nil
}
