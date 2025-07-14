package handler

import (
	"app/pkg"
	"app/pkg/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPurchaseOrderService is a mock implementation of ServicePurchaseOrder for testing
type MockPurchaseOrderService struct {
	mock.Mock
}

func (m *MockPurchaseOrderService) Create(purchaseOrder models.PurchaseOrder) (po models.PurchaseOrder, err error) {
	args := m.Called(purchaseOrder)
	return args.Get(0).(models.PurchaseOrder), args.Error(1)
}

// setupPurchaseOrderTestHandler creates a handler with mock service
func setupPurchaseOrderTestHandler() (*PurchaseOrderHandler, *MockPurchaseOrderService) {
	mockService := &MockPurchaseOrderService{}
	handler := NewPurchaseOrderHandler(mockService)
	return handler, mockService
}

func TestPurchaseOrderCreate(t *testing.T) {
	handler, mockService := setupPurchaseOrderTestHandler()

	inputPurchaseOrder := models.PurchaseOrder{
		PurchaseOrderAttributes: models.PurchaseOrderAttributes{
			OrderNumber:     "order#2",
			OrderDate:       "2021-04-05",
			TrackingCode:    "xyz789",
			BuyerID:         1,
			ProductRecordID: 2,
		},
	}

	expectedPurchaseOrder := models.PurchaseOrder{
		ID: 2,
		PurchaseOrderAttributes: models.PurchaseOrderAttributes{
			OrderNumber:     "order#2",
			OrderDate:       "2021-04-05",
			TrackingCode:    "xyz789",
			BuyerID:         1,
			ProductRecordID: 2,
		},
	}

	mockService.On("Create", inputPurchaseOrder).Return(expectedPurchaseOrder, nil)

	// Create request body
	purchaseOrderData := map[string]interface{}{
		"order_number":      "order#2",
		"order_date":        "2021-04-05",
		"tracking_code":     "xyz789",
		"buyer_id":          1,
		"product_record_id": 2,
	}
	body, _ := json.Marshal(purchaseOrderData)

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBuffer(body))
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

	purchaseOrderResponse := response["data"].(map[string]interface{})
	assert.Equal(t, "order#2", purchaseOrderResponse["order_number"])
	assert.Equal(t, "xyz789", purchaseOrderResponse["tracking_code"])

	mockService.AssertExpectations(t)
}

func TestPurchaseOrderCreateMissingOrderNumber(t *testing.T) {
	handler, _ := setupPurchaseOrderTestHandler()

	// Create request body with missing order_number
	purchaseOrderData := map[string]interface{}{
		"order_date":        "2021-04-05",
		"tracking_code":     "xyz789",
		"buyer_id":          1,
		"product_record_id": 2,
	}
	body, _ := json.Marshal(purchaseOrderData)

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPurchaseOrderCreateMissingOrderDate(t *testing.T) {
	handler, _ := setupPurchaseOrderTestHandler()

	// Create request body with missing order_date
	purchaseOrderData := map[string]interface{}{
		"order_number":      "order#2",
		"tracking_code":     "xyz789",
		"buyer_id":          1,
		"product_record_id": 2,
	}
	body, _ := json.Marshal(purchaseOrderData)

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPurchaseOrderCreateMissingTrackingCode(t *testing.T) {
	handler, _ := setupPurchaseOrderTestHandler()

	// Create request body with missing tracking_code
	purchaseOrderData := map[string]interface{}{
		"order_number":      "order#2",
		"order_date":        "2021-04-05",
		"buyer_id":          1,
		"product_record_id": 2,
	}
	body, _ := json.Marshal(purchaseOrderData)

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPurchaseOrderCreateMissingBuyerID(t *testing.T) {
	handler, _ := setupPurchaseOrderTestHandler()

	// Create request body with missing buyer_id
	purchaseOrderData := map[string]interface{}{
		"order_number":      "order#2",
		"order_date":        "2021-04-05",
		"tracking_code":     "xyz789",
		"product_record_id": 2,
	}
	body, _ := json.Marshal(purchaseOrderData)

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPurchaseOrderCreateMissingProductRecordID(t *testing.T) {
	handler, _ := setupPurchaseOrderTestHandler()

	// Create request body with missing product_record_id
	purchaseOrderData := map[string]interface{}{
		"order_number":  "order#2",
		"order_date":    "2021-04-05",
		"tracking_code": "xyz789",
		"buyer_id":      1,
	}
	body, _ := json.Marshal(purchaseOrderData)

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPurchaseOrderCreateInvalidJSON(t *testing.T) {
	handler, _ := setupPurchaseOrderTestHandler()

	// Create invalid JSON
	body := bytes.NewBuffer([]byte("invalid json"))

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPurchaseOrderCreateDuplicateOrderNumber(t *testing.T) {
	handler, mockService := setupPurchaseOrderTestHandler()

	inputPurchaseOrder := models.PurchaseOrder{
		PurchaseOrderAttributes: models.PurchaseOrderAttributes{
			OrderNumber:     "order#1",
			OrderDate:       "2021-04-05",
			TrackingCode:    "xyz789",
			BuyerID:         1,
			ProductRecordID: 2,
		},
	}

	mockService.On("Create", inputPurchaseOrder).Return(models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrBadRequest])

	// Create request body with duplicate order_number
	purchaseOrderData := map[string]interface{}{
		"order_number":      "order#1",
		"order_date":        "2021-04-05",
		"tracking_code":     "xyz789",
		"buyer_id":          1,
		"product_record_id": 2,
	}
	body, _ := json.Marshal(purchaseOrderData)

	// Create request
	req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.Create()(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}
