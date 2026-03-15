package domain

import errs "github.com/nlsnnn/berezhok/internal/shared/errors"

// GeoPoint represents a geographic coordinate used across modules.
type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

func NewGeoPoint(lat, lon float64) (GeoPoint, error) {
	p := GeoPoint{Latitude: lat, Longitude: lon}
	if !p.IsValid() {
		return GeoPoint{}, errs.ErrInvalidGeoPoint
	}
	return p, nil
}

func (p GeoPoint) IsValid() bool {
	return p.Latitude >= -90 && p.Latitude <= 90 && p.Longitude >= -180 && p.Longitude <= 180
}
