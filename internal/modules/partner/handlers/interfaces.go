package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/domain"
)

type appSvc interface {
	Create(ctx context.Context, arg sqlc.CreateApplicationParams) (sqlc.PartnerApplication, error)
	GetByID(ctx context.Context, id uuid.UUID) (sqlc.PartnerApplication, error)
	List(ctx context.Context) ([]sqlc.PartnerApplication, error)
	Approve(ctx context.Context, id uuid.UUID) error
	Reject(ctx context.Context, id uuid.UUID, reason string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type partnerSvc interface {
	List(ctx context.Context) ([]sqlc.Partner, error)
	FindByID(ctx context.Context, id uuid.UUID) (sqlc.Partner, error)
	Create(ctx context.Context, arg sqlc.CreatePartnerParams) (sqlc.Partner, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
	Profile(ctx context.Context, userID uuid.UUID) (sqlc.GetPartnerProfileRow, error)
}

type locationSvc interface {
	Create(ctx context.Context, partnerID uuid.UUID, code string, name string, address string, latitude float64, longitude float64) (domain.Location, error)
}
