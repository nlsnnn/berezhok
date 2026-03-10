package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	"github.com/nlsnnn/berezhok/internal/shared/generator"
)

type appService struct {
	repo            sqlc.Querier
	partService     partnerSvc
	employeeService employeeSvc
}

type partnerSvc interface {
	Create(ctx context.Context, arg sqlc.CreatePartnerParams) (sqlc.Partner, error)
}

type employeeSvc interface {
	Create(ctx context.Context, arg sqlc.CreatePartnerEmployeeParams) (sqlc.PartnerEmployee, error)
}

func NewApplicationService(appRepo sqlc.Querier, partSvc partnerSvc, employeeSvc employeeSvc) *appService {
	return &appService{
		repo:            appRepo,
		partService:     partSvc,
		employeeService: employeeSvc,
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

	partner, err := a.partService.Create(ctx, sqlc.CreatePartnerParams{
		LegalName: app.BusinessName,
		Status:    "pending_documents",
	})
	if err != nil {
		return err
	}

	password := generator.GeneratePassword()
	passwordHash, err := auth.Hash(password)
	if err != nil {
		return err
	}

	if _, err := a.employeeService.Create(ctx, sqlc.CreatePartnerEmployeeParams{
		PartnerID:    partner.ID,
		Name:         pgtype.Text{String: app.ContactName, Valid: true},
		Email:        app.ContactEmail,
		Role:         "owner",
		PasswordHash: passwordHash,
	}); err != nil {
		return err
	}

	// send email with credentials
	fmt.Printf("Partner approved. Contact: %s, Password: %s\n", app.ContactEmail, password)

	return a.Update(ctx, sqlc.UpdateApplicationParams{
		ID:         id,
		Status:     "approved",
		ReviewedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
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
		ReviewedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
}
