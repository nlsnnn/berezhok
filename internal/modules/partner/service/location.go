package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner"
)

type locationService struct {
	repo         *sqlc.Queries
	locationRepo locationRepo
}

type locationRepo interface {
	Create(ctx context.Context, location domain.Location) (domain.Location, error)
	FindCategoryByCode(ctx context.Context, code string) (domain.LocationCategory, error)
}

func NewLocationService(repo *sqlc.Queries, locationRepo locationRepo) locationService {
	return locationService{
		repo:         repo,
		locationRepo: locationRepo,
	}
}

func (s *locationService) List(ctx context.Context) ([]sqlc.Location, error) {
	return s.repo.ListLocations(ctx)
}

func (s *locationService) Create(ctx context.Context, partnerID uuid.UUID, code string, name string, address string, latitude float64, longitude float64) (domain.Location, error) {
	// TODO: check if partner exists

	category, err := s.locationRepo.FindCategoryByCode(ctx, code)
	if err != nil {
		return domain.Location{}, partner.ErrLocationCategoryNotFound
	}

	location, err := domain.NewLocation(partnerID.String(), name, address, category,
		domain.GeoPoint{
			Latitude:  latitude,
			Longitude: longitude,
		})
	if err != nil {
		return domain.Location{}, err
	}

	return s.locationRepo.Create(ctx, location)
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
