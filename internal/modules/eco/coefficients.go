package eco

import (
	"math"

	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/eco/domain"
)

// Average weight of one rescued surprise box per location category, in grams.
// Values are rough estimates based on Too Good To Go's public methodology and
// internal product judgement — they're meant to give honest order-of-magnitude
// numbers, not precise nutrition data.
var avgWeightGramsByCategory = map[string]int{
	"bakery":     800,
	"cafe":       1200,
	"restaurant": 1500,
	"grocery":    2000,
	"hotel":      1300,
}

const (
	// fallbackWeightGrams is used when a box comes from a category that's not
	// in the map (defence against schema drift).
	fallbackWeightGrams = 1200

	// co2PerKgFood is the CO₂-equivalent (kg) prevented per 1 kg of food
	// rescued from waste — a widely-cited FAO-derived average.
	co2PerKgFood = 2.5

	// mealKg is the weight of one "meal equivalent" used in marketing copy.
	mealKg = 1.5
)

// tierStep is one rung of the gamification ladder. Thresholds are the lower
// bound (inclusive) of total kilograms required to reach the tier.
type tierStep struct {
	tier        domain.Tier
	thresholdKg float64
}

var tierLadder = []tierStep{
	{domain.TierNewcomer, 0},
	{domain.TierHelper, 5},
	{domain.TierGuardian, 20},
	{domain.TierKeeper, 50},
	{domain.TierHero, 150},
}

// ComputeStats turns raw per-category counts into the read-model returned to
// the customer. Pure function — easy to unit-test without DB/Redis.
func ComputeStats(categoryCounts map[string]int, savings decimal.Decimal) domain.EcoStats {
	totalGrams := 0
	totalBoxes := 0
	for category, count := range categoryCounts {
		w, ok := avgWeightGramsByCategory[category]
		if !ok {
			w = fallbackWeightGrams
		}
		totalGrams += w * count
		totalBoxes += count
	}

	totalKg := float64(totalGrams) / 1000.0
	co2 := totalKg * co2PerKgFood
	meals := int(math.Round(totalKg / mealKg))

	tier, nextTier, kgToNext, progress := classifyTier(totalKg)

	return domain.EcoStats{
		BoxesPickedUp:   totalBoxes,
		TotalKg:         roundTo(totalKg, 2),
		CO2SavedKg:      roundTo(co2, 2),
		SavingsRub:      savings,
		MealsEquivalent: meals,
		Tier:            tier,
		TierProgress:    roundTo(progress, 4),
		NextTier:        nextTier,
		KgToNextTier:    roundTo(kgToNext, 2),
	}
}

// classifyTier returns the current tier plus progress to the next one.
// On the top tier nextTier is nil and progress is 1.0.
func classifyTier(totalKg float64) (current domain.Tier, next *domain.Tier, kgToNext, progress float64) {
	current = tierLadder[0].tier
	currentIdx := 0
	for i, step := range tierLadder {
		if totalKg >= step.thresholdKg {
			current = step.tier
			currentIdx = i
		}
	}

	if currentIdx == len(tierLadder)-1 {
		return current, nil, 0, 1.0
	}

	currentThreshold := tierLadder[currentIdx].thresholdKg
	nextThreshold := tierLadder[currentIdx+1].thresholdKg
	span := nextThreshold - currentThreshold

	kgToNext = nextThreshold - totalKg
	if kgToNext < 0 {
		kgToNext = 0
	}
	progress = (totalKg - currentThreshold) / span
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	nt := tierLadder[currentIdx+1].tier
	return current, &nt, kgToNext, progress
}

func roundTo(v float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(v*pow) / pow
}
