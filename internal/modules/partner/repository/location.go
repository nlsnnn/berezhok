package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type LocationRepo struct {
	q *sqlc.Queries
}

func NewLocationRepo(q *sqlc.Queries) *LocationRepo {
	return &LocationRepo{q: q}
}

func (r *LocationRepo) Create(ctx context.Context, location domain.Location) (domain.Location, error) {
	partnerID := uuid.MustParse(location.PartnerID)
	l, err := r.q.CreateLocation(ctx, sqlc.CreateLocationParams{
		PartnerID:     partnerID,
		Name:          location.Name,
		Address:       location.Address,
		CategoryCode:  location.Category.Code,
		Status:        string(location.Status),
		StMakepoint:   location.Coords.Longitude,
		StMakepoint_2: location.Coords.Latitude,
	})
	if err != nil {
		return domain.Location{}, err
	}

	return locationToDomain(l), nil
}

func (r *LocationRepo) FindByPartnerID(ctx context.Context, partnerID string) ([]domain.Location, error) {
	uid := uuid.MustParse(partnerID)
	rows, err := r.q.FindLocationsByPartnerID(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Location, len(rows))
	for i, l := range rows {
		result[i] = locationToDomain(l)
	}
	return result, nil
}

func (r *LocationRepo) Delete(ctx context.Context, id string) error {
	uid := uuid.MustParse(id)
	return r.q.DeleteLocation(ctx, uid)
}

func (r *LocationRepo) FindCategoryByCode(ctx context.Context, code string) (domain.LocationCategory, error) {
	cat, err := r.q.FindCategoryByCode(ctx, code)
	if err != nil {
		return domain.LocationCategory{}, err
	}

	return domain.LocationCategory{
		Code:    cat.Code,
		Name:    cat.NameRu,
		IconURL: cat.IconUrl.String,
		Color:   cat.Color.String,
		Sort:    int(cat.SortOrder.Int32),
	}, nil
}

func (r *LocationRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.Location, error) {
	row, err := r.q.FindLocationByID(ctx, id)
	if err != nil {
		return domain.Location{}, err
	}
	return locationToDomain(row), nil
}

func locationToDomain(l sqlc.Location) domain.Location {
	return domain.Location{
		ID:            l.ID.String(),
		PartnerID:     l.PartnerID.String(),
		Name:          l.Name,
		Address:       l.Address,
		Phone:         l.Phone.String,
		LogoURL:       l.LogoUrl.String,
		CoverImageURL: l.CoverImageUrl.String,
		GalleryURLs:   l.GalleryUrls,
		Status:        domain.LocationStatus(l.Status),
		Category:      domain.LocationCategory{Code: l.CategoryCode},
		Coords:        sharedDomain.GeoPoint{}, // PostGIS geometry не парсится напрямую из sqlc
	}
}
