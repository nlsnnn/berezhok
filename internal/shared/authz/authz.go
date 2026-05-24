package authz

import (
	"errors"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner    Role = "owner"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
)

type Permission string

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleAdmin      AdminRole = "admin"
	AdminRoleSupport    AdminRole = "support"
)

type AdminPermission string

const (
	PermissionPartnerPasswordChange  Permission = "partner.password.change"
	PermissionPartnerOrdersView      Permission = "partner.orders.view"
	PermissionPartnerOrdersPickup    Permission = "partner.orders.pickup"
	PermissionPartnerBoxesManage     Permission = "partner.boxes.manage"
	PermissionPartnerStatsView       Permission = "partner.stats.view"
	PermissionPartnerEmployeesManage Permission = "partner.employees.manage"
	PermissionPartnerLocationsManage Permission = "partner.locations.manage"
	PermissionPartnerProfileManage   Permission = "partner.profile.manage"
	PermissionPartnerLegalInfoManage Permission = "partner.legal_info.manage"
	PermissionPartnerDashboardView   Permission = "partner.dashboard.view"
	PermissionPartnerMediaUpload     Permission = "partner.media.upload"
	PermissionPartnerPayoutsView     Permission = "partner.payouts.view"
	PermissionPartnerPayoutsManage   Permission = "partner.payouts.manage"

	PermissionAdminOpsView      AdminPermission = "admin.ops.view"
	PermissionAdminOpsManage    AdminPermission = "admin.ops.manage"
	PermissionAdminAdminsManage AdminPermission = "admin.admins.manage"
	PermissionAdminAuditView    AdminPermission = "admin.audit.view"
)

var (
	ErrForbidden           = errors.New("access denied")
	ErrLocationScopeDenied = errors.New("location scope denied")
)

var RolePermissions = map[Role]map[Permission]struct{}{
	RoleOwner: {
		PermissionPartnerPasswordChange:  {},
		PermissionPartnerOrdersView:      {},
		PermissionPartnerOrdersPickup:    {},
		PermissionPartnerBoxesManage:     {},
		PermissionPartnerStatsView:       {},
		PermissionPartnerEmployeesManage: {},
		PermissionPartnerLocationsManage: {},
		PermissionPartnerProfileManage:   {},
		PermissionPartnerLegalInfoManage: {},
		PermissionPartnerDashboardView:   {},
		PermissionPartnerMediaUpload:     {},
		PermissionPartnerPayoutsView:     {},
		PermissionPartnerPayoutsManage:   {},
	},
	RoleManager: {
		PermissionPartnerPasswordChange: {},
		PermissionPartnerOrdersView:     {},
		PermissionPartnerOrdersPickup:   {},
		PermissionPartnerBoxesManage:    {},
		PermissionPartnerStatsView:      {},
	},
	RoleEmployee: {
		PermissionPartnerPasswordChange: {},
		PermissionPartnerOrdersView:     {},
		PermissionPartnerOrdersPickup:   {},
	},
}

var AdminRolePermissions = map[AdminRole]map[AdminPermission]struct{}{
	AdminRoleSuperAdmin: {
		PermissionAdminOpsView:      {},
		PermissionAdminOpsManage:    {},
		PermissionAdminAdminsManage: {},
		PermissionAdminAuditView:    {},
	},
	AdminRoleAdmin: {
		PermissionAdminOpsView:   {},
		PermissionAdminOpsManage: {},
		PermissionAdminAuditView: {},
	},
	AdminRoleSupport: {
		PermissionAdminOpsView:   {},
		PermissionAdminAuditView: {},
	},
}

type PartnerActor struct {
	PartnerID  uuid.UUID
	EmployeeID uuid.UUID
	Role       Role
	LocationID *uuid.UUID
}

type AdminActor struct {
	AdminID uuid.UUID
	Role    AdminRole
}

func (a AdminActor) Can(permission AdminPermission) bool {
	permissions, ok := AdminRolePermissions[a.Role]
	if !ok {
		return false
	}

	_, ok = permissions[permission]
	return ok
}

func (a AdminActor) EnsureCan(permission AdminPermission) error {
	if !a.Can(permission) {
		return ErrForbidden
	}

	return nil
}

func (a PartnerActor) Can(permission Permission) bool {
	permissions, ok := RolePermissions[a.Role]
	if !ok {
		return false
	}

	_, ok = permissions[permission]
	return ok
}

func (a PartnerActor) EnsureCan(permission Permission) error {
	if !a.Can(permission) {
		return ErrForbidden
	}

	return nil
}

func (a PartnerActor) EnsureLocation(locationID uuid.UUID) error {
	if a.Role == RoleOwner {
		return nil
	}

	if a.LocationID == nil || *a.LocationID != locationID {
		return ErrLocationScopeDenied
	}

	return nil
}

func (a PartnerActor) ScopedLocationID() *uuid.UUID {
	if a.Role == RoleOwner {
		return nil
	}

	return a.LocationID
}
