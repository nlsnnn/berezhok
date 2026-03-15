package errors

import "errors"

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrEmailAlreadyInUse       = errors.New("email is already in use")

	ErrPartnerNotFound          = errors.New("partner not found")
	ErrPasswordUnchanged        = errors.New("password must be different from the current one")
	ErrLocationCategoryNotFound = errors.New("location category not found")
	ErrInvalidCommissionRate    = errors.New("commission rate must be between 0 and 100")
)
