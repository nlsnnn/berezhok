package yookassa

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	yk "github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yooerrors "github.com/rvinnie/yookassa-sdk-go/yookassa/errors"
	yooopts "github.com/rvinnie/yookassa-sdk-go/yookassa/opts"
	yoopayout "github.com/rvinnie/yookassa-sdk-go/yookassa/payout"

	"github.com/nlsnnn/berezhok/internal/modules/payout/service"
	"github.com/nlsnnn/berezhok/internal/shared/config"
)

func NewPayoutHandlerClient(cfg config.YookassaPayout) *yk.PayoutHandler {
	client := yk.NewClient(
		cfg.AgentID,
		cfg.SecretKey,
		yooopts.WithHTTPClient(http.Client{Timeout: 30 * time.Second}),
	)

	return yk.NewPayoutHandler(client)
}

type PayoutAdapter struct {
	client *yk.PayoutHandler
}

func NewPayoutAdapter(client *yk.PayoutHandler) *PayoutAdapter {
	return &PayoutAdapter{client: client}
}

func (a *PayoutAdapter) SendSBP(ctx context.Context, in service.SendSBPInput) (service.SendSBPResult, error) {
	p := &yoopayout.Payout{
		Amount: &yoocommon.Amount{
			Value:    in.Amount.StringFixed(2),
			Currency: "RUB",
		},
		PayoutDestinationData: yoopayout.PayoutDestinationData{
			Type:   yoopayout.PayoutTypeSBP,
			Phone:  in.PhoneE164,
			BankId: in.BankID,
		},
		Description: in.Description,
		Metadata: yoopayout.Metadata{
			OrderId: in.IdempotencyKey,
		},
	}

	created, err := a.client.WithIdempotencyKey(in.IdempotencyKey).CreatePayout(ctx, p)
	if err != nil {
		return service.SendSBPResult{}, wrapYooError(err)
	}

	return service.SendSBPResult{
		ProviderPayoutID: created.Id,
		Status:           string(created.Status),
	}, nil
}

func (a *PayoutAdapter) GetSBPBanks(ctx context.Context) ([]service.SBPBank, error) {
	banks, err := a.client.GetSbpBanks(ctx)
	if err != nil {
		return nil, wrapYooError(err)
	}

	result := make([]service.SBPBank, len(banks))
	for i, b := range banks {
		result[i] = service.SBPBank{BankID: b.BankId, Name: b.Name, BIC: b.Bic}
	}

	return result, nil
}

func (a *PayoutAdapter) GetPayoutStatus(ctx context.Context, providerPayoutID string) (service.SendSBPResult, error) {
	p, err := a.client.GetPayout(ctx, providerPayoutID)
	if err != nil {
		return service.SendSBPResult{}, wrapYooError(err)
	}

	return service.SendSBPResult{
		ProviderPayoutID: p.Id,
		Status:           string(p.Status),
	}, nil
}

// wrapYooError turns a *yooerrors.YoomoneyError into a richer message that
// includes the API error code and the offending parameter name — without these
// the caller only sees "error: <description>" and has no idea what field failed.
func wrapYooError(err error) error {
	var yooErr *yooerrors.YoomoneyError
	if errors.As(err, &yooErr) {
		return fmt.Errorf("yookassa [code=%s param=%s]: %s", yooErr.Code, yooErr.Parameter, yooErr.Description)
	}

	return err
}
