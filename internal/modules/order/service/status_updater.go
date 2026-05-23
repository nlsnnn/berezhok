package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
)

type orderStatusUpdater struct {
	repo           orderRepository
	eventPublisher orderProjectionPublisher
	log            *slog.Logger
}

func NewOrderStatusUpdater(repo orderRepository, log *slog.Logger, publishers ...orderProjectionPublisher) *orderStatusUpdater {
	u := &orderStatusUpdater{
		repo: repo,
		log:  log,
	}
	if len(publishers) > 0 {
		u.eventPublisher = publishers[0]
	}
	return u
}

func (u *orderStatusUpdater) MarkOrderPaid(ctx context.Context, orderID uuid.UUID) error {
	const op = "orderStatusUpdater.MarkOrderPaid"
	log := u.log.With(slog.String("op", op), slog.String("order_id", orderID.String()))

	err := u.repo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusConfirmed)
	if err != nil {
		log.Error("failed to update order status to paid", slog.String("error", err.Error()))
		return err
	}

	log.Info("order marked as paid")
	u.publishProjection(ctx, orderID)
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
	u.publishProjection(ctx, orderID)
	return nil
}

func (u *orderStatusUpdater) publishProjection(ctx context.Context, orderID uuid.UUID) {
	if u.eventPublisher == nil {
		return
	}
	projection, err := u.repo.GetOrderProjection(ctx, orderID)
	if err != nil {
		u.log.Warn("failed to load order projection", slog.String("order_id", orderID.String()), slog.String("err", err.Error()))
		return
	}
	if err = u.eventPublisher.PublishOrderStatusChanged(ctx, *projection); err != nil {
		u.log.Warn("failed to publish order projection", slog.String("order_id", orderID.String()), slog.String("err", err.Error()))
	}
}
