package errors

import "errors"

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrInvalidPaymentID = errors.New("invalid payment id")
)
