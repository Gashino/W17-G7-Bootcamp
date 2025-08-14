package handler

import (
	"app/pkg"
	"app/pkg/models"
	warehouse "app/test/carry"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCarry_SearchByLocality(t *testing.T) {
	t.Run("Cuando la petición sea exitosa con id específico el backend devolverá los carries de esa localidad", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		expectedCarries := map[string]models.CarryByLocality{
			"1": {
				LocalityId:   "1",
				LocalityName: "Buenos Aires",
				CarriesCount: 5,
			},
		}

		mockService.On("SearchByLocality", 1).Return(expectedCarries, nil)

		// Act
		hd := NewCarryDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/carry?id=1", nil)
		res := httptest.NewRecorder()
		hd.SearchByLocality()(res, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{
			"data": [
				{
					"locality_id": "1",
					"locality_name": "Buenos Aires",
					"carries_count": 5
				}
			]
		}`

		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "SearchByLocality", 1)
		require.JSONEq(t, expectedBody, actualBody)
	})

	t.Run("Cuando la petición sea exitosa sin id específico el backend devolverá todos los carries por localidad", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		expectedCarries := map[string]models.CarryByLocality{
			"1": {
				LocalityId:   "1",
				LocalityName: "Buenos Aires",
				CarriesCount: 5,
			},
			"2": {
				LocalityId:   "2",
				LocalityName: "Córdoba",
				CarriesCount: 3,
			},
		}

		mockService.On("SearchByLocality", -1).Return(expectedCarries, nil)

		// Act
		hd := NewCarryDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/carry", nil)
		res := httptest.NewRecorder()
		hd.SearchByLocality()(res, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{
			"data": [
				{
					"locality_id": "1",
					"locality_name": "Buenos Aires",
					"carries_count": 5
				},
				{
					"locality_id": "2",
					"locality_name": "Córdoba",
					"carries_count": 3
				}
			]
		}`

		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "SearchByLocality", -1)
		require.JSONEq(t, expectedBody, actualBody)
	})

	t.Run("Cuando no se encuentren carries se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		expectedCarries := map[string]models.CarryByLocality{}

		mockService.On("SearchByLocality", 999).Return(expectedCarries, nil)

		// Act
		hd := NewCarryDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/carry?id=999", nil)
		res := httptest.NewRecorder()
		hd.SearchByLocality()(res, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{
			"status": "Not Found",
			"message": "Not found"
		}`

		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "SearchByLocality", 999)
		require.JSONEq(t, expectedBody, actualBody)
	})

	t.Run("Cuando el id no sea un número válido se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)

		// Act
		hd := NewCarryDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/carry?id=invalid", nil)
		res := httptest.NewRecorder()
		hd.SearchByLocality()(res, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{
			"status": "Not Found",
			"message": "error: Not found"
		}`

		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "SearchByLocality")
		require.JSONEq(t, expectedBody, actualBody)
	})

	t.Run("Cuando ocurra un error en el servicio se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		mockService.On("SearchByLocality", 1).Return(map[string]models.CarryByLocality{}, svcErr)

		// Act
		hd := NewCarryDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/carry?id=1", nil)
		res := httptest.NewRecorder()
		hd.SearchByLocality()(res, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{
			"status": "Not Found",
			"message": "error: Not found"
		}`

		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "SearchByLocality", 1)
		require.JSONEq(t, expectedBody, actualBody)
	})
}

func TestCarry_Create(t *testing.T) {
	t.Run("Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		carryInput := models.Carry{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Street 123",
			Telephone:   "555-0123",
			LocalityId:  1,
		}
		expectedCarry := models.Carry{
			ID:          1,
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Street 123",
			Telephone:   "555-0123",
			LocalityId:  1,
		}
		mockService.On("Create", carryInput).Return(expectedCarry, nil)

		hd := NewCarryDefault(mockService)
		requestBody := `{
			"cid": "CID123",
			"company_name": "Test Company",
			"address": "Test Street 123",
			"telephone": "555-0123",
			"locality_id": 1
		}`

		// Act
		req := httptest.NewRequest(http.MethodPost, "/api/v1/carry", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{
			"data": {
				"cid": "CID123",
				"company_name": "Test Company",
				"address": "Test Street 123",
				"telephone": "555-0123",
				"locality_id": 1
			}
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Create", carryInput)
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el objeto JSON no contiene los campos necesarios se devolverá un código 422", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		hd := NewCarryDefault(mockService)
		requestBody := `{
			"cid": "",
			"company_name": "Test Company",
			"address": "Test Street 123",
			"telephone": "555-0123",
			"locality_id": 1
		}`

		// Act
		req := httptest.NewRequest(http.MethodPost, "/api/v1/carry", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{
			"status": "Unprocessable Entity",
			"message": "error: Validation error"
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "Create")
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el JSON es inválido se devolverá un código 422", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		hd := NewCarryDefault(mockService)
		requestBody := `{
			"cid": "CID123",
			"company_name": "Test Company"
			"address": "Test Street 123",
		}` // JSON malformado

		// Act
		req := httptest.NewRequest(http.MethodPost, "/api/v1/carry", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{
			"status": "Unprocessable Entity",
			"message": "error: Validation error"
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "Create")
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el cid ya existe devuelve un error 409 (Conflict)", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockCarryService)
		carryInput := models.Carry{
			Cid:         "EXISTING123",
			CompanyName: "Test Company",
			Address:     "Test Street 123",
			Telephone:   "555-0123",
			LocalityId:  1,
		}
		conflictErr := pkg.ServiceErrors[pkg.ErrConflict]
		mockService.On("Create", carryInput).Return(models.Carry{}, conflictErr)

		hd := NewCarryDefault(mockService)
		requestBody := `{
			"cid": "EXISTING123",
			"company_name": "Test Company",
			"address": "Test Street 123",
			"telephone": "555-0123",
			"locality_id": 1
		}`

		// Act
		req := httptest.NewRequest(http.MethodPost, "/api/v1/carry", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{
			"status": "Conflict",
			"message": "error: error: Resource conflict"
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Create", carryInput)
		require.JSONEq(t, expectedBody, res.Body.String())
	})
}
