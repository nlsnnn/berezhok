package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	catalogDomain "github.com/nlsnnn/berezhok/internal/modules/catalog/domain"
	partnerDomain "github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	sharedDomain "github.com/nlsnnn/berezhok/internal/shared/domain"
)

type LocationRepo struct {
	q *sqlc.Queries
}

func NewLocationRepo(q *sqlc.Queries) *LocationRepo {
	return &LocationRepo{q: q}
}

// SearchLocationsParams holds search parameters
type SearchLocationsParams struct {
	CategoryCode *string
	Limit        int
	Offset       int
}

// SearchLocations returns list of locations with filter and pagination
func (r *LocationRepo) SearchLocations(ctx context.Context, params SearchLocationsParams) ([]LocationWithDetails, error) {
	var categoryCode pgtype.Text
	if params.CategoryCode != nil {
		categoryCode = pgtype.Text{String: *params.CategoryCode, Valid: true}
	}

	rows, err := r.q.SearchLocations(ctx, sqlc.SearchLocationsParams{
		CategoryCode: categoryCode,
		Limit:        int32(params.Limit),
		Offset:       int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]LocationWithDetails, len(rows))
	for i, row := range rows {
		result[i] = searchRowToLocationWithDetails(row)
	}

	return result, nil
}

// CountActiveLocations returns total count for pagination
func (r *LocationRepo) CountActiveLocations(ctx context.Context, categoryCode *string) (int64, error) {
	var catCode pgtype.Text
	if categoryCode != nil {
		catCode = pgtype.Text{String: *categoryCode, Valid: true}
	}
	return r.q.CountActiveLocations(ctx, catCode)
}

// GetLocationDetailsByID returns full location details including category
func (r *LocationRepo) GetLocationDetailsByID(ctx context.Context, id uuid.UUID) (LocationWithDetails, error) {
	row, err := r.q.GetLocationDetailsByID(ctx, id)
	if err != nil {
		return LocationWithDetails{}, err
	}

	return detailsRowToLocationWithDetails(row), nil
}

// CountActiveBoxesByLocationID returns count of active boxes
func (r *LocationRepo) CountActiveBoxesByLocationID(ctx context.Context, locationID uuid.UUID) (int64, error) {
	return r.q.CountActiveBoxesByLocationID(ctx, locationID)
}

// GetActiveBoxesByLocationID returns active boxes for location
func (r *LocationRepo) GetActiveBoxesByLocationID(ctx context.Context, locationID uuid.UUID) ([]catalogDomain.SurpriseBox, error) {
	boxes, err := r.q.ListActiveBoxesByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	result := make([]catalogDomain.SurpriseBox, len(boxes))
	for i, box := range boxes {
		result[i] = catalogDomain.SurpriseBox{
			ID:          box.ID,
			LocationID:  box.LocationID,
			Name:        box.Name,
			Description: box.Description.String,
			Price: catalogDomain.Price{
				Original: pgconverter.NumericToDecimalOrZero(box.OriginalPrice),
				Discount: pgconverter.NumericToDecimalOrZero(box.DiscountPrice),
			},
			PickupTime: sharedDomain.PickupTime{
				Start: pgconverter.TimeValue(box.PickupTimeStart),
				End:   pgconverter.TimeValue(box.PickupTimeEnd),
			},
			Quantity:  int(box.QuantityAvailable),
			Status:    catalogDomain.BoxStatus(box.Status),
			Image:     box.ImageUrl.String,
			CreatedAt: box.CreatedAt,
		}
	}

	return result, nil
}

// LocationWithDetails contains location data with category info
type LocationWithDetails struct {
	ID            uuid.UUID
	PartnerID     uuid.UUID
	Name          string
	Address       string
	Phone         string
	LogoURL       string
	CoverImageURL string
	GalleryURLs   []string
	WorkingHours  map[string]string
	Status        string
	Category      partnerDomain.LocationCategory
	Coords        sharedDomain.GeoPoint
	ActiveBoxes   int
	CreatedAt     interface{}
	UpdatedAt     interface{}
}

// searchRowToLocationWithDetails converts search row to domain model
func searchRowToLocationWithDetails(r sqlc.SearchLocationsRow) LocationWithDetails {
	var workingHours map[string]string
	if r.WorkingHours != nil {
		_ = json.Unmarshal(r.WorkingHours, &workingHours)
	}

	// Convert interface{} coordinates to float64
	lat, _ := r.Latitude.(float64)
	lng, _ := r.Longitude.(float64)

	return LocationWithDetails{
		ID:            r.ID,
		PartnerID:     r.PartnerID,
		Name:          r.Name,
		Address:       r.Address,
		Phone:         r.Phone.String,
		LogoURL:       r.LogoUrl.String,
		CoverImageURL: r.CoverImageUrl.String,
		GalleryURLs:   r.GalleryUrls,
		WorkingHours:  workingHours,
		Status:        r.Status,
		Category: partnerDomain.LocationCategory{
			Code:    r.CategoryCode,
			Name:    r.CategoryName,
			IconURL: r.CategoryIconUrl.String,
			Color:   r.CategoryColor.String,
		},
		Coords: sharedDomain.GeoPoint{
			Latitude:  lat,
			Longitude: lng,
		},
		ActiveBoxes: int(r.ActiveBoxesCount),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// detailsRowToLocationWithDetails converts details row to domain model
func detailsRowToLocationWithDetails(r sqlc.GetLocationDetailsByIDRow) LocationWithDetails {
	var workingHours map[string]string
	if r.WorkingHours != nil {
		_ = json.Unmarshal(r.WorkingHours, &workingHours)
	}

	// Convert interface{} coordinates to float64
	lat, _ := r.Latitude.(float64)
	lng, _ := r.Longitude.(float64)

	return LocationWithDetails{
		ID:            r.ID,
		PartnerID:     r.PartnerID,
		Name:          r.Name,
		Address:       r.Address,
		Phone:         r.Phone.String,
		LogoURL:       r.LogoUrl.String,
		CoverImageURL: r.CoverImageUrl.String,
		GalleryURLs:   r.GalleryUrls,
		WorkingHours:  workingHours,
		Status:        r.Status,
		Category: partnerDomain.LocationCategory{
			Code:    r.CategoryCode,
			Name:    r.CategoryName,
			IconURL: r.CategoryIconUrl.String,
			Color:   r.CategoryColor.String,
		},
		Coords: sharedDomain.GeoPoint{
			Latitude:  lat,
			Longitude: lng,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
