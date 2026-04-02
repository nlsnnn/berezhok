package contextx

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	sharedErrors "github.com/nlsnnn/berezhok/internal/shared/errors"
)

type key string

const (
	customerIDKey key = "customer_id"
	partnerIDKey  key = "partner_id"
	employeeIDKey key = "employee_id"
)

// Retrieves the customer ID from the request context
func CustomerID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), customerIDKey)
}

// Retrieves the customer ID from a context.Context
func CustomerIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, customerIDKey)
}

// Retrieves the partner ID from the request context
func PartnerID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), partnerIDKey)
}

// Retrieves the partner ID from a context.Context
func PartnerIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, partnerIDKey)
}

// EmployeeID retrieves the employee ID from the request context
func EmployeeID(r *http.Request) (uuid.UUID, error) {
	return getIDFromCtx(r.Context(), employeeIDKey)
}

// EmployeeID retrieves the employee ID from the request context
func EmployeeIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return getIDFromCtx(ctx, employeeIDKey)
}

// helper function to extract a UUID from the context using the provided key
func getIDFromCtx(ctx context.Context, key key) (uuid.UUID, error) {
	if ctx == nil {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}
	val := ctx.Value(key)
	if val == nil {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}
	return id, nil
}
