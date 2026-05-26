package dto

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
)

func (r CreateApplicationRequest) ToInput() service.CreateApplicationInput {
	return service.CreateApplicationInput{
		ContactName:  r.ContactName,
		ContactEmail: r.ContactEmail,
		ContactPhone: r.ContactPhone,
		BusinessName: r.BusinessName,
		CategoryCode: r.CategoryCode,
		Address:      r.Address,
		Longitude:    r.Longitude,
		Latitude:     r.Latitude,
		Description:  r.Description,
	}
}

func (r CreateLocationRequest) ToInput(partnerID string) service.CreateLocationInput {
	return service.CreateLocationInput{
		PartnerID:    partnerID,
		CategoryCode: r.CategoryCode,
		Name:         r.Name,
		Address:      r.Address,
		Latitude:     r.Latitude,
		Longitude:    r.Longitude,
	}
}

func (r AddLegalInfoRequest) ToInput(partnerID string) service.AddLegalInfoInput {
	return service.AddLegalInfoInput{
		PartnerID:    partnerID,
		Inn:          r.Inn,
		Ogrn:         r.Ogrn,
		Kpp:          r.Kpp,
		LegalAddress: r.LegalAddress,
		LegalName:    r.LegalName,
	}
}

func (r CreateEmployeeRequest) ToInput(actor authz.PartnerActor) service.CreateManagedEmployeeInput {
	return service.CreateManagedEmployeeInput{
		Actor:      actor,
		LocationID: r.LocationID,
		Email:      r.Email,
		Name:       r.Name,
	}
}

func FromApplication(a domain.Application) ApplicationResponse {
	return ApplicationResponse{
		ID:              a.ID,
		ContactName:     a.ContactName,
		ContactEmail:    a.ContactEmail,
		ContactPhone:    a.ContactPhone,
		BusinessName:    a.BusinessName,
		CategoryCode:    a.CategoryCode,
		Address:         a.Address,
		Description:     a.Description,
		Status:          string(a.Status),
		ReviewedAt:      a.ReviewedAt,
		RejectionReason: a.RejectionReason,
		CreatedAt:       a.CreatedAt,
	}
}

func FromLocation(l domain.Location) LocationResponse {
	var workingHours map[string]string
	if l.WorkingHours != "" {
		_ = json.Unmarshal([]byte(l.WorkingHours), &workingHours)
	}

	return LocationResponse{
		ID:            l.ID,
		Name:          l.Name,
		Address:       l.Address,
		Phone:         l.Phone,
		CategoryCode:  l.Category.Code,
		Status:        string(l.Status),
		LogoURL:       l.LogoURL,
		CoverImageURL: l.CoverImageURL,
		WorkingHours:  workingHours,
		CreatedAt:     l.CreatedAt,
	}
}

func (r UpdateLocationRequest) ToInput(locationID, partnerID uuid.UUID) service.UpdateLocationInput {
	return service.UpdateLocationInput{
		ID:            locationID,
		PartnerID:     partnerID,
		Phone:         r.Phone,
		LogoURL:       r.LogoURL,
		CoverImageURL: r.CoverImageURL,
		WorkingHours:  r.WorkingHours,
	}
}

func FromManagedEmployee(employee domain.ManagedEmployee) ManagedEmployeeResponse {
	return ManagedEmployeeResponse{
		ID:                 employee.ID,
		Email:              employee.Email,
		Name:               employee.Name,
		Role:               string(employee.Role),
		LocationID:         employee.LocationID,
		LocationName:       employee.LocationName,
		MustChangePassword: employee.MustChangePassword,
		CreatedAt:          employee.CreatedAt,
	}
}
