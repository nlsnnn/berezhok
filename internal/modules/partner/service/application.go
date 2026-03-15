package service

import (
	"context"
	"fmt"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	types "github.com/nlsnnn/berezhok/internal/shared/domain"
	"github.com/nlsnnn/berezhok/internal/shared/generator"
)

type CreateApplicationInput struct {
	ContactName  string
	ContactEmail string
	ContactPhone string
	BusinessName string
	CategoryCode string
	Address      string
	Longitude    float64
	Latitude     float64
	Description  string
}

type appService struct {
	repo             appRepo
	partnerSvc       partnerProvider
	employeeSvc      employeeCreator
	locationProvider locationProvider
}

type appRepo interface {
	FindByID(ctx context.Context, id string) (domain.Application, error)
	List(ctx context.Context) ([]domain.Application, error)
	Create(ctx context.Context, app domain.Application) (domain.Application, error)
	UpdateStatus(ctx context.Context, id string, status domain.ApplicationStatus, rejectionReason string) error
	Delete(ctx context.Context, id string) error
}

type partnerProvider interface {
	Create(ctx context.Context, name string) (domain.Partner, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type employeeCreator interface {
	Create(ctx context.Context, partnerID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error)
}

type locationProvider interface {
	Create(ctx context.Context, input CreateLocationInput) (domain.Location, error)
	FindCategoryByCode(ctx context.Context, code string) (domain.LocationCategory, error)
}

func NewApplicationService(repo appRepo, partnerSvc partnerProvider, employeeSvc employeeCreator, locationProvider locationProvider) *appService {
	return &appService{
		repo:             repo,
		partnerSvc:       partnerSvc,
		employeeSvc:      employeeSvc,
		locationProvider: locationProvider,
	}
}

func (s *appService) Create(ctx context.Context, input CreateApplicationInput) (domain.Application, error) {
	coords, err := types.NewGeoPoint(input.Latitude, input.Longitude)
	if err != nil {
		return domain.Application{}, err
	}

	if _, err := s.locationProvider.FindCategoryByCode(ctx, input.CategoryCode); err != nil {
		return domain.Application{}, errors.ErrLocationCategoryNotFound
	}

	emailExists, err := s.partnerSvc.CheckEmailExists(ctx, input.ContactEmail)
	if err == nil && emailExists {
		return domain.Application{}, errors.ErrEmailAlreadyInUse
	}
	if err != nil {
		return domain.Application{}, err
	}

	application, err := domain.NewApplication(
		input.ContactName,
		input.ContactEmail,
		input.ContactPhone,
		input.BusinessName,
		input.CategoryCode,
		input.Address,
		input.Description,
		coords,
	)
	if err != nil {
		return domain.Application{}, err
	}

	return s.repo.Create(ctx, application)
}

func (s *appService) GetByID(ctx context.Context, id string) (domain.Application, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *appService) List(ctx context.Context) ([]domain.Application, error) {
	return s.repo.List(ctx)
}

func (s *appService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *appService) Approve(ctx context.Context, id string) error {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !app.CanTransitionTo(domain.ApplicationStatusApproved) {
		return errors.ErrInvalidStatusTransition
	}

	password := generator.GeneratePassword()
	passwordHash, err := auth.Hash(password)
	if err != nil {
		return err
	}

	partner, err := s.partnerSvc.Create(ctx, app.BusinessName)
	if err != nil {
		return err
	}

	if _, err := s.employeeSvc.Create(ctx, partner.ID, app.ContactEmail, passwordHash, app.ContactName, domain.EmployeeRoleOwner); err != nil {
		return err
	}

	if _, err := s.locationProvider.Create(ctx, CreateLocationInput{
		PartnerID:    partner.ID,
		CategoryCode: app.CategoryCode,
		Name:         app.BusinessName,
		Address:      app.Address,
		Latitude:     app.Coords.Latitude,
		Longitude:    app.Coords.Longitude,
		Status:       domain.LocationStatusDraft,
	}); err != nil {
		return err
	}

	// TODO: send email with credentials
	fmt.Printf("Partner approved. Contact: %s, Password: %s\n", app.ContactEmail, password)

	return s.repo.UpdateStatus(ctx, id, domain.ApplicationStatusApproved, "")
}

func (s *appService) Reject(ctx context.Context, id string, reason string) error {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !app.CanTransitionTo(domain.ApplicationStatusRejected) {
		return errors.ErrInvalidStatusTransition
	}

	return s.repo.UpdateStatus(ctx, id, domain.ApplicationStatusRejected, reason)
}
