package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/admin/domain"
	"github.com/nlsnnn/berezhok/internal/modules/auth"
	hasher "github.com/nlsnnn/berezhok/internal/shared/auth"
)

type adminFinder interface {
	FindByEmail(ctx context.Context, email string) (domain.AdminUser, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

type adminAuthenticator struct {
	repo         adminFinder
	tokenService TokenService
}

func NewAdminAuthenticator(repo adminFinder, tokenService TokenService) *adminAuthenticator {
	return &adminAuthenticator{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (a *adminAuthenticator) Authenticate(ctx context.Context, email, password string) (*auth.TokenClaims, error) {
	admin, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !admin.IsActive || !hasher.Compare(admin.PasswordHash, password) {
		return nil, auth.ErrInvalidCredentials
	}

	claims := auth.TokenClaims{
		UserID:   admin.ID,
		UserType: "admin",
		Role:     string(admin.Role),
		UserData: admin,
	}

	token, err := a.tokenService.Generate(claims)
	if err != nil {
		return nil, err
	}

	if err := a.repo.UpdateLastLogin(ctx, admin.ID); err != nil {
		return nil, err
	}

	claims.Access = token
	return &claims, nil
}
