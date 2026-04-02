package errors

import "errors"

var (
	ErrOrderNotCompleted      = errors.New("order is not completed")
	ErrOrderOwnershipMismatch = errors.New("order does not belong to customer")
	ErrReviewAlreadyExists    = errors.New("review already exists")
	ErrRatingOutOfRange       = errors.New("rating must be between 1 and 5")
)
