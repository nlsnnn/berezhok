package service

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/modules/auth"
	"github.com/nlsnnn/berezhok/internal/modules/customer/domain"
	"github.com/nlsnnn/berezhok/internal/shared/generator"
)

type userProvider interface {
	FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error)
}

type smsProvider interface {
	SendCode(ctx context.Context, phone, code string) error
	ValidateCode(ctx context.Context, phone, code string) (bool, error)
}

type customerAuthenticator struct {
	repo         userProvider
	tokenService TokenService
	smsProvider  smsProvider
}

func NewCustomerAuthenticator(repo userProvider, tokenService TokenService, smsProvider smsProvider) *customerAuthenticator {
	return &customerAuthenticator{
		repo:         repo,
		tokenService: tokenService,
		smsProvider:  smsProvider,
	}
}

func (a *customerAuthenticator) Authenticate(ctx context.Context, phone, code string) (*auth.TokenClaims, error) {
	isValid, err := a.smsProvider.ValidateCode(ctx, phone, code)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, auth.ErrInvalidCredentials
	}

	user, err := a.repo.FindOrCreateByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	claims := auth.TokenClaims{
		UserID:   user.ID,
		UserType: "customer",
	}
	token, err := a.tokenService.Generate(claims)
	if err != nil {
		return nil, err
	}

	claims.Access = token
	claims.UserData = user
	return &claims, nil
}

func (a *customerAuthenticator) SendCode(ctx context.Context, phone string) error {
	code := generator.GenerateOTP()
	if err := a.smsProvider.SendCode(ctx, phone, code); err != nil {
		return err
	}
	return nil
}
