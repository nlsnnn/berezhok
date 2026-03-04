package dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ToText converts optional string to pgtype.Text
func ToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// ToTimestamptz converts time.Time to pgtype.Timestamptz
func ToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: !t.IsZero()}
}

// MapSlice applies a mapper function to a slice
func MapSlice[T, U any](data []T, mapper func(T) U) []U {
	if data == nil {
		return nil
	}
	result := make([]U, len(data))
	for i, v := range data {
		result[i] = mapper(v)
	}
	return result
}
