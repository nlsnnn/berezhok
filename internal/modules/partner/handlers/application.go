package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	partnerErrors "github.com/nlsnnn/berezhok/internal/modules/partner/errors"
	"github.com/nlsnnn/berezhok/internal/modules/partner/handlers/dto"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type appHandler struct {
	log        *slog.Logger
	validator  *validator.Validator
	appService appSvc
}

func NewApplicationHandler(log *slog.Logger, svc appSvc) appHandler {
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

	app, err := a.appService.Create(r.Context(), req.ToInput())
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrLocationCategoryNotFound):
			a.log.Warn("invalid location category code", sl.Err(err))
			response.BadRequest(w, "invalid location category code")
		case errors.Is(err, partnerErrors.ErrEmailAlreadyInUse):
			a.log.Warn("email already in use", sl.Err(err))
			response.BadRequest(w, "email already in use")
		default:
			a.log.Error("failed to create application", sl.Err(err))
			response.InternalError(w, err)
		}
		return
	}

	response.Success(w, dto.FromApplication(app))
}

func (a *appHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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
	id := chi.URLParam(r, "id")

	if err := a.appService.Delete(r.Context(), id); err != nil {
		a.log.Error("failed to delete application", sl.Err(err))
		response.InternalError(w, err)
		return
	}
	response.Success(w, map[string]string{"message": "application deleted successfully"})
}

func (a *appHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := a.appService.Approve(r.Context(), id); err != nil {
		if errors.Is(err, partnerErrors.ErrInvalidStatusTransition) {
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
	id := chi.URLParam(r, "id")

	var req dto.RejectApplicationRequest
	if errs := a.validator.DecodeAndValidate(r, &req); errs != nil {
		a.log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	if err := a.appService.Reject(r.Context(), id, req.RejectionReason); err != nil {
		if errors.Is(err, partnerErrors.ErrInvalidStatusTransition) {
			response.BadRequest(w, "only pending applications can be rejected")
			return
		}
		a.log.Error("failed to reject application", sl.Err(err))
		response.InternalError(w, err)
		return
	}

	response.Success(w, map[string]string{"message": "application rejected successfully"})
}
