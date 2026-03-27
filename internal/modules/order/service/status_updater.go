package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
)

type orderStatusUpdater struct {
	repo orderRepository
	log  *slog.Logger
}

func NewOrderStatusUpdater(repo orderRepository, log *slog.Logger) *orderStatusUpdater {
	return &orderStatusUpdater{
		repo: repo,
		log:  log,
	}
}

func (u *orderStatusUpdater) MarkOrderPaid(ctx context.Context, orderID uuid.UUID) error {
	const op = "orderStatusUpdater.MarkOrderPaid"
	log := u.log.With(slog.String("op", op), slog.String("order_id", orderID.String()))

	err := u.repo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusPaid)
	if err != nil {
		log.Error("failed to update order status to paid", slog.String("error", err.Error()))
		return err
	}

	log.Info("order marked as paid")
	return nil
}

func (u *orderStatusUpdater) MarkOrderCanceled(ctx context.Context, orderID uuid.UUID) error {
	const op = "orderStatusUpdater.MarkOrderCanceled"
	log := u.log.With(slog.String("op", op), slog.String("order_id", orderID.String()))

	err := u.repo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusCancelled)
	if err != nil {
		log.Error("failed to update order status to canceled", slog.String("error", err.Error()))
		return err
	}

	log.Info("order marked as canceled")
	return nil
}
