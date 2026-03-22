package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	"github.com/shopspring/decimal"
)

type BoxRepo struct {
	q *sqlc.Queries
}

func NewBoxRepo(q *sqlc.Queries) *BoxRepo {
	return &BoxRepo{q: q}
}

func (r *BoxRepo) CreateBox(ctx context.Context, box *domain.SurpriseBox) error {
	status := box.Status
	if status == "" {
		status = domain.BoxStatusDraft
	}

	b, err := r.q.CreateBox(ctx, sqlc.CreateBoxParams{
		LocationID:        box.LocationID,
		Name:              box.Name,
		Description:       stringToText(box.Description),
		OriginalPrice:     decimalToNumeric(box.Price.Original, false),
		DiscountPrice:     decimalToNumeric(box.Price.Discount, true),
		QuantityAvailable: int32(box.Quantity),
		PickupTimeStart:   timeToPGTime(box.PickupTime.Start),
		PickupTimeEnd:     timeToPGTime(box.PickupTime.End),
		ImageUrl:          stringToText(box.Image),
		Status:            string(status),
	})
	if err != nil {
		return err
	}
	box.ID = b.ID
	return nil
}

func (r *BoxRepo) GetBoxesByLocationID(ctx context.Context, locationID uuid.UUID) ([]domain.SurpriseBox, error) {
	boxes, err := r.q.ListBoxesByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	domainBoxes := make([]domain.SurpriseBox, len(boxes))
	for i, b := range boxes {
		domainBoxes[i] = boxToDomain(b)
	}

	return domainBoxes, nil
}

func (r *BoxRepo) GetBoxesByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]domain.SurpriseBox, error) {
	boxes, err := r.q.ListBoxesByPartnerID(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	domainBoxes := make([]domain.SurpriseBox, len(boxes))
	for i, b := range boxes {
		domainBoxes[i] = boxToDomain(b)
	}

	return domainBoxes, nil
}

func (r *BoxRepo) GetBoxByID(ctx context.Context, id string) (*domain.SurpriseBox, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, catalogErrors.ErrInvalidBoxID
	}

	b, err := r.q.FindBoxByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, catalogErrors.ErrBoxNotFound
		}
		return nil, err
	}
	box := boxToDomain(b)
	return &box, nil
}

func (r *BoxRepo) UpdateBox(ctx context.Context, box *domain.SurpriseBox) error {
	_, err := r.q.UpdateBox(ctx, sqlc.UpdateBoxParams{
		ID:                box.ID,
		Name:              box.Name,
		Description:       stringToText(box.Description),
		OriginalPrice:     decimalToNumeric(box.Price.Original, false),
		DiscountPrice:     decimalToNumeric(box.Price.Discount, true),
		QuantityAvailable: int32(box.Quantity),
		PickupTimeStart:   timeToPGTime(box.PickupTime.Start),
		PickupTimeEnd:     timeToPGTime(box.PickupTime.End),
		ImageUrl:          stringToText(box.Image),
		Status:            string(box.Status),
	})
	return err
}

func (r *BoxRepo) DeleteBox(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return catalogErrors.ErrInvalidBoxID
	}

	_, err = r.q.FindBoxByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return catalogErrors.ErrBoxNotFound
		}
		return err
	}

	return r.q.DeleteBox(ctx, uid)
}

func boxToDomain(b sqlc.SurpriseBox) domain.SurpriseBox {
	originalPrice := numericToDecimalOrZero(b.OriginalPrice)
	discountPrice := numericToDecimalOrZero(b.DiscountPrice)

	return domain.SurpriseBox{
		ID:          b.ID,
		LocationID:  b.LocationID,
		Name:        b.Name,
		Description: textToString(b.Description),
		Price: domain.Price{
			Original: originalPrice,
			Discount: discountPrice,
		},
		PickupTime: domain.PickupTime{
			Start: timeValue(b.PickupTimeStart),
			End:   timeValue(b.PickupTimeEnd),
		},
		Quantity:  int(b.QuantityAvailable),
		Status:    domain.BoxStatus(b.Status),
		Image:     textToString(b.ImageUrl),
		CreatedAt: b.CreatedAt,
	}
}

func numericToDecimalOrZero(v pgtype.Numeric) decimal.Decimal {
	if !v.Valid || v.Int == nil {
		return decimal.Zero
	}

	return decimal.NewFromBigInt(v.Int, v.Exp)
}

func textToString(v pgtype.Text) string {
	if !v.Valid {
		return ""
	}

	return v.String
}

func timeValue(v pgtype.Time) (result time.Time) {
	if !v.Valid {
		return result
	}

	return time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(v.Microseconds) * time.Microsecond)
}

func stringToText(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}

func decimalToNumeric(value decimal.Decimal, required bool) pgtype.Numeric {
	if !required && value.IsZero() {
		return pgtype.Numeric{}
	}

	return pgtype.Numeric{
		Int:   value.Coefficient(),
		Exp:   value.Exponent(),
		Valid: true,
	}
}

func timeToPGTime(value time.Time) pgtype.Time {
	if value.IsZero() {
		return pgtype.Time{}
	}

	microseconds := int64(value.Hour())*int64(time.Hour/time.Microsecond) +
		int64(value.Minute())*int64(time.Minute/time.Microsecond) +
		int64(value.Second())*int64(time.Second/time.Microsecond) +
		int64(value.Nanosecond())/int64(time.Microsecond)

	return pgtype.Time{Microseconds: microseconds, Valid: true}
}
