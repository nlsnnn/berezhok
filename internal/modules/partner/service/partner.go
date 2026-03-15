package service

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
)

type ChangePasswordInput struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

type partService struct {
	repo    partnerRepo
	empRepo employeeRepoForPartner
}

type partnerRepo interface {
	FindByID(ctx context.Context, id string) (domain.Partner, error)
	List(ctx context.Context) ([]domain.Partner, error)
	Create(ctx context.Context, legalName string) (domain.Partner, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	GetProfile(ctx context.Context, employeeID string) (domain.PartnerProfile, error)
	UpdateEmployeePassword(ctx context.Context, employeeID, newHash string) error
}

type employeeRepoForPartner interface {
	FindByID(ctx context.Context, id string) (domain.Employee, error)
}

func NewPartnerService(repo partnerRepo, empRepo employeeRepoForPartner) *partService {
	return &partService{repo: repo, empRepo: empRepo}
}

func (s *partService) List(ctx context.Context) ([]domain.Partner, error) {
	return s.repo.List(ctx)
}

func (s *partService) FindByID(ctx context.Context, id string) (domain.Partner, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *partService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	emailExists, err := s.repo.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}

	return emailExists, nil
}

func (s *partService) Create(ctx context.Context, legalName string) (domain.Partner, error) {
	return s.repo.Create(ctx, legalName)
}

func (s *partService) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	employee, err := s.empRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if !auth.Compare(employee.PasswordHash, input.CurrentPassword) {
		return errors.ErrInvalidCredentials
	}

	if input.CurrentPassword == input.NewPassword {
		return errors.ErrPasswordUnchanged
	}

	newHash, err := auth.Hash(input.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdateEmployeePassword(ctx, input.UserID, newHash)
}

func (s *partService) Profile(ctx context.Context, userID string) (domain.PartnerProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}
