package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nlsnnn/berezhok/internal/modules/payout/domain"
	payoutErrors "github.com/nlsnnn/berezhok/internal/modules/payout/errors"
)

// ----- mocks -----

type mockDispatchRepo struct {
	pendingIDs     []uuid.UUID
	processingRows []ProcessingPayout
	payout         domain.Payout
	destination    domain.PayoutDestination
	destErr        error

	completedID         uuid.UUID
	completedProvider   string
	failedID            uuid.UUID
	failedMsg           string
	savedProviderID     string
	savedProviderPayout uuid.UUID
}

func (m *mockDispatchRepo) LockPendingForDispatch(_ context.Context, limit int32) ([]uuid.UUID, error) {
	if int(limit) < len(m.pendingIDs) {
		return m.pendingIDs[:limit], nil
	}

	return m.pendingIDs, nil
}

func (m *mockDispatchRepo) GetByID(_ context.Context, _ uuid.UUID) (domain.Payout, error) {
	return m.payout, nil
}

func (m *mockDispatchRepo) GetDestination(_ context.Context, _ uuid.UUID) (domain.PayoutDestination, error) {
	return m.destination, m.destErr
}

func (m *mockDispatchRepo) MarkCompleted(_ context.Context, id uuid.UUID, providerID string) error {
	m.completedID = id
	m.completedProvider = providerID

	return nil
}

func (m *mockDispatchRepo) MarkFailed(_ context.Context, id uuid.UUID, msg string) error {
	m.failedID = id
	m.failedMsg = msg

	return nil
}

func (m *mockDispatchRepo) SetProviderPayoutID(_ context.Context, id uuid.UUID, providerID string) error {
	m.savedProviderPayout = id
	m.savedProviderID = providerID

	return nil
}

func (m *mockDispatchRepo) ListProcessingPayouts(_ context.Context, _ int32) ([]ProcessingPayout, error) {
	return m.processingRows, nil
}

type mockProvider struct {
	result     SendSBPResult
	err        error
	pollResult SendSBPResult
	pollErr    error
}

func (m *mockProvider) SendSBP(_ context.Context, _ SendSBPInput) (SendSBPResult, error) {
	return m.result, m.err
}

func (m *mockProvider) GetPayoutStatus(_ context.Context, _ string) (SendSBPResult, error) {
	return m.pollResult, m.pollErr
}

func newDispatchSvc(repo *mockDispatchRepo, provider PayoutProvider) *DispatchService {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	return NewDispatchService(repo, provider, log)
}

// ----- tests -----

func TestDispatchPending_SuccessPath(t *testing.T) {
	payoutID := uuid.New()
	partnerID := uuid.New()

	repo := &mockDispatchRepo{
		pendingIDs: []uuid.UUID{payoutID},
		payout: domain.Payout{
			ID:             payoutID,
			PartnerID:      partnerID,
			Net:            decimal.NewFromFloat(800),
			IdempotencyKey: payoutID.String(),
		},
		destination: domain.PayoutDestination{
			PartnerID: partnerID,
			Type:      domain.DestinationTypeSBP,
			SBPPhone:  "+79991234567",
			SBPBankID: "100000000111",
		},
	}
	provider := &mockProvider{result: SendSBPResult{ProviderPayoutID: "prov-123", Status: "succeeded"}}

	svc := newDispatchSvc(repo, provider)
	err := svc.DispatchPending(context.Background(), 10)
	require.NoError(t, err)

	assert.Equal(t, payoutID, repo.completedID)
	assert.Equal(t, "prov-123", repo.completedProvider)
	assert.Equal(t, uuid.Nil, repo.failedID)
}

func TestDispatchPending_FailsWhenNoDestination(t *testing.T) {
	payoutID := uuid.New()

	repo := &mockDispatchRepo{
		pendingIDs: []uuid.UUID{payoutID},
		payout:     domain.Payout{ID: payoutID, PartnerID: uuid.New()},
		destErr:    payoutErrors.ErrDestinationNotFound,
	}
	provider := &mockProvider{}

	svc := newDispatchSvc(repo, provider)
	err := svc.DispatchPending(context.Background(), 10)
	require.NoError(t, err)

	assert.Equal(t, payoutID, repo.failedID)
	assert.Contains(t, repo.failedMsg, "destination")
}

func TestDispatchPending_FailsWhenProviderErrors(t *testing.T) {
	payoutID := uuid.New()
	partnerID := uuid.New()

	repo := &mockDispatchRepo{
		pendingIDs: []uuid.UUID{payoutID},
		payout:     domain.Payout{ID: payoutID, PartnerID: partnerID, Net: decimal.NewFromFloat(500)},
		destination: domain.PayoutDestination{
			PartnerID: partnerID,
			Type:      domain.DestinationTypeSBP,
			SBPPhone:  "+79001234567",
			SBPBankID: "100000000111",
		},
	}
	provider := &mockProvider{err: errors.New("api timeout")}

	svc := newDispatchSvc(repo, provider)
	err := svc.DispatchPending(context.Background(), 10)
	require.NoError(t, err)

	assert.Equal(t, payoutID, repo.failedID)
	assert.Contains(t, repo.failedMsg, "timeout")
}

func TestDispatchPending_NoPendingDoesNothing(t *testing.T) {
	repo := &mockDispatchRepo{pendingIDs: []uuid.UUID{}}
	provider := &mockProvider{}

	svc := newDispatchSvc(repo, provider)
	err := svc.DispatchPending(context.Background(), 10)
	require.NoError(t, err)

	assert.Equal(t, uuid.Nil, repo.completedID)
	assert.Equal(t, uuid.Nil, repo.failedID)
}

func TestDispatchPending_PendingStatusSavesProviderID(t *testing.T) {
	payoutID := uuid.New()
	partnerID := uuid.New()

	repo := &mockDispatchRepo{
		pendingIDs: []uuid.UUID{payoutID},
		payout: domain.Payout{
			ID: payoutID, PartnerID: partnerID,
			Net: decimal.NewFromFloat(500), IdempotencyKey: payoutID.String(),
		},
		destination: domain.PayoutDestination{
			PartnerID: partnerID, Type: domain.DestinationTypeSBP,
			SBPPhone: "+79001234567", SBPBankID: "100000000111",
		},
	}
	provider := &mockProvider{result: SendSBPResult{ProviderPayoutID: "prov-456", Status: "pending"}}

	svc := newDispatchSvc(repo, provider)
	err := svc.DispatchPending(context.Background(), 10)
	require.NoError(t, err)

	// must NOT mark completed or failed
	assert.Equal(t, uuid.Nil, repo.completedID)
	assert.Equal(t, uuid.Nil, repo.failedID)
	// must save provider ID for future polling
	assert.Equal(t, payoutID, repo.savedProviderPayout)
	assert.Equal(t, "prov-456", repo.savedProviderID)
}

func TestDispatchPending_CanceledStatusMarksFailed(t *testing.T) {
	payoutID := uuid.New()
	partnerID := uuid.New()

	repo := &mockDispatchRepo{
		pendingIDs: []uuid.UUID{payoutID},
		payout: domain.Payout{
			ID: payoutID, PartnerID: partnerID,
			Net: decimal.NewFromFloat(500), IdempotencyKey: payoutID.String(),
		},
		destination: domain.PayoutDestination{
			PartnerID: partnerID, Type: domain.DestinationTypeSBP,
			SBPPhone: "+79001234567", SBPBankID: "100000000111",
		},
	}
	provider := &mockProvider{result: SendSBPResult{ProviderPayoutID: "prov-789", Status: "canceled"}}

	svc := newDispatchSvc(repo, provider)
	err := svc.DispatchPending(context.Background(), 10)
	require.NoError(t, err)

	assert.Equal(t, payoutID, repo.failedID)
	assert.Contains(t, repo.failedMsg, "canceled")
}

func TestPollProcessing_SucceededTransitionsToCompleted(t *testing.T) {
	payoutID := uuid.New()

	repo := &mockDispatchRepo{
		processingRows: []ProcessingPayout{{ID: payoutID, ProviderPayoutID: "prov-abc"}},
	}
	provider := &mockProvider{pollResult: SendSBPResult{ProviderPayoutID: "prov-abc", Status: "succeeded"}}

	svc := newDispatchSvc(repo, provider)
	err := svc.PollProcessing(context.Background(), 10)
	require.NoError(t, err)

	assert.Equal(t, payoutID, repo.completedID)
	assert.Equal(t, "prov-abc", repo.completedProvider)
}

func TestPollProcessing_CanceledTransitionsToFailed(t *testing.T) {
	payoutID := uuid.New()

	repo := &mockDispatchRepo{
		processingRows: []ProcessingPayout{{ID: payoutID, ProviderPayoutID: "prov-def"}},
	}
	provider := &mockProvider{pollResult: SendSBPResult{ProviderPayoutID: "prov-def", Status: "canceled"}}

	svc := newDispatchSvc(repo, provider)
	err := svc.PollProcessing(context.Background(), 10)
	require.NoError(t, err)

	assert.Equal(t, payoutID, repo.failedID)
}
