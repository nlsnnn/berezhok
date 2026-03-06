package service

import (
	"context"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
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
