package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
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
	partnerSvc  partnerChecker
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
	// GetBoxesByLocationID retrieves all surprise boxes for a given location ID.
	GetBoxesByLocationID(ctx context.Context, locationID uuid.UUID) ([]domain.SurpriseBox, error)
	// GetBoxesByPartnerID retrieves all surprise boxes for a given partner ID.
	GetBoxesByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]domain.SurpriseBox, error)
}

type locationFinder interface {
	// LocationExists checks if a location with the given ID exists.
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	PartnerOwnsLocation(ctx context.Context, partnerID, locationID uuid.UUID) (bool, error)
}

type partnerChecker interface {
	CanActivateBoxes(ctx context.Context, partnerID uuid.UUID) (bool, error)
}

func NewBoxService(boxRepo BoxRepository, locationSvc locationFinder, partnerSvc partnerChecker) *boxService {
	return &boxService{
		boxRepo:     boxRepo,
		locationSvc: locationSvc,
		partnerSvc:  partnerSvc,
	}
}

func (s *boxService) CreateBox(ctx context.Context, actor authz.PartnerActor, input CreateBoxInput) (domain.SurpriseBox, error) {
	if err := actor.EnsureCan(authz.PermissionPartnerBoxesManage); err != nil {
		return domain.SurpriseBox{}, err
	}
	if err := actor.EnsureLocation(input.LocationID); err != nil {
		return domain.SurpriseBox{}, err
	}

	// Validate location
	exists, err := s.locationSvc.Exists(ctx, input.LocationID)
	if err != nil {
		return domain.SurpriseBox{}, fmt.Errorf("checking location existence: %w", err)
	}
	if !exists {
		return domain.SurpriseBox{}, catalogErrors.ErrLocationNotFound
	}

	// Validate partner ownership of the location
	owns, err := s.locationSvc.PartnerOwnsLocation(ctx, actor.PartnerID, input.LocationID)
	if err != nil {
		return domain.SurpriseBox{}, fmt.Errorf("checking location ownership: %w", err)
	}
	if !owns {
		return domain.SurpriseBox{}, catalogErrors.ErrUnauthorizedLocation
	}

	// Validate pickup time
	pickupTime, err := sharedDomain.NewPickupTimeFromStrings(input.PickupTimeStart, input.PickupTimeEnd)
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

	canActivate, err := s.partnerSvc.CanActivateBoxes(ctx, actor.PartnerID)
	if err != nil {
		return domain.SurpriseBox{}, fmt.Errorf("checking partner legal info: %w", err)
	}

	if !canActivate && box.Status == domain.BoxStatusActive {
		box.Status = domain.BoxStatusInactive
	}

	err = s.boxRepo.CreateBox(ctx, &box)
	if err != nil {
		return domain.SurpriseBox{}, err
	}

	return box, nil
}

func (s *boxService) GetBoxesByLocationID(ctx context.Context, locationID uuid.UUID) ([]domain.SurpriseBox, error) {
	return s.boxRepo.GetBoxesByLocationID(ctx, locationID)
}

func (s *boxService) GetBoxesByPartnerID(ctx context.Context, actor authz.PartnerActor) ([]domain.SurpriseBox, error) {
	if err := actor.EnsureCan(authz.PermissionPartnerBoxesManage); err != nil {
		return nil, err
	}
	if actor.Role != authz.RoleOwner {
		if actor.LocationID == nil {
			return nil, authz.ErrLocationScopeDenied
		}
		return s.boxRepo.GetBoxesByLocationID(ctx, *actor.LocationID)
	}

	return s.boxRepo.GetBoxesByPartnerID(ctx, actor.PartnerID)
}

func (s *boxService) GetBoxByID(ctx context.Context, id string) (*domain.SurpriseBox, error) {
	return s.boxRepo.GetBoxByID(ctx, id)
}

func (s *boxService) GetPartnerBoxByID(ctx context.Context, actor authz.PartnerActor, id string) (*domain.SurpriseBox, error) {
	if err := actor.EnsureCan(authz.PermissionPartnerBoxesManage); err != nil {
		return nil, err
	}

	box, err := s.boxRepo.GetBoxByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err = actor.EnsureLocation(box.LocationID); err != nil {
		return nil, err
	}

	return box, nil
}

func (s *boxService) UpdateBox(ctx context.Context, actor authz.PartnerActor, input UpdateBoxInput) (domain.SurpriseBox, error) {
	if err := actor.EnsureCan(authz.PermissionPartnerBoxesManage); err != nil {
		return domain.SurpriseBox{}, err
	}

	pickupTime, err := sharedDomain.NewPickupTimeFromStrings(input.PickupTimeStart, input.PickupTimeEnd)
	if err != nil {
		return domain.SurpriseBox{}, mapPickupTimeErr(err)
	}

	box, err := s.boxRepo.GetBoxByID(ctx, input.ID)
	if err != nil {
		return domain.SurpriseBox{}, err
	}
	if err = actor.EnsureLocation(box.LocationID); err != nil {
		return domain.SurpriseBox{}, err
	}

	if domain.BoxStatus(input.Status) == domain.BoxStatusActive {
		canActivate, err := s.partnerSvc.CanActivateBoxes(ctx, actor.PartnerID)
		if err != nil {
			return domain.SurpriseBox{}, fmt.Errorf("checking partner legal info: %w", err)
		}

		if !canActivate {
			return domain.SurpriseBox{}, catalogErrors.ErrPartnerDocumentsRequired
		}
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

func (s *boxService) DeleteBox(ctx context.Context, actor authz.PartnerActor, id string) error {
	if err := actor.EnsureCan(authz.PermissionPartnerBoxesManage); err != nil {
		return err
	}
	box, err := s.boxRepo.GetBoxByID(ctx, id)
	if err != nil {
		return err
	}
	if err = actor.EnsureLocation(box.LocationID); err != nil {
		return err
	}

	return s.boxRepo.DeleteBox(ctx, id)
}

func mapPickupTimeErr(err error) error {
	if errors.Is(err, catalogErrors.ErrInvalidPickupTimeRange) {
		return catalogErrors.ErrInvalidPickupTimeRange
	}

	return fmt.Errorf("pickup time: %w", catalogErrors.ErrInvalidPickupTimeFormat)
}
