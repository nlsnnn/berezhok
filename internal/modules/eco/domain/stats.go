package domain

import "github.com/shopspring/decimal"

// Tier ranks a customer based on total kilograms of food saved.
type Tier string

const (
	TierNewcomer Tier = "newcomer"
	TierHelper   Tier = "helper"
	TierGuardian Tier = "guardian"
	TierKeeper   Tier = "keeper"
	TierHero     Tier = "hero"
)

// EcoStats is a read model returned by GET /customer/eco-stats.
type EcoStats struct {
	BoxesPickedUp   int
	TotalKg         float64
	CO2SavedKg      float64
	SavingsRub      decimal.Decimal
	MealsEquivalent int
	Tier            Tier
	TierProgress    float64 // 0..1 progress to NextTier; 1.0 when on top tier
	NextTier        *Tier   // nil when on top tier
	KgToNextTier    float64 // 0 when on top tier
}

// CategoryAggregate is the raw per-location-category aggregation returned by
// the repository. The service multiplies counts by the per-category weight
// coefficient from coefficients.go.
type CategoryAggregate struct {
	CategoryCode string
	PickedCount  int
	Savings      decimal.Decimal
}
