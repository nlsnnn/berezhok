package domain

import (
	"time"

	"github.com/nlsnnn/berezhok/internal/shared/domain"
	"github.com/nlsnnn/berezhok/internal/shared/errors"
)

type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusApproved ApplicationStatus = "approved"
	ApplicationStatusRejected ApplicationStatus = "rejected"
)

type Application struct {
	ID              string
	ContactName     string
	ContactEmail    string
	ContactPhone    string
	BusinessName    string
	CategoryCode    string
	Address         string
	Description     string
	Coords          domain.GeoPoint
	Status          ApplicationStatus
	ReviewedAt      *time.Time
	RejectionReason string
	CreatedAt       time.Time
}

func NewApplication(contactName, contactEmail, contactPhone, businessName, categoryCode, address, description string, coords domain.GeoPoint) (Application, error) {
	if contactName == "" {
		return Application{}, errors.ErrInvalidInput
	}

	return Application{
		ContactName:  contactName,
		ContactEmail: contactEmail,
		ContactPhone: contactPhone,
		BusinessName: businessName,
		CategoryCode: categoryCode,
		Address:      address,
		Description:  description,
		Coords:       coords,
		Status:       ApplicationStatusPending,
	}, nil
}

func (a Application) CanTransitionTo(newStatus ApplicationStatus) bool {
	switch a.Status {
	case ApplicationStatusPending:
		return newStatus == ApplicationStatusApproved || newStatus == ApplicationStatusRejected
	default:
		return false
	}
}
