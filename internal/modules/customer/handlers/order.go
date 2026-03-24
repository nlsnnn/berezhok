package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type orderHandler struct {
	log *slog.Logger
}

func NewOrderHandler(log *slog.Logger) *orderHandler {
	return &orderHandler{log: log}
}

// CreateOrder handles POST /customer/orders
func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Order module not implemented yet", http.StatusNotImplemented)
}

// GetOrder handles GET /customer/orders/{order_id}
func (h *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Order module not implemented yet", http.StatusNotImplemented)
}

// ListOrders handles GET /customer/orders
func (h *orderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Order module not implemented yet", http.StatusNotImplemented)
}

// ConfirmPickup handles POST /customer/orders/{order_id}/confirm-pickup
func (h *orderHandler) ConfirmPickup(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Order module not implemented yet", http.StatusNotImplemented)
}

// CreateDispute handles POST /customer/orders/{order_id}/dispute
func (h *orderHandler) CreateDispute(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Order module not implemented yet", http.StatusNotImplemented)
}
