package service

import (
	"context"

	"github.com/google/uuid"
	catalogDomain "github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	"github.com/nlsnnn/berezhok/internal/modules/customer/repository"
)

type locationService struct {
	repo locationRepo
}

type locationRepo interface {
	SearchLocations(ctx context.Context, params repository.SearchLocationsParams) ([]repository.LocationWithDetails, error)
	CountActiveLocations(ctx context.Context, categoryCode *string) (int64, error)
	GetLocationDetailsByID(ctx context.Context, id uuid.UUID) (repository.LocationWithDetails, error)
	CountActiveBoxesByLocationID(ctx context.Context, locationID uuid.UUID) (int64, error)
	GetActiveBoxesByLocationID(ctx context.Context, locationID uuid.UUID) ([]catalogDomain.SurpriseBox, error)
}

func NewLocationService(repo locationRepo) *locationService {
	return &locationService{repo: repo}
}

// SearchLocationsInput holds search parameters
type SearchLocationsInput struct {
	CategoryCode *string
	Limit        int
	Offset       int
}

// SearchLocationsResult holds search results with pagination
type SearchLocationsResult struct {
	Locations []repository.LocationWithDetails
	Total     int64
	Limit     int
	Offset    int
	HasMore   bool
}

func (s *locationService) SearchLocations(ctx context.Context, input SearchLocationsInput) (SearchLocationsResult, error) {
	// Get total count for pagination
	total, err := s.repo.CountActiveLocations(ctx, input.CategoryCode)
	if err != nil {
		return SearchLocationsResult{}, err
	}

	// Get locations
	locations, err := s.repo.SearchLocations(ctx, repository.SearchLocationsParams{
		CategoryCode: input.CategoryCode,
		Limit:        input.Limit,
		Offset:       input.Offset,
	})
	if err != nil {
		return SearchLocationsResult{}, err
	}

	hasMore := int64(input.Offset+input.Limit) < total

	return SearchLocationsResult{
		Locations: locations,
		Total:     total,
		Limit:     input.Limit,
		Offset:    input.Offset,
		HasMore:   hasMore,
	}, nil
}

// LocationDetails holds full location details with boxes
type LocationDetails struct {
	Location    repository.LocationWithDetails
	ActiveBoxes []catalogDomain.SurpriseBox
}

func (s *locationService) GetLocationDetails(ctx context.Context, locationID uuid.UUID) (LocationDetails, error) {
	// Get location details
	location, err := s.repo.GetLocationDetailsByID(ctx, locationID)
	if err != nil {
		return LocationDetails{}, err
	}

	// Get active boxes
	boxes, err := s.repo.GetActiveBoxesByLocationID(ctx, locationID)
	if err != nil {
		return LocationDetails{}, err
	}

	return LocationDetails{
		Location:    location,
		ActiveBoxes: boxes,
	}, nil
}
