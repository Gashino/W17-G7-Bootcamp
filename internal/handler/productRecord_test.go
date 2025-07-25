package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/productRecord"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Helper functions for creating pointers
func timePtr(t time.Time) *time.Time { return &t }

func TestProductRecordHandler_Create(t *testing.T) {
	// Create a sample product record for testing
	now := time.Now()
	purchasePrice := 100.50
	salePrice := 150.75
	productId := 1

	productRecordMock := models.ProductRecord{
		LastUpdateDate: timePtr(now),
		PurchasePrice:  float64Ptr(purchasePrice),
		SalePrice:      float64Ptr(salePrice),
		ProductId:      models.IntPtr(productId),
	}

	t.Run("create_ok", func(t *testing.T) {
		// Arrange
		mockProductRecordService := new(productRecord.MockProductRecordService)
		hd := NewProductRecordDefault(mockProductRecordService)

		// Setup the mock to return the product record with an ID
		expectedRecord := productRecordMock
		expectedRecord.ID = 1

		// Use AnythingOfType to ignore time comparison issues
		mockProductRecordService.On("Create", mock.AnythingOfType("models.ProductRecord")).Return(&expectedRecord, nil).Run(func(args mock.Arguments) {
			actual := args.Get(0).(models.ProductRecord)
			// Verify important fields
			assert.Equal(t, *productRecordMock.PurchasePrice, *actual.PurchasePrice)
			assert.Equal(t, *productRecordMock.SalePrice, *actual.SalePrice)
			assert.Equal(t, *productRecordMock.ProductId, *actual.ProductId)
			// For time, use Equal() which ignores monotonic clock
			assert.True(t, actual.LastUpdateDate.Equal(*productRecordMock.LastUpdateDate))
		})

		// Create request body
		body, errParsing := json.Marshal(productRecordMock)
		require.NoError(t, errParsing)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/productRecords", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		mockProductRecordService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, res.Code)

		// Verify response body
		var response map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response, "data")
	})

	t.Run("create_invalid_request_body", func(t *testing.T) {
		// Arrange
		mockProductRecordService := new(productRecord.MockProductRecordService)
		hd := NewProductRecordDefault(mockProductRecordService)

		// Invalid JSON to trigger parse error
		invalidJSON := []byte(`{"invalid":json}`)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/productRecords", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		mockProductRecordService.AssertNotCalled(t, "Create")
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("create_invalid_product_record", func(t *testing.T) {
		// Arrange
		mockProductRecordService := new(productRecord.MockProductRecordService)
		hd := NewProductRecordDefault(mockProductRecordService)

		// Create an invalid product record (missing required fields)
		invalidRecord := models.ProductRecord{
			// Missing LastUpdateDate
			PurchasePrice: float64Ptr(100.0),
			SalePrice:     float64Ptr(150.0),
			ProductId:     models.IntPtr(1),
		}

		body, errParsing := json.Marshal(invalidRecord)
		require.NoError(t, errParsing)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/productRecords", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		mockProductRecordService.AssertNotCalled(t, "Create")
		require.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})

	t.Run("create_service_error", func(t *testing.T) {
		// Arrange
		mockProductRecordService := new(productRecord.MockProductRecordService)
		hd := NewProductRecordDefault(mockProductRecordService)

		// Setup the mock to return a service error
		errorResponse := pkg.ServiceErrors[pkg.ErrNotFound]
		errorResponse.Message = "invalid product_id"

		// Use AnythingOfType to ignore time comparison issues
		mockProductRecordService.On("Create", mock.AnythingOfType("models.ProductRecord")).Return(nil, errorResponse).Run(func(args mock.Arguments) {
			actual := args.Get(0).(models.ProductRecord)
			// Verify important fields manually
			assert.Equal(t, *productRecordMock.PurchasePrice, *actual.PurchasePrice)
			assert.Equal(t, *productRecordMock.SalePrice, *actual.SalePrice)
			assert.Equal(t, *productRecordMock.ProductId, *actual.ProductId)
			// For time, use Equal() which ignores monotonic clock
			assert.True(t, actual.LastUpdateDate.Equal(*productRecordMock.LastUpdateDate))
		})

		body, errParsing := json.Marshal(productRecordMock)
		require.NoError(t, errParsing)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/productRecords", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		mockProductRecordService.AssertExpectations(t)
		require.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("create_internal_server_error", func(t *testing.T) {
		// Arrange
		mockProductRecordService := new(productRecord.MockProductRecordService)
		hd := NewProductRecordDefault(mockProductRecordService)

		// Setup the mock to return a generic error
		// Use AnythingOfType to ignore time comparison issues
		mockProductRecordService.On("Create", mock.AnythingOfType("models.ProductRecord")).Return(nil, errors.New("database connection error")).Run(func(args mock.Arguments) {
			actual := args.Get(0).(models.ProductRecord)
			// Verify important fields manually
			assert.Equal(t, *productRecordMock.PurchasePrice, *actual.PurchasePrice)
			assert.Equal(t, *productRecordMock.SalePrice, *actual.SalePrice)
			assert.Equal(t, *productRecordMock.ProductId, *actual.ProductId)
			// For time, use Equal() which ignores monotonic clock
			assert.True(t, actual.LastUpdateDate.Equal(*productRecordMock.LastUpdateDate))
		})

		body, errParsing := json.Marshal(productRecordMock)
		require.NoError(t, errParsing)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/productRecords", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		mockProductRecordService.AssertExpectations(t)
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
