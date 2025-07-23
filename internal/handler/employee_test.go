package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/employee"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", reqBody)
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", reqBody)
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", reqBody)
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

		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
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
	t.Run("find_by_id_non_existent 404", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		mockService.On("FindById", 3).Return(
			models.Employee{},
			pkg.ServiceErrors[pkg.ErrNotFound],
		)
		hd := NewEmployeeHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/3", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		expected := `{"message":"error: Not found", "status":"Not Found"}`
		expectedCode := 404
		hd.GetEmployee(res, req)
		// then
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("find_by_id_existent", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 3

		mockService.On("FindById", 3).Return(
			models.Employee{
				ID:           &id,
				CardNumberID: &cardNumberID,
				FirstName:    &firstName,
				LastName:     &lastName,
				WarehouseID:  &warehouseID,
			},
			nil,
		)
		hd := NewEmployeeHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/3", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		expected := `{
	"data": {
	    "id": 3,
	    "card_number_id": "11223342",
	    "first_name": "Carlos",
	    "last_name": "López",
	    "warehouse_id": 3
	}
	}`
		expectedCode := 200
		hd.GetEmployee(res, req)
		// then
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("update_ok", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 3

		mockService.On("Update", models.Employee{
			ID:           nil,
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}, 3).Return(
			models.Employee{
				ID:           &id,
				CardNumberID: &cardNumberID,
				FirstName:    &firstName,
				LastName:     &lastName,
				WarehouseID:  &warehouseID,
			},
			nil,
		)
		hd := NewEmployeeHandler(mockService)

		body := `{
		"card_number_id": "11223342",
		"first_name": "Carlos",
		"last_name": "López",
		"warehouse_id": 3
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/employees/3", reqBody)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		expected := `{
    "data": {
        "id": 3,
        "card_number_id": "11223342",
        "first_name": "Carlos",
        "last_name": "López",
        "warehouse_id": 3
    }
}`
		expectedCode := 200
		hd.UpdateEmployee(res, req)
		// then
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
	t.Run("update_non_existent", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3

		mockService.On("Update", models.Employee{
			ID:           nil,
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}, 3).Return(
			models.Employee{},
			pkg.ServiceErrors[pkg.ErrNotFound],
		)
		hd := NewEmployeeHandler(mockService)

		body := `{
		"card_number_id": "11223342",
		"first_name": "Carlos",
		"last_name": "López",
		"warehouse_id": 3
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/employees/3", reqBody)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		expected := `{"message":"error: Not found", "status":"Not Found"}`
		expectedCode := 404
		hd.UpdateEmployee(res, req)
		// then
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}

func TestDeleteEmployee(t *testing.T) {
	t.Run("delete_non_existent", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		mockService.On("Delete", 3).Return(
			pkg.ServiceErrors[pkg.ErrNotFound],
		)
		hd := NewEmployeeHandler(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/3", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		expected := `{"message":"error: Not found", "status":"Not Found"}`
		expectedCode := 404
		hd.DeleteEmployee(res, req)
		// then
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("delete_ok", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		mockService.On("Delete", 3).Return(
			nil,
		)
		hd := NewEmployeeHandler(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/employees/3", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 204
		hd.DeleteEmployee(res, req)
		// then
		require.Equal(t, expectedCode, res.Code)
	})
}
