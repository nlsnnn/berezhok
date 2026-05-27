package dto

// EcoStatsResponse is the customer-facing eco-score payload.
// SavingsRub is sent as an integer ruble amount — the frontend formats it
// with Intl.NumberFormat('ru-RU').
type EcoStatsResponse struct {
	BoxesPickedUp   int     `json:"boxes_picked_up"`
	TotalKg         float64 `json:"total_kg"`
	CO2SavedKg      float64 `json:"co2_saved_kg"`
	SavingsRub      int64   `json:"savings_rub"`
	MealsEquivalent int     `json:"meals_equivalent"`
	Tier            string  `json:"tier"`
	TierProgress    float64 `json:"tier_progress"`
	NextTier        *string `json:"next_tier,omitempty"`
	KgToNextTier    float64 `json:"kg_to_next_tier"`
}
