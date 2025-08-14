package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/employee"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetEmployee(t *testing.T) {
	t.Run("find_by_id_invalid_format", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		hd := NewEmployeeHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/abc", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		hd.GetEmployee(res, req)

		expected := `{"message":"error: Invalid ID format", "status":"Bad Request"}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("find_by_id_missing", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		hd := NewEmployeeHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/", nil)
		routeCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		hd.GetEmployee(res, req)

		expected := `{"message":"error: ID is required", "status":"Bad Request"}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}

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
	t.Run("create_fail_invalid_json", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		hd := NewEmployeeHandler(mockService)

		body := `{
			invalid json
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/employees", reqBody)
		res := httptest.NewRecorder()

		hd.CreateEmployee(res, req)

		expected := `{"message":"error: Validation error", "status":"Unprocessable Entity"}`
		expectedCode := 422
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

	t.Run("find_all_internal_server_error", func(t *testing.T) {
		// Arrange
		mockService := new(employee.MockEmployeeService)
		mockService.On("FindAll").Return(map[int]models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer])
		hd := NewEmployeeHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
		res := httptest.NewRecorder()

		// Act
		hd.GetAllEmployees(res, req)

		// Assert
		expected := `{"message":"error: Internal server error", "status":"Internal Server Error"}`
		expectedCode := 500
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
		mockService.AssertExpectations(t)
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
	t.Run("update_invalid_json", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		hd := NewEmployeeHandler(mockService)

		body := `{
			invalid json
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/employees/3", reqBody)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		hd.UpdateEmployee(res, req)

		expected := `{"message":"error: Validation error", "status":"Unprocessable Entity"}`
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("update_empty_fields", func(t *testing.T) {
		mockService := new(employee.MockEmployeeService)
		hd := NewEmployeeHandler(mockService)

		body := `{}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/employees/3", reqBody)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		hd.UpdateEmployee(res, req)

		expected := `{"message":"error: Validation error", "status":"Unprocessable Entity"}`
		expectedCode := 422
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

func TestGetEmployeeInboundOrdersReport(t *testing.T) {
	t.Run("success_get_all_employees_report", func(t *testing.T) {
		// Arrange
		mockService := new(employee.MockEmployeeService)
		handler := NewEmployeeHandler(mockService)

		reports := []models.EmployeeReport{
			{
				ID:                 1,
				CardNumberID:       "E001",
				FirstName:          "John",
				LastName:           "Doe",
				InboundOrdersCount: 5,
			},
			{
				ID:                 2,
				CardNumberID:       "E002",
				FirstName:          "Jane",
				LastName:           "Smith",
				InboundOrdersCount: 3,
			},
		}

		mockService.On("ReportInboundOrdersCountByEmployee", (*int)(nil)).Return(reports, nil)

		req := httptest.NewRequest("GET", "/employees/report/inbound-orders", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetEmployeeInboundOrdersReport(w, req)

		// Assert
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		var response models.EmployeeReportsResponse
		json.Unmarshal(body, &response)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, reports, response.Data)
		mockService.AssertExpectations(t)
	})

	t.Run("success_get_specific_employee_report", func(t *testing.T) {
		// Arrange
		mockService := new(employee.MockEmployeeService)
		handler := NewEmployeeHandler(mockService)

		employeeID := 1
		reports := []models.EmployeeReport{
			{
				ID:                 employeeID,
				CardNumberID:       "E001",
				FirstName:          "John",
				LastName:           "Doe",
				InboundOrdersCount: 5,
			},
		}

		mockService.On("ReportInboundOrdersCountByEmployee", &employeeID).Return(reports, nil)

		req := httptest.NewRequest("GET", fmt.Sprintf("/employees/report/inbound-orders?id=%d", employeeID), nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetEmployeeInboundOrdersReport(w, req)

		// Assert
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		var response models.EmployeeReportsResponse
		json.Unmarshal(body, &response)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, reports, response.Data)
		mockService.AssertExpectations(t)
	})

	t.Run("error_invalid_id_format", func(t *testing.T) {
		// Arrange
		mockService := new(employee.MockEmployeeService)
		handler := NewEmployeeHandler(mockService)

		req := httptest.NewRequest("GET", "/employees/report/inbound-orders?id=invalid", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetEmployeeInboundOrdersReport(w, req)

		// Assert
		resp := w.Result()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error_employee_not_found", func(t *testing.T) {
		// Arrange
		mockService := new(employee.MockEmployeeService)
		handler := NewEmployeeHandler(mockService)

		employeeID := 999
		srvError := pkg.ServiceErrors[pkg.ErrNotFound]
		mockService.On("ReportInboundOrdersCountByEmployee", &employeeID).Return([]models.EmployeeReport{}, srvError)

		req := httptest.NewRequest("GET", fmt.Sprintf("/employees/report/inbound-orders?id=%d", employeeID), nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetEmployeeInboundOrdersReport(w, req)

		// Assert
		resp := w.Result()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("error_internal_server", func(t *testing.T) {
		// Arrange
		mockService := new(employee.MockEmployeeService)
		handler := NewEmployeeHandler(mockService)

		srvError := pkg.ServiceErrors[pkg.ErrInternalServer]
		mockService.On("ReportInboundOrdersCountByEmployee", (*int)(nil)).Return([]models.EmployeeReport{}, srvError)

		req := httptest.NewRequest("GET", "/employees/report/inbound-orders", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetEmployeeInboundOrdersReport(w, req)

		// Assert
		resp := w.Result()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestWriteResponse(t *testing.T) {
	t.Run("write_success_response", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()
		data := struct {
			Message string `json:"message"`
		}{
			Message: "test message",
		}

		// act
		writeResponse(w, http.StatusOK, data, nil)

		// assert
		// Verificar el status code
		require.Equal(t, http.StatusOK, w.Code)

		// Verificar el Content-Type
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		// Verificar el body
		expected := `{"message":"test message"}`
		require.JSONEq(t, expected, w.Body.String())
	})

	t.Run("write_error_response", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()
		testError := fmt.Errorf("test error")

		// act
		writeResponse(w, http.StatusBadRequest, nil, testError)

		// assert
		// Verificar el status code
		require.Equal(t, http.StatusBadRequest, w.Code)

		// Verificar el Content-Type
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		// Verificar el body
		expected := `{"error":"test error"}`
		require.JSONEq(t, expected, w.Body.String())
	})

	t.Run("write_nil_data_response", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()

		// act
		writeResponse(w, http.StatusNoContent, nil, nil)

		// assert
		// Verificar el status code
		require.Equal(t, http.StatusNoContent, w.Code)

		// Verificar el Content-Type
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		// Verificar que el body está vacío o es "null" (dependiendo de cómo json.Encoder maneje nil)
		require.Contains(t, []string{"", "null\n"}, w.Body.String())
	})

	t.Run("write_complex_data_response", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()
		data := struct {
			Items []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"items"`
			Total int `json:"total"`
		}{
			Items: []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				{ID: 1, Name: "Item 1"},
				{ID: 2, Name: "Item 2"},
			},
			Total: 2,
		}

		// act
		writeResponse(w, http.StatusOK, data, nil)

		// assert
		// Verificar el status code
		require.Equal(t, http.StatusOK, w.Code)

		// Verificar el Content-Type
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		// Verificar el body
		expected := `{
			"items": [
				{"id": 1, "name": "Item 1"},
				{"id": 2, "name": "Item 2"}
			],
			"total": 2
		}`
		require.JSONEq(t, expected, w.Body.String())
	})

	t.Run("write_response_with_empty_struct", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()
		data := struct{}{}

		// act
		writeResponse(w, http.StatusOK, data, nil)

		// assert
		// Verificar el status code
		require.Equal(t, http.StatusOK, w.Code)

		// Verificar el Content-Type
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		// Verificar el body
		expected := `{}`
		require.JSONEq(t, expected, w.Body.String())
	})
}

func TestHandleServiceError(t *testing.T) {
	t.Run("handle_service_error", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()
		err := pkg.ServiceErrors[pkg.ErrNotFound]

		// act
		handleServiceError(w, err)

		// assert
		require.Equal(t, http.StatusNotFound, w.Code)
		expected := `{"message":"error: Not found", "status":"Not Found"}`
		require.JSONEq(t, expected, w.Body.String())
	})

	t.Run("handle_non_service_error", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()
		err := fmt.Errorf("un error cualquiera que no es ServiceError")

		// act
		handleServiceError(w, err)

		// assert
		require.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"message":"error: Internal server error", "status":"Internal Server Error"}`
		require.JSONEq(t, expected, w.Body.String())
	})

	t.Run("handle_nil_error", func(t *testing.T) {
		// arrange
		w := httptest.NewRecorder()

		// act
		handleServiceError(w, nil)

		// assert
		require.Equal(t, http.StatusInternalServerError, w.Code)
		expected := `{"message":"error: Internal server error", "status":"Internal Server Error"}`
		require.JSONEq(t, expected, w.Body.String())
	})
}
