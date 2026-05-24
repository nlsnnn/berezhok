package handlers

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
)

func TestVerifyPartnerLegalInfoRequiresAdminID(t *testing.T) {
	t.Parallel()

	handler := NewOpsHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New(), nil, fakeApplicationAuditRepo{})
	request := httptest.NewRequest(http.MethodPost, "/admin/partners/11111111-1111-1111-1111-111111111111/legal-info/verify", nil)
	recorder := httptest.NewRecorder()

	handler.VerifyPartnerLegalInfo(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestRejectPartnerLegalInfoRequiresComment(t *testing.T) {
	t.Parallel()

	adminID := uuid.New()
	handler := NewOpsHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), validator.New(), nil, fakeApplicationAuditRepo{})
	request := httptest.NewRequest(http.MethodPost, "/admin/partners/11111111-1111-1111-1111-111111111111/legal-info/reject", bytes.NewReader([]byte(`{"verification_comment":""}`)))
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("partner_id", "11111111-1111-1111-1111-111111111111")
	request = request.WithContext(context.WithValue(request.Context(), contextx.AdminIDKey, adminID))
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
	recorder := httptest.NewRecorder()

	handler.RejectPartnerLegalInfo(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestVerifyPartnerLegalInfoActivatesDraftLocations(t *testing.T) {
	t.Parallel()

	partnerID := uuid.New()
	adminID := uuid.New()
	queries := &fakeLegalInfoReviewQueries{}

	_, _, err := verifyPartnerLegalInfo(context.Background(), queries, partnerID, adminID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !queries.legalInfoVerified {
		t.Fatal("expected legal info to be verified")
	}
	if !queries.partnerActivated {
		t.Fatal("expected partner to be activated")
	}
	if !queries.locationsActivated {
		t.Fatal("expected draft locations to be activated")
	}
}

type fakeLegalInfoReviewQueries struct {
	legalInfoVerified  bool
	partnerActivated   bool
	locationsActivated bool
}

func (q *fakeLegalInfoReviewQueries) VerifyPartnerLegalInfo(ctx context.Context, arg sqlc.VerifyPartnerLegalInfoParams) (sqlc.PartnerLegalInfo, error) {
	q.legalInfoVerified = true
	return sqlc.PartnerLegalInfo{PartnerID: arg.PartnerID, VerificationStatus: pgtype.Text{String: "verified", Valid: true}}, nil
}

func (q *fakeLegalInfoReviewQueries) ActivatePartnerIfPendingDocuments(ctx context.Context, id uuid.UUID) (sqlc.Partner, error) {
	q.partnerActivated = true
	return sqlc.Partner{ID: id, Status: "active"}, nil
}

func (q *fakeLegalInfoReviewQueries) ActivatePartnerDraftLocations(ctx context.Context, partnerID uuid.UUID) (int64, error) {
	q.locationsActivated = true
	return 1, nil
}
