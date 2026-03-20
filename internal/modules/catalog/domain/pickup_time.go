package domain

import (
	"fmt"
	"time"

	"github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
)

const pickupTimeLayout = "15:04"

type PickupTime struct {
	Start time.Time
	End   time.Time
}

func NewPickupTime(start, end time.Time) (PickupTime, error) {
	if !end.After(start) {
		return PickupTime{}, errors.ErrInvalidPickupTimeRange
	}

	return PickupTime{
		Start: start,
		End:   end,
	}, nil
}

func NewPickupTimeFromStrings(startStr, endStr string) (PickupTime, error) {
	start, err := time.Parse(pickupTimeLayout, startStr)
	if err != nil {
		return PickupTime{}, fmt.Errorf("invalid pickup time format for start time: %w", err)
	}

	end, err := time.Parse(pickupTimeLayout, endStr)
	if err != nil {
		return PickupTime{}, fmt.Errorf("invalid pickup time format for end time: %w", err)
	}

	return NewPickupTime(start, end)
}
