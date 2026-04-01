package errors

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderID     = errors.New("invalid order id")
	ErrBoxNotAvailable    = errors.New("box is not available or out of stock")
	ErrInvalidBoxStatus   = errors.New("box status is not active")
	ErrInvalidOrderStatus = errors.New("invalid order status transition")
	ErrOrderNotReady      = errors.New("order is not ready for pickup")
	ErrInvalidCustomerID  = errors.New("invalid customer id")
	ErrPaymentFailed      = errors.New("payment creation failed")
)
