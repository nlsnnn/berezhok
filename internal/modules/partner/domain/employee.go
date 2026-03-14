package domain

import "time"

type EmployeeRole string

const (
	EmployeeRoleOwner    EmployeeRole = "owner"
	EmployeeRoleManager  EmployeeRole = "manager"
	EmployeeRoleEmployee EmployeeRole = "employee"
)

type Employee struct {
	ID                 string
	PartnerID          string
	LocationID         string
	Email              string
	PasswordHash       string
	Role               EmployeeRole
	Name               string
	MustChangePassword bool
	CreatedAt          time.Time
}
