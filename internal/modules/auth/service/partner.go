package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/auth"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	hasher "github.com/nlsnnn/berezhok/internal/shared/auth"
)

type employeeFinder interface {
	FindByEmail(ctx context.Context, email string) (domain.Employee, error)
}

type partnerAuthenticator struct {
	repo         employeeFinder
	tokenService TokenService
}

func NewPartnerAuthenticator(repo employeeFinder, tokenService TokenService) *partnerAuthenticator {
	return &partnerAuthenticator{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (a *partnerAuthenticator) Authenticate(ctx context.Context, email, password string) (*auth.TokenClaims, error) {
	emp, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !hasher.Compare(emp.PasswordHash, password) {
		return nil, auth.ErrInvalidCredentials
	}

	userID, err := uuid.Parse(emp.ID)
	if err != nil {
		return nil, err
	}
	partnerID, err := uuid.Parse(emp.PartnerID)
	if err != nil {
		return nil, err
	}

	claims := auth.TokenClaims{
		UserID:   userID,
		UserType: "partner",
		Role:     string(emp.Role),
		UserData: partnerID,
	}
	token, err := a.tokenService.Generate(claims)
	if err != nil {
		return nil, err
	}

	claims.Access = token
	claims.UserData = emp
	return &claims, nil
}
