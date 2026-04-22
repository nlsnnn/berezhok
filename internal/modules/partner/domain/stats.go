package domain

import "time"

type StatsFilter struct {
	Period           string
	DateFrom         time.Time
	DateTo           time.Time
	LocationID       string
	Status           string
	TopLocationsSort string
	TopBoxesSort     string
	OrdersSort       string
	Limit            int
	Offset           int
}

type PartnerStats struct {
	Summary         PartnerStatsSummary
	Timeline        []PartnerStatsTimelinePoint
	StatusBreakdown []PartnerStatsStatusBreakdownItem
	TopLocations    []PartnerStatsLocation
	TopBoxes        []PartnerStatsBox
	Orders          []PartnerStatsOrder
	Meta            PartnerStatsMeta
}

type PartnerStatsSummary struct {
	OrdersTotal               int
	OrdersCompleted           int
	OrdersCancelled           int
	OrdersPendingConfirmation int
	GrossRevenue              int
	NetRevenue                int
	AvgOrderValue             float64
	AvgRating                 float64
	ReviewsCount              int
}

type PartnerStatsTimelinePoint struct {
	Date            time.Time
	OrdersTotal     int
	OrdersCompleted int
	GrossRevenue    int
	NetRevenue      int
}

type PartnerStatsStatusBreakdownItem struct {
	Status string
	Count  int
	Share  float64
}

type PartnerStatsLocation struct {
	LocationID      string
	Name            string
	Address         string
	OrdersTotal     int
	OrdersCompleted int
	GrossRevenue    int
	NetRevenue      int
	AvgRating       float64
}

type PartnerStatsBox struct {
	BoxID           string
	Name            string
	ImageURL        string
	LocationName    string
	OrdersTotal     int
	OrdersCompleted int
	GrossRevenue    int
	NetRevenue      int
}

type PartnerStatsOrder struct {
	ID              string
	Status          string
	PickupCode      string
	Amount          float64
	BoxName         string
	BoxImageURL     string
	CustomerPhone   string
	CustomerName    string
	LocationID      string
	LocationName    string
	LocationAddress string
	PickupTimeStart time.Time
	PickupTimeEnd   time.Time
	CreatedAt       time.Time
	CanPickup       bool
}

type PartnerStatsMeta struct {
	Period           string
	DateFrom         time.Time
	DateTo           time.Time
	LocationID       string
	Status           string
	TopLocationsSort string
	TopBoxesSort     string
	OrdersSort       string
	Pagination       PartnerStatsPagination
}

type PartnerStatsPagination struct {
	Total   int
	Limit   int
	Offset  int
	HasMore bool
}
