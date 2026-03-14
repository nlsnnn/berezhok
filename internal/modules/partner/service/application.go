package service

import (
	"context"
	"fmt"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	"github.com/nlsnnn/berezhok/internal/shared/generator"
)

type CreateApplicationInput struct {
	ContactName  string
	ContactEmail string
	ContactPhone string
	BusinessName string
	CategoryCode string
	Address      string
	Description  string
}

type appService struct {
	repo        appRepo
	partnerSvc  partnerCreator
	employeeSvc employeeCreator
}

type appRepo interface {
	FindByID(ctx context.Context, id string) (domain.Application, error)
	List(ctx context.Context) ([]domain.Application, error)
	Create(ctx context.Context, app domain.Application) (domain.Application, error)
	UpdateStatus(ctx context.Context, id string, status domain.ApplicationStatus, rejectionReason string) error
	Delete(ctx context.Context, id string) error
}

type partnerCreator interface {
	Create(ctx context.Context, legalName string) (domain.Partner, error)
}

type employeeCreator interface {
	Create(ctx context.Context, partnerID, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error)
}

func NewApplicationService(repo appRepo, partnerSvc partnerCreator, employeeSvc employeeCreator) *appService {
	return &appService{
		repo:        repo,
		partnerSvc:  partnerSvc,
		employeeSvc: employeeSvc,
	}
}

func (s *appService) Create(ctx context.Context, input CreateApplicationInput) (domain.Application, error) {
	if input.ContactName == "" {
		return domain.Application{}, errors.ErrInvalidInput
	}
	return s.repo.Create(ctx, domain.Application{
		ContactName:  input.ContactName,
		ContactEmail: input.ContactEmail,
		ContactPhone: input.ContactPhone,
		BusinessName: input.BusinessName,
		CategoryCode: input.CategoryCode,
		Address:      input.Address,
		Description:  input.Description,
	})
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

	if app.Status != domain.ApplicationStatusPending {
		return errors.ErrInvalidStatusTransition
	}

	partner, err := s.partnerSvc.Create(ctx, app.BusinessName)
	if err != nil {
		return err
	}

	password := generator.GeneratePassword()
	passwordHash, err := auth.Hash(password)
	if err != nil {
		return err
	}

	if _, err := s.employeeSvc.Create(ctx, partner.ID, app.ContactEmail, passwordHash, app.ContactName, domain.EmployeeRoleOwner); err != nil {
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

	if app.Status != domain.ApplicationStatusPending {
		return errors.ErrInvalidStatusTransition
	}

	return s.repo.UpdateStatus(ctx, id, domain.ApplicationStatusRejected, reason)
}
