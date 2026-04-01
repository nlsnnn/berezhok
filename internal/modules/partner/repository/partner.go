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

func (r *PartnerRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return r.q.CheckEmailExists(ctx, email)
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

func (r *PartnerRepo) Create(ctx context.Context, name string) (domain.Partner, error) {
	p, err := r.q.CreatePartner(ctx, sqlc.CreatePartnerParams{
		BrandName:      name,
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

	commission, err := getCommissionFromRow(row)
	if err != nil {
		return domain.PartnerProfile{}, err
	}

	profile := domain.PartnerProfile{
		Partner: domain.Partner{
			ID:         row.PartnerID.String(),
			BrandName:  row.BrandName,
			Status:     domain.PartnerStatus(row.PartnerStatus),
			Commission: commission,
			CreatedAt:  row.PartnerCreatedAt,
		},
		Employee: domain.Employee{
			ID:                 row.EmployeeID.String(),
			Email:              row.Email,
			Name:               row.EmployeeName.String,
			Role:               domain.EmployeeRole(row.Role),
			MustChangePassword: row.MustChangePassword.Bool,
			CreatedAt:          row.EmployeeCreatedAt,
		},
	}

	if row.LocationID.Valid {
		profile.Location = &domain.LocationSummary{
			ID:        row.LocationID.String(),
			Name:      row.LocationName.String,
			Address:   row.LocationAddress.String,
			CreatedAt: row.LocationCreatedAt.Time,
		}
	}

	// Get all locations for the partner
	locationRows, err := r.q.FindLocationsByPartnerID(ctx, row.PartnerID)
	if err != nil {
		return domain.PartnerProfile{}, err
	}

	locations := make([]domain.LocationSummary, len(locationRows))
	for i, loc := range locationRows {
		locations[i] = domain.LocationSummary{
			ID:        loc.ID.String(),
			Name:      loc.Name,
			Address:   loc.Address,
			CreatedAt: time.Now(), // TODO: Add CreatedAt to the SQL query and use it here
		}
	}
	profile.Locations = locations

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
	commission, err := getCommission(p)
	if err != nil {
		commission = domain.Commission{}
	}

	return domain.Partner{
		ID:         p.ID.String(),
		BrandName:  p.BrandName,
		LogoURL:    p.LogoUrl.String,
		Status:     domain.PartnerStatus(p.Status),
		Commission: commission,
		CreatedAt:  p.CreatedAt,
	}
}

func getCommissionFromRow(row sqlc.GetPartnerProfileRow) (domain.Commission, error) {
	commissionRate, _ := numericToFloat64(row.CommissionRate.(pgtype.Numeric))

	var promoUntil *time.Time
	if row.PromoCommissionUntil.Valid {
		promoUntil = &row.PromoCommissionUntil.Time
	}

	commission, err := domain.NewCommission(commissionRate, promoUntil)
	if err != nil {
		return domain.Commission{}, err
	}

	return commission, nil
}

func getCommission(p sqlc.Partner) (domain.Commission, error) {
	commissionRate, _ := p.CommissionRate.Int.Float64()

	var promoUntil *time.Time
	if p.PromoCommissionUntil.Valid {
		promoUntil = &p.PromoCommissionUntil.Time
	}

	commission, err := domain.NewCommission(commissionRate, promoUntil)
	if err != nil {
		return domain.Commission{}, err
	}

	return commission, nil
}

func numericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid || n.Int == nil {
		return 0, nil
	}

	f := new(big.Float).SetInt(n.Int)
	if n.Exp > 0 {
		mul := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil))
		f.Mul(f, mul)
	} else if n.Exp < 0 {
		div := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil))
		f.Quo(f, div)
	}

	result, _ := f.Float64()
	return result, nil
}
