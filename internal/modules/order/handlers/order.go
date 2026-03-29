package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
	"github.com/nlsnnn/berezhok/internal/lib/validator"
	catalogErrors "github.com/nlsnnn/berezhok/internal/modules/catalog/errors"
	"github.com/nlsnnn/berezhok/internal/modules/order/domain"
	orderErrors "github.com/nlsnnn/berezhok/internal/modules/order/errors"
	"github.com/nlsnnn/berezhok/internal/modules/order/handlers/dto"
	orderService "github.com/nlsnnn/berezhok/internal/modules/order/service"
	sharedErrors "github.com/nlsnnn/berezhok/internal/shared/errors"
	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type orderServiceInterface interface {
	CreateOrder(ctx context.Context, boxID, customerID uuid.UUID) (*orderService.CreateOrderResult, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
	GetOrderDetailsByID(ctx context.Context, orderID uuid.UUID) (*orderService.OrderDetailsResult, error)
	ListOrdersByCustomerID(ctx context.Context, customerID uuid.UUID, status string, limit, offset int) (*orderService.ListOrdersResult, error)
}

type orderHandler struct {
	log       *slog.Logger
	validator *validator.Validator
	service   orderServiceInterface
}

func NewOrderHandler(service orderServiceInterface, log *slog.Logger, v *validator.Validator) *orderHandler {
	return &orderHandler{
		service:   service,
		log:       log,
		validator: v,
	}
}

// CreateOrder handles POST /customer/orders
func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const op = "order.handler.CreateOrder"
	log := h.log.With(slog.String("op", op))

	// Parse and validate request
	var req dto.CreateOrderRequest
	if validationErrs := h.validator.DecodeAndValidate(r, &req); validationErrs != nil {
		log.Warn("invalid request", sl.Errs(validationErrs))
		response.ValidationError(w, "validation failed", validationErrs)
		return
	}

	// Extract customer ID from JWT context
	customerID, err := h.getCustomerIDFromContext(r)
	if err != nil {
		log.Error("failed to get customer_id from context", sl.Err(err))
		response.Unauthorized(w, "authentication required")
		return
	}

	// Parse box ID
	boxID, err := uuid.Parse(req.BoxID)
	if err != nil {
		log.Warn("invalid box_id format", slog.String("box_id", req.BoxID))
		response.BadRequest(w, "invalid box_id format")
		return
	}

	// Create order
	result, err := h.service.CreateOrder(r.Context(), boxID, customerID)
	if err != nil {
		switch {
		case errors.Is(err, orderErrors.ErrBoxNotAvailable):
			response.NotFound(w, "box is not available or out of stock")
		case errors.Is(err, orderErrors.ErrInvalidBoxStatus):
			response.BadRequest(w, "box is not active")
		case errors.Is(err, catalogErrors.ErrBoxNotFound):
			response.NotFound(w, "box not found")
		case errors.Is(err, orderErrors.ErrPaymentFailed):
			log.Error("payment creation failed", sl.Err(err))
			response.InternalErrorWithMessage(w, "failed to create payment")
		default:
			log.Error("failed to create order", sl.Err(err))
			response.InternalError(w, nil)
		}
		return
	}

	// Build response
	// Payment link typically expires in 15 minutes
	expiresAt := time.Now().Add(15 * time.Minute)

	resp := dto.ToCreateOrderResponse(
		result.OrderID.String(),
		result.PaymentLink,
		0, // Amount will be fetched from order if needed
		expiresAt,
	)

	log.Info("order created successfully", slog.String("order_id", result.OrderID.String()))
	response.Created(w, resp)
}

// GetOrder handles GET /customer/orders/{order_id}
func (h *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	const op = "order.handler.GetOrder"
	log := h.log.With(slog.String("op", op))

	// Extract order ID from URL
	orderIDStr := chi.URLParam(r, "order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		log.Warn("invalid order_id format", slog.String("order_id", orderIDStr))
		response.BadRequest(w, "invalid order_id format")
		return
	}

	// Extract customer ID from JWT context
	customerID, err := h.getCustomerIDFromContext(r)
	if err != nil {
		log.Error("failed to get customer_id from context", sl.Err(err))
		response.Unauthorized(w, "authentication required")
		return
	}

	// Get order details
	order, err := h.service.GetOrderDetailsByID(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, orderErrors.ErrOrderNotFound) {
			response.NotFound(w, "order not found")
			return
		}
		log.Error("failed to get order", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	// Verify ownership
	if order.CustomerID != customerID {
		log.Warn("customer tried to access another customer's order",
			slog.String("customer_id", customerID.String()),
			slog.String("order_customer_id", order.CustomerID.String()),
		)
		response.Forbidden(w, "access denied")
		return
	}

	resp := dto.ToOrderDetailResponse(order)
	response.Success(w, resp)
}

// ListOrders handles GET /customer/orders
func (h *orderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	const op = "order.handler.ListOrders"
	log := h.log.With(slog.String("op", op))

	// Extract customer ID from JWT context
	customerID, err := h.getCustomerIDFromContext(r)
	if err != nil {
		log.Error("failed to get customer_id from context", sl.Err(err))
		response.Unauthorized(w, "authentication required")
		return
	}

	// Parse query params
	status := r.URL.Query().Get("status")

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// List orders
	result, err := h.service.ListOrdersByCustomerID(r.Context(), customerID, status, limit, offset)
	if err != nil {
		log.Error("failed to list orders", sl.Err(err))
		response.InternalError(w, nil)
		return
	}

	// Build response items
	items := make([]dto.OrderListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = dto.ToOrderListItem(
			item.ID.String(),
			string(item.Status),
			item.PickupCode,
			item.Amount,
			item.BoxName,
			item.LocationName,
			item.PickupTimeStart,
			item.CreatedAt,
			item.HasReview,
		)
	}

	resp := dto.ToOrderListResponse(items, result.Total, result.Limit, result.Offset)
	response.Success(w, resp)
}

// ConfirmPickup handles POST /customer/orders/{order_id}/confirm-pickup
func (h *orderHandler) ConfirmPickup(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Order pickup confirmation not implemented yet", http.StatusNotImplemented)
}

// CreateDispute handles POST /customer/orders/{order_id}/dispute
func (h *orderHandler) CreateDispute(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Dispute creation not implemented yet", http.StatusNotImplemented)
}

func (h *orderHandler) getCustomerIDFromContext(r *http.Request) (uuid.UUID, error) {
	customerIDRaw := r.Context().Value("customer_id")
	if customerIDRaw == nil {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}

	customerID, ok := customerIDRaw.(uuid.UUID)
	if !ok {
		return uuid.Nil, sharedErrors.ErrNotFoundContextValue
	}

	return customerID, nil
}
