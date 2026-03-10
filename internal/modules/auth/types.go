package auth

import "github.com/google/uuid"

type TokenClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	UserType string    `json:"user_type"` // "customer" | "partner" | "admin"
	Role     string    `json:"role,omitempty"`
	Access   string    `json:"access,omitempty"`
	Refresh  string    `json:"refresh,omitempty"`
	UserData any       `json:"user_data,omitempty"`
}
