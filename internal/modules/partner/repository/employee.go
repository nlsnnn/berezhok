package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
)

type EmployeeRepo struct {
	q *sqlc.Queries
}

func NewEmployeeRepo(q *sqlc.Queries) *EmployeeRepo {
	return &EmployeeRepo{q: q}
}

func (r *EmployeeRepo) FindByID(ctx context.Context, id string) (domain.Employee, error) {
	uid := uuid.MustParse(id)
	e, err := r.q.FindPartnerEmployeeByID(ctx, uid)
	if err != nil {
		return domain.Employee{}, err
	}
	return employeeToDomain(e), nil
}

func (r *EmployeeRepo) FindByEmail(ctx context.Context, email string) (domain.Employee, error) {
	e, err := r.q.FindPartnerEmployeeByEmail(ctx, email)
	if err != nil {
		return domain.Employee{}, err
	}
	return employeeToDomain(e), nil
}

func (r *EmployeeRepo) List(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.q.ListPartnerEmployees(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Employee, len(rows))
	for i, e := range rows {
		result[i] = employeeToDomain(e)
	}
	return result, nil
}

func (r *EmployeeRepo) ListByPartnerID(ctx context.Context, partnerID string) ([]domain.Employee, error) {
	uid := uuid.MustParse(partnerID)
	rows, err := r.q.ListEmployeesByPartnerID(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Employee, len(rows))
	for i, e := range rows {
		result[i] = employeeToDomain(e)
	}
	return result, nil
}

func (r *EmployeeRepo) Create(ctx context.Context, partnerID string, email, passwordHash, name string, role domain.EmployeeRole) (domain.Employee, error) {
	uid := uuid.MustParse(partnerID)
	e, err := r.q.CreatePartnerEmployee(ctx, sqlc.CreatePartnerEmployeeParams{
		PartnerID:    uid,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         string(role),
		Name:         pgtype.Text{String: name, Valid: name != ""},
	})
	if err != nil {
		return domain.Employee{}, err
	}
	return employeeToDomain(e), nil
}

func (r *EmployeeRepo) Delete(ctx context.Context, id string) error {
	uid := uuid.MustParse(id)
	return r.q.DeletePartnerEmployee(ctx, uid)
}

func employeeToDomain(e sqlc.PartnerEmployee) domain.Employee {
	locationID := ""
	if e.LocationID.Valid {
		locationID = e.LocationID.String()
	}
	return domain.Employee{
		ID:                 e.ID.String(),
		PartnerID:          e.PartnerID.String(),
		LocationID:         locationID,
		Email:              e.Email,
		PasswordHash:       e.PasswordHash,
		Role:               domain.EmployeeRole(e.Role),
		Name:               e.Name.String,
		MustChangePassword: e.MustChangePassword.Bool,
		CreatedAt:          e.CreatedAt,
	}
}
