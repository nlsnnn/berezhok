package yookassa

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/modules/payment/domain"
	yk "github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type YookassaAdapter struct {
	client *yk.PaymentHandler
}

func NewAdapter(client *yk.PaymentHandler) *YookassaAdapter {
	return &YookassaAdapter{client: client}
}

func (a *YookassaAdapter) Create(ctx context.Context, amount string, description string, returnURL string) (domain.ProviderPaymentResult, error) {
	payment, err := a.client.CreatePayment(ctx, &yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    amount,
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType(yoopayment.PaymentTypeBankCard),
		Confirmation: yoopayment.MobileApplication{
			Type:            yoopayment.TypeMobileApplication,
			ConfirmationURL: returnURL,
		},
		Description: description,
	})

	if err != nil {
		return domain.ProviderPaymentResult{}, err
	}
	// Confirmation:map[confirmation_url:https://yoomoney.ru/checkout/payments/v2/contract?orderId=315661f0-000f-5001-8000-17f1c044277d type:redirect]
	confirmationURL := payment.Confirmation.(map[string]interface{})["confirmation_url"].(string)

	return domain.ProviderPaymentResult{
		PaymentLink:       confirmationURL,
		ProviderPaymentID: payment.ID,
	}, nil
}
