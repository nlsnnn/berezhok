package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	"github.com/shopspring/decimal"
)

type CreateBoxInput struct {
	LocationID      uuid.UUID
	Name            string
	Description     string
	DiscountPrice   decimal.Decimal
	OriginalPrice   decimal.Decimal
	PickupTimeStart string
	PickupTimeEnd   string
	Quantity        int
	Image           string
	Status          string
}

type UpdateBoxInput struct {
	ID              string
	Name            string
	Description     string
	DiscountPrice   decimal.Decimal
	OriginalPrice   decimal.Decimal
	PickupTimeStart string
	PickupTimeEnd   string
	Quantity        int
	Image           string
	Status          string
}

type boxService struct {
	boxRepo     BoxRepository
	locationSvc locationFinder
}

type BoxRepository interface {
	// CreateBox creates a new surprise box in the database.
	CreateBox(ctx context.Context, box *domain.SurpriseBox) error
	// GetBoxByID retrieves a surprise box by its ID.
	GetBoxByID(ctx context.Context, id string) (*domain.SurpriseBox, error)
	// UpdateBox updates the details of an existing surprise box.
	UpdateBox(ctx context.Context, box *domain.SurpriseBox) error
	// DeleteBox removes a surprise box from the database.
	DeleteBox(ctx context.Context, id string) error
	// ListBoxes retrieves a list of surprise boxes, optionally filtered by location or status.
	// ListBoxes(ctx context.Context, locationID string, status domain.BoxStatus) ([]*domain.SurpriseBox, error)
}

type locationFinder interface {
	// LocationExists checks if a location with the given ID exists.
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	PartnerOwnsLocation(ctx context.Context, partnerID uuid.UUID, locationID uuid.UUID) (bool, error)
}

func NewBoxService(boxRepo BoxRepository, locationSvc locationFinder) *boxService {
	return &boxService{
		boxRepo:     boxRepo,
		locationSvc: locationSvc,
	}
}

func (s *boxService) CreateBox(ctx context.Context, partnerID uuid.UUID, input CreateBoxInput) (domain.SurpriseBox, error) {
	// Validate location
	exists, err := s.locationSvc.Exists(ctx, input.LocationID)
	if err != nil {
		return domain.SurpriseBox{}, fmt.Errorf("checking location existence: %w", err)
	}
	if !exists {
		return domain.SurpriseBox{}, catalogErrors.ErrLocationNotFound
	}

	// Validate partner ownership of the location
	owns, err := s.locationSvc.PartnerOwnsLocation(ctx, partnerID, input.LocationID)
	if err != nil {
		return domain.SurpriseBox{}, fmt.Errorf("checking location ownership: %w", err)
	}
	if !owns {
		return domain.SurpriseBox{}, catalogErrors.ErrUnauthorizedLocation
	}

	// Validate pickup time
	pickupTime, err := domain.NewPickupTimeFromStrings(input.PickupTimeStart, input.PickupTimeEnd)
	if err != nil {
		return domain.SurpriseBox{}, mapPickupTimeErr(err)
	}

	box := domain.SurpriseBox{
		LocationID:  input.LocationID,
		Name:        input.Name,
		Description: input.Description,
		Price: domain.Price{
			Original: input.OriginalPrice,
			Discount: input.DiscountPrice,
		},
		PickupTime: pickupTime,
		Quantity:   input.Quantity,
		Status:     domain.BoxStatus(input.Status),
		Image:      input.Image,
	}

	err = s.boxRepo.CreateBox(ctx, &box)
	if err != nil {
		return domain.SurpriseBox{}, err
	}

	return box, nil
}

func (s *boxService) GetBoxByID(ctx context.Context, id string) (*domain.SurpriseBox, error) {
	return s.boxRepo.GetBoxByID(ctx, id)
}

func (s *boxService) UpdateBox(ctx context.Context, input UpdateBoxInput) (domain.SurpriseBox, error) {
	pickupTime, err := domain.NewPickupTimeFromStrings(input.PickupTimeStart, input.PickupTimeEnd)
	if err != nil {
		return domain.SurpriseBox{}, mapPickupTimeErr(err)
	}

	box, err := s.boxRepo.GetBoxByID(ctx, input.ID)
	if err != nil {
		return domain.SurpriseBox{}, err
	}

	box.Name = input.Name
	box.Description = input.Description
	box.Price.Original = input.OriginalPrice
	box.Price.Discount = input.DiscountPrice
	box.PickupTime = pickupTime
	box.Quantity = input.Quantity
	box.Image = input.Image
	box.Status = domain.BoxStatus(input.Status)

	if err = s.boxRepo.UpdateBox(ctx, box); err != nil {
		return domain.SurpriseBox{}, err
	}

	return *box, nil
}

func (s *boxService) DeleteBox(ctx context.Context, id string) error {
	return s.boxRepo.DeleteBox(ctx, id)
}

func mapPickupTimeErr(err error) error {
	if errors.Is(err, catalogErrors.ErrInvalidPickupTimeRange) {
		return catalogErrors.ErrInvalidPickupTimeRange
	}

	return fmt.Errorf("pickup time: %w", catalogErrors.ErrInvalidPickupTimeFormat)
}
