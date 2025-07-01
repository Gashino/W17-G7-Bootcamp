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

// MockSellerService es un mock del servicio de sellers
type MockSellerService struct {
	mock.Mock
}

func (m *MockSellerService) FindAll() (map[int]models.Seller, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Seller), args.Error(1)
}

func (m *MockSellerService) GetById(id int) (models.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) Create(seller models.Seller) (models.Seller, error) {
	args := m.Called(seller)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) UpdateFields(id int, request models.SellerCreateRequest) (models.Seller, error) {
	args := m.Called(id, request)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) DeleteSeller(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test para GetById - Caso exitoso
func TestSellerDefault_GetById_Success(t *testing.T) {
	// Arrange
	mockService := new(MockSellerService)
	handler := NewSellerDefault(mockService)

	// Seller de prueba
	expectedSeller := models.Seller{
		Id: 1,
		SellerAttributes: models.SellerAttributes{
			CId:         "12345",
			CompanyName: "Test Company",
			Adress:      "Test Address",
			Telephone:   "123456789",
		},
	}

	// Configurar mock para retornar el seller
	mockService.On("GetById", 1).Return(expectedSeller, nil)

	// Crear request y response recorder
	req := httptest.NewRequest("GET", "/sellers/1", nil)
	w := httptest.NewRecorder()

	// Configurar chi router para simular parámetro de URL
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handlerFunc := handler.GetById()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "success", response["message"])
	assert.NotNil(t, response["data"])

	// Verificar que el mock fue llamado correctamente
	mockService.AssertExpectations(t)
}

// Test para GetById - Caso de error (seller no encontrado)
func TestSellerDefault_GetById_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockSellerService)
	handler := NewSellerDefault(mockService)

	// Configurar mock para retornar error
	mockService.On("GetById", 999).Return(models.Seller{}, fmt.Errorf("seller not found"))

	// Crear request y response recorder
	req := httptest.NewRequest("GET", "/sellers/999", nil)
	w := httptest.NewRecorder()

	// Configurar chi router para simular parámetro de URL
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handlerFunc := handler.GetById()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Not found", response["error"])

	// Verificar que el mock fue llamado correctamente
	mockService.AssertExpectations(t)
}

// Test para GetById - Caso de error (ID inválido)
func TestSellerDefault_GetById_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockSellerService)
	handler := NewSellerDefault(mockService)

	// Crear request y response recorder con ID inválido
	req := httptest.NewRequest("GET", "/sellers/invalid", nil)
	w := httptest.NewRecorder()

	// Configurar chi router para simular parámetro de URL
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handlerFunc := handler.GetById()
	handlerFunc(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Bad request", response["error"])

	// No se debe llamar al service con ID inválido
	mockService.AssertNotCalled(t, "GetById")
}

// Test para Create - Caso exitoso
func TestSellerDefault_Create_Success(t *testing.T) {
	// Arrange
	mockService := new(MockSellerService)
	handler := NewSellerDefault(mockService)

	// Request de creación válido
	createRequest := models.SellerCreateRequest{
		CId:         stringPtr("12345"),
		CompanyName: stringPtr("New Company"),
		Address:     stringPtr("New Address"),
		Telephone:   stringPtr("987654321"),
	}

	// Seller esperado que retornará el service
	expectedSeller := models.Seller{
		Id: 2,
		SellerAttributes: models.SellerAttributes{
			CId:         "12345",
			CompanyName: "New Company",
			Adress:      "New Address",
			Telephone:   "987654321",
		},
	}

	// Configurar mock
	mockService.On("Create", mock.AnythingOfType("models.Seller")).Return(expectedSeller, nil)

	// Crear request HTTP
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/sellers", bytes.NewBuffer(reqBody))
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

	assert.Equal(t, "success", response["message"])
	assert.NotNil(t, response["data"])

	// Verificar que el mock fue llamado correctamente
	mockService.AssertExpectations(t)
}

// Test para Create - Caso de error (datos inválidos)
func TestSellerDefault_Create_InvalidData(t *testing.T) {
	// Arrange
	mockService := new(MockSellerService)
	handler := NewSellerDefault(mockService)

	// Request de creación inválido (faltan campos obligatorios)
	createRequest := models.SellerCreateRequest{
		CId: stringPtr("12345"),
		// Faltan CompanyName, Address y Telephone
	}

	// Crear request HTTP
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/sellers", bytes.NewBuffer(reqBody))
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

	assert.Equal(t, "Validation error", response["error"])

	// No se debe llamar al service con datos inválidos
	mockService.AssertNotCalled(t, "Create")
}

// Test para Create - Caso de error (JSON malformado)
func TestSellerDefault_Create_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockSellerService)
	handler := NewSellerDefault(mockService)

	// JSON malformado
	invalidJSON := `{"cid": "12345", "company_name": }`

	// Crear request HTTP
	req := httptest.NewRequest("POST", "/sellers", bytes.NewBuffer([]byte(invalidJSON)))
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

	assert.Equal(t, "Validation error", response["error"])

	// No se debe llamar al service con JSON inválido
	mockService.AssertNotCalled(t, "Create")
}

// Helper function para crear punteros a string
func stringPtr(s string) *string {
	return &s
} 