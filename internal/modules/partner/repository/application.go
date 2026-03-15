package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type ApplicationRepo struct {
	q *sqlc.Queries
}

func NewApplicationRepo(q *sqlc.Queries) *ApplicationRepo {
	return &ApplicationRepo{q: q}
}

func (r *ApplicationRepo) FindByID(ctx context.Context, id string) (domain.Application, error) {
	uid := uuid.MustParse(id)
	a, err := r.q.FindApplicationByID(ctx, uid)
	if err != nil {
		return domain.Application{}, err
	}
	return applicationToDomain(a), nil
}

func (r *ApplicationRepo) List(ctx context.Context) ([]domain.Application, error) {
	rows, err := r.q.ListApplications(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Application, len(rows))
	for i, a := range rows {
		result[i] = applicationToDomain(a)
	}
	return result, nil
}

func (r *ApplicationRepo) Create(ctx context.Context, app domain.Application) (domain.Application, error) {
	a, err := r.q.CreateApplication(ctx, sqlc.CreateApplicationParams{
		ContactName:  app.ContactName,
		ContactEmail: app.ContactEmail,
		ContactPhone: app.ContactPhone,
		BusinessName: app.BusinessName,
		CategoryCode: pgtype.Text{String: app.CategoryCode, Valid: app.CategoryCode != ""},
		Address:      pgtype.Text{String: app.Address, Valid: app.Address != ""},
		Description:  pgtype.Text{String: app.Description, Valid: app.Description != ""},
		Status:       string(domain.ApplicationStatusPending),
		Latitude:     pgtype.Float8{Float64: app.Coords.Latitude, Valid: app.Coords.Latitude != 0},
		Longitude:    pgtype.Float8{Float64: app.Coords.Longitude, Valid: app.Coords.Longitude != 0},
	})
	if err != nil {
		return domain.Application{}, err
	}
	return applicationToDomain(a), nil
}

func (r *ApplicationRepo) UpdateStatus(ctx context.Context, id string, status domain.ApplicationStatus, rejectionReason string) error {
	uid := uuid.MustParse(id)
	return r.q.UpdateApplication(ctx, sqlc.UpdateApplicationParams{
		ID:              uid,
		Status:          string(status),
		ReviewedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		RejectionReason: pgtype.Text{String: rejectionReason, Valid: rejectionReason != ""},
	})
}

func (r *ApplicationRepo) Delete(ctx context.Context, id string) error {
	uid := uuid.MustParse(id)
	return r.q.DeleteApplication(ctx, uid)
}

func applicationToDomain(a sqlc.PartnerApplication) domain.Application {
	var reviewedAt *time.Time
	if a.ReviewedAt.Valid {
		t := a.ReviewedAt.Time
		reviewedAt = &t
	}
	return domain.Application{
		ID:              a.ID.String(),
		ContactName:     a.ContactName,
		ContactEmail:    a.ContactEmail,
		ContactPhone:    a.ContactPhone,
		BusinessName:    a.BusinessName,
		CategoryCode:    a.CategoryCode.String,
		Address:         a.Address.String,
		Description:     a.Description.String,
		Status:          domain.ApplicationStatus(a.Status),
		ReviewedAt:      reviewedAt,
		RejectionReason: a.RejectionReason.String,
		CreatedAt:       a.CreatedAt,
		Coords: sharedDomain.GeoPoint{
			Latitude:  a.Latitude.Float64,
			Longitude: a.Longitude.Float64,
		},
	}
}
