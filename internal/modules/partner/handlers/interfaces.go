package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
)

type appSvc interface {
	Create(ctx context.Context, arg sqlc.CreateApplicationParams) (sqlc.PartnerApplication, error)
	GetByID(ctx context.Context, id uuid.UUID) (sqlc.PartnerApplication, error)
	List(ctx context.Context) ([]sqlc.PartnerApplication, error)
	Approve(ctx context.Context, id uuid.UUID) error
	Reject(ctx context.Context, id uuid.UUID, reason string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
