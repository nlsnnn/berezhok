package eco

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/eco/domain"
)

func TestComputeStats_TiersAndFormulas(t *testing.T) {
	tests := []struct {
		name         string
		counts       map[string]int
		savings      decimal.Decimal
		wantBoxes    int
		wantKg       float64
		wantCO2      float64
		wantTier     domain.Tier
		wantNextTier *domain.Tier
	}{
		{
			name:         "empty",
			counts:       map[string]int{},
			savings:      decimal.Zero,
			wantBoxes:    0,
			wantKg:       0,
			wantCO2:      0,
			wantTier:     domain.TierNewcomer,
			wantNextTier: tierPtr(domain.TierHelper),
		},
		{
			// 3 bakery (800g) + 2 grocery (2000g) = 2.4 + 4.0 = 6.4 kg → helper
			name: "helper boundary mix",
			counts: map[string]int{
				"bakery":  3,
				"grocery": 2,
			},
			savings:      decimal.NewFromInt(1820),
			wantBoxes:    5,
			wantKg:       6.4,
			wantCO2:      16,
			wantTier:     domain.TierHelper,
			wantNextTier: tierPtr(domain.TierGuardian),
		},
		{
			// Exactly at threshold — must include it (>= semantics).
			name:         "exactly 5 kg",
			counts:       map[string]int{"cafe": 5}, // 5 * 1200g = 6 kg, still helper
			savings:      decimal.NewFromInt(500),
			wantBoxes:    5,
			wantKg:       6,
			wantCO2:      15,
			wantTier:     domain.TierHelper,
			wantNextTier: tierPtr(domain.TierGuardian),
		},
		{
			name:         "guardian",
			counts:       map[string]int{"restaurant": 20}, // 30 kg
			savings:      decimal.NewFromInt(0),
			wantBoxes:    20,
			wantKg:       30,
			wantCO2:      75,
			wantTier:     domain.TierGuardian,
			wantNextTier: tierPtr(domain.TierKeeper),
		},
		{
			name:         "keeper",
			counts:       map[string]int{"grocery": 30}, // 60 kg
			savings:      decimal.NewFromInt(0),
			wantBoxes:    30,
			wantKg:       60,
			wantCO2:      150,
			wantTier:     domain.TierKeeper,
			wantNextTier: tierPtr(domain.TierHero),
		},
		{
			name:         "hero (top tier, no next)",
			counts:       map[string]int{"grocery": 100}, // 200 kg
			savings:      decimal.NewFromInt(0),
			wantBoxes:    100,
			wantKg:       200,
			wantCO2:      500,
			wantTier:     domain.TierHero,
			wantNextTier: nil,
		},
		{
			name:         "unknown category uses fallback",
			counts:       map[string]int{"food_truck": 2}, // 2 * 1200g = 2.4 kg
			savings:      decimal.NewFromInt(0),
			wantBoxes:    2,
			wantKg:       2.4,
			wantCO2:      6,
			wantTier:     domain.TierNewcomer,
			wantNextTier: tierPtr(domain.TierHelper),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ComputeStats(tc.counts, tc.savings)

			if got.BoxesPickedUp != tc.wantBoxes {
				t.Errorf("BoxesPickedUp = %d, want %d", got.BoxesPickedUp, tc.wantBoxes)
			}
			if got.TotalKg != tc.wantKg {
				t.Errorf("TotalKg = %v, want %v", got.TotalKg, tc.wantKg)
			}
			if got.CO2SavedKg != tc.wantCO2 {
				t.Errorf("CO2SavedKg = %v, want %v", got.CO2SavedKg, tc.wantCO2)
			}
			if got.Tier != tc.wantTier {
				t.Errorf("Tier = %q, want %q", got.Tier, tc.wantTier)
			}
			if (got.NextTier == nil) != (tc.wantNextTier == nil) {
				t.Errorf("NextTier nilness = %v, want %v", got.NextTier == nil, tc.wantNextTier == nil)
			} else if got.NextTier != nil && *got.NextTier != *tc.wantNextTier {
				t.Errorf("NextTier = %q, want %q", *got.NextTier, *tc.wantNextTier)
			}
			if !got.SavingsRub.Equal(tc.savings) {
				t.Errorf("SavingsRub = %s, want %s", got.SavingsRub.String(), tc.savings.String())
			}
		})
	}
}

func TestClassifyTier_ProgressAndKgToNext(t *testing.T) {
	// 30 kg: tier guardian (>=20), next keeper (>=50), 20 kg to go, progress 10/30 ≈ 0.333
	current, next, kgToNext, progress := classifyTier(30)

	if current != domain.TierGuardian {
		t.Fatalf("current = %q, want %q", current, domain.TierGuardian)
	}
	if next == nil || *next != domain.TierKeeper {
		t.Fatalf("next = %v, want %q", next, domain.TierKeeper)
	}
	if kgToNext != 20 {
		t.Errorf("kgToNext = %v, want 20", kgToNext)
	}
	// 10 of 30 = 0.3333...
	if progress < 0.3333 || progress > 0.3334 {
		t.Errorf("progress = %v, want ~0.3333", progress)
	}
}

func tierPtr(t domain.Tier) *domain.Tier { return &t }
