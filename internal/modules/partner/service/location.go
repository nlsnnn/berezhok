package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner"
)

type locationService struct {
	repo *sqlc.Queries
}

func NewLocationService(repo *sqlc.Queries) locationService {
	return locationService{
		repo: repo,
	}
}

func (s *locationService) List(ctx context.Context) ([]sqlc.Location, error) {
	return s.repo.ListLocations(ctx)
}

func (s *locationService) Create(ctx context.Context, arg sqlc.CreateLocationParams) (sqlc.Location, error) {
	// func (s *locationService) Create(ctx context.Context, partnerID uuid.UUID) (sqlc.Location, error) {
	// TODO: check if partner exists

	_, err := s.repo.FindCategoryByCode(ctx, arg.CategoryCode)
	if err != nil {
		return sqlc.Location{}, partner.ErrLocationCategoryNotFound
	}

	return s.repo.CreateLocation(ctx, arg)
}

func (s *locationService) Update(ctx context.Context, arg sqlc.UpdateLocationParams) error {
	location, err := s.repo.UpdateLocation(ctx, arg)
	if err != nil {
		return err
	}

	_ = location

	return nil
}

func (s *locationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLocation(ctx, id)
}
