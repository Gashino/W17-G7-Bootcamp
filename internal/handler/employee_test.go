package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/employee"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEmployee(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 21

		mockService.On("Save", models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}).Return(models.Employee{
			ID:           &id,
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}, nil)
		hd := NewEmployeeHandler(mockService)
		body := `{
		"card_number_id": "11223342",
		"first_name": "Carlos",
		"last_name": "López",
		"warehouse_id": 3
		}`

		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employee", reqBody)
		res := httptest.NewRecorder()
		hd.CreateEmployee(res, req)
		expected := `{
			    "data": {
					"id": 21,
					"card_number_id": "11223342",
					"first_name": "Carlos",
					"last_name": "López",
					"warehouse_id": 3
    		}
		}`
		expectedCode := 201
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
	t.Run("create_fail 422", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		hd := NewEmployeeHandler(mockService)
		body := `{
		"card_number_id": "11223342",
		"first_name": "Carlos"
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employee", reqBody)
		res := httptest.NewRecorder()
		hd.CreateEmployee(res, req)
		expected := `{
				"message":"error: last_name is required", "status":"Unprocessable Entity"
				}`
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
	t.Run("create_fail 409", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		mockService.On("Save", models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}).Return(models.Employee{}, pkg.ServiceErrors[pkg.ErrConflict])

		hd := NewEmployeeHandler(mockService)

		body := `{
		"card_number_id": "11223342",
		"first_name": "Carlos",
		"last_name": "López",
		"warehouse_id": 3
		}`

		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employee", reqBody)
		res := httptest.NewRecorder()

		hd.CreateEmployee(res, req)

		expected := `{"message":"error: Resource conflict", "status":"Conflict"}`
		expectedCode := 409

		// then
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}

func TestFindEmployee(t *testing.T) {
	t.Run("find_all", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 21

		mockService.On("FindAll").Return(
			map[int]models.Employee{
				1: models.Employee{
					ID:           &id,
					CardNumberID: &cardNumberID,
					FirstName:    &firstName,
					LastName:     &lastName,
					WarehouseID:  &warehouseID,
				},
			},
			nil,
		)
		hd := NewEmployeeHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/employee", nil)
		res := httptest.NewRecorder()
		expected := `{
    "data": [
        {
            "id": 21,
            "card_number_id": "11223342",
            "first_name": "Carlos",
            "last_name": "López",
            "warehouse_id": 3
        }
			]
		}`
		expectedCode := 200
		hd.GetAllEmployees(res, req)
		// then
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}
