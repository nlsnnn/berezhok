package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/auth"
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
}

func TestRequirePartnerRolesRejectsNonOwner(t *testing.T) {
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

	handler := mw.RequireAuth("partner")(mw.RequirePartnerRoles("owner")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
