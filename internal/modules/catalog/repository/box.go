package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
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
		Description:       pgconverter.StringToText(box.Description),
		OriginalPrice:     pgconverter.DecimalToNumeric(box.Price.Original, false),
		DiscountPrice:     pgconverter.DecimalToNumeric(box.Price.Discount, true),
		QuantityAvailable: int32(box.Quantity),
		PickupTimeStart:   pgconverter.TimeToPGTime(box.PickupTime.Start),
		PickupTimeEnd:     pgconverter.TimeToPGTime(box.PickupTime.End),
		ImageUrl:          pgconverter.StringToText(box.Image),
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
		Description:       pgconverter.StringToText(box.Description),
		OriginalPrice:     pgconverter.DecimalToNumeric(box.Price.Original, false),
		DiscountPrice:     pgconverter.DecimalToNumeric(box.Price.Discount, true),
		QuantityAvailable: int32(box.Quantity),
		PickupTimeStart:   pgconverter.TimeToPGTime(box.PickupTime.Start),
		PickupTimeEnd:     pgconverter.TimeToPGTime(box.PickupTime.End),
		ImageUrl:          pgconverter.StringToText(box.Image),
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
	originalPrice := pgconverter.NumericToDecimalOrZero(b.OriginalPrice)
	discountPrice := pgconverter.NumericToDecimalOrZero(b.DiscountPrice)

	return domain.SurpriseBox{
		ID:          b.ID,
		LocationID:  b.LocationID,
		Name:        b.Name,
		Description: pgconverter.TextToString(b.Description),
		Price: domain.Price{
			Original: originalPrice,
			Discount: discountPrice,
		},
		PickupTime: domain.PickupTime{
			Start: pgconverter.TimeValue(b.PickupTimeStart),
			End:   pgconverter.TimeValue(b.PickupTimeEnd),
		},
		Quantity:  int(b.QuantityAvailable),
		Status:    domain.BoxStatus(b.Status),
		Image:     pgconverter.TextToString(b.ImageUrl),
		CreatedAt: b.CreatedAt,
	}
}
