package domain

import "time"

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
	Status          ApplicationStatus
	ReviewedAt      *time.Time
	RejectionReason string
	CreatedAt       time.Time
}
