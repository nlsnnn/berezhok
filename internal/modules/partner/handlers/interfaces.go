package handlers

import (
	"context"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
)

type appSvc interface {
	Create(ctx context.Context, input service.CreateApplicationInput) (domain.Application, error)
	GetByID(ctx context.Context, id string) (domain.Application, error)
	List(ctx context.Context) ([]domain.Application, error)
	Approve(ctx context.Context, id string) error
	Reject(ctx context.Context, id, reason string) error
	Delete(ctx context.Context, id string) error
}

type partnerSvc interface {
	ChangePassword(ctx context.Context, input service.ChangePasswordInput) error
	Profile(ctx context.Context, userID string) (domain.PartnerProfile, error)
}

type locationSvc interface {
	Create(ctx context.Context, input service.CreateLocationInput) (domain.Location, error)
	ListByPartner(ctx context.Context, partnerID string) ([]domain.Location, error)
}
