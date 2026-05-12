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
	"github.com/nlsnnn/berezhok/internal/modules/partner/service"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type employeeHandler struct {
	log       *slog.Logger
	validator *validator.Validator
	svc       employeeSvc
}

func NewEmployeeHandler(log *slog.Logger, svc employeeSvc) employeeHandler {
	return employeeHandler{
		log:       log,
		validator: validator.New(),
		svc:       svc,
	}
}

func (h *employeeHandler) List(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.employee.List"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	employees, err := h.svc.ListManagedByPartnerID(r.Context(), partnerID.String())
	if err != nil {
		log.Error("failed to list employees", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	response.Success(w, dto.MapSlice(employees, dto.FromManagedEmployee))
}

func (h *employeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.employee.Create"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	var req dto.CreateEmployeeRequest
	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		log.Error("validation failed", sl.Errs(errs))
		response.ValidationError(w, "validation failed", errs)
		return
	}

	employee, err := h.svc.CreateManaged(r.Context(), req.ToInput(partnerID.String()))
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrEmailAlreadyInUse),
			errors.Is(err, partnerErrors.ErrLocationNotOwnedByPartner):
			log.Warn("failed to create employee", sl.Err(err))
			response.BadRequest(w, err.Error())
		case errors.Is(err, partnerErrors.ErrLocationNotFound):
			log.Warn("location not found", sl.Err(err))
			response.NotFound(w, err.Error())
		default:
			log.Error("failed to create employee", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Created(w, dto.FromManagedEmployee(employee))
}

func (h *employeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "partner.handler.employee.Delete"
	log := h.log.With(slog.String("op", op))

	partnerID, err := contextx.PartnerID(r)
	if err != nil {
		log.Error("partner_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	actorEmployeeID, err := contextx.EmployeeID(r)
	if err != nil {
		log.Error("employee_id not found in context", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	employeeID := chi.URLParam(r, "id")
	err = h.svc.DeleteManaged(r.Context(), service.DeleteManagedEmployeeInput{
		PartnerID:       partnerID.String(),
		ActorEmployeeID: actorEmployeeID.String(),
		EmployeeID:      employeeID,
	})
	if err != nil {
		switch {
		case errors.Is(err, partnerErrors.ErrEmployeeNotFound):
			log.Warn("employee not found", sl.Err(err))
			response.NotFound(w, err.Error())
		case errors.Is(err, partnerErrors.ErrCannotDeleteOwner),
			errors.Is(err, partnerErrors.ErrCannotDeleteSelf):
			log.Warn("delete employee rejected", sl.Err(err))
			response.BadRequest(w, err.Error())
		default:
			log.Error("failed to delete employee", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	response.Success(w, map[string]string{"message": "employee deleted successfully"})
}
