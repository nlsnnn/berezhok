package authz

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestRolePermissionsMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		role       Role
		permission Permission
		want       bool
	}{
		{name: "owner manages employees", role: RoleOwner, permission: PermissionPartnerEmployeesManage, want: true},
		{name: "owner manages boxes", role: RoleOwner, permission: PermissionPartnerBoxesManage, want: true},
		{name: "manager manages boxes", role: RoleManager, permission: PermissionPartnerBoxesManage, want: true},
		{name: "manager cannot manage employees", role: RoleManager, permission: PermissionPartnerEmployeesManage, want: false},
		{name: "employee can pickup orders", role: RoleEmployee, permission: PermissionPartnerOrdersPickup, want: true},
		{name: "employee cannot manage boxes", role: RoleEmployee, permission: PermissionPartnerBoxesManage, want: false},
		{name: "employee changes password", role: RoleEmployee, permission: PermissionPartnerPasswordChange, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := PartnerActor{Role: tt.role}
			if got := actor.Can(tt.permission); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestPartnerActorCanAccessLocation(t *testing.T) {
	t.Parallel()

	locationID := uuid.New()
	otherLocationID := uuid.New()

	tests := []struct {
		name       string
		actor      PartnerActor
		locationID uuid.UUID
		wantErr    error
	}{
		{name: "owner can access any location", actor: PartnerActor{Role: RoleOwner}, locationID: locationID},
		{name: "manager can access assigned location", actor: PartnerActor{Role: RoleManager, LocationID: &locationID}, locationID: locationID},
		{name: "manager without location is denied", actor: PartnerActor{Role: RoleManager}, locationID: locationID, wantErr: ErrLocationScopeDenied},
		{name: "manager cannot access other location", actor: PartnerActor{Role: RoleManager, LocationID: &otherLocationID}, locationID: locationID, wantErr: ErrLocationScopeDenied},
		{name: "employee can access assigned location", actor: PartnerActor{Role: RoleEmployee, LocationID: &locationID}, locationID: locationID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.actor.EnsureLocation(tt.locationID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
