package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nlsnnn/berezhok/internal/shared/response"
)

type customerOrderHandler struct {
	log *slog.Logger
}

func NewCustomerOrderHandler(log *slog.Logger) *customerOrderHandler {
	return &customerOrderHandler{log: log}
}

// CreateOrder handles POST /customer/orders (stub - redirects to order module)
func (h *customerOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Use order module handlers", http.StatusNotImplemented)
}

// GetOrder handles GET /customer/orders/{order_id} (stub - redirects to order module)
func (h *customerOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Use order module handlers", http.StatusNotImplemented)
}

// ListOrders handles GET /customer/orders (stub - redirects to order module)
func (h *customerOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Use order module handlers", http.StatusNotImplemented)
}

// ConfirmPickup handles POST /customer/orders/{order_id}/confirm-pickup
func (h *customerOrderHandler) ConfirmPickup(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Order pickup confirmation not implemented yet", http.StatusNotImplemented)
}

// CreateDispute handles POST /customer/orders/{order_id}/dispute
func (h *customerOrderHandler) CreateDispute(w http.ResponseWriter, r *http.Request) {
	response.Error(w, "Dispute creation not implemented yet", http.StatusNotImplemented)
}
