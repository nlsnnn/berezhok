package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
)

type boxRepoStub struct {
	createFn func(ctx context.Context, box *domain.SurpriseBox) error
	getFn    func(ctx context.Context, id string) (*domain.SurpriseBox, error)
}

func (r *boxRepoStub) CreateBox(ctx context.Context, box *domain.SurpriseBox) error {
	if r.createFn != nil {
		return r.createFn(ctx, box)
	}
	return nil
}

func (r *boxRepoStub) GetBoxByID(ctx context.Context, id string) (*domain.SurpriseBox, error) {
	if r.getFn != nil {
		return r.getFn(ctx, id)
	}
	return nil, catalogErrors.ErrBoxNotFound
}

func (r *boxRepoStub) UpdateBox(ctx context.Context, box *domain.SurpriseBox) error { return nil }
func (r *boxRepoStub) DeleteBox(ctx context.Context, id string) error               { return nil }
func (r *boxRepoStub) GetBoxesByLocationID(ctx context.Context, locationID uuid.UUID) ([]domain.SurpriseBox, error) {
	return nil, nil
}

func (r *boxRepoStub) GetBoxesByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]domain.SurpriseBox, error) {
	return nil, nil
}

type boxLocationStub struct {
	partnerID uuid.UUID
}

func (s *boxLocationStub) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return true, nil
}

func (s *boxLocationStub) PartnerOwnsLocation(ctx context.Context, partnerID, locationID uuid.UUID) (bool, error) {
	return partnerID == s.partnerID, nil
}

type boxPartnerStub struct{}

func (s *boxPartnerStub) CanActivateBoxes(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	return true, nil
}

func TestCreateBoxRejectsEmployee(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()
	locationID := uuid.New()
	actor := authz.PartnerActor{PartnerID: partnerID, EmployeeID: uuid.New(), Role: authz.RoleEmployee, LocationID: &locationID}
	svc := NewBoxService(&boxRepoStub{}, &boxLocationStub{partnerID: partnerID}, &boxPartnerStub{})

	_, err := svc.CreateBox(context.Background(), actor, validCreateBoxInput(locationID))
	if !errors.Is(err, authz.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestCreateBoxRejectsManagerOutsideAssignedLocation(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()
	assignedLocationID := uuid.New()
	requestedLocationID := uuid.New()
	actor := authz.PartnerActor{PartnerID: partnerID, EmployeeID: uuid.New(), Role: authz.RoleManager, LocationID: &assignedLocationID}
	svc := NewBoxService(&boxRepoStub{}, &boxLocationStub{partnerID: partnerID}, &boxPartnerStub{})

	_, err := svc.CreateBox(context.Background(), actor, validCreateBoxInput(requestedLocationID))
	if !errors.Is(err, authz.ErrLocationScopeDenied) {
		t.Fatalf("expected location scope denied, got %v", err)
	}
}

func TestCreateBoxAllowsManagerAssignedLocation(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()
	locationID := uuid.New()
	actor := authz.PartnerActor{PartnerID: partnerID, EmployeeID: uuid.New(), Role: authz.RoleManager, LocationID: &locationID}
	created := false
	svc := NewBoxService(&boxRepoStub{
		createFn: func(ctx context.Context, box *domain.SurpriseBox) error {
			created = true
			if box.LocationID != locationID {
				t.Fatalf("expected location %s, got %s", locationID, box.LocationID)
			}
			return nil
		},
	}, &boxLocationStub{partnerID: partnerID}, &boxPartnerStub{})

	_, err := svc.CreateBox(context.Background(), actor, validCreateBoxInput(locationID))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !created {
		t.Fatalf("expected box to be created")
	}
}

func TestGetPartnerBoxByIDRejectsManagerOutsideAssignedLocation(t *testing.T) {
	t.Parallel()

	assignedLocationID := uuid.New()
	boxLocationID := uuid.New()
	actor := authz.PartnerActor{PartnerID: uuid.New(), EmployeeID: uuid.New(), Role: authz.RoleManager, LocationID: &assignedLocationID}
	svc := NewBoxService(&boxRepoStub{
		getFn: func(ctx context.Context, id string) (*domain.SurpriseBox, error) {
			return &domain.SurpriseBox{ID: uuid.New(), LocationID: boxLocationID}, nil
		},
	}, &boxLocationStub{}, &boxPartnerStub{})

	_, err := svc.GetPartnerBoxByID(context.Background(), actor, uuid.New().String())
	if !errors.Is(err, authz.ErrLocationScopeDenied) {
		t.Fatalf("expected location scope denied, got %v", err)
	}
}

func validCreateBoxInput(locationID uuid.UUID) CreateBoxInput {
	return CreateBoxInput{
		LocationID:      locationID,
		Name:            "Dinner box",
		Description:     "Surprise",
		DiscountPrice:   decimal.NewFromInt(500),
		OriginalPrice:   decimal.NewFromInt(1000),
		PickupTimeStart: "18:00",
		PickupTimeEnd:   "20:00",
		Quantity:        3,
		Image:           "https://example.com/box.jpg",
		Status:          string(domain.BoxStatusActive),
	}
}
