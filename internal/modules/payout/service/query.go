package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
)

type QueryRepository interface {
	ListByPartner(ctx context.Context, partnerID uuid.UUID, limit, offset int32) ([]domain.Payout, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Payout, error)
	ListOrdersByPayoutID(ctx context.Context, payoutID uuid.UUID) ([]domain.PayoutOrder, error)
}

type QueryService struct {
	repo QueryRepository
}

func NewQueryService(repo QueryRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) ListByPartner(ctx context.Context, partnerID uuid.UUID, limit, offset int32) ([]domain.Payout, int64, error) {
	return s.repo.ListByPartner(ctx, partnerID, limit, offset)
}

func (s *QueryService) GetWithOrders(ctx context.Context, id uuid.UUID) (domain.Payout, []domain.PayoutOrder, error) {
	payout, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Payout{}, nil, err
	}

	orders, err := s.repo.ListOrdersByPayoutID(ctx, id)
	if err != nil {
		return domain.Payout{}, nil, err
	}

	return payout, orders, nil
}
