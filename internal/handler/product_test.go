package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/product"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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
		mockProductService.On("Create", productMock).Return(&productMock, nil)

		body, errParsing := json.Marshal(productMock)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd := NewProductDefault(mockProductService)

		//act
		hd.Create()(res, req)

		//assert
		require.NoError(t, errParsing)
		mockProductService.AssertExpectations(t)
		require.Equal(t, res.Code, http.StatusCreated)
		require.NotEmpty(t, res.Body)
		require.Equal(t, productMock.ID, 1)

	})
	t.Run("create_failure", func(t *testing.T) {
		//arrange
		productMock := models.Product{}

		mockProductService := new(product.MockProductService)

		body, errParsing := json.Marshal(productMock)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd := NewProductDefault(mockProductService)

		//act
		hd.Create()(res, req)

		//assert
		require.NoError(t, errParsing)
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
		mockProductService.AssertExpectations(t)
		require.Equal(t, http.StatusConflict, res.Code)
	})
}

func TestHandler_Find(t *testing.T) {

	t.Run("find_all", func(t *testing.T) {

	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {

	})

	t.Run("find_by_id_existent", func(t *testing.T) {

	})
}
