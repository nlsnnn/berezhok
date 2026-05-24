package errors

import "errors"

var (
	ErrNoDestination       = errors.New("payout destination not configured")
	ErrNoOrders            = errors.New("no unsettled orders for payout period")
	ErrBelowMinAmount      = errors.New("payout amount below minimum threshold")
	ErrPayoutNotFound      = errors.New("payout not found")
	ErrDestinationNotFound = errors.New("payout destination not found")
)
