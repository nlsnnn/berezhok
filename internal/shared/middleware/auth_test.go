package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/auth"
	"github.com/nlsnnn/berezhok/internal/shared/authz"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
)

type tokenServiceStub struct {
	validateFn func(token string) (*auth.TokenClaims, error)
}

func (s *tokenServiceStub) Validate(token string) (*auth.TokenClaims, error) {
	return s.validateFn(token)
}

func TestRequireAuthStoresPartnerEmployeeContext(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	partnerID := uuid.New()
	locationID := uuid.New()

	mw := NewAuthMiddleware(&tokenServiceStub{
		validateFn: func(token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID:     userID,
				UserType:   "partner",
				Role:       "employee",
				PartnerID:  &partnerID,
				LocationID: &locationID,
			}, nil
		},
	})

	handler := mw.RequireAuth("partner")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{}

		employeeID, _ := contextx.EmployeeID(r)
		payload["employee_id"] = employeeID.String()

		gotPartnerID, _ := contextx.PartnerID(r)
		payload["partner_id"] = gotPartnerID.String()

		gotLocationID, _ := contextx.LocationID(r)
		payload["location_id"] = gotLocationID.String()

		role, _ := contextx.UserRole(r)
		payload["role"] = role

		actor, _ := contextx.PartnerActor(r)
		payload["actor_role"] = string(actor.Role)
		payload["actor_partner_id"] = actor.PartnerID.String()

		_ = json.NewEncoder(w).Encode(payload)
	}))

	req := httptest.NewRequest(http.MethodGet, "/partner/orders", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if payload["employee_id"] != userID.String() {
		t.Fatalf("expected employee_id %s, got %s", userID, payload["employee_id"])
	}
	if payload["partner_id"] != partnerID.String() {
		t.Fatalf("expected partner_id %s, got %s", partnerID, payload["partner_id"])
	}
	if payload["location_id"] != locationID.String() {
		t.Fatalf("expected location_id %s, got %s", locationID, payload["location_id"])
	}
	if payload["role"] != "employee" {
		t.Fatalf("expected role employee, got %s", payload["role"])
	}
	if payload["actor_role"] != "employee" {
		t.Fatalf("expected actor role employee, got %s", payload["actor_role"])
	}
	if payload["actor_partner_id"] != partnerID.String() {
		t.Fatalf("expected actor partner_id %s, got %s", partnerID, payload["actor_partner_id"])
	}
}

func TestRequirePermissionRejectsMissingPermission(t *testing.T) {
	t.Parallel()

	mw := NewAuthMiddleware(&tokenServiceStub{
		validateFn: func(token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID:   uuid.New(),
				UserType: "partner",
				Role:     "employee",
			}, nil
		},
	})

	handler := mw.RequireAuth("partner")(mw.RequirePermission(authz.PermissionPartnerEmployeesManage)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodGet, "/partner/employees", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rr.Code)
	}
}

func TestRequirePermissionAllowsManagerBoxes(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()
	locationID := uuid.New()
	mw := NewAuthMiddleware(&tokenServiceStub{
		validateFn: func(token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID:     uuid.New(),
				UserType:   "partner",
				Role:       "manager",
				PartnerID:  &partnerID,
				LocationID: &locationID,
			}, nil
		},
	})

	handler := mw.RequireAuth("partner")(mw.RequirePermission(authz.PermissionPartnerBoxesManage)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodGet, "/partner/boxes", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestRequireAuthStoresAdminContext(t *testing.T) {
	t.Parallel()

	adminID := uuid.New()
	mw := NewAuthMiddleware(&tokenServiceStub{
		validateFn: func(token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID:   adminID,
				UserType: "admin",
				Role:     "super_admin",
			}, nil
		},
	})

	handler := mw.RequireAuth("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{}

		gotAdminID, _ := contextx.AdminID(r)
		payload["admin_id"] = gotAdminID.String()

		actor, _ := contextx.AdminActor(r)
		payload["actor_role"] = string(actor.Role)
		payload["actor_admin_id"] = actor.AdminID.String()

		_ = json.NewEncoder(w).Encode(payload)
	}))

	req := httptest.NewRequest(http.MethodGet, "/admin/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if payload["admin_id"] != adminID.String() {
		t.Fatalf("expected admin_id %s, got %s", adminID, payload["admin_id"])
	}
	if payload["actor_role"] != "super_admin" {
		t.Fatalf("expected actor role super_admin, got %s", payload["actor_role"])
	}
	if payload["actor_admin_id"] != adminID.String() {
		t.Fatalf("expected actor admin_id %s, got %s", adminID, payload["actor_admin_id"])
	}
}

func TestRequireAdminPermissionRejectsSupportManage(t *testing.T) {
	t.Parallel()

	mw := NewAuthMiddleware(&tokenServiceStub{
		validateFn: func(token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID:   uuid.New(),
				UserType: "admin",
				Role:     "support",
			}, nil
		},
	})

	handler := mw.RequireAuth("admin")(mw.RequireAdminPermission(authz.PermissionAdminOpsManage)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodPatch, "/admin/partners/id", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rr.Code)
	}
}

func TestRequireAdminPermissionAllowsAdminManage(t *testing.T) {
	t.Parallel()

	mw := NewAuthMiddleware(&tokenServiceStub{
		validateFn: func(token string) (*auth.TokenClaims, error) {
			return &auth.TokenClaims{
				UserID:   uuid.New(),
				UserType: "admin",
				Role:     "admin",
			}, nil
		},
	})

	handler := mw.RequireAuth("admin")(mw.RequireAdminPermission(authz.PermissionAdminOpsManage)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodPatch, "/admin/partners/id", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}
