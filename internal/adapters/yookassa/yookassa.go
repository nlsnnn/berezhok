package yookassa

import (
	"net/http"
	"time"

	"github.com/nlsnnn/berezhok/internal/shared/config"
	yk "github.com/rvinnie/yookassa-sdk-go/yookassa"
	yooopts "github.com/rvinnie/yookassa-sdk-go/yookassa/opts"
)

func New(cfg config.Yookassa) *yk.PaymentHandler {
	client := yk.NewClient(
		cfg.AccountID,
		cfg.SecretKey,
		yooopts.WithHTTPClient(http.Client{Timeout: 30 * time.Second}),
	)

	paymentHandler := yk.NewPaymentHandler(client)

	return paymentHandler
}
