package partner

import "errors"

var (
	ErrInvalidInput            = errors.New("invalid input")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
