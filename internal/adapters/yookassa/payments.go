package yookassa

import (
	"context"
	"fmt"

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

func (a *YookassaAdapter) CreatePayment(ctx context.Context, amount string, description string, method string, returnURL string) (domain.ProviderPaymentResult, error) {
	payment, err := a.client.CreatePayment(ctx, &yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    amount,
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType(method),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: returnURL,
		},
		Description: description,
	})
	if err != nil {
		return domain.ProviderPaymentResult{}, err
	}
	redirect, ok := payment.Confirmation.(yoopayment.Redirect)
	if !ok {
		return domain.ProviderPaymentResult{}, fmt.Errorf("unexpected confirmation type")
	}
	return domain.ProviderPaymentResult{
		PaymentLink:       redirect.ConfirmationURL,
		ProviderPaymentID: payment.ID,
	}, nil
}
