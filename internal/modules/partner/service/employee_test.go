package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
)

type employeeRepoStub struct {
	findByIDFn        func(ctx context.Context, id string) (domain.Employee, error)
	listByPartnerIDFn func(ctx context.Context, partnerID string) ([]domain.Employee, error)
	createFn          func(ctx context.Context, partnerID, locationID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error)
	deleteFn          func(ctx context.Context, id string) error
}

func (s *employeeRepoStub) FindByID(ctx context.Context, id string) (domain.Employee, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return domain.Employee{}, nil
}

func (s *employeeRepoStub) List(ctx context.Context) ([]domain.Employee, error) {
	return nil, nil
}

func (s *employeeRepoStub) ListByPartnerID(ctx context.Context, partnerID string) ([]domain.Employee, error) {
	if s.listByPartnerIDFn != nil {
		return s.listByPartnerIDFn(ctx, partnerID)
	}
	return nil, nil
}

func (s *employeeRepoStub) Create(ctx context.Context, partnerID, locationID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error) {
	if s.createFn != nil {
		return s.createFn(ctx, partnerID, locationID, email, passwordHash, name, role)
	}
	return domain.Employee{}, nil
}

func (s *employeeRepoStub) Delete(ctx context.Context, id string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type employeeEmailCheckerStub struct {
	checkFn func(ctx context.Context, email string) (bool, error)
}

func (s *employeeEmailCheckerStub) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	if s.checkFn != nil {
		return s.checkFn(ctx, email)
	}
	return false, nil
}

type employeeLocationProviderStub struct {
	findByIDFn        func(ctx context.Context, id uuid.UUID) (domain.Location, error)
	findByPartnerIDFn func(ctx context.Context, partnerID string) ([]domain.Location, error)
}

func (s *employeeLocationProviderStub) FindByID(ctx context.Context, id uuid.UUID) (domain.Location, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return domain.Location{}, nil
}

func (s *employeeLocationProviderStub) FindByPartnerID(ctx context.Context, partnerID string) ([]domain.Location, error) {
	if s.findByPartnerIDFn != nil {
		return s.findByPartnerIDFn(ctx, partnerID)
	}
	return nil, nil
}

type employeeNotificationProviderStub struct {
	sendInviteFn func(ctx context.Context, email, name, password string) error
}

func (s *employeeNotificationProviderStub) SendEmployeeInviteNotification(ctx context.Context, email, name, password string) error {
	if s.sendInviteFn != nil {
		return s.sendInviteFn(ctx, email, name, password)
	}
	return nil
}

func TestCreateManagedEmployeeSuccess(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New().String()
	locationID := uuid.New()
	createdAt := time.Now()

	svc := NewEmployeeService(
		&employeeRepoStub{
			createFn: func(ctx context.Context, gotPartnerID, gotLocationID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error) {
				if gotPartnerID != partnerID {
					t.Fatalf("expected partner id %s, got %s", partnerID, gotPartnerID)
				}
				if gotLocationID != locationID.String() {
					t.Fatalf("expected location id %s, got %s", locationID, gotLocationID)
				}
				if email != "staff@example.com" {
					t.Fatalf("unexpected email %s", email)
				}
				if passwordHash == "" {
					t.Fatal("expected password hash to be set")
				}
				if role != domain.EmployeeRoleEmployee {
					t.Fatalf("expected role employee, got %s", role)
				}

				return domain.Employee{
					ID:                 uuid.New().String(),
					PartnerID:          gotPartnerID,
					LocationID:         gotLocationID,
					Email:              email,
					Name:               name,
					Role:               role,
					MustChangePassword: true,
					CreatedAt:          createdAt,
				}, nil
			},
		},
		&employeeEmailCheckerStub{
			checkFn: func(ctx context.Context, email string) (bool, error) {
				return false, nil
			},
		},
		&employeeLocationProviderStub{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (domain.Location, error) {
				return domain.Location{
					ID:        id.String(),
					PartnerID: partnerID,
					Name:      "Coffee Point",
				}, nil
			},
		},
		&employeeNotificationProviderStub{
			sendInviteFn: func(ctx context.Context, email, name, password string) error {
				if email != "staff@example.com" {
					t.Fatalf("unexpected email %s", email)
				}
				if password == "" {
					t.Fatal("expected generated password")
				}
				return nil
			},
		},
	)

	employee, err := svc.CreateManaged(context.Background(), CreateManagedEmployeeInput{
		PartnerID:  partnerID,
		LocationID: locationID.String(),
		Email:      "staff@example.com",
		Name:       "Ivan",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if employee.LocationName != "Coffee Point" {
		t.Fatalf("expected location name Coffee Point, got %s", employee.LocationName)
	}
	if employee.Role != domain.EmployeeRoleEmployee {
		t.Fatalf("expected role employee, got %s", employee.Role)
	}
}

func TestCreateManagedEmployeeRejectsForeignLocation(t *testing.T) {
	t.Parallel()

	svc := NewEmployeeService(
		&employeeRepoStub{},
		&employeeEmailCheckerStub{},
		&employeeLocationProviderStub{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (domain.Location, error) {
				return domain.Location{ID: id.String(), PartnerID: uuid.New().String()}, nil
			},
		},
		&employeeNotificationProviderStub{},
	)

	_, err := svc.CreateManaged(context.Background(), CreateManagedEmployeeInput{
		PartnerID:  uuid.New().String(),
		LocationID: uuid.New().String(),
		Email:      "staff@example.com",
		Name:       "Ivan",
	})
	if !errors.Is(err, partnerErrors.ErrLocationNotOwnedByPartner) {
		t.Fatalf("expected ErrLocationNotOwnedByPartner, got %v", err)
	}
}

func TestCreateManagedEmployeeRollsBackWhenInviteFails(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New().String()
	locationID := uuid.New()
	createdEmployeeID := uuid.New().String()
	deleted := false

	svc := NewEmployeeService(
		&employeeRepoStub{
			createFn: func(ctx context.Context, gotPartnerID, gotLocationID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error) {
				return domain.Employee{
					ID:         createdEmployeeID,
					PartnerID:  gotPartnerID,
					LocationID: gotLocationID,
					Email:      email,
					Name:       name,
					Role:       role,
				}, nil
			},
			deleteFn: func(ctx context.Context, id string) error {
				if id != createdEmployeeID {
					t.Fatalf("expected rollback delete for %s, got %s", createdEmployeeID, id)
				}
				deleted = true
				return nil
			},
		},
		&employeeEmailCheckerStub{},
		&employeeLocationProviderStub{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (domain.Location, error) {
				return domain.Location{ID: id.String(), PartnerID: partnerID, Name: "Coffee Point"}, nil
			},
		},
		&employeeNotificationProviderStub{
			sendInviteFn: func(ctx context.Context, email, name, password string) error {
				return errors.New("queue unavailable")
			},
		},
	)

	_, err := svc.CreateManaged(context.Background(), CreateManagedEmployeeInput{
		PartnerID:  partnerID,
		LocationID: locationID.String(),
		Email:      "staff@example.com",
		Name:       "Ivan",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !deleted {
		t.Fatal("expected created employee to be rolled back")
	}
}

func TestDeleteManagedEmployeeRejectsOwnerAndSelf(t *testing.T) {
	t.Parallel()

	svc := NewEmployeeService(
		&employeeRepoStub{
			findByIDFn: func(ctx context.Context, id string) (domain.Employee, error) {
				return domain.Employee{
					ID:        id,
					PartnerID: "partner-1",
					Role:      domain.EmployeeRoleOwner,
				}, nil
			},
		},
		&employeeEmailCheckerStub{},
		&employeeLocationProviderStub{},
		&employeeNotificationProviderStub{},
	)

	err := svc.DeleteManaged(context.Background(), DeleteManagedEmployeeInput{
		PartnerID:       "partner-1",
		ActorEmployeeID: "employee-1",
		EmployeeID:      "employee-2",
	})
	if !errors.Is(err, partnerErrors.ErrCannotDeleteOwner) {
		t.Fatalf("expected ErrCannotDeleteOwner, got %v", err)
	}

	selfSvc := NewEmployeeService(
		&employeeRepoStub{
			findByIDFn: func(ctx context.Context, id string) (domain.Employee, error) {
				return domain.Employee{
					ID:        id,
					PartnerID: "partner-1",
					Role:      domain.EmployeeRoleEmployee,
				}, nil
			},
		},
		&employeeEmailCheckerStub{},
		&employeeLocationProviderStub{},
		&employeeNotificationProviderStub{},
	)

	err = selfSvc.DeleteManaged(context.Background(), DeleteManagedEmployeeInput{
		PartnerID:       "partner-1",
		ActorEmployeeID: "employee-1",
		EmployeeID:      "employee-1",
	})
	if !errors.Is(err, partnerErrors.ErrCannotDeleteSelf) {
		t.Fatalf("expected ErrCannotDeleteSelf, got %v", err)
	}
}
