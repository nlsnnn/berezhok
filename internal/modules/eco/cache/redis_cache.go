package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/modules/eco/domain"
)

const (
	keyPrefix = "eco:stats:"
	cacheTTL  = 10 * time.Minute
)

type EcoStatsCache struct {
	client *redis.Client
}

func NewEcoStatsCache(client *redis.Client) *EcoStatsCache {
	return &EcoStatsCache{client: client}
}

// cacheEntry is the wire format stored in Redis. We don't reuse domain.EcoStats
// directly because decimal.Decimal marshals as a string and we want stable,
// explicit fields.
type cacheEntry struct {
	BoxesPickedUp   int     `json:"boxes_picked_up"`
	TotalKg         float64 `json:"total_kg"`
	CO2SavedKg      float64 `json:"co2_saved_kg"`
	SavingsRub      string  `json:"savings_rub"`
	MealsEquivalent int     `json:"meals_equivalent"`
	Tier            string  `json:"tier"`
	TierProgress    float64 `json:"tier_progress"`
	NextTier        *string `json:"next_tier,omitempty"`
	KgToNextTier    float64 `json:"kg_to_next_tier"`
}

func (c *EcoStatsCache) Get(ctx context.Context, userID uuid.UUID) (*domain.EcoStats, bool, error) {
	raw, err := c.client.Get(ctx, keyPrefix+userID.String()).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var entry cacheEntry
	if err := json.Unmarshal(raw, &entry); err != nil {
		return nil, false, fmt.Errorf("unmarshal eco cache entry: %w", err)
	}

	savings, err := decimal.NewFromString(entry.SavingsRub)
	if err != nil {
		savings = decimal.Zero
	}

	var nextTier *domain.Tier
	if entry.NextTier != nil {
		nt := domain.Tier(*entry.NextTier)
		nextTier = &nt
	}

	return &domain.EcoStats{
		BoxesPickedUp:   entry.BoxesPickedUp,
		TotalKg:         entry.TotalKg,
		CO2SavedKg:      entry.CO2SavedKg,
		SavingsRub:      savings,
		MealsEquivalent: entry.MealsEquivalent,
		Tier:            domain.Tier(entry.Tier),
		TierProgress:    entry.TierProgress,
		NextTier:        nextTier,
		KgToNextTier:    entry.KgToNextTier,
	}, true, nil
}

func (c *EcoStatsCache) Set(ctx context.Context, userID uuid.UUID, stats domain.EcoStats) error {
	var nextTier *string
	if stats.NextTier != nil {
		s := string(*stats.NextTier)
		nextTier = &s
	}

	entry := cacheEntry{
		BoxesPickedUp:   stats.BoxesPickedUp,
		TotalKg:         stats.TotalKg,
		CO2SavedKg:      stats.CO2SavedKg,
		SavingsRub:      stats.SavingsRub.String(),
		MealsEquivalent: stats.MealsEquivalent,
		Tier:            string(stats.Tier),
		TierProgress:    stats.TierProgress,
		NextTier:        nextTier,
		KgToNextTier:    stats.KgToNextTier,
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal eco cache entry: %w", err)
	}
	return c.client.Set(ctx, keyPrefix+userID.String(), payload, cacheTTL).Err()
}

// Invalidate drops the cached entry for a user. Safe to call when the key
// doesn't exist.
func (c *EcoStatsCache) Invalidate(ctx context.Context, userID uuid.UUID) error {
	return c.client.Del(ctx, keyPrefix+userID.String()).Err()
}
