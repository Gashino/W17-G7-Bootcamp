package handler

import (
	"app/pkg/models"
	"app/test/locality"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type localityResponseStruct struct {
	Data models.LocalitiesDoc `json:"data"`
}

type localityErrorResponseStruct struct {
	Message string `json:"message"`
}

type localityReportResponseStruct struct {
	Data models.LocalityBySellerResponse `json:"data"`
}

func TestLocalityDefault_Create(t *testing.T) {
	t.Run("Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado", func(t *testing.T) {
		// Arrange
		mockService := new(locality.MockLocalityService)
		expectedLocality := models.Locality{
			ID: 1,
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "Buenos Aires",
				ProvinceName: "Buenos Aires",
				CountryName:  "Argentina",
			},
		}
		mockService.On("Create", mock.AnythingOfType("models.Locality")).Return(expectedLocality, nil)

		createRequest := models.LocalityCreateRequest{
			LocalityName: models.StringPtr("Buenos Aires"),
			ProvinceName: models.StringPtr("Buenos Aires"),
			CountryName:  models.StringPtr("Argentina"),
		}

		hd := NewLocalityDefault(mockService)
		reqBody, _ := json.Marshal(createRequest)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/localities", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		actualResp := localityResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusCreated
		expectedResp := localityResponseStruct{
			Data: models.LocalitiesDoc{
				ID:           1,
				LocalityName: "Buenos Aires",
				ProvinceName: "Buenos Aires",
				CountryName:  "Argentina",
			},
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Create", mock.AnythingOfType("models.Locality"))
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Si el objeto JSON no contiene los campos necesarios se devolverá un código 422", func(t *testing.T) {
		// Arrange
		mockService := new(locality.MockLocalityService)

		// Datos incompletos - faltan campos requeridos
		createRequest := models.LocalityCreateRequest{
			LocalityName: models.StringPtr("Buenos Aires"),
			// Faltan: ProvinceName y CountryName
		}

		hd := NewLocalityDefault(mockService)
		reqBody, _ := json.Marshal(createRequest)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/localities", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		actualResp := localityErrorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusUnprocessableEntity
		expectedResp := localityErrorResponseStruct{
			Message: "error: Validation error",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "Create")
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Si la locality ya existe devuelve un error 409 Conflict", func(t *testing.T) {
		// Arrange
		mockService := new(locality.MockLocalityService)
		mockService.On("Create", mock.AnythingOfType("models.Locality")).Return(models.Locality{}, fmt.Errorf("locality already exists"))

		createRequest := models.LocalityCreateRequest{
			LocalityName: models.StringPtr("Buenos Aires"),
			ProvinceName: models.StringPtr("Buenos Aires"),
			CountryName:  models.StringPtr("Argentina"),
		}

		hd := NewLocalityDefault(mockService)
		reqBody, _ := json.Marshal(createRequest)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/localities", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		actualResp := localityErrorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusConflict
		expectedResp := localityErrorResponseStruct{
			Message: "error: Resource conflict",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Create", mock.AnythingOfType("models.Locality"))
		require.Equal(t, expectedResp, actualResp)
	})
}

func TestLocalityDefault_SellersByLocality(t *testing.T) {
	t.Run("Cuando la petición sea exitosa el backend devolverá la información del reporte solicitado", func(t *testing.T) {
		// Arrange
		mockService := new(locality.MockLocalityService)
		expectedResponse := models.LocalityBySellerResponse{
			ID:           1,
			LocalityName: models.StringPtr("Buenos Aires"),
			SellerCount:  models.StringPtr("5"),
		}
		mockService.On("GetCantSellersByLocality", 1).Return(expectedResponse, nil)

		hd := NewLocalityDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportSellers/1", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.SellersByLocality()(res, req)

		// Assert
		actualResp := localityReportResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusCreated
		expectedResp := localityReportResponseStruct{
			Data: models.LocalityBySellerResponse{
				ID:           1,
				LocalityName: models.StringPtr("Buenos Aires"),
				SellerCount:  models.StringPtr("5"),
			},
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "GetCantSellersByLocality", 1)
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando el ID no sea válido se devolverá un código 400", func(t *testing.T) {
		// Arrange
		mockService := new(locality.MockLocalityService)

		hd := NewLocalityDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportSellers/invalid", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL inválido
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.SellersByLocality()(res, req)

		// Assert
		actualResp := localityErrorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusBadRequest
		expectedResp := localityErrorResponseStruct{
			Message: "error: Bad request",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "GetCantSellersByLocality")
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando la locality no exista se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(locality.MockLocalityService)
		mockService.On("GetCantSellersByLocality", 999).Return(models.LocalityBySellerResponse{}, fmt.Errorf("locality not found"))

		hd := NewLocalityDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportSellers/999", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.SellersByLocality()(res, req)

		// Assert
		actualResp := localityErrorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusNotFound
		expectedResp := localityErrorResponseStruct{
			Message: "error: Not found",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "GetCantSellersByLocality", 999)
		require.Equal(t, expectedResp, actualResp)
	})
}
