package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner"
)

type appService struct {
	repo sqlc.Querier
}

func NewApplicationService(appRepo sqlc.Querier) *appService {
	return &appService{
		repo: appRepo,
	}
}

func (a *appService) Create(ctx context.Context, arg sqlc.CreateApplicationParams) (sqlc.PartnerApplication, error) {
	if arg.ContactName == "" {
		return sqlc.PartnerApplication{}, partner.ErrInvalidInput
	}
	return a.repo.CreateApplication(ctx, arg)
}

func (a *appService) GetByID(ctx context.Context, id uuid.UUID) (sqlc.PartnerApplication, error) {
	return a.repo.FindApplicationByID(ctx, id)
}

func (a *appService) List(ctx context.Context) ([]sqlc.PartnerApplication, error) {
	return a.repo.ListApplications(ctx)
}

func (a *appService) Update(ctx context.Context, arg sqlc.UpdateApplicationParams) error {
	return a.repo.UpdateApplication(ctx, arg)
}

func (a *appService) Delete(ctx context.Context, id uuid.UUID) error {
	return a.repo.DeleteApplication(ctx, id)
}

func (a *appService) Approve(ctx context.Context, id uuid.UUID) error {
	app, err := a.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if app.Status != "pending" {
		return partner.ErrInvalidStatusTransition
	}

	return a.Update(ctx, sqlc.UpdateApplicationParams{
		ID:     id,
		Status: "approved",
	})
}

func (a *appService) Reject(ctx context.Context, id uuid.UUID, reason string) error {
	app, err := a.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if app.Status != "pending" {
		return partner.ErrInvalidStatusTransition
	}

	return a.Update(ctx, sqlc.UpdateApplicationParams{
		ID:              id,
		Status:          "rejected",
		RejectionReason: pgtype.Text{String: reason, Valid: true},
	})
}
