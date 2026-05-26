package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
)

type pinsRepo interface {
	ListAvailable(ctx context.Context) ([]domain.LocationPin, error)
	GetForLocation(ctx context.Context, locationID uuid.UUID) ([]domain.LocationPin, error)
	SetForLocation(ctx context.Context, locationID uuid.UUID, pinCodes []string) error
}

type PinsService struct {
	pins     pinsRepo
	location locationRepo
}

func NewPinsService(pins pinsRepo, location locationRepo) *PinsService {
	return &PinsService{pins: pins, location: location}
}

func (s *PinsService) ListAvailable(ctx context.Context) ([]domain.LocationPin, error) {
	return s.pins.ListAvailable(ctx)
}

func (s *PinsService) UpdateLocationPins(ctx context.Context, partnerID string, locationID uuid.UUID, pinCodes []string) error {
	loc, err := s.location.FindByID(ctx, locationID)
	if err != nil {
		return err
	}
	if loc.PartnerID != partnerID {
		return fmt.Errorf("location does not belong to partner")
	}
	return s.pins.SetForLocation(ctx, locationID, pinCodes)
}

func (s *PinsService) GetForLocation(ctx context.Context, locationID uuid.UUID) ([]domain.LocationPin, error) {
	return s.pins.GetForLocation(ctx, locationID)
}
