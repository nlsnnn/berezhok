package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nlsnnn/berezhok/internal/modules/partner"
	"github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/logger/sl"
	"github.com/nlsnnn/berezhok/internal/shared/response"
	"github.com/nlsnnn/berezhok/internal/shared/validator"
)

type appHandler struct {
	log        *slog.Logger
	appService appSvc
	validator  *validator.Validator
}

func NewApplicationHandler(
	log *slog.Logger,
	svc appSvc,
) appHandler {
	return appHandler{
		log:        log,
		appService: svc,
		validator:  validator.New(),
	}
}

func (a *appHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateApplicationRequest

	if errs := a.validator.DecodeAndValidate(r, &req); errs != nil {
		a.log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	app, err := a.appService.Create(r.Context(), req.ToModel())
	if err != nil {
		a.log.Error("failed to create application", sl.Err(err))
		response.InternalError(w, err)
		return
	}

	response.Success(w, dto.FromApplication(app))
}

func (a *appHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		a.log.Error("failed to parse id", sl.Err(err))
		response.BadRequest(w, "invalid id")
		return
	}

	app, err := a.appService.GetByID(r.Context(), id)
	if err != nil {
		a.log.Error("failed to get application by id", sl.Err(err))
		response.InternalError(w, err)
		return
	}
	response.Success(w, dto.FromApplication(app))
}

func (a *appHandler) List(w http.ResponseWriter, r *http.Request) {
	apps, err := a.appService.List(r.Context())
	if err != nil {
		a.log.Error("failed to list applications", sl.Err(err))
		response.InternalError(w, err)
		return
	}

	response.Success(w, dto.MapSlice(apps, dto.FromApplication))
}

func (a *appHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		a.log.Error("failed to parse id", sl.Err(err))
		response.BadRequest(w, "invalid id")
		return
	}
	err = a.appService.Delete(r.Context(), id)
	if err != nil {
		a.log.Error("failed to delete application", sl.Err(err))
		response.InternalError(w, err)
		return
	}
	response.Success(w, map[string]string{"message": "application deleted successfully"})
}

func (a *appHandler) Approve(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		a.log.Error("failed to parse id", sl.Err(err))
		response.BadRequest(w, "invalid id")
		return
	}

	err = a.appService.Approve(r.Context(), id)
	if err != nil {
		if errors.Is(err, partner.ErrInvalidStatusTransition) {
			response.BadRequest(w, "only pending applications can be approved")
			return
		}
		a.log.Error("failed to approve application", sl.Err(err))
		response.InternalError(w, err)
		return
	}

	response.Success(w, map[string]string{"message": "application approved successfully"})
}

func (a *appHandler) Reject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		a.log.Error("failed to parse id", sl.Err(err))
		response.BadRequest(w, "invalid id")
		return
	}

	var req dto.RejectApplicationRequest
	if errs := a.validator.DecodeAndValidate(r, &req); errs != nil {
		a.log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	err = a.appService.Reject(r.Context(), id, req.RejectionReason)
	if err != nil {
		if errors.Is(err, partner.ErrInvalidStatusTransition) {
			response.BadRequest(w, "only pending applications can be rejected")
			return
		}
		a.log.Error("failed to reject application", sl.Err(err))
		response.InternalError(w, err)
		return
	}

	response.Success(w, map[string]string{"message": "application rejected successfully"})
}
