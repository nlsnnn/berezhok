package repository

import (
	"context"
	"encoding/json"
	"errors"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

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

func (r *AdminRepo) MarkApplicationReviewed(ctx context.Context, applicationID, adminID uuid.UUID) error {
	_, err := r.q.MarkApplicationReviewed(ctx, sqlc.MarkApplicationReviewedParams{
		ID: applicationID,
		ReviewedBy: pgtype.UUID{
			Bytes: adminID,
			Valid: true,
		},
	})
	return err
}

func (r *AdminRepo) CreateAuditLog(ctx context.Context, input domain.AuditLog) (domain.AuditLog, error) {
	entityType := pgtype.Text{String: input.EntityType, Valid: input.EntityType != ""}
	entityID := pgtype.UUID{}
	if input.EntityID != nil {
		entityID = pgtype.UUID{Bytes: *input.EntityID, Valid: true}
	}

	details := input.Details
	if details == nil {
		details = json.RawMessage(`{}`)
	}

	row, err := r.q.CreateAdminAuditLog(ctx, sqlc.CreateAdminAuditLogParams{
		AdminUserID: input.AdminID,
		Action:      input.Action,
		EntityType:  entityType,
		EntityID:    entityID,
		Details:     details,
		IpAddress:   input.IPAddress,
	})
	if err != nil {
		return domain.AuditLog{}, err
	}

	return toDomainAudit(row), nil
}

func toDomainAudit(row sqlc.AdminAuditLog) domain.AuditLog {
	var entityID *uuid.UUID
	if row.EntityID.Valid {
		value := uuid.UUID(row.EntityID.Bytes)
		entityID = &value
	}

	var ipAddress *netip.Addr
	if row.IpAddress != nil {
		value := *row.IpAddress
		ipAddress = &value
	}

	return domain.AuditLog{
		ID:         row.ID,
		AdminID:    row.AdminUserID,
		Action:     row.Action,
		EntityType: row.EntityType.String,
		EntityID:   entityID,
		Details:    row.Details,
		IPAddress:  ipAddress,
		CreatedAt:  row.CreatedAt,
	}
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
