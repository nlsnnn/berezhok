package domain

import (
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID
	AdminID    uuid.UUID
	Action     string
	EntityType string
	EntityID   *uuid.UUID
	Details    []byte
	IPAddress  *netip.Addr
	CreatedAt  time.Time
	AdminEmail string
	AdminName  string
}
