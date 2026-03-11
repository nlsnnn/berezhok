package service

import (
	"context"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
)

type partService struct {
	repo sqlc.Querier
}

func NewPartnerService(repo sqlc.Querier) *partService {
	return &partService{repo: repo}
}

func (s *partService) List(ctx context.Context) ([]sqlc.Partner, error) {
	return s.repo.ListPartners(ctx)
}

func (s *partService) FindByID(ctx context.Context, id uuid.UUID) (sqlc.Partner, error) {
	return s.repo.FindPartnerByID(ctx, id)
}

func (s *partService) Create(ctx context.Context, arg sqlc.CreatePartnerParams) (sqlc.Partner, error) {
	arg.CommissionRate = pgtype.Numeric{Int: big.NewInt(20), Exp: -2, Valid: true} // default commission rate
	return s.repo.CreatePartner(ctx, arg)
}

func (s *partService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.repo.FindPartnerEmployeeByID(ctx, userID)
	if err != nil {
		return err
	}

	if !auth.Compare(user.PasswordHash, currentPassword) {
		return partner.ErrInvalidCredentials
	}

	if currentPassword == newPassword {
		return partner.ErrPasswordUnchanged
	}

	newPasswordHash, err := auth.Hash(newPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePartnerEmployeePassword(ctx, sqlc.UpdatePartnerEmployeePasswordParams{
		ID:                 userID,
		PasswordHash:       newPasswordHash,
		MustChangePassword: pgtype.Bool{Valid: true, Bool: false},
	})
}

func (s *partService) Profile(ctx context.Context, userID uuid.UUID) (sqlc.GetPartnerProfileRow, error) {
	return s.repo.GetPartnerProfile(ctx, userID)
}
