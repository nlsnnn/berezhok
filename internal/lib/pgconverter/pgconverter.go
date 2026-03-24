package pgconverter

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func NumericToDecimalOrZero(v pgtype.Numeric) decimal.Decimal {
	if !v.Valid || v.Int == nil {
		return decimal.Zero
	}

	return decimal.NewFromBigInt(v.Int, v.Exp)
}

func TextToString(v pgtype.Text) string {
	if !v.Valid {
		return ""
	}

	return v.String
}

func TimeValue(v pgtype.Time) (result time.Time) {
	if !v.Valid {
		return result
	}

	return time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(v.Microseconds) * time.Microsecond)
}

func StringToText(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}

func DecimalToNumeric(value decimal.Decimal, required bool) pgtype.Numeric {
	if !required && value.IsZero() {
		return pgtype.Numeric{}
	}

	return pgtype.Numeric{
		Int:   value.Coefficient(),
		Exp:   value.Exponent(),
		Valid: true,
	}
}

func TimeToPGTime(value time.Time) pgtype.Time {
	if value.IsZero() {
		return pgtype.Time{}
	}

	microseconds := int64(value.Hour())*int64(time.Hour/time.Microsecond) +
		int64(value.Minute())*int64(time.Minute/time.Microsecond) +
		int64(value.Second())*int64(time.Second/time.Microsecond) +
		int64(value.Nanosecond())/int64(time.Microsecond)

	return pgtype.Time{Microseconds: microseconds, Valid: true}
}
