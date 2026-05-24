package domain

import (
	"time"

	"github.com/google/uuid"
)

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleAdmin      AdminRole = "admin"
	AdminRoleSupport    AdminRole = "support"
)

type AdminUser struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	Role         AdminRole
	IsActive     bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
