package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
	payoutErrors "github.com/nlsnnn/berezhok/internal/modules/payout/errors"
)

type DispatchRepository interface {
	LockPendingForDispatch(ctx context.Context, limit int32) ([]uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Payout, error)
	GetDestination(ctx context.Context, partnerID uuid.UUID) (domain.PayoutDestination, error)
	MarkCompleted(ctx context.Context, id uuid.UUID, providerPayoutID string) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error
	SetProviderPayoutID(ctx context.Context, id uuid.UUID, providerPayoutID string) error
	ListProcessingPayouts(ctx context.Context, limit int32) ([]ProcessingPayout, error)
}

type DispatchService struct {
	repo     DispatchRepository
	provider PayoutProvider
	log      *slog.Logger
}

func NewDispatchService(repo DispatchRepository, provider PayoutProvider, log *slog.Logger) *DispatchService {
	return &DispatchService{
		repo:     repo,
		provider: provider,
		log:      log,
	}
}

func (s *DispatchService) DispatchPending(ctx context.Context, limit int32) error {
	ids, err := s.repo.LockPendingForDispatch(ctx, limit)
	if err != nil {
		return fmt.Errorf("lock pending payouts: %w", err)
	}

	if len(ids) == 0 {
		return nil
	}

	s.log.Info("dispatching payouts", slog.Int("count", len(ids)))

	for _, id := range ids {
		if err := s.dispatchOne(ctx, id); err != nil {
			s.log.Error("dispatch payout failed",
				slog.String("payout_id", id.String()),
				slog.String("err", err.Error()),
			)
		}
	}

	s.log.Info("dispatch completed", slog.Int("dispatched", len(ids)))

	return nil
}

func (s *DispatchService) dispatchOne(ctx context.Context, id uuid.UUID) error {
	payout, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get payout: %w", err)
	}

	dest, err := s.repo.GetDestination(ctx, payout.PartnerID)
	if err != nil {
		if err == payoutErrors.ErrDestinationNotFound {
			s.log.Warn("payout destination not configured, marking failed",
				slog.String("payout_id", id.String()),
				slog.String("partner_id", payout.PartnerID.String()),
			)
			_ = s.repo.MarkFailed(ctx, id, "payout destination not configured")

			return nil
		}

		return fmt.Errorf("get destination: %w", err)
	}

	result, err := s.provider.SendSBP(ctx, SendSBPInput{
		Amount:         payout.Net,
		PhoneE164:      dest.SBPPhone,
		BankID:         dest.SBPBankID,
		RecipientName:  dest.RecipientName,
		Description:    fmt.Sprintf("Выплата партнёру за период %s–%s", payout.PeriodStart.Format("02.01"), payout.PeriodEnd.Format("02.01.2006")),
		IdempotencyKey: payout.IdempotencyKey,
		Metadata: map[string]string{
			"payout_id":  payout.ID.String(),
			"partner_id": payout.PartnerID.String(),
		},
	})
	if err != nil {
		s.log.Error("SendSBP failed",
			slog.String("payout_id", id.String()),
			slog.String("partner_id", payout.PartnerID.String()),
			slog.String("provider_error", err.Error()),
		)
		_ = s.repo.MarkFailed(ctx, id, err.Error())

		return nil
	}

	switch result.Status {
	case "succeeded":
		if err := s.repo.MarkCompleted(ctx, id, result.ProviderPayoutID); err != nil {
			return fmt.Errorf("mark completed: %w", err)
		}
		s.log.Info("payout completed",
			slog.String("payout_id", id.String()),
			slog.String("provider_id", result.ProviderPayoutID),
		)

	case "pending":
		// YooKassa processes SBP payouts asynchronously; save the provider ID and poll later.
		if err := s.repo.SetProviderPayoutID(ctx, id, result.ProviderPayoutID); err != nil {
			return fmt.Errorf("set provider payout id: %w", err)
		}
		s.log.Info("payout pending at provider, will poll",
			slog.String("payout_id", id.String()),
			slog.String("provider_id", result.ProviderPayoutID),
		)

	case "canceled":
		s.log.Warn("payout canceled by provider",
			slog.String("payout_id", id.String()),
			slog.String("provider_id", result.ProviderPayoutID),
		)
		_ = s.repo.MarkFailed(ctx, id, "payout canceled by provider")

	default:
		s.log.Warn("unknown payout status from provider",
			slog.String("payout_id", id.String()),
			slog.String("status", result.Status),
		)
	}

	return nil
}

// PollProcessing checks payouts that are in processing state and have a known
// provider ID, fetching their current status from YooKassa and updating the DB.
func (s *DispatchService) PollProcessing(ctx context.Context, limit int32) error {
	payouts, err := s.repo.ListProcessingPayouts(ctx, limit)
	if err != nil {
		return fmt.Errorf("list processing payouts: %w", err)
	}

	if len(payouts) == 0 {
		return nil
	}

	s.log.Info("polling processing payouts", slog.Int("count", len(payouts)))

	for _, p := range payouts {
		result, err := s.provider.GetPayoutStatus(ctx, p.ProviderPayoutID)
		if err != nil {
			s.log.Error("GetPayoutStatus failed",
				slog.String("payout_id", p.ID.String()),
				slog.String("provider_id", p.ProviderPayoutID),
				slog.String("err", err.Error()),
			)
			continue
		}

		switch result.Status {
		case "succeeded":
			if err := s.repo.MarkCompleted(ctx, p.ID, result.ProviderPayoutID); err != nil {
				s.log.Error("mark completed failed", slog.String("payout_id", p.ID.String()), slog.String("err", err.Error()))
			} else {
				s.log.Info("payout completed via poll", slog.String("payout_id", p.ID.String()))
			}

		case "canceled":
			_ = s.repo.MarkFailed(ctx, p.ID, "payout canceled by provider")
			s.log.Warn("payout canceled (polled)", slog.String("payout_id", p.ID.String()))

		case "pending":
			s.log.Info("payout still pending at provider", slog.String("payout_id", p.ID.String()))
		}
	}

	return nil
}
