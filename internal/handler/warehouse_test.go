package handler

import (
	"app/pkg/models"
	"app/test/warehouse"
	"net/http"
	"net/http/httptest"
	"testing"

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

		// Act

		// Assert
	})
	t.Run("Cuando el warehouse no exista se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouse_Add(t *testing.T) {
	t.Run("Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado.", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si el objeto JSON no contiene los campos necesarios se devolverá un código 422", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si el warehouse_code ya existe devuelve un error 409 Conflict", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouse_Update(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá el warehouse con la información actualizada junto con un código 200", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si el warehouse que se desea actualizar no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
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
