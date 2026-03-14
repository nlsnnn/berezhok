package dto

import (
	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
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

func FromLocation(m sqlc.Location) LocationResponse {
	return LocationResponse{
		ID:      m.ID,
		Name:    m.Name,
		Address: m.Address,
	}
}
