package service

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/auth"
	hasher "github.com/nlsnnn/berezhok/internal/shared/auth"
)

type TokenService interface {
	Generate(claims auth.TokenClaims) (string, error)
	Validate(tokenString string) (*auth.TokenClaims, error)
}

type partnerAuthenticator struct {
	repo         *sqlc.Queries
	tokenService TokenService
}

func NewPartnerAuthenticator(repo *sqlc.Queries, tokenService TokenService) *partnerAuthenticator {
	return &partnerAuthenticator{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (a *partnerAuthenticator) Authenticate(ctx context.Context, email, password string) (*auth.TokenClaims, error) {
	part, err := a.repo.FindPartnerEmployeeByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !hasher.Compare(part.PasswordHash, password) {
		return nil, auth.ErrInvalidCredentials
	}

	claims := auth.TokenClaims{
		UserID:   part.ID,
		UserType: "partner",
		Role:     part.Role,
	}
	token, err := a.tokenService.Generate(claims)
	if err != nil {
		return nil, err
	}

	claims.Access = token
	claims.UserData = part
	return &claims, nil
}
