package service

import (
	"context"
	"unicode"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
)

type ChangePasswordInput struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

type AddLegalInfoInput struct {
	PartnerID    string
	Inn          string
	Ogrn         string
	Kpp          string
	LegalAddress string
	LegalName    string
}

type partService struct {
	repo    partnerRepo
	empRepo employeeRepoForPartner
}

type partnerRepo interface {
	FindByID(ctx context.Context, id string) (domain.Partner, error)
	List(ctx context.Context) ([]domain.Partner, error)
	Create(ctx context.Context, name string) (domain.Partner, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	GetProfile(ctx context.Context, employeeID string) (domain.PartnerProfile, error)
	GetDashboard(ctx context.Context, employeeID string) (domain.PartnerDashboard, error)
	UpdateEmployeePassword(ctx context.Context, employeeID, newHash string) error
	UpsertLegalInfo(ctx context.Context, info domain.LegalInfo) error
	UpdateStatus(ctx context.Context, partnerID string, status domain.PartnerStatus) error
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

func (s *partService) Create(ctx context.Context, name string) (domain.Partner, error) {
	return s.repo.Create(ctx, name)
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

func (s *partService) Dashboard(ctx context.Context, userID string) (domain.PartnerDashboard, error) {
	return s.repo.GetDashboard(ctx, userID)
}

func (s *partService) AddLegalInfo(ctx context.Context, input AddLegalInfoInput) error {
	// TODO: move validation to handler (?)
	if !isDigitsOnly(input.Inn) || (len(input.Inn) != 10 && len(input.Inn) != 12) {
		return errors.ErrInvalidINN
	}

	if input.Ogrn != "" && (!isDigitsOnly(input.Ogrn) || (len(input.Ogrn) != 13 && len(input.Ogrn) != 15)) {
		return errors.ErrInvalidOGRN
	}

	if input.Kpp != "" && (!isDigitsOnly(input.Kpp) || len(input.Kpp) != 9) {
		return errors.ErrInvalidKPP
	}

	partner, err := s.repo.FindByID(ctx, input.PartnerID)
	if err != nil {
		return err
	}

	if partner.Status != domain.PartnerStatusPendingDocuments && partner.Status != domain.PartnerStatusActive {
		return errors.ErrPartnerStatusInvalid
	}

	err = s.repo.UpsertLegalInfo(ctx, domain.LegalInfo{
		PartnerID:    input.PartnerID,
		Inn:          input.Inn,
		Ogrn:         input.Ogrn,
		Kpp:          input.Kpp,
		LegalAddress: input.LegalAddress,
		LegalName:    input.LegalName,
	})
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, input.PartnerID, domain.PartnerStatusActive)
}

func (s *partService) CanActivateBoxes(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	partner, err := s.repo.FindByID(ctx, partnerID.String())
	if err != nil {
		return false, err
	}

	return partner.Status != domain.PartnerStatusPendingDocuments, nil
}

func isDigitsOnly(value string) bool {
	for _, r := range value {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
