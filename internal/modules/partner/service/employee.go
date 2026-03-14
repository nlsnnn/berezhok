package service

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
)

type empService struct {
	repo empRepo
}

type empRepo interface {
	FindByID(ctx context.Context, id string) (domain.Employee, error)
	List(ctx context.Context) ([]domain.Employee, error)
	ListByPartnerID(ctx context.Context, partnerID string) ([]domain.Employee, error)
	Create(ctx context.Context, partnerID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error)
	Delete(ctx context.Context, id string) error
}

func NewEmployeeService(repo empRepo) *empService {
	return &empService{repo: repo}
}

func (s *empService) List(ctx context.Context) ([]domain.Employee, error) {
	return s.repo.List(ctx)
}

func (s *empService) ListByPartnerID(ctx context.Context, partnerID string) ([]domain.Employee, error) {
	return s.repo.ListByPartnerID(ctx, partnerID)
}

func (s *empService) FindByID(ctx context.Context, id string) (domain.Employee, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *empService) Create(ctx context.Context, partnerID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error) {
	return s.repo.Create(ctx, partnerID, email, passwordHash, name, role)
}

func (s *empService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
