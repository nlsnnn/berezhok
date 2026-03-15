package dto

import (
	"github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
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
	return LocationResponse{
		ID:      l.ID,
		Name:    l.Name,
		Address: l.Address,
	}
}
