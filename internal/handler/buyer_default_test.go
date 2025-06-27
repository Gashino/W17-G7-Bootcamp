package handler

import (
	"app/internal/repository"
	"app/internal/service"
	"app/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// setupTestHandler creates a handler with test data
func setupTestHandler() *BuyerDefault {
	// Create test data
	testBuyers := map[int]models.Buyer{
		1: {
			Id: 1,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "402323",
				FirstName:    "John",
				LastName:     "Doe",
			},
		},
		2: {
			Id: 2,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "402324",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		},
	}

	// Create repository with test data
	repo := repository.NewBuyerMap(testBuyers)

	// Create service
	svc := service.NewBuyerDefault(repo)

	// Create handler
	return NewBuyerDefault(svc)
}

func TestGetAll(t *testing.T) {
	handler := setupTestHandler()

	// Create request
	req := httptest.NewRequest("GET", "/buyers", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetAll()(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check response structure
	if _, exists := response["data"]; !exists {
		t.Error("Response should contain 'data' field")
	}

	// Check that we have buyers in the response
	data := response["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("Expected 2 buyers, got %d", len(data))
	}
}

func TestGetByID(t *testing.T) {
	handler := setupTestHandler()

	// Create request with ID parameter
	req := httptest.NewRequest("GET", "/buyers/1", nil)
	w := httptest.NewRecorder()

	// Set up chi context with URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.GetByID()(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check response structure
	if _, exists := response["data"]; !exists {
		t.Error("Response should contain 'data' field")
	}

	// Check buyer data
	buyerData := response["data"].(map[string]interface{})
	if buyerData["Id"].(float64) != 1 {
		t.Errorf("Expected buyer ID 1, got %v", buyerData["Id"])
	}
	if buyerData["first_name"] != "John" {
		t.Errorf("Expected first name 'John', got %v", buyerData["first_name"])
	}
}

func TestGetByIDNotFound(t *testing.T) {
	handler := setupTestHandler()

	// Create request with non-existent ID
	req := httptest.NewRequest("GET", "/buyers/999", nil)
	w := httptest.NewRecorder()

	// Set up chi context with URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.GetByID()(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestCreate(t *testing.T) {
	handler := setupTestHandler()

	// Create request body
	buyerData := map[string]string{
		"card_number_id": "402325",
		"first_name":     "Bob",
		"last_name":      "Johnson",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check response structure
	if _, exists := response["data"]; !exists {
		t.Error("Response should contain 'data' field")
	}

	// Check buyer data
	buyerDataResponse := response["data"].(map[string]interface{})
	if buyerDataResponse["first_name"] != "Bob" {
		t.Errorf("Expected first name 'Bob', got %v", buyerDataResponse["first_name"])
	}
}

func TestCreateMissingFields(t *testing.T) {
	handler := setupTestHandler()

	// Create request body with missing fields
	buyerData := map[string]string{
		"card_number_id": "402325",
		// missing first_name and last_name
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateDuplicateCardNumber(t *testing.T) {
	handler := setupTestHandler()

	// Create request body with duplicate card_number_id
	buyerData := map[string]string{
		"card_number_id": "402323", // This already exists
		"first_name":     "Alice",
		"last_name":      "Brown",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdate(t *testing.T) {
	handler := setupTestHandler()

	// Create request body for partial update
	updateData := map[string]string{
		"first_name": "Johnny",
	}
	body, _ := json.Marshal(updateData)

	// Create request
	req := httptest.NewRequest("PATCH", "/buyers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Set up chi context with URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Update()(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check response structure
	if _, exists := response["data"]; !exists {
		t.Error("Response should contain 'data' field")
	}

	// Check buyer data
	buyerData := response["data"].(map[string]interface{})
	if buyerData["first_name"] != "Johnny" {
		t.Errorf("Expected first name 'Johnny', got %v", buyerData["first_name"])
	}
	// Check that other fields remain unchanged
	if buyerData["last_name"] != "Doe" {
		t.Errorf("Expected last name 'Doe', got %v", buyerData["last_name"])
	}
}

func TestUpdateDuplicateCardNumber(t *testing.T) {
	handler := setupTestHandler()

	// Create request body with duplicate card_number_id
	updateData := map[string]string{
		"card_number_id": "402324", // This belongs to buyer 2
	}
	body, _ := json.Marshal(updateData)

	// Create request
	req := httptest.NewRequest("PATCH", "/buyers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Set up chi context with URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Update()(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDelete(t *testing.T) {
	handler := setupTestHandler()

	// Create request
	req := httptest.NewRequest("DELETE", "/buyers/1", nil)
	w := httptest.NewRecorder()

	// Set up chi context with URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Delete()(w, req)

	// Check status code
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestDeleteNotFound(t *testing.T) {
	handler := setupTestHandler()

	// Create request with non-existent ID
	req := httptest.NewRequest("DELETE", "/buyers/999", nil)
	w := httptest.NewRecorder()

	// Set up chi context with URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Delete()(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
