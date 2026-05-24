package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
)

type DestinationRepository interface {
	GetDestination(ctx context.Context, partnerID uuid.UUID) (domain.PayoutDestination, error)
	UpsertDestination(ctx context.Context, d domain.PayoutDestination) (domain.PayoutDestination, error)
}

type DestinationService struct {
	repo  DestinationRepository
	banks BanksProvider
}

func NewDestinationService(repo DestinationRepository, banks BanksProvider) *DestinationService {
	return &DestinationService{repo: repo, banks: banks}
}

func (s *DestinationService) Get(ctx context.Context, partnerID uuid.UUID) (domain.PayoutDestination, error) {
	return s.repo.GetDestination(ctx, partnerID)
}

func (s *DestinationService) GetSBPBanks(ctx context.Context) ([]SBPBank, error) {
	return s.banks.GetSBPBanks(ctx)
}

func (s *DestinationService) Upsert(ctx context.Context, partnerID uuid.UUID, destType domain.DestinationType, phone, bankID, recipientName string) (domain.PayoutDestination, error) {
	d := domain.PayoutDestination{
		PartnerID:     partnerID,
		Type:          destType,
		SBPPhone:      phone,
		SBPBankID:     bankID,
		RecipientName: recipientName,
	}

	return s.repo.UpsertDestination(ctx, d)
}
