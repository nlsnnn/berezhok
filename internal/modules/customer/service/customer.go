package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/customer/domain"
)

type customerService struct {
	repo userRepo
}

type userRepo interface {
	GetProfile(ctx context.Context, id uuid.UUID) (domain.Profile, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, name string) (domain.User, error)
}

func NewCustomerService(repo userRepo) *customerService {
	return &customerService{repo: repo}
}

func (s *customerService) GetProfile(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	return s.repo.GetProfile(ctx, userID)
}

func (s *customerService) UpdateProfile(ctx context.Context, userID uuid.UUID, name string) (domain.User, error) {
	return s.repo.UpdateProfile(ctx, userID, name)
}
