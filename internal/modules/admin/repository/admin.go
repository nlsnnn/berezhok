package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/admin/domain"
	authErrors "github.com/nlsnnn/berezhok/internal/modules/auth"
)

type AdminRepo struct {
	q *sqlc.Queries
}

func NewAdminRepo(q *sqlc.Queries) *AdminRepo {
	return &AdminRepo{q: q}
}

func (r *AdminRepo) FindByEmail(ctx context.Context, email string) (domain.AdminUser, error) {
	admin, err := r.q.FindAdminByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AdminUser{}, authErrors.ErrUserNotFound
		}
		return domain.AdminUser{}, err
	}

	return toDomainAdmin(admin), nil
}

func (r *AdminRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.AdminUser, error) {
	admin, err := r.q.FindAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AdminUser{}, authErrors.ErrUserNotFound
		}
		return domain.AdminUser{}, err
	}

	return toDomainAdmin(admin), nil
}

func (r *AdminRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	return r.q.UpdateAdminLastLogin(ctx, id)
}

func toDomainAdmin(admin sqlc.AdminUser) domain.AdminUser {
	var lastLoginAt *time.Time
	if admin.LastLoginAt.Valid {
		value := admin.LastLoginAt.Time
		lastLoginAt = &value
	}

	return domain.AdminUser{
		ID:           admin.ID,
		Email:        admin.Email,
		PasswordHash: admin.PasswordHash,
		Name:         admin.Name,
		Role:         domain.AdminRole(admin.Role),
		IsActive:     admin.IsActive,
		LastLoginAt:  lastLoginAt,
		CreatedAt:    admin.CreatedAt,
		UpdatedAt:    admin.UpdatedAt,
	}
}
