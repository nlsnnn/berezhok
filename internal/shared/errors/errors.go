package errors

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrNotFound        = errors.New("not found")
	ErrInvalidGeoPoint = errors.New("invalid geo point")

	// Authentication
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrNotFoundContextValue = errors.New("required value not found in context")

	// Phone
	ErrInvalidPhoneNumber = errors.New("invalid phone number")

	// Pickup time
	ErrInvalidPickupTimeRange = errors.New("invalid pickup time range")
)
