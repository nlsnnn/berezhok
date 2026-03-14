package partner

import (
	"context"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/domain"
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
		StMakepoint:   location.Location.Latitude,
		StMakepoint_2: location.Location.Longitude,
	})

	if err != nil {
		return domain.Location{}, err
	}

	return locationToDomain(l), nil
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

func locationToDomain(l sqlc.Location) domain.Location {
	category := domain.LocationCategory{
		Code: l.CategoryCode,
	}

	return domain.Location{
		ID:            l.ID.String(),
		PartnerID:     l.PartnerID.String(),
		Category:      category,
		Name:          l.Name,
		Address:       l.Address,
		Phone:         l.Phone.String,
		LogoURL:       l.LogoUrl.String,
		CoverImageURL: l.CoverImageUrl.String,
		GalleryURLs:   l.GalleryUrls,
	}
}
