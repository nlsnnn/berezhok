package errors

import "errors"

var (
	ErrInvalidPickupTimeFormat = errors.New("invalid pickup time format")
	ErrInvalidPickupTimeRange  = errors.New("pickup_time_end must be after pickup_time_start")
	ErrInvalidBoxID            = errors.New("invalid box id")
	ErrBoxNotFound             = errors.New("box not found")
)
