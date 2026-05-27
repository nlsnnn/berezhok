package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/modules/eco"
	"github.com/nlsnnn/berezhok/internal/modules/eco/domain"
)

type ecoRepo interface {
	GetAggregate(ctx context.Context, userID uuid.UUID) ([]domain.CategoryAggregate, error)
}

type statsCache interface {
	Get(ctx context.Context, userID uuid.UUID) (*domain.EcoStats, bool, error)
	Set(ctx context.Context, userID uuid.UUID, stats domain.EcoStats) error
}

type ecoService struct {
	repo  ecoRepo
	cache statsCache
	log   *slog.Logger
}

// NewEcoService wires the eco read-path. `cache` is optional — pass nil to
// disable caching (useful in tests).
func NewEcoService(repo ecoRepo, cache statsCache, log *slog.Logger) *ecoService {
	return &ecoService{repo: repo, cache: cache, log: log}
}

func (s *ecoService) GetForUser(ctx context.Context, userID uuid.UUID) (domain.EcoStats, error) {
	const op = "eco.service.GetForUser"

	if s.cache != nil {
		cached, hit, err := s.cache.Get(ctx, userID)
		if err != nil {
			s.log.Warn("eco cache read failed", sl.Err(err))
		} else if hit {
			return *cached, nil
		}
	}

	aggregate, err := s.repo.GetAggregate(ctx, userID)
	if err != nil {
		return domain.EcoStats{}, fmt.Errorf("%s: %w", op, err)
	}

	counts := make(map[string]int, len(aggregate))
	savings := decimal.Zero
	for _, row := range aggregate {
		counts[row.CategoryCode] = row.PickedCount
		savings = savings.Add(row.Savings)
	}

	stats := eco.ComputeStats(counts, savings)

	if s.cache != nil {
		if err := s.cache.Set(ctx, userID, stats); err != nil {
			s.log.Warn("eco cache write failed", sl.Err(err))
		}
	}

	return stats, nil
}
