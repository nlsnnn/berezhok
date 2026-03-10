package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
)

type empService struct {
	repo sqlc.Querier
}

func NewEmployeeService(repo sqlc.Querier) *empService {
	return &empService{repo: repo}
}

func (s *empService) List(ctx context.Context) ([]sqlc.PartnerEmployee, error) {
	return s.repo.ListPartnerEmployees(ctx)
}

func (s *empService) ListByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]sqlc.PartnerEmployee, error) {
	return s.repo.ListEmployeesByPartnerID(ctx, partnerID)
}

func (s *empService) FindByID(ctx context.Context, id uuid.UUID) (sqlc.PartnerEmployee, error) {
	return s.repo.FindPartnerEmployeeByID(ctx, id)
}

func (s *empService) Create(ctx context.Context, arg sqlc.CreatePartnerEmployeeParams) (sqlc.PartnerEmployee, error) {
	return s.repo.CreatePartnerEmployee(ctx, arg)
}

func (s *empService) Update(ctx context.Context, arg sqlc.UpdatePartnerEmployeeParams) error {
	return s.repo.UpdatePartnerEmployee(ctx, arg)
}

func (s *empService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePartnerEmployee(ctx, id)
}
