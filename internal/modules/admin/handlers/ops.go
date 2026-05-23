package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"

	"github.com/nlsnnn/berezhok/internal/adapters/postgresql/sqlc"
	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/pgconverter"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	"github.com/nlsnnn/berezhok/internal/modules/admin/domain"
	"github.com/nlsnnn/berezhok/internal/shared/auth"
	"github.com/nlsnnn/berezhok/internal/shared/contextx"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

const (
	defaultAdminLimit = 20
	maxAdminLimit     = 100
)

type OpsHandler struct {
	log       *slog.Logger
	validator *validator.Validator
	q         *sqlc.Queries
	auditRepo applicationAuditRepo
}

func NewOpsHandler(log *slog.Logger, validator *validator.Validator, q *sqlc.Queries, auditRepo applicationAuditRepo) *OpsHandler {
	return &OpsHandler{log: log, validator: validator, q: q, auditRepo: auditRepo}
}

type createAdminRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=super_admin admin support"`
}

type updateAdminRequest struct {
	Name     *string `json:"name"`
	Role     *string `json:"role" validate:"omitempty,oneof=super_admin admin support"`
	IsActive *bool   `json:"is_active"`
}

type updatePartnerRequest struct {
	Status               *string  `json:"status" validate:"omitempty,oneof=pending_documents active suspended blocked"`
	CommissionRate       *float64 `json:"commission_rate" validate:"omitempty,gte=0,lte=1"`
	PromoCommissionRate  *float64 `json:"promo_commission_rate" validate:"omitempty,gte=0,lte=1"`
	PromoCommissionUntil *string  `json:"promo_commission_until"`
}

type updateStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

func (h *OpsHandler) Me(w http.ResponseWriter, r *http.Request) {
	adminID, err := contextx.AdminID(r)
	if err != nil {
		response.Unauthorized(w, "authentication required")
		return
	}

	admin, err := h.q.FindAdminByID(r.Context(), adminID)
	if err != nil {
		h.writeError(w, "failed to get admin", err)
		return
	}

	response.Success(w, adminResponse(admin))
}

func (h *OpsHandler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	search := r.URL.Query().Get("search")

	items, err := h.q.ListAdminUsers(r.Context(), sqlc.ListAdminUsersParams{Search: search, PageLimit: int32(p.Limit), PageOffset: int32(p.Offset)})
	if err != nil {
		h.writeError(w, "failed to list admins", err)
		return
	}
	total, err := h.q.CountAdminUsers(r.Context(), search)
	if err != nil {
		h.writeError(w, "failed to count admins", err)
		return
	}

	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = adminResponse(item)
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req createAdminRequest
	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		response.ValidationError(w, "validation failed", errs)
		return
	}

	passwordHash, err := auth.Hash(req.Password)
	if err != nil {
		h.writeError(w, "failed to hash password", err)
		return
	}

	admin, err := h.q.CreateAdminUser(r.Context(), sqlc.CreateAdminUserParams{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Role:         req.Role,
		IsActive:     true,
	})
	if err != nil {
		h.writeError(w, "failed to create admin", err)
		return
	}
	h.auditMutation(r, "admin.admin.create", "admin_user", admin.ID, map[string]any{"email": admin.Email, "role": admin.Role})

	response.Created(w, adminResponse(admin))
}

func (h *OpsHandler) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	var req updateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid json")
		return
	}

	admin, err := h.q.UpdateAdminUser(r.Context(), sqlc.UpdateAdminUserParams{
		ID:       id,
		Name:     optionalText(req.Name),
		Role:     optionalText(req.Role),
		IsActive: optionalBool(req.IsActive),
	})
	if err != nil {
		h.writeError(w, "failed to update admin", err)
		return
	}
	h.auditMutation(r, "admin.admin.update", "admin_user", admin.ID, nil)

	response.Success(w, adminResponse(admin))
}

func (h *OpsHandler) DeactivateAdmin(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	if _, err := h.q.DeactivateAdminUser(r.Context(), id); err != nil {
		h.writeError(w, "failed to deactivate admin", err)
		return
	}
	h.auditMutation(r, "admin.admin.deactivate", "admin_user", id, nil)
	response.Success(w, map[string]string{"message": "admin deactivated"})
}

func (h *OpsHandler) ListAudit(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	params := sqlc.ListAdminAuditLogParams{
		ActionFilter:     r.URL.Query().Get("action"),
		EntityTypeFilter: r.URL.Query().Get("entity_type"),
		PageLimit:        int32(p.Limit),
		PageOffset:       int32(p.Offset),
	}
	items, err := h.q.ListAdminAuditLog(r.Context(), params)
	if err != nil {
		h.writeError(w, "failed to list audit log", err)
		return
	}
	total, err := h.q.CountAdminAuditLog(r.Context(), sqlc.CountAdminAuditLogParams{
		ActionFilter:     params.ActionFilter,
		EntityTypeFilter: params.EntityTypeFilter,
	})
	if err != nil {
		h.writeError(w, "failed to count audit log", err)
		return
	}
	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = map[string]any{
			"id": item.ID, "admin_user_id": item.AdminUserID, "admin_email": item.AdminEmail, "admin_name": item.AdminName,
			"action": item.Action, "entity_type": item.EntityType.String, "entity_id": nullableUUID(item.EntityID),
			"details": json.RawMessage(item.Details), "ip_address": item.IpAddress, "created_at": item.CreatedAt,
		}
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) ListPartners(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	params := sqlc.ListAdminPartnersParams{StatusFilter: r.URL.Query().Get("status"), Search: r.URL.Query().Get("search"), PageLimit: int32(p.Limit), PageOffset: int32(p.Offset)}
	items, err := h.q.ListAdminPartners(r.Context(), params)
	if err != nil {
		h.writeError(w, "failed to list partners", err)
		return
	}
	total, err := h.q.CountAdminPartners(r.Context(), sqlc.CountAdminPartnersParams{StatusFilter: params.StatusFilter, Search: params.Search})
	if err != nil {
		h.writeError(w, "failed to count partners", err)
		return
	}
	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = partnerListResponse(item)
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) GetPartner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := h.q.GetAdminPartnerByID(r.Context(), id)
	if err != nil {
		h.writeError(w, "failed to get partner", err)
		return
	}
	response.Success(w, partnerDetailResponse(item))
}

func (h *OpsHandler) UpdatePartner(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	var req updatePartnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid json")
		return
	}
	partner, err := h.q.UpdateAdminPartner(r.Context(), sqlc.UpdateAdminPartnerParams{
		ID:                   id,
		Status:               optionalText(req.Status),
		CommissionRate:       optionalNumeric(req.CommissionRate),
		PromoCommissionRate:  optionalNumeric(req.PromoCommissionRate),
		PromoCommissionUntil: optionalDate(req.PromoCommissionUntil),
	})
	if err != nil {
		h.writeError(w, "failed to update partner", err)
		return
	}
	h.auditMutation(r, "admin.partner.update", "partner", id, nil)
	response.Success(w, map[string]any{"id": partner.ID, "status": partner.Status})
}

func (h *OpsHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	params := sqlc.ListAdminLocationsParams{StatusFilter: r.URL.Query().Get("status"), PartnerIDFilter: optionalUUIDString(r.URL.Query().Get("partner_id")), Search: r.URL.Query().Get("search"), PageLimit: int32(p.Limit), PageOffset: int32(p.Offset)}
	items, err := h.q.ListAdminLocations(r.Context(), params)
	if err != nil {
		h.writeError(w, "failed to list locations", err)
		return
	}
	total, err := h.q.CountAdminLocations(r.Context(), sqlc.CountAdminLocationsParams{StatusFilter: params.StatusFilter, PartnerIDFilter: params.PartnerIDFilter, Search: params.Search})
	if err != nil {
		h.writeError(w, "failed to count locations", err)
		return
	}
	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = locationListResponse(item)
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := h.q.GetAdminLocationByID(r.Context(), id)
	if err != nil {
		h.writeError(w, "failed to get location", err)
		return
	}
	response.Success(w, locationDetailResponse(item))
}

func (h *OpsHandler) UpdateLocationStatus(w http.ResponseWriter, r *http.Request) {
	h.updateStatus(w, r, "admin.location.status.update", "location", "location", func(id uuid.UUID, status string) error {
		_, err := h.q.UpdateAdminLocationStatus(r.Context(), sqlc.UpdateAdminLocationStatusParams{ID: id, Status: status})
		return err
	})
}

func (h *OpsHandler) ListBoxes(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	params := sqlc.ListAdminBoxesParams{StatusFilter: r.URL.Query().Get("status"), LocationIDFilter: optionalUUIDString(r.URL.Query().Get("location_id")), Search: r.URL.Query().Get("search"), PageLimit: int32(p.Limit), PageOffset: int32(p.Offset)}
	items, err := h.q.ListAdminBoxes(r.Context(), params)
	if err != nil {
		h.writeError(w, "failed to list boxes", err)
		return
	}
	total, err := h.q.CountAdminBoxes(r.Context(), sqlc.CountAdminBoxesParams{StatusFilter: params.StatusFilter, LocationIDFilter: params.LocationIDFilter, Search: params.Search})
	if err != nil {
		h.writeError(w, "failed to count boxes", err)
		return
	}
	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = boxListResponse(item)
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) GetBox(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := h.q.GetAdminBoxByID(r.Context(), id)
	if err != nil {
		h.writeError(w, "failed to get box", err)
		return
	}
	response.Success(w, boxDetailResponse(item))
}

func (h *OpsHandler) UpdateBoxStatus(w http.ResponseWriter, r *http.Request) {
	h.updateStatus(w, r, "admin.box.status.update", "surprise_box", "box", func(id uuid.UUID, status string) error {
		_, err := h.q.UpdateAdminBoxStatus(r.Context(), sqlc.UpdateAdminBoxStatusParams{ID: id, Status: status})
		return err
	})
}

func (h *OpsHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	search := r.URL.Query().Get("search")
	items, err := h.q.ListAdminCustomers(r.Context(), sqlc.ListAdminCustomersParams{Search: search, PageLimit: int32(p.Limit), PageOffset: int32(p.Offset)})
	if err != nil {
		h.writeError(w, "failed to list customers", err)
		return
	}
	total, err := h.q.CountAdminCustomers(r.Context(), search)
	if err != nil {
		h.writeError(w, "failed to count customers", err)
		return
	}
	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = customerListResponse(item)
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := h.q.GetAdminCustomerByID(r.Context(), id)
	if err != nil {
		h.writeError(w, "failed to get customer", err)
		return
	}
	response.Success(w, customerDetailResponse(item))
}

func (h *OpsHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	params := sqlc.ListAdminOrdersParams{StatusFilter: r.URL.Query().Get("status"), Search: r.URL.Query().Get("search"), PageLimit: int32(p.Limit), PageOffset: int32(p.Offset)}
	items, err := h.q.ListAdminOrders(r.Context(), params)
	if err != nil {
		h.writeError(w, "failed to list orders", err)
		return
	}
	total, err := h.q.CountAdminOrders(r.Context(), sqlc.CountAdminOrdersParams{StatusFilter: params.StatusFilter, Search: params.Search})
	if err != nil {
		h.writeError(w, "failed to count orders", err)
		return
	}
	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = orderListResponse(item)
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := h.q.GetAdminOrderByID(r.Context(), id)
	if err != nil {
		h.writeError(w, "failed to get order", err)
		return
	}
	response.Success(w, orderDetailResponse(item))
}

func (h *OpsHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	p := paginationFromRequest(r)
	params := sqlc.ListAdminPaymentsParams{StatusFilter: r.URL.Query().Get("status"), Search: r.URL.Query().Get("search"), PageLimit: int32(p.Limit), PageOffset: int32(p.Offset)}
	items, err := h.q.ListAdminPayments(r.Context(), params)
	if err != nil {
		h.writeError(w, "failed to list payments", err)
		return
	}
	total, err := h.q.CountAdminPayments(r.Context(), sqlc.CountAdminPaymentsParams{StatusFilter: params.StatusFilter, Search: params.Search})
	if err != nil {
		h.writeError(w, "failed to count payments", err)
		return
	}
	mapped := make([]map[string]any, len(items))
	for i, item := range items {
		mapped[i] = paymentListResponse(item)
	}
	response.Success(w, paginated(mapped, int(total), p))
}

func (h *OpsHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	item, err := h.q.GetAdminPaymentByID(r.Context(), id)
	if err != nil {
		h.writeError(w, "failed to get payment", err)
		return
	}
	response.Success(w, paymentDetailResponse(item))
}

func (h *OpsHandler) ListPaymentEvents(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	p := paginationFromRequest(r)
	items, err := h.q.ListAdminPaymentEvents(r.Context(), sqlc.ListAdminPaymentEventsParams{PaymentID: id, Limit: int32(p.Limit), Offset: int32(p.Offset)})
	if err != nil {
		h.writeError(w, "failed to list payment events", err)
		return
	}
	total, err := h.q.CountAdminPaymentEvents(r.Context(), id)
	if err != nil {
		h.writeError(w, "failed to count payment events", err)
		return
	}
	response.Success(w, paginated(items, int(total), p))
}

func (h *OpsHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.q.GetAdminStats(r.Context())
	if err != nil {
		h.writeError(w, "failed to get stats", err)
		return
	}
	response.Success(w, map[string]any{
		"customers": map[string]any{"total": stats.CustomersTotal},
		"partners":  map[string]any{"total": stats.PartnersTotal, "active": stats.PartnersActive},
		"locations": map[string]any{"total": stats.LocationsTotal},
		"boxes":     map[string]any{"total": stats.BoxesTotal},
		"orders":    map[string]any{"total": stats.OrdersTotal, "completed": stats.OrdersCompleted, "cancelled": stats.OrdersCancelled, "disputed": stats.OrdersDisputed},
		"revenue":   map[string]any{"gross": numeric(stats.GrossRevenue)},
		"payments":  map[string]any{"total": stats.PaymentsTotal, "succeeded": stats.PaymentsSucceeded},
	})
}

func (h *OpsHandler) updateStatus(w http.ResponseWriter, r *http.Request, action, entityType, paramName string, update func(uuid.UUID, string) error) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	var req updateStatusRequest
	if errs := h.validator.DecodeAndValidate(r, &req); errs != nil {
		response.ValidationError(w, "validation failed", errs)
		return
	}
	if err := update(id, req.Status); err != nil {
		h.writeError(w, "failed to update status", err)
		return
	}
	h.auditMutation(r, action, entityType, id, map[string]any{"status": req.Status})
	response.Success(w, map[string]any{paramName + "_id": id, "status": req.Status})
}

func (h *OpsHandler) auditMutation(r *http.Request, action, entityType string, entityID uuid.UUID, details map[string]any) {
	adminID, err := contextx.AdminID(r)
	if err != nil {
		h.log.Warn("failed to audit admin mutation: missing admin id", sl.Err(err))
		return
	}
	payload, err := json.Marshal(details)
	if err != nil {
		h.log.Warn("failed to marshal audit details", sl.Err(err))
		return
	}
	if _, err := h.auditRepo.CreateAuditLog(r.Context(), domain.AuditLog{
		AdminID: adminID, Action: action, EntityType: entityType, EntityID: &entityID, Details: payload, IPAddress: requestIP(r),
	}); err != nil {
		h.log.Warn("failed to write audit log", sl.Err(err))
	}
}

func (h *OpsHandler) writeError(w http.ResponseWriter, message string, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		response.NotFound(w, "not found")
		return
	}
	h.log.Error(message, sl.Err(err))
	response.InternalError(w, nil)
}

type page struct {
	Limit  int
	Offset int
}

func paginationFromRequest(r *http.Request) page {
	limit := intQuery(r, "limit", defaultAdminLimit)
	if limit <= 0 {
		limit = defaultAdminLimit
	}
	if limit > maxAdminLimit {
		limit = maxAdminLimit
	}
	offset := intQuery(r, "offset", 0)
	if offset < 0 {
		offset = 0
	}
	return page{Limit: limit, Offset: offset}
}

func intQuery(r *http.Request, name string, fallback int) int {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func paginated(items any, total int, p page) map[string]any {
	return map[string]any{
		"items": items,
		"pagination": map[string]any{
			"total": total, "limit": p.Limit, "offset": p.Offset, "has_more": p.Offset+p.Limit < total,
		},
	}
}

func parseIDParam(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		response.BadRequest(w, "invalid id")
		return uuid.Nil, false
	}
	return id, true
}

func adminResponse(admin sqlc.AdminUser) map[string]any {
	return map[string]any{"id": admin.ID, "email": admin.Email, "name": admin.Name, "role": admin.Role, "is_active": admin.IsActive, "last_login_at": timestamptz(admin.LastLoginAt), "created_at": admin.CreatedAt, "updated_at": admin.UpdatedAt}
}

func partnerListResponse(row sqlc.ListAdminPartnersRow) map[string]any {
	return map[string]any{"id": row.ID, "brand_name": row.BrandName, "legal_name": row.LegalName, "logo_url": text(row.LogoUrl), "account_type": text(row.AccountType), "commission_rate": numeric(row.CommissionRate), "promo_commission_rate": nullableNumeric(row.PromoCommissionRate), "promo_commission_until": date(row.PromoCommissionUntil), "status": row.Status, "locations_count": row.LocationsCount, "total_orders": row.TotalOrders, "total_revenue": numeric(row.TotalRevenue), "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func partnerDetailResponse(row sqlc.GetAdminPartnerByIDRow) map[string]any {
	return map[string]any{"id": row.ID, "brand_name": row.BrandName, "legal_name": row.LegalName, "logo_url": text(row.LogoUrl), "account_type": text(row.AccountType), "commission_rate": numeric(row.CommissionRate), "promo_commission_rate": nullableNumeric(row.PromoCommissionRate), "promo_commission_until": date(row.PromoCommissionUntil), "status": row.Status, "legal_info": map[string]any{"inn": row.Inn, "ogrn": row.Ogrn, "kpp": row.Kpp, "legal_address": row.LegalAddress}, "stats": map[string]any{"locations_count": row.LocationsCount, "total_orders": row.TotalOrders, "total_revenue": numeric(row.TotalRevenue)}, "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func locationListResponse(row sqlc.ListAdminLocationsRow) map[string]any {
	return map[string]any{"id": row.ID, "partner_id": row.PartnerID, "partner_name": row.PartnerName, "category_code": row.CategoryCode, "name": row.Name, "address": row.Address, "phone": row.Phone, "logo_url": row.LogoUrl, "cover_image_url": row.CoverImageUrl, "status": row.Status, "boxes_count": row.BoxesCount, "total_orders": row.TotalOrders, "created_at": timestamp(row.CreatedAt), "updated_at": timestamp(row.UpdatedAt)}
}

func locationDetailResponse(row sqlc.GetAdminLocationByIDRow) map[string]any {
	item := locationListResponse(sqlc.ListAdminLocationsRow{ID: row.ID, PartnerID: row.PartnerID, PartnerName: row.PartnerName, CategoryCode: row.CategoryCode, Name: row.Name, Address: row.Address, Phone: row.Phone, LogoUrl: row.LogoUrl, CoverImageUrl: row.CoverImageUrl, Status: row.Status, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt, BoxesCount: row.BoxesCount, TotalOrders: row.TotalOrders})
	item["gallery_urls"] = row.GalleryUrls
	item["working_hours"] = json.RawMessage(row.WorkingHours)
	return item
}

func boxListResponse(row sqlc.ListAdminBoxesRow) map[string]any {
	return map[string]any{"id": row.ID, "location_id": row.LocationID, "location_name": row.LocationName, "partner_id": row.PartnerID, "partner_name": row.PartnerName, "name": row.Name, "description": row.Description, "original_price": nullableNumeric(row.OriginalPrice), "discount_price": numeric(row.DiscountPrice), "quantity_available": row.QuantityAvailable, "pickup_time_start": pgTime(row.PickupTimeStart), "pickup_time_end": pgTime(row.PickupTimeEnd), "image_url": row.ImageUrl, "status": row.Status, "total_orders": row.TotalOrders, "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func boxDetailResponse(row sqlc.GetAdminBoxByIDRow) map[string]any {
	return boxListResponse(sqlc.ListAdminBoxesRow(row))
}

func customerListResponse(row sqlc.ListAdminCustomersRow) map[string]any {
	return map[string]any{"id": row.ID, "phone": row.Phone, "name": row.Name, "total_orders": row.TotalOrders, "total_spent": numeric(row.TotalSpent), "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func customerDetailResponse(row sqlc.GetAdminCustomerByIDRow) map[string]any {
	return map[string]any{"id": row.ID, "phone": row.Phone, "name": row.Name, "total_orders": row.TotalOrders, "total_spent": numeric(row.TotalSpent), "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func orderListResponse(row sqlc.ListAdminOrdersRow) map[string]any {
	return map[string]any{"id": row.ID, "user_id": row.UserID, "customer_phone": row.CustomerPhone, "customer_name": row.CustomerName, "box_id": row.BoxID, "box_name": row.BoxName, "location_id": row.LocationID, "location_name": row.LocationName, "partner_id": row.PartnerID, "partner_name": row.PartnerName, "pickup_code": row.PickupCode, "amount": numeric(row.Amount), "pickup_time_start": row.PickupTimeStart, "pickup_time_end": row.PickupTimeEnd, "status": row.Status, "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func orderDetailResponse(row sqlc.GetAdminOrderByIDRow) map[string]any {
	return map[string]any{"id": row.ID, "user_id": row.UserID, "customer_phone": row.CustomerPhone, "customer_name": row.CustomerName, "box_id": row.BoxID, "box_name": row.BoxName, "location_id": row.LocationID, "location_name": row.LocationName, "location_address": row.LocationAddress, "partner_id": row.PartnerID, "partner_name": row.PartnerName, "pickup_code": row.PickupCode, "qr_code_url": text(row.QrCodeUrl), "amount": numeric(row.Amount), "pickup_time_start": row.PickupTimeStart, "pickup_time_end": row.PickupTimeEnd, "status": row.Status, "partner_confirmation_deadline": row.PartnerConfirmationDeadline, "partner_confirmed_at": timestamptz(row.PartnerConfirmedAt), "partner_confirmed_by": nullableUUID(row.PartnerConfirmedBy), "cancellation_reason": text(row.CancellationReason), "cancelled_at": timestamptz(row.CancelledAt), "picked_up_at": timestamptz(row.PickedUpAt), "picked_up_confirmed_by": nullableUUID(row.PickedUpConfirmedBy), "user_confirmed_at": timestamptz(row.UserConfirmedAt), "auto_completed_at": timestamptz(row.AutoCompletedAt), "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func paymentListResponse(row sqlc.ListAdminPaymentsRow) map[string]any {
	return map[string]any{"id": row.ID, "order_id": row.OrderID, "provider_payment_id": text(row.ProviderPaymentID), "payment_url": text(row.PaymentUrl), "method": nullablePaymentMethod(row.Method), "provider": nullablePaymentProvider(row.Provider), "amount": numeric(row.Amount), "status": row.Status, "paid_at": timestamptz(row.PaidAt), "pickup_code": row.PickupCode, "customer_phone": row.CustomerPhone, "location_name": row.LocationName, "partner_name": row.PartnerName, "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func paymentDetailResponse(row sqlc.GetAdminPaymentByIDRow) map[string]any {
	return map[string]any{"id": row.ID, "order_id": row.OrderID, "provider_payment_id": text(row.ProviderPaymentID), "payment_url": text(row.PaymentUrl), "method": nullablePaymentMethod(row.Method), "provider": nullablePaymentProvider(row.Provider), "amount": numeric(row.Amount), "status": row.Status, "paid_at": timestamptz(row.PaidAt), "pickup_code": row.PickupCode, "customer_phone": row.CustomerPhone, "location_name": row.LocationName, "partner_name": row.PartnerName, "created_at": row.CreatedAt, "updated_at": row.UpdatedAt}
}

func text(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func numeric(value pgtype.Numeric) float64 {
	return pgconverter.NumericToDecimalOrZero(value).InexactFloat64()
}

func nullableNumeric(value pgtype.Numeric) any {
	if !value.Valid {
		return nil
	}
	return numeric(value)
}

func date(value pgtype.Date) any {
	if !value.Valid {
		return nil
	}
	return value.Time.Format(time.DateOnly)
}

func timestamp(value pgtype.Timestamp) any {
	if !value.Valid {
		return nil
	}
	return value.Time
}

func timestamptz(value pgtype.Timestamptz) any {
	if !value.Valid {
		return nil
	}
	return value.Time
}

func nullableUUID(value pgtype.UUID) any {
	if !value.Valid {
		return nil
	}
	return uuid.UUID(value.Bytes)
}

func pgTime(value pgtype.Time) string {
	if !value.Valid {
		return ""
	}
	return time.UnixMicro(value.Microseconds).UTC().Format("15:04:05")
}

func nullablePaymentMethod(value sqlc.NullPaymentMethod) any {
	if !value.Valid {
		return nil
	}
	return value.PaymentMethod
}

func nullablePaymentProvider(value sqlc.NullPaymentProvider) any {
	if !value.Valid {
		return nil
	}
	return value.PaymentProvider
}

func optionalText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func optionalBool(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *value, Valid: true}
}

func optionalNumeric(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{}
	}
	return pgconverter.DecimalToNumeric(decimal.NewFromFloat(*value), true)
}

func optionalDate(value *string) pgtype.Date {
	if value == nil || *value == "" {
		return pgtype.Date{}
	}
	parsed, err := time.Parse(time.DateOnly, *value)
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: parsed, Valid: true}
}

func optionalUUIDString(value string) string {
	if value == "" {
		return ""
	}
	if _, err := uuid.Parse(value); err != nil {
		return ""
	}
	return value
}
