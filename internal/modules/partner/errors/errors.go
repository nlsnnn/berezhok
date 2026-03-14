package errors

import "errors"

var (
	ErrInvalidInput             = errors.New("invalid input")
	ErrInvalidStatusTransition  = errors.New("invalid status transition")
	ErrPartnerNotFound          = errors.New("partner not found")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrPasswordUnchanged        = errors.New("password must be different from the current one")
	ErrLocationCategoryNotFound = errors.New("location category not found")
)
