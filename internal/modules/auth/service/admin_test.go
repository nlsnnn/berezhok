package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/admin/domain"
	"github.com/nlsnnn/berezhok/internal/modules/auth"
	hasher "github.com/nlsnnn/berezhok/internal/shared/auth"
)

type adminFinderStub struct {
	admin domain.AdminUser
	err   error
}

type tokenGeneratorStub struct {
	token string
}

func (s *tokenGeneratorStub) Generate(claims auth.TokenClaims) (string, error) {
	if s.token == "" {
		return "token", nil
	}
	return s.token, nil
}

func (s *tokenGeneratorStub) Validate(tokenString string) (*auth.TokenClaims, error) {
	return nil, nil
}

func (s adminFinderStub) FindByEmail(ctx context.Context, email string) (domain.AdminUser, error) {
	return s.admin, s.err
}

func (s adminFinderStub) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestAdminAuthenticatorRejectsInactiveAdmin(t *testing.T) {
	t.Parallel()

	passwordHash, err := hasher.Hash("secret123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	authenticator := NewAdminAuthenticator(adminFinderStub{
		admin: domain.AdminUser{
			ID:           uuid.New(),
			Email:        "admin@berezhok.local",
			PasswordHash: passwordHash,
			Role:         domain.AdminRoleAdmin,
			IsActive:     false,
		},
	}, &tokenGeneratorStub{})

	_, err = authenticator.Authenticate(context.Background(), "admin@berezhok.local", "secret123")
	if !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAdminAuthenticatorReturnsAdminClaims(t *testing.T) {
	t.Parallel()

	adminID := uuid.New()
	passwordHash, err := hasher.Hash("secret123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	authenticator := NewAdminAuthenticator(adminFinderStub{
		admin: domain.AdminUser{
			ID:           adminID,
			Email:        "admin@berezhok.local",
			PasswordHash: passwordHash,
			Role:         domain.AdminRoleSuperAdmin,
			IsActive:     true,
		},
	}, &tokenGeneratorStub{token: "admin-token"})

	claims, err := authenticator.Authenticate(context.Background(), "admin@berezhok.local", "secret123")
	if err != nil {
		t.Fatalf("authenticate: %v", err)
	}

	if claims.UserID != adminID {
		t.Fatalf("expected user id %s, got %s", adminID, claims.UserID)
	}
	if claims.UserType != "admin" {
		t.Fatalf("expected user type admin, got %s", claims.UserType)
	}
	if claims.Role != "super_admin" {
		t.Fatalf("expected role super_admin, got %s", claims.Role)
	}
	if claims.Access != "admin-token" {
		t.Fatalf("expected token admin-token, got %s", claims.Access)
	}
}
