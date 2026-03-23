package domain

import "time"

// PartnerProfile is a read model for the partner profile endpoint.
// It aggregates partner, employee and optional location data.
type PartnerProfile struct {
	Partner   Partner
	Employee  Employee
	Location  *LocationSummary  // Employee's assigned location (for backwards compatibility)
	Locations []LocationSummary // All partner locations
}

// LocationSummary is a lightweight read-only view of a location
// used when full location data is not needed.
type LocationSummary struct {
	ID        string
	Name      string
	Address   string
	CreatedAt time.Time
}
