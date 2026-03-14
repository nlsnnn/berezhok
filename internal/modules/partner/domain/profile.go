package domain

// PartnerProfile is a read model for the partner profile endpoint.
// It aggregates partner, employee and optional location data.
type PartnerProfile struct {
	Partner  Partner
	Employee Employee
	Location *LocationSummary
}

// LocationSummary is a lightweight read-only view of a location
// used when full location data is not needed.
type LocationSummary struct {
	ID      string
	Name    string
	Address string
}
