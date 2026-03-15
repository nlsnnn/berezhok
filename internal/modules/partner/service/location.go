package service

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type CreateLocationInput struct {
	PartnerID    string
	CategoryCode string
	Name         string
	Address      string
	Latitude     float64
	Longitude    float64
	Status       domain.LocationStatus
}

type UpdateLocationInput struct {
	ID            string
	Name          string
	Address       string
	CategoryCode  string
	LogoURL       string
	CoverImageURL string
}

type locationService struct {
	repo locationRepo
}

type locationRepo interface {
	Create(ctx context.Context, location domain.Location) (domain.Location, error)
	FindByPartnerID(ctx context.Context, partnerID string) ([]domain.Location, error)
	FindCategoryByCode(ctx context.Context, code string) (domain.LocationCategory, error)
	Delete(ctx context.Context, id string) error
}

func NewLocationService(repo locationRepo) *locationService {
	return &locationService{repo: repo}
}

func (s *locationService) ListByPartner(ctx context.Context, partnerID string) ([]domain.Location, error) {
	return s.repo.FindByPartnerID(ctx, partnerID)
}

func (s *locationService) FindCategoryByCode(ctx context.Context, code string) (domain.LocationCategory, error) {
	return s.repo.FindCategoryByCode(ctx, code)
}

func (s *locationService) Create(ctx context.Context, input CreateLocationInput) (domain.Location, error) {
	category, err := s.repo.FindCategoryByCode(ctx, input.CategoryCode)
	if err != nil {
		return domain.Location{}, errors.ErrLocationCategoryNotFound
	}

	location, err := domain.NewLocation(input.PartnerID, input.Name, input.Address, category,
		input.Status,
		sharedDomain.GeoPoint{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		})
	if err != nil {
		return domain.Location{}, err
	}

	return s.repo.Create(ctx, location)
}

func (s *locationService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
