package domain

import "fmt"

type LocationStatus string

const (
	LocationStatusActive   LocationStatus = "active"
	LocationStatusInactive LocationStatus = "inactive"
	LocationStatusClosed   LocationStatus = "closed"
)

type Location struct {
	ID            string
	PartnerID     string
	Name          string
	Address       string
	Phone         string
	LogoURL       string
	CoverImageURL string
	GalleryURLs   []string
	WorkingHours  string
	Status        LocationStatus
	Category      LocationCategory
	Location      GeoPoint
}

type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

type LocationCategory struct {
	Code    string
	Name    string
	IconURL string
	Color   string
	Sort    int
}

func NewLocation(partnerID, name, address string, category LocationCategory, location GeoPoint) (Location, error) {
	if partnerID == "" {
		return Location{}, fmt.Errorf("partner ID is required")
	}
	if name == "" {
		return Location{}, fmt.Errorf("name is required")
	}
	if address == "" {
		return Location{}, fmt.Errorf("address is required")
	}

	return Location{
		PartnerID: partnerID,
		Name:      name,
		Address:   address,
		Status:    LocationStatusInactive,
		Category:  category,
		Location:  location,
	}, nil
}
