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
			response.Error(w, http.StatusBadRequest, "card_number_id is required")
			return
		}
		if buyer.FirstName == "" {
			response.Error(w, http.StatusBadRequest, "first_name is required")
			return
		}
		if buyer.LastName == "" {
			response.Error(w, http.StatusBadRequest, "last_name is required")
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

// Delete is a method that handles DELETE requests to delete a buyer
func (h *BuyerHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		// Delete buyer
		err = h.sv.Delete(id)
		if err != nil {
			handleError(w, err)
			return
		}

		// Send response using the response library
		response.JSON(w, http.StatusNoContent, nil)
	}
}

// handleError is a helper function that handles errors and sends appropriate HTTP responses
func handleError(w http.ResponseWriter, err error) {
	if serviceError, ok := err.(pkg.ServiceError); ok {
		response.Error(w, serviceError.ResponseCode, serviceError.Message)
		return
	}

	// Default error response
	response.Error(w, http.StatusInternalServerError, "Internal server error")
}
