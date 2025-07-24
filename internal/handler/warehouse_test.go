package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/warehouse"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestWarehouse_GetAll(t *testing.T) {
	t.Run("Cuando la petición sea exitosa el backend devolverá un listado de todas los warehouses existentes", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		expectedWarehouses := map[int]models.Warehouse{
			1: {
				ID:             1,
				WarehouseCode:  "TEST1",
				Address:        "Fake Street",
				Telephone:      "123456789",
				MinCapacity:    1,
				MinTemperature: 2,
			},
			2: {
				ID:             2,
				WarehouseCode:  "TEST2",
				Address:        "Fake 2 Street",
				Telephone:      "987654321",
				MinCapacity:    10,
				MinTemperature: 20,
			},
		}

		mockService.On("FindAll").Return(expectedWarehouses, nil)

		// Act
		hd := NewWarehouseDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/warehouse", nil)
		res := httptest.NewRecorder()
		hd.GetAll()(res, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `{
			"data": {
				"1": {
					"id": 1,
					"warehouse_code": "TEST1",
					"address": "Fake Street",
					"telephone": "123456789",
					"minimun_capacity": 1,
					"minimun_temperature": 2
				},
				"2": {
					"id": 2,
					"warehouse_code": "TEST2",
					"address": "Fake 2 Street",
					"telephone": "987654321",
					"minimun_capacity": 10,
					"minimun_temperature": 20
				}
			}
		}`

		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)

		mockService.AssertCalled(t, "FindAll")

		require.JSONEq(t, expectedBody, actualBody)
	})
}

func TestWarehouse_GetOne(t *testing.T) {
	t.Run("Cuando la petición sea exitosa el backend devolverá la información del warehouse solicitado", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		expectedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "TEST1",
			Address:        "Fake Street",
			Telephone:      "123456789",
			MinCapacity:    1,
			MinTemperature: 2,
		}
		mockService.On("FindByID", 1).Return(expectedWarehouse, nil)

		hd := NewWarehouseDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/warehouse/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		// Act
		hd.GetOne()(res, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `{
							"data": {
									"id": 1,
									"warehouse_code": "TEST1",
									"address": "Fake Street",
									"telephone": "123456789",
									"minimun_capacity": 1,
									"minimun_temperature": 2
							}
						}`
		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "FindByID", 1)
		require.JSONEq(t, expectedBody, actualBody)
	})
	t.Run("Cuando el warehouse no exista se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		mockService.On("FindByID", 1).Return(models.Warehouse{}, svcErr)

		hd := NewWarehouseDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/warehouse/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		// Act
		hd.GetOne()(res, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{
							"status": "Not Found",
							"message": "error: Not found"
						}`
		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "FindByID", 1)
		require.JSONEq(t, expectedBody, actualBody)
	})
}

func TestWarehouse_Add(t *testing.T) {
	t.Run("Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado.", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "TEST123",
			Address:        "Test Street 123",
			Telephone:      "555-0123",
			MinCapacity:    100,
			MinTemperature: 5,
		}
		expectedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "TEST123",
			Address:        "Test Street 123",
			Telephone:      "555-0123",
			MinCapacity:    100,
			MinTemperature: 5,
		}
		mockService.On("Add", warehouseDoc).Return(expectedWarehouse, nil)

		hd := NewWarehouseDefault(mockService)
		requestBody := `{
			"warehouse_code": "TEST123",
			"address": "Test Street 123",
			"telephone": "555-0123",
			"minimun_capacity": 100,
			"minimun_temperature": 5
		}`

		// Act
		req := httptest.NewRequest(http.MethodPost, "/api/v1/warehouse", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd.Add()(res, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{
			"data": {
				"id": 1,
				"warehouse_code": "TEST123",
				"address": "Test Street 123",
				"telephone": "555-0123",
				"minimun_capacity": 100,
				"minimun_temperature": 5
			}
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Add", warehouseDoc)
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el objeto JSON no contiene los campos necesarios se devolverá un código 422", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)
		requestBody := `{
			"warehouse_code": "",
			"address": "Test Street 123",
			"telephone": "555-0123",
			"minimun_capacity": 100,
			"minimun_temperature": 5
		}`

		// Act
		req := httptest.NewRequest(http.MethodPost, "/api/v1/warehouse", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd.Add()(res, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{
			"status": "Unprocessable Entity",
			"message": "error: Validation error"
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertNotCalled(t, "Add")
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el warehouse_code ya existe devuelve un error 409 Conflict", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "EXISTING123",
			Address:        "Test Street 123",
			Telephone:      "555-0123",
			MinCapacity:    100,
			MinTemperature: 5,
		}
		conflictErr := pkg.ServiceErrors[pkg.ErrConflict]
		mockService.On("Add", warehouseDoc).Return(models.Warehouse{}, conflictErr)

		hd := NewWarehouseDefault(mockService)
		requestBody := `{
			"warehouse_code": "EXISTING123",
			"address": "Test Street 123",
			"telephone": "555-0123",
			"minimun_capacity": 100,
			"minimun_temperature": 5
		}`

		// Act
		req := httptest.NewRequest(http.MethodPost, "/api/v1/warehouse", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		hd.Add()(res, req)

		// Assert
		expectedCode := http.StatusInternalServerError
		expectedBody := `{
			"status": "Internal Server Error",
			"message": "error: error: Resource conflict"
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Add", warehouseDoc)
		require.JSONEq(t, expectedBody, res.Body.String())
	})
}

func TestWarehouse_Update(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá el warehouse con la información actualizada junto con un código 200", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "UPDATED123",
			Address:        "Updated Street 456",
			Telephone:      "555-9999",
			MinCapacity:    200,
			MinTemperature: 10,
		}
		expectedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "UPDATED123",
			Address:        "Updated Street 456",
			Telephone:      "555-9999",
			MinCapacity:    200,
			MinTemperature: 10,
		}
		mockService.On("Update", 1, warehouseDoc).Return(expectedWarehouse, nil)

		hd := NewWarehouseDefault(mockService)
		requestBody := `{
			"warehouse_code": "UPDATED123",
			"address": "Updated Street 456",
			"telephone": "555-9999",
			"minimun_capacity": 200,
			"minimun_temperature": 10
		}`

		// Act
		req := httptest.NewRequest(http.MethodPut, "/api/v1/warehouse/1", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		hd.Update()(res, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{
			"data": {
				"id": 1,
				"warehouse_code": "UPDATED123",
				"address": "Updated Street 456",
				"telephone": "555-9999",
				"minimun_capacity": 200,
				"minimun_temperature": 10
			}
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Update", 1, warehouseDoc)
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el warehouse que se desea actualizar no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange
		mockService := new(warehouse.MockWarehouseService)
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "NONEXISTENT123",
			Address:        "Test Street 789",
			Telephone:      "555-1111",
			MinCapacity:    50,
			MinTemperature: 15,
		}
		notFoundErr := pkg.ServiceErrors[pkg.ErrNotFound]
		mockService.On("Update", 999, warehouseDoc).Return(models.Warehouse{}, notFoundErr)

		hd := NewWarehouseDefault(mockService)
		requestBody := `{
			"warehouse_code": "NONEXISTENT123",
			"address": "Test Street 789",
			"telephone": "555-1111",
			"minimun_capacity": 50,
			"minimun_temperature": 15
		}`

		// Act
		req := httptest.NewRequest(http.MethodPut, "/api/v1/warehouse/999", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		hd.Update()(res, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{
			"status": "Not Found",
			"message": "error: error: Not found"
		}`

		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "Update", 999, warehouseDoc)
		require.JSONEq(t, expectedBody, res.Body.String())
	})
}

func TestWarehouse_Delete(t *testing.T) {
	t.Run("Cuando el warehouse no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Cuando la eliminación sea exitosa se devolverá un código 204", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}
