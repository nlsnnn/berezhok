package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	"github.com/nlsnnn/berezhok/internal/shared/generator"
)

type empService struct {
	repo                 empRepo
	emailChecker         employeeEmailChecker
	locationProvider     employeeLocationProvider
	notificationProvider employeeNotificationProvider
}

type empRepo interface {
	FindByID(ctx context.Context, id string) (domain.Employee, error)
	List(ctx context.Context) ([]domain.Employee, error)
	ListByPartnerID(ctx context.Context, partnerID string) ([]domain.Employee, error)
	Create(ctx context.Context, partnerID, locationID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error)
	Delete(ctx context.Context, id string) error
}

type employeeEmailChecker interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type employeeLocationProvider interface {
	FindByID(ctx context.Context, id uuid.UUID) (domain.Location, error)
	FindByPartnerID(ctx context.Context, partnerID string) ([]domain.Location, error)
}

type employeeNotificationProvider interface {
	SendEmployeeInviteNotification(ctx context.Context, email, name, password string) error
}

type CreateManagedEmployeeInput struct {
	PartnerID  string
	LocationID string
	Email      string
	Name       string
}

type DeleteManagedEmployeeInput struct {
	PartnerID       string
	ActorEmployeeID string
	EmployeeID      string
}

func NewEmployeeService(
	repo empRepo,
	emailChecker employeeEmailChecker,
	locationProvider employeeLocationProvider,
	notificationProvider employeeNotificationProvider,
) *empService {
	return &empService{
		repo:                 repo,
		emailChecker:         emailChecker,
		locationProvider:     locationProvider,
		notificationProvider: notificationProvider,
	}
}

func (s *empService) List(ctx context.Context) ([]domain.Employee, error) {
	return s.repo.List(ctx)
}

func (s *empService) ListByPartnerID(ctx context.Context, partnerID string) ([]domain.Employee, error) {
	return s.repo.ListByPartnerID(ctx, partnerID)
}

func (s *empService) FindByID(ctx context.Context, id string) (domain.Employee, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *empService) Create(ctx context.Context, partnerID, locationID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error) {
	return s.repo.Create(ctx, partnerID, locationID, email, passwordHash, name, role)
}

func (s *empService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *empService) ListManagedByPartnerID(ctx context.Context, partnerID string) ([]domain.ManagedEmployee, error) {
	employees, err := s.repo.ListByPartnerID(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	locations, err := s.locationProvider.FindByPartnerID(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	locationNames := make(map[string]string, len(locations))
	for _, location := range locations {
		locationNames[location.ID] = location.Name
	}

	result := make([]domain.ManagedEmployee, len(employees))
	for i, employee := range employees {
		result[i] = managedEmployeeFromDomain(employee, locationNames[employee.LocationID])
	}

	return result, nil
}

func (s *empService) CreateManaged(ctx context.Context, input CreateManagedEmployeeInput) (domain.ManagedEmployee, error) {
	emailExists, err := s.emailChecker.CheckEmailExists(ctx, input.Email)
	if err != nil {
		return domain.ManagedEmployee{}, err
	}
	if emailExists {
		return domain.ManagedEmployee{}, partnerErrors.ErrEmailAlreadyInUse
	}

	locationID, err := uuid.Parse(input.LocationID)
	if err != nil {
		return domain.ManagedEmployee{}, partnerErrors.ErrLocationNotFound
	}

	location, err := s.locationProvider.FindByID(ctx, locationID)
	if err != nil {
		return domain.ManagedEmployee{}, err
	}
	if location.PartnerID != input.PartnerID {
		return domain.ManagedEmployee{}, partnerErrors.ErrLocationNotOwnedByPartner
	}

	password := generator.GeneratePassword()
	passwordHash, err := auth.Hash(password)
	if err != nil {
		return domain.ManagedEmployee{}, err
	}

	employee, err := s.repo.Create(
		ctx,
		input.PartnerID,
		input.LocationID,
		input.Email,
		passwordHash,
		input.Name,
		domain.EmployeeRoleEmployee,
	)
	if err != nil {
		return domain.ManagedEmployee{}, err
	}

	if err := s.notificationProvider.SendEmployeeInviteNotification(ctx, input.Email, input.Name, password); err != nil {
		if rollbackErr := s.repo.Delete(ctx, employee.ID); rollbackErr != nil {
			return domain.ManagedEmployee{}, fmt.Errorf("send invite: %w; rollback create employee: %v", err, rollbackErr)
		}

		return domain.ManagedEmployee{}, err
	}

	return managedEmployeeFromDomain(employee, location.Name), nil
}

func (s *empService) DeleteManaged(ctx context.Context, input DeleteManagedEmployeeInput) error {
	employee, err := s.repo.FindByID(ctx, input.EmployeeID)
	if err != nil {
		return err
	}

	if employee.PartnerID != input.PartnerID {
		return partnerErrors.ErrEmployeeNotFound
	}
	if employee.Role == domain.EmployeeRoleOwner {
		return partnerErrors.ErrCannotDeleteOwner
	}
	if employee.ID == input.ActorEmployeeID {
		return partnerErrors.ErrCannotDeleteSelf
	}

	return s.repo.Delete(ctx, input.EmployeeID)
}

func managedEmployeeFromDomain(employee domain.Employee, locationName string) domain.ManagedEmployee {
	return domain.ManagedEmployee{
		ID:                 employee.ID,
		PartnerID:          employee.PartnerID,
		LocationID:         employee.LocationID,
		LocationName:       locationName,
		Email:              employee.Email,
		Role:               employee.Role,
		Name:               employee.Name,
		MustChangePassword: employee.MustChangePassword,
		CreatedAt:          employee.CreatedAt,
	}
}
