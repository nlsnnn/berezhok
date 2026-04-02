package contextx

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	sharedErrors "github.com/nlsnnn/berezhok/internal/shared/errors"
)

type key string

const (
	UserIDKey     key = "user_id"
	UserTypeKey   key = "user_type"
	CustomerIDKey key = "customer_id"
	PartnerIDKey  key = "partner_id"
	EmployeeIDKey key = "employee_id"
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

func CustomerID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), CustomerIDKey)
}

func CustomerIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, CustomerIDKey)
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
