package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

// NewPurchaseOrderHandler is a function that returns a new instance of PurchaseOrderHandler
func NewPurchaseOrderHandler(sv service.ServicePurchaseOrder) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{
		sv: sv,
	}
}

// PurchaseOrderHandler is a struct that handles purchase order HTTP requests
type PurchaseOrderHandler struct {
	// sv is the service for purchase order business logic operations
	sv service.ServicePurchaseOrder
}

// Create is a method that handles POST requests to create a new purchase order
func (h *PurchaseOrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var purchaseOrder models.PurchaseOrder
		err := json.NewDecoder(r.Body).Decode(&purchaseOrder)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		// Validate required fields
		if purchaseOrder.OrderNumber == "" {
			response.Error(w, http.StatusBadRequest, "order_number is required")
			return
		}
		if purchaseOrder.OrderDate == "" {
			response.Error(w, http.StatusBadRequest, "order_date is required")
			return
		}
		if purchaseOrder.TrackingCode == "" {
			response.Error(w, http.StatusBadRequest, "tracking_code is required")
			return
		}
		if purchaseOrder.BuyerID == 0 {
			response.Error(w, http.StatusBadRequest, "buyer_id is required")
			return
		}
		if purchaseOrder.ProductRecordID == 0 {
			response.Error(w, http.StatusBadRequest, "product_record_id is required")
			return
		}

		// Create purchase order
		createdPurchaseOrder, err := h.sv.Create(purchaseOrder)
		if err != nil {
			if serviceError, ok := err.(pkg.ServiceError); ok {
				response.Error(w, serviceError.ResponseCode, serviceError.Message)
				return
			}
			response.Error(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Create response
		responseData := map[string]interface{}{
			"data": createdPurchaseOrder,
		}

		// Send response using the response library
		response.JSON(w, http.StatusCreated, responseData)
	}
}
