package contextx

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/shared/authz"
	sharedErrors "github.com/nlsnnn/berezhok/internal/shared/errors"
)

type key string

const (
	UserIDKey       key = "user_id"
	UserTypeKey     key = "user_type"
	UserRoleKey     key = "user_role"
	CustomerIDKey   key = "customer_id"
	AdminIDKey      key = "admin_id"
	PartnerIDKey    key = "partner_id"
	EmployeeIDKey   key = "employee_id"
	LocationIDKey   key = "location_id"
	PartnerActorKey key = "partner_actor"
	AdminActorKey   key = "admin_actor"
)

func UserID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), UserIDKey)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, UserIDKey)
}

func UserType(r *http.Request) (string, error) {
	return getStringFromCtx(r.Context(), UserTypeKey)
}

func UserRole(r *http.Request) (string, error) {
	return getStringFromCtx(r.Context(), UserRoleKey)
}

func CustomerID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), CustomerIDKey)
}

func CustomerIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, CustomerIDKey)
}

func AdminID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), AdminIDKey)
}

func AdminIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, AdminIDKey)
}

func PartnerID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), PartnerIDKey)
}

func PartnerIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, PartnerIDKey)
}

func EmployeeID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), EmployeeIDKey)
}

func EmployeeIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, EmployeeIDKey)
}

func LocationID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), LocationIDKey)
}

func PartnerActor(r *http.Request) (authz.PartnerActor, error) {
	return PartnerActorFromContext(r.Context())
}

func PartnerActorFromContext(ctx context.Context) (authz.PartnerActor, error) {
	if ctx == nil {
		return authz.PartnerActor{}, sharedErrors.ErrNotFoundContextValue
	}
	val := ctx.Value(PartnerActorKey)
	if val == nil {
		return authz.PartnerActor{}, sharedErrors.ErrNotFoundContextValue
	}
	actor, ok := val.(authz.PartnerActor)
	if !ok {
		return authz.PartnerActor{}, sharedErrors.ErrNotFoundContextValue
	}
	return actor, nil
}

func AdminActor(r *http.Request) (authz.AdminActor, error) {
	return AdminActorFromContext(r.Context())
}

func AdminActorFromContext(ctx context.Context) (authz.AdminActor, error) {
	if ctx == nil {
		return authz.AdminActor{}, sharedErrors.ErrNotFoundContextValue
	}
	val := ctx.Value(AdminActorKey)
	if val == nil {
		return authz.AdminActor{}, sharedErrors.ErrNotFoundContextValue
	}
	actor, ok := val.(authz.AdminActor)
	if !ok {
		return authz.AdminActor{}, sharedErrors.ErrNotFoundContextValue
	}
	return actor, nil
}

func getIDFromCtx(ctx context.Context, k key) (uuid.UUID, error) {
	if ctx == nil {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}
	val := ctx.Value(k)
	if val == nil {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}
	return id, nil
}

func getStringFromCtx(ctx context.Context, k key) (string, error) {
	if ctx == nil {
		return "", sharedErrors.ErrNotFoundContextValue
	}
	val := ctx.Value(k)
	if val == nil {
		return "", sharedErrors.ErrNotFoundContextValue
	}
	s, ok := val.(string)
	if !ok {
		return "", sharedErrors.ErrNotFoundContextValue
	}
	return s, nil
}
