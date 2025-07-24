package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/product"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestHandler_Find(t *testing.T) {
	// helpers for pointers
	str := func(s string) *string { return &s }
	f64 := func(f float64) *float64 { return &f }
	i := func(x int) *int { return &x }

	testProducts := map[int]models.Product{
		1: {
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    str("P001"),
				Description:                    str("Product test 1"),
				NetWeight:                      f64(10.5),
				ExpirationRate:                 i(30),
				RecommendedFreezingTemperature: f64(-18.0),
				FreezingRate:                   i(5),
				ProductTypeId:                  i(1),
				SellerId:                       i(1),
			},
			Dimensions: models.Dimensions{
				Width:  f64(2.0),
				Height: f64(3.0),
				Length: f64(4.0),
			},
		},
		2: {
			ID: 2,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    str("P002"),
				Description:                    str("Product test 2"),
				NetWeight:                      f64(20.0),
				ExpirationRate:                 i(60),
				RecommendedFreezingTemperature: f64(-20.0),
				FreezingRate:                   i(10),
				ProductTypeId:                  i(2),
				SellerId:                       i(2),
			},
			Dimensions: models.Dimensions{
				Width:  f64(2.5),
				Height: f64(3.5),
				Length: f64(4.5),
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
	// helpers for pointers
	str := func(s string) *string { return &s }
	f64 := func(f float64) *float64 { return &f }
	i := func(x int) *int { return &x }

	testProducts := map[int]models.Product{
		1: {
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    str("P001"),
				Description:                    str("Product test 1"),
				NetWeight:                      f64(10.5),
				ExpirationRate:                 i(30),
				RecommendedFreezingTemperature: f64(-18.0),
				FreezingRate:                   i(5),
				ProductTypeId:                  i(1),
				SellerId:                       i(1),
			},
			Dimensions: models.Dimensions{
				Width:  f64(2.0),
				Height: f64(3.0),
				Length: f64(4.0),
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
