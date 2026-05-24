package handlers

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

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
	request = request.WithContext(context.WithValue(request.Context(), contextx.AdminIDKey, adminID))
	recorder := httptest.NewRecorder()

	handler.RejectPartnerLegalInfo(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}
