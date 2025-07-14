package handler

import (
	"app/pkg"
	"app/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBuyerService is a mock implementation of ServiceBuyer for testing
type MockBuyerService struct {
	mock.Mock
}

func (m *MockBuyerService) GetAll() (b map[int]models.Buyer, err error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Buyer), args.Error(1)
}

func (m *MockBuyerService) GetByID(id int) (b models.Buyer, err error) {
	args := m.Called(id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) Create(buyer models.Buyer) (b models.Buyer, err error) {
	args := m.Called(buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) Update(id int, buyer models.Buyer) (b models.Buyer, err error) {
	args := m.Called(id, buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) Delete(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

// setupBuyerTestHandler creates a handler with mock service
func setupBuyerTestHandler() (*BuyerHandler, *MockBuyerService) {
	mockService := &MockBuyerService{}
	handler := NewBuyerHandler(mockService)
	return handler, mockService
}

func TestBuyerGetAll(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	expectedBuyers := map[int]models.Buyer{
		1: {
			ID: 1,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "12345678",
				FirstName:    "John",
				LastName:     "Doe",
			},
		},
		2: {
			ID: 2,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		},
	}

	mockService.On("GetAll").Return(expectedBuyers, nil)

	// Create request
	req := httptest.NewRequest("GET", "/buyers", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetAll()(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")

	buyersData := response["data"].([]interface{})
	assert.Len(t, buyersData, 2)

	mockService.AssertExpectations(t)
}

func TestBuyerGetByID(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	expectedBuyer := models.Buyer{
		ID: 1,
		BuyerAttributes: models.BuyerAttributes{
			CardNumberID: "12345678",
			FirstName:    "John",
			LastName:     "Doe",
		},
	}

	mockService.On("GetByID", 1).Return(expectedBuyer, nil)

	// Create request with ID parameter
	req := httptest.NewRequest("GET", "/buyers/1", nil)
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.GetByID()(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")

	buyerData := response["data"].(map[string]interface{})
	assert.Equal(t, "John", buyerData["first_name"])

	mockService.AssertExpectations(t)
}

func TestBuyerGetByIDNotFound(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	mockService.On("GetByID", 999).Return(models.Buyer{}, pkg.ServiceErrors[pkg.ErrNotFound])

	// Create request with non-existent ID
	req := httptest.NewRequest("GET", "/buyers/999", nil)
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.GetByID()(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestBuyerGetByIDInvalidID(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create request with invalid ID
	req := httptest.NewRequest("GET", "/buyers/invalid", nil)
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.GetByID()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyerCreate(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	inputBuyer := models.Buyer{
		BuyerAttributes: models.BuyerAttributes{
			CardNumberID: "99999999",
			FirstName:    "Alice",
			LastName:     "Johnson",
		},
	}

	expectedBuyer := models.Buyer{
		ID: 3,
		BuyerAttributes: models.BuyerAttributes{
			CardNumberID: "99999999",
			FirstName:    "Alice",
			LastName:     "Johnson",
		},
	}

	mockService.On("Create", inputBuyer).Return(expectedBuyer, nil)

	// Create request body
	buyerData := map[string]interface{}{
		"card_number_id": "99999999",
		"first_name":     "Alice",
		"last_name":      "Johnson",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")

	buyerResponse := response["data"].(map[string]interface{})
	assert.Equal(t, "Alice", buyerResponse["first_name"])

	mockService.AssertExpectations(t)
}

func TestBuyerCreateMissingCardNumberID(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create request body with missing card_number_id
	buyerData := map[string]interface{}{
		"first_name": "Alice",
		"last_name":  "Johnson",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyerCreateMissingFirstName(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create request body with missing first_name
	buyerData := map[string]interface{}{
		"card_number_id": "99999999",
		"last_name":      "Johnson",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyerCreateMissingLastName(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create request body with missing last_name
	buyerData := map[string]interface{}{
		"card_number_id": "99999999",
		"first_name":     "Alice",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyerCreateInvalidJSON(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create invalid JSON
	body := bytes.NewBuffer([]byte("invalid json"))

	// Create request
	req := httptest.NewRequest("POST", "/buyers", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyerCreateDuplicateCardNumberID(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	inputBuyer := models.Buyer{
		BuyerAttributes: models.BuyerAttributes{
			CardNumberID: "12345678",
			FirstName:    "Alice",
			LastName:     "Johnson",
		},
	}

	mockService.On("Create", inputBuyer).Return(models.Buyer{}, pkg.ServiceErrors[pkg.ErrBadRequest])

	// Create request body with duplicate card_number_id
	buyerData := map[string]interface{}{
		"card_number_id": "12345678",
		"first_name":     "Alice",
		"last_name":      "Johnson",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("POST", "/buyers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestBuyerUpdate(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	inputBuyer := models.Buyer{
		BuyerAttributes: models.BuyerAttributes{
			FirstName: "UpdatedJohn",
			LastName:  "UpdatedDoe",
		},
	}

	expectedBuyer := models.Buyer{
		ID: 1,
		BuyerAttributes: models.BuyerAttributes{
			FirstName: "UpdatedJohn",
			LastName:  "UpdatedDoe",
		},
	}

	mockService.On("Update", 1, inputBuyer).Return(expectedBuyer, nil)

	// Create request body
	buyerData := map[string]interface{}{
		"first_name": "UpdatedJohn",
		"last_name":  "UpdatedDoe",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("PATCH", "/buyers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Update()(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")

	mockService.AssertExpectations(t)
}

func TestBuyerUpdateNotFound(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	inputBuyer := models.Buyer{
		BuyerAttributes: models.BuyerAttributes{
			FirstName: "UpdatedJohn",
		},
	}

	mockService.On("Update", 999, inputBuyer).Return(models.Buyer{}, pkg.ServiceErrors[pkg.ErrNotFound])

	// Create request body
	buyerData := map[string]interface{}{
		"first_name": "UpdatedJohn",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("PATCH", "/buyers/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Update()(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestBuyerUpdateInvalidID(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create request body
	buyerData := map[string]interface{}{
		"first_name": "UpdatedJohn",
	}
	body, _ := json.Marshal(buyerData)

	// Create request
	req := httptest.NewRequest("PATCH", "/buyers/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Update()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyerUpdateInvalidJSON(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create invalid JSON
	body := bytes.NewBuffer([]byte("invalid json"))

	// Create request
	req := httptest.NewRequest("PATCH", "/buyers/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Update()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBuyerDelete(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	mockService.On("Delete", 1).Return(nil)

	// Create request
	req := httptest.NewRequest("DELETE", "/buyers/1", nil)
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Delete()(w, req)

	// Assertions
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestBuyerDeleteNotFound(t *testing.T) {
	handler, mockService := setupBuyerTestHandler()

	mockService.On("Delete", 999).Return(pkg.ServiceErrors[pkg.ErrNotFound])

	// Create request
	req := httptest.NewRequest("DELETE", "/buyers/999", nil)
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Delete()(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestBuyerDeleteInvalidID(t *testing.T) {
	handler, _ := setupBuyerTestHandler()

	// Create request
	req := httptest.NewRequest("DELETE", "/buyers/invalid", nil)
	w := httptest.NewRecorder()

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call handler
	handler.Delete()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
