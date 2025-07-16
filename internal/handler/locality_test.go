package handler

import (
	"app/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocalityService is a mock of locality service
type MockLocalityService struct {
	mock.Mock
}

func (m *MockLocalityService) Create(locality models.Locality) (models.Locality, error) {
	args := m.Called(locality)
	return args.Get(0).(models.Locality), args.Error(1)
}

func (m *MockLocalityService) GetById(id int) (models.Locality, error) {
	args := m.Called(id)
	return args.Get(0).(models.Locality), args.Error(1)
}

func (m *MockLocalityService) GetCantSellersByLocality(id int) (models.LocalityBySellerResponse, error) {
	args := m.Called(id)
	return args.Get(0).(models.LocalityBySellerResponse), args.Error(1)
}

// Test for Create - Success case
func TestLocalityDefault_Create_Success(t *testing.T) {
	// Arrange
	mockService := new(MockLocalityService)
	handler := NewLocalityDefault(mockService)

	// Valid creation request
	createRequest := models.LocalityCreateRequest{
		LocalityName: stringPtr("Buenos Aires"),
		ProvinceName: stringPtr("Buenos Aires"),
		CountryName:  stringPtr("Argentina"),
	}

	// Expected locality returned by service
	expectedLocality := models.Locality{
		ID: 1,
		LocalitiesAttributes: models.LocalitiesAttributes{
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		},
	}

	// Configure mock
	mockService.On("Create", mock.AnythingOfType("models.Locality")).Return(expectedLocality, nil)

	// Create HTTP request
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/localities", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handlerFunc := handler.Create()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response["data"])

	// Verify mock was called correctly
	mockService.AssertExpectations(t)
}

// Test for Create - Invalid data case
func TestLocalityDefault_Create_InvalidData(t *testing.T) {
	// Arrange
	mockService := new(MockLocalityService)
	handler := NewLocalityDefault(mockService)

	// Invalid creation request (missing required fields)
	createRequest := models.LocalityCreateRequest{
		LocalityName: stringPtr("Buenos Aires"),
		// Missing ProvinceName and CountryName
	}

	// Create HTTP request
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/localities", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handlerFunc := handler.Create()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error: Validation error", response["message"])

	// Service should not be called with invalid data
	mockService.AssertNotCalled(t, "Create")
}

// Test for Create - JSON malformed case
func TestLocalityDefault_Create_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockLocalityService)
	handler := NewLocalityDefault(mockService)

	// Malformed JSON
	invalidJSON := `{"locality_name": "Buenos Aires", "province_name": }`

	// Create HTTP request
	req := httptest.NewRequest("POST", "/localities", bytes.NewBuffer([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handlerFunc := handler.Create()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error: Validation error", response["message"])

	// Service should not be called with invalid JSON
	mockService.AssertNotCalled(t, "Create")
}

// Test for Create - Conflict case
func TestLocalityDefault_Create_Conflict(t *testing.T) {
	// Arrange
	mockService := new(MockLocalityService)
	handler := NewLocalityDefault(mockService)

	// Valid creation request
	createRequest := models.LocalityCreateRequest{
		LocalityName: stringPtr("Buenos Aires"),
		ProvinceName: stringPtr("Buenos Aires"),
		CountryName:  stringPtr("Argentina"),
	}

	// Configure mock to return conflict error
	mockService.On("Create", mock.AnythingOfType("models.Locality")).Return(models.Locality{}, fmt.Errorf("locality already exists"))

	// Create HTTP request
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/localities", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handlerFunc := handler.Create()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error: Resource conflict", response["message"])

	// Verify mock was called correctly
	mockService.AssertExpectations(t)
}

// Test for SellersByLocality - Success case
func TestLocalityDefault_SellersByLocality_Success(t *testing.T) {
	// Arrange
	mockService := new(MockLocalityService)
	handler := NewLocalityDefault(mockService)

	// Expected response
	expectedResponse := models.LocalityBySellerResponse{
		ID:           1,
		LocalityName: stringPtr("Buenos Aires"),
		SellerCount:  stringPtr("5"),
	}

	// Configure mock to return response
	mockService.On("GetCantSellersByLocality", 1).Return(expectedResponse, nil)

	// Create request and response recorder
	req := httptest.NewRequest("GET", "/localities/reportSellers/1", nil)
	w := httptest.NewRecorder()

	// Configure chi router to simulate URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handlerFunc := handler.SellersByLocality()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response["data"])

	// Verify mock was called correctly
	mockService.AssertExpectations(t)
}

// Test for SellersByLocality - Not found case
func TestLocalityDefault_SellersByLocality_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockLocalityService)
	handler := NewLocalityDefault(mockService)

	// Configure mock to return error
	mockService.On("GetCantSellersByLocality", 999).Return(models.LocalityBySellerResponse{}, fmt.Errorf("locality not found"))

	// Create request and response recorder
	req := httptest.NewRequest("GET", "/localities/reportSellers/999", nil)
	w := httptest.NewRecorder()

	// Configure chi router to simulate URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handlerFunc := handler.SellersByLocality()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error: Not found", response["message"])

	// Verify mock was called correctly
	mockService.AssertExpectations(t)
}

// Test for SellersByLocality - Invalid ID case
func TestLocalityDefault_SellersByLocality_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockLocalityService)
	handler := NewLocalityDefault(mockService)

	// Create request and response recorder with invalid ID
	req := httptest.NewRequest("GET", "/localities/reportSellers/invalid", nil)
	w := httptest.NewRecorder()

	// Configure chi router to simulate URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handlerFunc := handler.SellersByLocality()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "error: Bad request", response["message"])

	// Service should not be called with invalid ID
	mockService.AssertNotCalled(t, "GetCantSellersByLocality")
}
