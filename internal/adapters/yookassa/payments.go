package yookassa

import (
	"context"
	"fmt"

	yk "github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"

	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
)

type YookassaAdapter struct {
	client *yk.PaymentHandler
}

func NewAdapter(client *yk.PaymentHandler) *YookassaAdapter {
	return &YookassaAdapter{client: client}
}

func (a *YookassaAdapter) Create(ctx context.Context, amount, description, returnURL string, metadata map[string]string) (domain.ProviderPaymentResult, error) {
	payment, err := a.client.CreatePayment(ctx, &yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    amount,
			Currency: "RUB",
		},
		Confirmation: yoopayment.Redirect{
			Type:      yoopayment.TypeRedirect,
			ReturnURL: returnURL,
		},
		Description: description,
		Capture:     true,
		Metadata:    metadata,
	})
	if err != nil {
		return domain.ProviderPaymentResult{}, err
	}
	// Confirmation:map[confirmation_url:https://yoomoney.ru/checkout/payments/v2/contract?orderId=315661f0-000f-5001-8000-17f1c044277d type:redirect]
	confirmationMap, ok := payment.Confirmation.(map[string]interface{})
	if !ok {
		return domain.ProviderPaymentResult{}, fmt.Errorf("unexpected confirmation format from yookassa")
	}

	confirmationURL, ok := confirmationMap["confirmation_url"].(string)
	if !ok {
		return domain.ProviderPaymentResult{}, fmt.Errorf("confirmation_url missing in yookassa response")
	}

	return domain.ProviderPaymentResult{
		PaymentLink:       confirmationURL,
		ProviderPaymentID: payment.ID,
	}, nil
}
