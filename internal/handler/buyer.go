package handler

import (
	"app/internal/service"
	"app/pkg"
	"app/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

// NewBuyerHandler is a function that returns a new instance of BuyerHandler
func NewBuyerHandler(sv service.ServiceBuyer) *BuyerHandler {
	return &BuyerHandler{
		sv: sv,
	}
}

// BuyerHandler is a struct that handles buyer HTTP requests
type BuyerHandler struct {
	// sv is the service for buyer business logic operations
	sv service.ServiceBuyer
}

// GetAll is a method that handles GET requests to retrieve all buyers
func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all buyers
		buyers, err := h.sv.GetAll()
		if err != nil {
			handleError(w, err)
			return
		}

		// Convert map to slice for JSON response
		var buyersSlice []models.Buyer
		for _, buyer := range buyers {
			buyersSlice = append(buyersSlice, buyer)
		}

		// Create response
		responseData := map[string]interface{}{
			"data": buyersSlice,
		}

		// Send response using the response library
		response.JSON(w, http.StatusOK, responseData)
	}
}

// GetByID is a method that handles GET requests to retrieve a buyer by ID
func (h *BuyerHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		// Get buyer by ID
		buyer, err := h.sv.GetByID(id)
		if err != nil {
			handleError(w, err)
			return
		}

		// Create response
		responseData := map[string]interface{}{
			"data": buyer,
		}

		// Send response using the response library
		response.JSON(w, http.StatusOK, responseData)
	}
}

// Create is a method that handles POST requests to create a new buyer
func (h *BuyerHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var buyer models.Buyer
		err := json.NewDecoder(r.Body).Decode(&buyer)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		// Validate required fields
		if buyer.CardNumberID == "" {
			response.Error(w, http.StatusUnprocessableEntity, "error: card_number_id is required")
			return
		}
		if buyer.FirstName == "" {
			response.Error(w, http.StatusUnprocessableEntity, "error: first_name is required")
			return
		}
		if buyer.LastName == "" {
			response.Error(w, http.StatusUnprocessableEntity, "error: last_name is required")
			return
		}

		// Create buyer
		createdBuyer, err := h.sv.Create(buyer)
		if err != nil {
			handleError(w, err)
			return
		}

		// Create response
		responseData := map[string]interface{}{
			"data": createdBuyer,
		}

		// Send response using the response library
		response.JSON(w, http.StatusCreated, responseData)
	}
}

// Update is a method that handles PATCH requests to update an existing buyer
func (h *BuyerHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		// Parse request body
		var buyer models.Buyer
		err = json.NewDecoder(r.Body).Decode(&buyer)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		// Update buyer (partial updates are handled in repository)
		updatedBuyer, err := h.sv.Update(id, buyer)
		if err != nil {
			handleError(w, err)
			return
		}

		// Create response
		responseData := map[string]interface{}{
			"data": updatedBuyer,
		}

		// Send response using the response library
		response.JSON(w, http.StatusOK, responseData)
	}
}

// Delete is a method that handles DELETE requests to remove a buyer
func (h *BuyerHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the buyer ID from the path parameter
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]string{
				"status":  "Bad Request",
				"message": "Invalid ID format",
			})
			return
		}

		// Delete the buyer
		err = h.sv.Delete(id)
		if err != nil {
			handleError(w, err)
			return
		}

		// Send success response
		response.JSON(w, http.StatusNoContent, nil)
	}
}

// GetPurchaseOrdersReport is a method that handles GET requests to retrieve purchase orders report for buyers
func (h *BuyerHandler) GetPurchaseOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the optional buyer ID query parameter
		var buyerID *int
		idStr := r.URL.Query().Get("id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				response.JSON(w, http.StatusBadRequest, map[string]string{
					"status":  "Bad Request",
					"message": "Invalid ID format",
				})
				return
			}
			buyerID = &id
		}

		// Get purchase orders report
		reports, err := h.sv.GetPurchaseOrdersReport(buyerID)
		if err != nil {
			handleError(w, err)
			return
		}

		// Create response
		responseData := map[string]interface{}{
			"data": reports,
		}

		// Send response
		response.JSON(w, http.StatusOK, responseData)
	}
}

// handleError is a helper function that handles error responses
func handleError(w http.ResponseWriter, err error) {
	// Check if it's a ServiceError
	serviceError, ok := err.(pkg.ServiceError)
	if !ok {
		// If it's not a ServiceError, return internal server error
		response.JSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "Internal Server Error",
			"message": "Internal server error",
		})
		return
	}

	// Get the error details based on the service error code
	var statusText string
	switch serviceError.ResponseCode {
	case http.StatusBadRequest:
		statusText = "Bad Request"
	case http.StatusNotFound:
		statusText = "Not Found"
	case http.StatusConflict:
		statusText = "Conflict"
	case http.StatusUnprocessableEntity:
		statusText = "Unprocessable Entity"
	default:
		statusText = "Internal Server Error"
	}

	// Return appropriate error response based on the service error
	response.JSON(w, serviceError.ResponseCode, map[string]string{
		"status":  statusText,
		"message": serviceError.Message,
	})
}
