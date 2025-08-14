package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/product"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func float64Ptr(f float64) *float64 { return &f }
func TestHandler_Create(t *testing.T) {
	productCode := "TEST001"
	description := "Test Product"
	netWeight := 10.5
	expirationRate := 30
	recommendedFreezingTemperature := -18.0
	freezingRate := 10
	productTypeId := 101
	sellerId := 1
	width := 100.4
	height := 50.1
	length := 100.2

	productMock := models.Product{
		ProductAttributes: models.ProductAttributes{
			ProductCode:                    &productCode,
			Description:                    &description,
			NetWeight:                      &netWeight,
			ExpirationRate:                 &expirationRate,
			RecommendedFreezingTemperature: &recommendedFreezingTemperature,
			FreezingRate:                   &freezingRate,
			ProductTypeId:                  &productTypeId,
			SellerId:                       &sellerId,
		},
		Dimensions: models.Dimensions{
			Width:  &width,
			Height: &height,
			Length: &length,
		},
	}

	t.Run("create_ok", func(t *testing.T) {
		//arrange

		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		mockProductService.On("Create", productMock).Return(&productMock, nil)

		body, errParsing := json.Marshal(productMock)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		//act
		hd.Create()(res, req)

		//assert
		assert.NoError(t, errParsing)
		mockProductService.AssertExpectations(t)
		require.Equal(t, res.Code, http.StatusCreated)
		require.NotEmpty(t, res.Body)
		require.Equal(t, productMock.ID, 1)

	})
	t.Run("create_failure", func(t *testing.T) {
		//arrange
		productMock := models.Product{}

		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		body, errParsing := json.Marshal(productMock)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		//act
		hd.Create()(res, req)

		//assert
		require.NoError(t, errParsing)
		mockProductService.AssertNotCalled(t, "Create")
		require.Equal(t, http.StatusUnprocessableEntity, res.Code)

	})
	t.Run("create_json_parse_error", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		// Invalid JSON to trigger parse error
		invalidJSON := []byte(`{"invalid":json}`)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		//act
		hd.Create()(res, req)

		//assert
		mockProductService.AssertNotCalled(t, "Create")
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("create_internal_server_error", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		// Return a generic error that isn't a ServiceError
		mockProductService.On("Create", productMock).Return(&models.Product{}, errors.New("database connection error"))

		body, errParsing := json.Marshal(productMock)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		//act
		hd.Create()(res, req)

		//assert
		assert.NoError(t, errParsing)
		mockProductService.AssertCalled(t, "Create", productMock)
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("create_conflict", func(t *testing.T) {
		//arranges
		mockProductService := new(product.MockProductService)

		body, errParsing := json.Marshal(productMock)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd := NewProductDefault(mockProductService)

		mockProductService.On("Create", productMock).Return(&models.Product{}, pkg.ServiceErrors[pkg.ErrConflict])

		//act
		hd.Create()(res, req)

		//assert
		require.NoError(t, errParsing)
		mockProductService.AssertCalled(t, "Create", productMock)
		require.Equal(t, http.StatusConflict, res.Code)
	})
}

func TestHandler_GetAll(t *testing.T) {

	testProducts := map[int]models.Product{
		1: {
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    models.StringPtr("P001"),
				Description:                    models.StringPtr("Product test 1"),
				NetWeight:                      float64Ptr(10.5),
				ExpirationRate:                 models.IntPtr(30),
				RecommendedFreezingTemperature: float64Ptr(-18.0),
				FreezingRate:                   models.IntPtr(5),
				ProductTypeId:                  models.IntPtr(1),
				SellerId:                       models.IntPtr(1),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(2.0),
				Height: float64Ptr(3.0),
				Length: float64Ptr(4.0),
			},
		},
		2: {
			ID: 2,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    models.StringPtr("P002"),
				Description:                    models.StringPtr("Product test 2"),
				NetWeight:                      float64Ptr(15.5),
				ExpirationRate:                 models.IntPtr(45),
				RecommendedFreezingTemperature: float64Ptr(-20.0),
				FreezingRate:                   models.IntPtr(8),
				ProductTypeId:                  models.IntPtr(2),
				SellerId:                       models.IntPtr(2),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(5.0),
				Height: float64Ptr(6.0),
				Length: float64Ptr(7.0),
			},
		},
	}

	t.Run("get_all_ok", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		mockProductService.On("GetAll").Return(testProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		res := httptest.NewRecorder()

		//act
		hd.GetAll()(res, req)

		//assert
		mockProductService.AssertExpectations(t)
		require.Equal(t, http.StatusOK, res.Code)

		// Verify response body
		var response map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response, "data")
	})

	t.Run("get_all_internal_error", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		mockProductService.On("GetAll").Return(map[int]models.Product{}, errors.New("database connection error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		res := httptest.NewRecorder()

		//act
		hd.GetAll()(res, req)

		//assert
		mockProductService.AssertExpectations(t)
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestHandler_Find(t *testing.T) {
	// Using models package helper functions for pointers

	testProducts := map[int]models.Product{
		1: {
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    models.StringPtr("P001"),
				Description:                    models.StringPtr("Product test 1"),
				NetWeight:                      float64Ptr(10.5),
				ExpirationRate:                 models.IntPtr(30),
				RecommendedFreezingTemperature: float64Ptr(-18.0),
				FreezingRate:                   models.IntPtr(5),
				ProductTypeId:                  models.IntPtr(1),
				SellerId:                       models.IntPtr(1),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(2.0),
				Height: float64Ptr(3.0),
				Length: float64Ptr(4.0),
			},
		},
		2: {
			ID: 2,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    models.StringPtr("P002"),
				Description:                    models.StringPtr("Product test 2"),
				NetWeight:                      float64Ptr(20.0),
				ExpirationRate:                 models.IntPtr(60),
				RecommendedFreezingTemperature: float64Ptr(-20.0),
				FreezingRate:                   models.IntPtr(10),
				ProductTypeId:                  models.IntPtr(2),
				SellerId:                       models.IntPtr(2),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(2.5),
				Height: float64Ptr(3.5),
				Length: float64Ptr(4.5),
			},
		},
	}

	t.Run("find_all", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		expectedBody := map[string]any{"data": testProducts}
		expectedJson, errParsing := json.Marshal(expectedBody)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)

		mockProductService.On("GetAll").Return(testProducts, nil)

		//act
		hd.GetAll()(res, req)

		//assert
		assert.NoError(t, errParsing)
		mockProductService.AssertCalled(t, "GetAll")
		require.Equal(t, expectedJson, res.Body.Bytes())
		require.Equal(t, http.StatusOK, res.Code)

	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		mockProductService.On("GetById", 1).Return(&models.Product{}, pkg.ServiceErrors[pkg.ErrNotFound])

		//act
		hd.GetById()(res, req)

		//assert
		mockProductService.AssertCalled(t, "GetById", 1)
		require.Equal(t, `{"status":"Not Found","message":"error: Not found"}`, string(res.Body.Bytes()))
		require.Equal(t, http.StatusNotFound, res.Code)

	})

	t.Run("find_by_id_existent", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		productTest := testProducts[1]

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		mockProductService.On("GetById", 1).Return(&productTest, nil)

		expectedBody := map[string]any{"data": productTest}
		expectedBodyJson, errParsing := json.Marshal(expectedBody)

		//act
		hd.GetById()(res, req)

		//assert
		assert.NoError(t, errParsing)
		mockProductService.AssertCalled(t, "GetById", 1)
		require.Equal(t, http.StatusOK, res.Code)
		require.Equal(t, expectedBodyJson, res.Body.Bytes())
	})
}

func TestHandler_Update(t *testing.T) {

	testProducts := map[int]models.Product{
		1: {
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    models.StringPtr("P001"),
				Description:                    models.StringPtr("Product test 1"),
				NetWeight:                      float64Ptr(10.5),
				ExpirationRate:                 models.IntPtr(30),
				RecommendedFreezingTemperature: float64Ptr(-18.0),
				FreezingRate:                   models.IntPtr(5),
				ProductTypeId:                  models.IntPtr(1),
				SellerId:                       models.IntPtr(1),
			},
			Dimensions: models.Dimensions{
				Width:  float64Ptr(2.0),
				Height: float64Ptr(3.0),
				Length: float64Ptr(4.0),
			},
		},
	}

	t.Run("update_ok", func(t *testing.T) {
		//arrange

		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		testProduct := testProducts[1]
		productJson, errParsing := json.Marshal(testProduct)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/", bytes.NewReader(productJson))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockProductService.On("Update", 1, testProduct).Return(&testProduct, nil)

		//act
		hd.Patch()(res, req)

		//assert
		assert.NoError(t, errParsing)
		mockProductService.AssertCalled(t, "Update", 1, testProduct)
		require.Equal(t, http.StatusOK, res.Code)
		require.NotEmpty(t, res.Body)
	})

	t.Run("update_json_parse_error", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		// Invalid JSON to trigger parse error
		invalidJSON := []byte(`{"invalid":json}`)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/", bytes.NewReader(invalidJSON))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		//act
		hd.Patch()(res, req)

		//assert
		mockProductService.AssertNotCalled(t, "Update")
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("update_internal_server_error", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		testProduct := testProducts[1]
		productJson, errParsing := json.Marshal(testProduct)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/", bytes.NewReader(productJson))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Return a generic error that isn't a ServiceError
		mockProductService.On("Update", 1, testProduct).Return(&models.Product{}, errors.New("database connection error"))

		//act
		hd.Patch()(res, req)

		//assert
		assert.NoError(t, errParsing)
		mockProductService.AssertCalled(t, "Update", 1, testProduct)
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("update_non_existent", func(t *testing.T) {
		//arrange

		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		testProduct := testProducts[1]
		productJson, errParsing := json.Marshal(testProduct)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/", bytes.NewReader(productJson))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockProductService.On("Update", 1, testProduct).Return(&models.Product{}, pkg.ServiceErrors[pkg.ErrNotFound])

		//act
		hd.Patch()(res, req)

		//assert
		assert.NoError(t, errParsing)
		mockProductService.AssertCalled(t, "Update", 1, testProduct)
		require.Equal(t, http.StatusNotFound, res.Code)
		require.NotEmpty(t, res.Body)
	})

}

func TestHandler_ReportRecords(t *testing.T) {
	// Create test product records
	productRecords := []models.ProductRecordResponse{
		{
			ProductId:    1,
			Description:  "Test Product 1",
			RecordsCount: 5,
		},
		{
			ProductId:    2,
			Description:  "Test Product 2",
			RecordsCount: 3,
		},
	}

	t.Run("get_all_product_records_ok", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		// Return all product records when id is nil
		mockProductService.On("GetProductRecords", (*int)(nil)).Return(productRecords, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/reportRecords", nil)
		res := httptest.NewRecorder()

		//act
		hd.ReportRecords()(res, req)

		//assert
		mockProductService.AssertExpectations(t)
		require.Equal(t, http.StatusOK, res.Code)

		// Verify response body
		var response map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response, "data")
	})

	t.Run("get_product_records_by_id_ok", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		// Filter by product ID
		productId := 1
		filtered := []models.ProductRecordResponse{productRecords[0]}
		mockProductService.On("GetProductRecords", &productId).Return(filtered, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/reportRecords?id=1", nil)
		res := httptest.NewRecorder()

		//act
		hd.ReportRecords()(res, req)

		//assert
		mockProductService.AssertExpectations(t)
		require.Equal(t, http.StatusOK, res.Code)

		// Verify response body
		var response map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response, "data")
	})

	t.Run("get_product_records_invalid_id", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		// Invalid ID format
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/reportRecords?id=invalid", nil)
		res := httptest.NewRecorder()

		//act
		hd.ReportRecords()(res, req)

		//assert
		mockProductService.AssertNotCalled(t, "GetProductRecords")
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("get_product_records_not_found", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		productId := 999 // Non-existent ID
		mockProductService.On("GetProductRecords", &productId).Return([]models.ProductRecordResponse(nil), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/reportRecords?id=999", nil)
		res := httptest.NewRecorder()

		//act
		hd.ReportRecords()(res, req)

		//assert
		mockProductService.AssertExpectations(t)
		require.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("get_product_records_internal_error", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		productId := 1
		mockProductService.On("GetProductRecords", &productId).Return([]models.ProductRecordResponse{}, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/reportRecords?id=1", nil)
		res := httptest.NewRecorder()

		//act
		hd.ReportRecords()(res, req)

		//assert
		mockProductService.AssertExpectations(t)
		require.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		mockProductService.On("Delete", 1).Return(nil)

		//act
		hd.Delete()(res, req)

		//assert
		mockProductService.AssertCalled(t, "Delete", 1)
		require.Equal(t, http.StatusNoContent, res.Code)
		require.Empty(t, res.Body)
	})

	t.Run("delete_internal_server_error", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		// Return a generic error that isn't a ServiceError
		mockProductService.On("Delete", 1).Return(errors.New("database connection error"))

		//act
		hd.Delete()(res, req)

		//assert
		mockProductService.AssertCalled(t, "Delete", 1)
		require.Equal(t, http.StatusInternalServerError, res.Code)
		require.NotEmpty(t, res.Body)
	})

	t.Run("delete_non_existent", func(t *testing.T) {
		//arrange
		mockProductService := new(product.MockProductService)
		hd := NewProductDefault(mockProductService)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		mockProductService.On("Delete", 1).Return(pkg.ServiceErrors[pkg.ErrNotFound])

		//act
		hd.Delete()(res, req)

		//assert
		mockProductService.AssertCalled(t, "Delete", 1)
		require.Equal(t, http.StatusNotFound, res.Code)
		require.NotEmpty(t, res.Body)
	})

}
