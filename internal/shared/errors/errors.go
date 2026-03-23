package errors

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrNotFound        = errors.New("not found")
	ErrInvalidGeoPoint = errors.New("invalid geo point")

	// Phone
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)
