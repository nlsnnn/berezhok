package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/admin/domain"
	partnerDomain "github.com/nlsnnn/berezhok/internal/modules/partner/domain"
)

type fakeApplicationService struct {
	lastInput partnerDomain.ApplicationListInput
}

func (f *fakeApplicationService) GetByID(context.Context, string) (partnerDomain.Application, error) {
	return partnerDomain.Application{}, nil
}

func (f *fakeApplicationService) List(_ context.Context, input partnerDomain.ApplicationListInput) (partnerDomain.ApplicationListResult, error) {
	f.lastInput = input
	return partnerDomain.ApplicationListResult{
		Items: []partnerDomain.Application{{
			ID:           "10000000-0000-0000-0000-000000000001",
			ContactName:  "Ирина Соколова",
			ContactEmail: "pending.partner@berezhok.local",
			BusinessName: "Булочка у дома",
			Status:       partnerDomain.ApplicationStatusPending,
		}},
		Total: 1,
	}, nil
}

func (f *fakeApplicationService) Approve(context.Context, string) error { return nil }
func (f *fakeApplicationService) Reject(context.Context, string, string) error {
	return nil
}
func (f *fakeApplicationService) Delete(context.Context, string) error { return nil }

type fakeApplicationAuditRepo struct{}

func (f fakeApplicationAuditRepo) MarkApplicationReviewed(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}

func (f fakeApplicationAuditRepo) CreateAuditLog(context.Context, domain.AuditLog) (domain.AuditLog, error) {
	return domain.AuditLog{}, nil
}

func TestApplicationHandlerListUsesSearchAndPagination(t *testing.T) {
	appSvc := &fakeApplicationService{}
	handler := NewApplicationHandler(slog.Default(), nil, appSvc, fakeApplicationAuditRepo{})
	request := httptest.NewRequest(http.MethodGet, "/admin/applications?search=jopa%40mail.ru&status=pending&limit=3&offset=6", nil)
	recorder := httptest.NewRecorder()

	handler.List(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if appSvc.lastInput.Search != "jopa@mail.ru" {
		t.Fatalf("expected search to be passed to service, got %q", appSvc.lastInput.Search)
	}
	if appSvc.lastInput.Status != "pending" {
		t.Fatalf("expected status to be passed to service, got %q", appSvc.lastInput.Status)
	}
	if appSvc.lastInput.Limit != 3 || appSvc.lastInput.Offset != 6 {
		t.Fatalf("expected limit/offset 3/6, got %d/%d", appSvc.lastInput.Limit, appSvc.lastInput.Offset)
	}

	var body struct {
		Data struct {
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Pagination.Total != 1 || body.Data.Pagination.Limit != 3 || body.Data.Pagination.Offset != 6 {
		t.Fatalf("unexpected pagination: %+v", body.Data.Pagination)
	}
}
