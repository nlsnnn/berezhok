package handlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
)

type appSvc interface {
	Create(ctx context.Context, input service.CreateApplicationInput) (domain.Application, error)
	GetByID(ctx context.Context, id string) (domain.Application, error)
	List(ctx context.Context, input domain.ApplicationListInput) (domain.ApplicationListResult, error)
	Approve(ctx context.Context, id string) error
	Reject(ctx context.Context, id, reason string) error
	Delete(ctx context.Context, id string) error
}

type partnerSvc interface {
	ChangePassword(ctx context.Context, input service.ChangePasswordInput) error
	UpdateEmployeeName(ctx context.Context, employeeID, name string) error
	Profile(ctx context.Context, userID string) (domain.PartnerProfile, error)
	Dashboard(ctx context.Context, userID string) (domain.PartnerDashboard, error)
	Stats(ctx context.Context, actor authz.PartnerActor, filter domain.StatsFilter) (domain.PartnerStats, error)
	AddLegalInfo(ctx context.Context, input service.AddLegalInfoInput) error
}

type locationSvc interface {
	Create(ctx context.Context, input service.CreateLocationInput) (domain.Location, error)
	ListByPartner(ctx context.Context, partnerID string) ([]domain.Location, error)
	GetByID(ctx context.Context, partnerID, locationID uuid.UUID) (domain.Location, error)
	Update(ctx context.Context, input service.UpdateLocationInput) (domain.Location, error)
}

type employeeSvc interface {
	ListManagedByPartnerID(ctx context.Context, actor authz.PartnerActor) ([]domain.ManagedEmployee, error)
	CreateManaged(ctx context.Context, input service.CreateManagedEmployeeInput) (domain.ManagedEmployee, error)
	DeleteManaged(ctx context.Context, input service.DeleteManagedEmployeeInput) error
}

type pinsSvc interface {
	ListAvailable(ctx context.Context) ([]domain.LocationPin, error)
	UpdateLocationPins(ctx context.Context, partnerID string, locationID uuid.UUID, pinCodes []string) error
	GetForLocation(ctx context.Context, locationID uuid.UUID) ([]domain.LocationPin, error)
}
