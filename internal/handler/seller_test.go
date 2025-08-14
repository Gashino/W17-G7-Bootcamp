package handler

import (
	"app/pkg/models"
	"app/test/seller"
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

type responseSellerStruct struct {
	Data    map[int]models.SellerDoc `json:"data"`
	Message string                   `json:"message"`
}

type singleResponseStruct struct {
	Data models.SellerDoc `json:"data"`
}

type errorResponseStruct struct {
	Message string `json:"message"`
}

type deleteResponseStruct struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func TestSellerDefault_GetAll(t *testing.T) {
	t.Run("Cuando la petición sea exitosa el backend devolverá un listado de todas los sections existentes", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		expectedSeller := map[int]models.Seller{
			1: {ID: 1, SellerAttributes: models.SellerAttributes{CId: "12345", CompanyName: "Test Company", Address: "Test Address", Telephone: "123456789", LocalityID: 1}},
			2: {ID: 2, SellerAttributes: models.SellerAttributes{CId: "123456", CompanyName: "Test Company 2", Address: "Test Address 2", Telephone: "987654321", LocalityID: 2}},
			3: {ID: 3, SellerAttributes: models.SellerAttributes{CId: "123457", CompanyName: "Test Company 3", Address: "Test Address 3", Telephone: "555666777", LocalityID: 3}},
		}
		mockService.On("FindAll").Return(expectedSeller, nil)

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sellers", nil)
		res := httptest.NewRecorder()

		// Act
		hd.GetAll()(res, req)

		// Assert
		actualResp := responseSellerStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusOK
		expectedResp := responseSellerStruct{
			Data: map[int]models.SellerDoc{
				1: {ID: 1, CId: "12345", CompanyName: "Test Company", Address: "Test Address", Telephone: "123456789", LocalityID: 1},
				2: {ID: 2, CId: "123456", CompanyName: "Test Company 2", Address: "Test Address 2", Telephone: "987654321", LocalityID: 2},
				3: {ID: 3, CId: "123457", CompanyName: "Test Company 3", Address: "Test Address 3", Telephone: "555666777", LocalityID: 3},
			},
			Message: "success",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "FindAll")
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("error", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		mockService.On("FindAll").Return(map[int]models.Seller{}, fmt.Errorf("database error"))

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sellers", nil)
		res := httptest.NewRecorder()

		// Act
		hd.GetAll()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusInternalServerError
		expectedResp := errorResponseStruct{
			Message: "error: Internal server error",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "FindAll")
		require.Equal(t, expectedResp, actualResp)
	})
}

func TestSellerDefault_GetByID(t *testing.T) {
	t.Run("Cuando la section no exista se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		mockService.On("GetById", 999).Return(models.Seller{}, fmt.Errorf("seller not found"))

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sellers/999", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.GetById()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusNotFound
		expectedResp := errorResponseStruct{
			Message: "error: Not found",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "GetById", 999)
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando la petición sea exitosa el backend devolverá la información de la section solicitada", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		expectedSeller := models.Seller{
			ID: 1,
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}
		mockService.On("GetById", 1).Return(expectedSeller, nil)

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sellers/1", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.GetById()(res, req)

		// Assert
		actualResp := singleResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusCreated
		expectedResp := singleResponseStruct{
			Data: models.SellerDoc{
				ID:          1,
				CId:         "12345",
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "GetById", 1)
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando el ID no es válido se devolverá un código 400", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sellers/invalid", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL con ID inválido
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.GetById()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusBadRequest
		expectedResp := errorResponseStruct{
			Message: "error: Bad request",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "GetById")
		require.Equal(t, expectedResp, actualResp)
	})
}

func TestSellerDefault_Post(t *testing.T) {
	t.Run("Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado.", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		expectedSeller := models.Seller{
			ID: 1,
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "New Company",
				Address:     "New Address",
				Telephone:   "987654321",
				LocalityID:  1,
			},
		}
		mockService.On("Create", mock.AnythingOfType("models.Seller")).Return(expectedSeller, nil)

		createRequest := models.SellerCreateRequest{
			CId:         models.StringPtr("12345"),
			CompanyName: models.StringPtr("New Company"),
			Address:     models.StringPtr("New Address"),
			Telephone:   models.StringPtr("987654321"),
			LocalityID:  models.IntPtr(1),
		}

		hd := NewSellerDefault(mockService)
		reqBody, _ := json.Marshal(createRequest)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		actualResp := singleResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusCreated
		expectedResp := singleResponseStruct{
			Data: models.SellerDoc{
				ID:          1,
				CId:         "12345",
				CompanyName: "New Company",
				Address:     "New Address",
				Telephone:   "987654321",
				LocalityID:  1,
			},
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Create", mock.AnythingOfType("models.Seller"))
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Si el objeto JSON no contiene los campos necesarios se devolverá un código 422", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)

		// Datos incompletos - faltan campos requeridos
		createRequest := models.SellerCreateRequest{
			CId: models.StringPtr("12345"),
			// Faltan: CompanyName, Address, Telephone, LocalityID
		}

		hd := NewSellerDefault(mockService)
		reqBody, _ := json.Marshal(createRequest)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusUnprocessableEntity
		expectedResp := errorResponseStruct{
			Message: "error: Validation error",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "Create")
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Si el section_number ya existe devuelve un error 409 Conflict", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		mockService.On("Create", mock.AnythingOfType("models.Seller")).Return(models.Seller{}, fmt.Errorf("seller already exists"))

		createRequest := models.SellerCreateRequest{
			CId:         models.StringPtr("12345"),
			CompanyName: models.StringPtr("Existing Company"),
			Address:     models.StringPtr("Existing Address"),
			Telephone:   models.StringPtr("987654321"),
			LocalityID:  models.IntPtr(1),
		}

		hd := NewSellerDefault(mockService)
		reqBody, _ := json.Marshal(createRequest)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Act
		hd.Create()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusConflict
		expectedResp := errorResponseStruct{
			Message: "error: Resource conflict",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Create", mock.AnythingOfType("models.Seller"))
		require.Equal(t, expectedResp, actualResp)
	})
}

func TestSellerDefault_Update(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá la section con la información actualizada junto con un código 200", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		expectedSeller := models.Seller{
			ID: 1,
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Updated Company",
				Address:     "Updated Address",
				Telephone:   "987654321",
				LocalityID:  1,
			},
		}
		mockService.On("UpdateFields", 1, mock.AnythingOfType("models.SellerCreateRequest")).Return(expectedSeller, nil)

		updateRequest := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
			Address:     models.StringPtr("Updated Address"),
		}

		hd := NewSellerDefault(mockService)
		reqBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.Update()(res, req)

		// Assert
		actualResp := singleResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusCreated
		expectedResp := singleResponseStruct{
			Data: models.SellerDoc{
				ID:          1,
				CId:         "12345",
				CompanyName: "Updated Company",
				Address:     "Updated Address",
				Telephone:   "987654321",
				LocalityID:  1,
			},
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "UpdateFields", 1, mock.AnythingOfType("models.SellerCreateRequest"))
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Si el section que se desea actualizar no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		mockService.On("UpdateFields", 999, mock.AnythingOfType("models.SellerCreateRequest")).Return(models.Seller{}, fmt.Errorf("seller not found"))

		updateRequest := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
		}

		hd := NewSellerDefault(mockService)
		reqBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/sellers/999", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.Update()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusNotFound
		expectedResp := errorResponseStruct{
			Message: "error: Not found",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "UpdateFields", 999, mock.AnythingOfType("models.SellerCreateRequest"))
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando el ID no es válido se devolverá un código 422", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)

		updateRequest := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
		}

		hd := NewSellerDefault(mockService)
		reqBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/sellers/invalid", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Simular parámetro de URL con ID inválido
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.Update()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusUnprocessableEntity
		expectedResp := errorResponseStruct{
			Message: "error: Validation error",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "UpdateFields")
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando el JSON es inválido se devolverá un código 422", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/sellers/1", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.Update()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusUnprocessableEntity
		expectedResp := errorResponseStruct{
			Message: "error: Validation error",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "UpdateFields")
		require.Equal(t, expectedResp, actualResp)
	})
}

func TestSellerDefault_Delete(t *testing.T) {
	t.Run("Cuando el section no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		mockService.On("DeleteSeller", 999).Return(fmt.Errorf("seller not found"))

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/sellers/999", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.Delete()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusNotFound
		expectedResp := errorResponseStruct{
			Message: "error: Not found",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "DeleteSeller", 999)
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando la eliminación sea exitosa se devolverá un código 204", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)
		mockService.On("DeleteSeller", 1).Return(nil)

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/sellers/1", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.Delete()(res, req)

		// Assert
		actualResp := deleteResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusNoContent
		expectedResp := deleteResponseStruct{
			Message: "success",
			Data:    nil,
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "DeleteSeller", 1)
		require.Equal(t, expectedResp, actualResp)
	})

	t.Run("Cuando el ID no es válido se devolverá un código 400", func(t *testing.T) {
		// Arrange
		mockService := new(seller.MockSellerService)

		hd := NewSellerDefault(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/sellers/invalid", nil)
		res := httptest.NewRecorder()

		// Simular parámetro de URL con ID inválido
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		hd.Delete()(res, req)

		// Assert
		actualResp := errorResponseStruct{}
		err := json.Unmarshal([]byte(res.Body.Bytes()), &actualResp)
		require.NoError(t, err)
		expectedCode := http.StatusBadRequest
		expectedResp := errorResponseStruct{
			Message: "error: Bad request",
		}
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "DeleteSeller")
		require.Equal(t, expectedResp, actualResp)
	})
}
