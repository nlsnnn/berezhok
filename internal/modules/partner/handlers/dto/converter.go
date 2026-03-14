package dto

import (
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/domain"
)

func (r CreateApplicationRequest) ToModel() sqlc.CreateApplicationParams {
	return sqlc.CreateApplicationParams{
		ContactName:  r.ContactName,
		ContactEmail: r.ContactEmail,
		ContactPhone: r.ContactPhone,
		BusinessName: r.BusinessName,
		CategoryCode: ToText(r.CategoryCode),
		Address:      ToText(r.Address),
		Description:  ToText(r.Description),
		Status:       "pending",
	}
}

func FromLocation(l domain.Location) LocationResponse {
	return LocationResponse{
		ID:      l.ID,
		Name:    l.Name,
		Address: l.Address,
	}
}
