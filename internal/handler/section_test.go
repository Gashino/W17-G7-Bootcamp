package handler

import (
	"app/pkg/models"
	"app/test/section"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSectionDefault_GetAll(t *testing.T) {
	t.Run("Cuando la petición sea exitosa el backend devolverá un listado de todas los sections existentes", func(t *testing.T) {
		mockService := new(section.MockSectionService)
		expectedSections := map[int]models.Section{
			3: {
				ID: 3,
				SectionAttributes: models.SectionAttributes{
					SectionNumber:      3,
					CurrentTemperature: 25,
					MinimumTemperature: 15,
					CurrentCapacity:    30,
					MinimumCapacity:    5,
					MaximumCapacity:    0,
					WarehouseID:        2,
					ProductTypeID:      102,
				},
			},
			4: {
				ID: 4,
				SectionAttributes: models.SectionAttributes{
					SectionNumber:      4,
					CurrentTemperature: 10,
					MinimumTemperature: 5,
					CurrentCapacity:    120,
					MinimumCapacity:    30,
					MaximumCapacity:    0,
					WarehouseID:        2,
					ProductTypeID:      107,
				},
			},
			5: {
				ID: 5,
				SectionAttributes: models.SectionAttributes{
					SectionNumber:      5,
					CurrentTemperature: -5,
					MinimumTemperature: -10,
					CurrentCapacity:    90,
					MinimumCapacity:    15,
					MaximumCapacity:    0,
					WarehouseID:        3,
					ProductTypeID:      101,
				},
			},
		}
		mockService.On("GetAll").Return(expectedSections, nil)

		hd := NewSectionDefault(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/sections", nil)
		res := httptest.NewRecorder()

		hd.GetAll()(res, req)

		expectedCode := http.StatusOK
		expectedBody := `{
							"data": [
								{
									"id": 3,
									"section_number": 3,
									"current_temperature": 25,
									"minimum_temperature": 15,
									"current_capacity": 30,
									"minimum_capacity": 5,
									"maximum_capacity": 0,
									"warehouse_id": 2,
									"product_type_id": 102
								},
								{
									"id": 4,
									"section_number": 4,
									"current_temperature": 10,
									"minimum_temperature": 5,
									"current_capacity": 120,
									"minimum_capacity": 30,
									"maximum_capacity": 0,
									"warehouse_id": 2,
									"product_type_id": 107
								},
								{
									"id": 5,
									"section_number": 5,
									"current_temperature": -5,
									"minimum_temperature": -10,
									"current_capacity": 90,
									"minimum_capacity": 15,
									"maximum_capacity": 0,
									"warehouse_id": 3,
									"product_type_id": 101
								}
							],
							"message": "success"
						}`
		actualBody := res.Body.String()
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "GetAll")
		require.JSONEq(t, expectedBody, actualBody)
	})
	t.Run("error", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestSectionDefault_GetByID(t *testing.T) {
	t.Run("Cuando la section no exista se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
	t.Run("Cuando la petición sea exitosa el backend devolverá la información de la section solicitada", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestSectionDefault_PostSection(t *testing.T) {
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
	t.Run("Si el section_number ya existe devuelve un error 409 Conflict", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestSectionDefault_Update(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá la section con la información actualizada junto con un código 200", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
	t.Run("Si el section que se desea actualizar no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestSectionDefault_Delete(t *testing.T) {
	t.Run("Cuando el section no existe se devolverá un código 404", func(t *testing.T) {
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
