package errors

import "errors"

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrInvalidPaymentID = errors.New("invalid payment id")

	ErrPaymentAlreadyExists  = errors.New("payment already exists for this order")
	ErrPaymentCreationFailed = errors.New("failed to create payment")
	ErrPaymentUpdateFailed   = errors.New("failed to update payment status")
	ErrPaymentAlreadyHandled = errors.New("payment event already handled")
)
