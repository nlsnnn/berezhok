package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/netip"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/admin/domain"
	partnerDomain "github.com/nlsnnn/berezhok/internal/modules/partner/domain"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	partnerDTO "github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type applicationService interface {
	GetByID(ctx context.Context, id string) (partnerDomain.Application, error)
	List(ctx context.Context) ([]partnerDomain.Application, error)
	Approve(ctx context.Context, id string) error
	Reject(ctx context.Context, id, reason string) error
	Delete(ctx context.Context, id string) error
}

type applicationAuditRepo interface {
	MarkApplicationReviewed(ctx context.Context, applicationID, adminID uuid.UUID) error
	CreateAuditLog(ctx context.Context, input domain.AuditLog) (domain.AuditLog, error)
}

type ApplicationHandler struct {
	log       *slog.Logger
	validator *validator.Validator
	appSvc    applicationService
	auditRepo applicationAuditRepo
}

func NewApplicationHandler(log *slog.Logger, validator *validator.Validator, appSvc applicationService, auditRepo applicationAuditRepo) *ApplicationHandler {
	return &ApplicationHandler{
		log:       log,
		validator: validator,
		appSvc:    appSvc,
		auditRepo: auditRepo,
	}
}

func (h *ApplicationHandler) List(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appSvc.List(r.Context())
	if err != nil {
		h.log.Error("failed to list applications", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	items := partnerDTO.MapSlice(apps, partnerDTO.FromApplication)
	response.Success(w, map[string]any{
		"items": items,
		"pagination": map[string]any{
			"total":    len(items),
			"limit":    len(items),
			"offset":   0,
			"has_more": false,
		},
	})
}

func (h *ApplicationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	app, err := h.appSvc.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get application", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, partnerDTO.FromApplication(app))
}

func (h *ApplicationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	h.review(w, r, "admin.application.approve", func(ctx context.Context, id string) error {
		return h.appSvc.Approve(ctx, id)
	})
}

func (h *ApplicationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	var req partnerDTO.RejectApplicationRequest
	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		h.log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	h.review(w, r, "admin.application.reject", func(ctx context.Context, id string) error {
		return h.appSvc.Reject(ctx, id, req.RejectionReason)
	})
}

func (h *ApplicationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, applicationID, adminID, ok := h.adminActionContext(w, r)
	if !ok {
		return
	}

	if err := h.appSvc.Delete(r.Context(), id); err != nil {
		h.log.Error("failed to delete application", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	if err := h.audit(r, adminID, "admin.application.delete", "partner_application", applicationID, nil); err != nil {
		h.log.Error("failed to audit application delete", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, map[string]string{"message": "application deleted successfully"})
}

func (h *ApplicationHandler) review(w http.ResponseWriter, r *http.Request, action string, fn func(context.Context, string) error) {
	id, applicationID, adminID, ok := h.adminActionContext(w, r)
	if !ok {
		return
	}

	if err := fn(r.Context(), id); err != nil {
		if errors.Is(err, partnerErrors.ErrInvalidStatusTransition) {
			response.BadRequest(w, "only pending applications can be reviewed")
			return
		}
		h.log.Error("failed to review application", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	if err := h.auditRepo.MarkApplicationReviewed(r.Context(), applicationID, adminID); err != nil {
		h.log.Error("failed to mark application reviewed by admin", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	if err := h.audit(r, adminID, action, "partner_application", applicationID, nil); err != nil {
		h.log.Error("failed to audit application review", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, map[string]string{"message": "application reviewed successfully"})
}

func (h *ApplicationHandler) adminActionContext(w http.ResponseWriter, r *http.Request) (string, uuid.UUID, uuid.UUID, bool) {
	id := chi.URLParam(r, "id")
	applicationID, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return "", uuid.Nil, uuid.Nil, false
	}

	adminID, err := contextx.AdminID(r)
	if err != nil {
		response.Unauthorized(w, "authentication required")
		return "", uuid.Nil, uuid.Nil, false
	}

	return id, applicationID, adminID, true
}

func (h *ApplicationHandler) audit(r *http.Request, adminID uuid.UUID, action, entityType string, entityID uuid.UUID, details map[string]any) error {
	payload, err := json.Marshal(details)
	if err != nil {
		return err
	}

	_, err = h.auditRepo.CreateAuditLog(r.Context(), domain.AuditLog{
		AdminID:    adminID,
		Action:     action,
		EntityType: entityType,
		EntityID:   &entityID,
		Details:    payload,
		IPAddress:  requestIP(r),
	})
	return err
}

func requestIP(r *http.Request) *netip.Addr {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return nil
	}

	return &addr
}
