package repository

import (
	"context"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
)

type PartnerRepo struct {
	q *sqlc.Queries
}

func NewPartnerRepo(q *sqlc.Queries) *PartnerRepo {
	return &PartnerRepo{q: q}
}

func (r *PartnerRepo) FindByID(ctx context.Context, id string) (domain.Partner, error) {
	uid := uuid.MustParse(id)
	p, err := r.q.FindPartnerByID(ctx, uid)
	if err != nil {
		return domain.Partner{}, err
	}
	return partnerToDomain(p), nil
}

func (r *PartnerRepo) List(ctx context.Context) ([]domain.Partner, error) {
	rows, err := r.q.ListPartners(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Partner, len(rows))
	for i, p := range rows {
		result[i] = partnerToDomain(p)
	}
	return result, nil
}

func (r *PartnerRepo) Create(ctx context.Context, legalName string) (domain.Partner, error) {
	p, err := r.q.CreatePartner(ctx, sqlc.CreatePartnerParams{
		LegalName:      legalName,
		Status:         string(domain.PartnerStatusPendingDocuments),
		CommissionRate: pgtype.Numeric{Int: big.NewInt(20), Exp: -2, Valid: true},
	})
	if err != nil {
		return domain.Partner{}, err
	}
	return partnerToDomain(p), nil
}

func (r *PartnerRepo) GetProfile(ctx context.Context, employeeID string) (domain.PartnerProfile, error) {
	uid := uuid.MustParse(employeeID)
	row, err := r.q.GetPartnerProfile(ctx, uid)
	if err != nil {
		return domain.PartnerProfile{}, err
	}

	commissionRate, _ := row.CommissionRate.Float64Value()

	profile := domain.PartnerProfile{
		Partner: domain.Partner{
			ID:             row.PartnerID.String(),
			LegalName:      row.LegalName,
			BrandName:      row.BrandName.String,
			Status:         domain.PartnerStatus(row.PartnerStatus),
			CommissionRate: commissionRate.Float64,
		},
		Employee: domain.Employee{
			ID:    row.EmployeeID.String(),
			Email: row.Email,
			Name:  row.EmployeeName.String,
			Role:  domain.EmployeeRole(row.Role),
		},
	}

	if row.PromoCommissionUntil.Valid {
		t := row.PromoCommissionUntil.Time
		profile.Partner.PromoCommissionUntil = &t
	}

	if row.LocationID.Valid {
		profile.Location = &domain.LocationSummary{
			ID:      row.LocationID.String(),
			Name:    row.LocationName.String,
			Address: row.LocationAddress.String,
		}
	}

	return profile, nil
}

func (r *PartnerRepo) UpdateEmployeePassword(ctx context.Context, employeeID, newHash string) error {
	uid := uuid.MustParse(employeeID)
	return r.q.UpdatePartnerEmployeePassword(ctx, sqlc.UpdatePartnerEmployeePasswordParams{
		ID:                 uid,
		PasswordHash:       newHash,
		MustChangePassword: pgtype.Bool{Valid: true, Bool: false},
	})
}

func partnerToDomain(p sqlc.Partner) domain.Partner {
	var promoUntil *time.Time
	if p.PromoCommissionUntil.Valid {
		t := p.PromoCommissionUntil.Time
		promoUntil = &t
	}

	commissionRate := 0.0
	if f, err := p.CommissionRate.Float64Value(); err == nil {
		commissionRate = f.Float64
	}

	return domain.Partner{
		ID:                   p.ID.String(),
		LegalName:            p.LegalName,
		BrandName:            p.BrandName.String,
		LogoURL:              p.LogoUrl.String,
		Status:               domain.PartnerStatus(p.Status),
		CommissionRate:       commissionRate,
		PromoCommissionUntil: promoUntil,
		CreatedAt:            p.CreatedAt,
	}
}
