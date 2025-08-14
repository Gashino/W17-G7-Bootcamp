package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/inbound_order"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateInboundOrder(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		mockService := new(inbound_order.MockInboundOrderService)
		orderDate, err := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
		if err != nil {
			log.Fatal(err)
		}
		input := models.InboundOrder{
			OrderDate:      orderDate,
			EmployeeID:     3,
			OrderNumber:    "IO-2024-001",
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		mockService.On("Create", input).Return(input, nil)
		inputJson, errInputMarshal := json.Marshal(input)
		if errInputMarshal != nil {
			log.Fatal("Error")
		}
		hd := NewInboundOrderHandler(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inboundOrders", bytes.NewReader(inputJson))
		res := httptest.NewRecorder()
		hd.CreateInboundOrder(res, req)
		expected := `{
    "data": {
        "order_date": "2024-01-15T10:30:00Z",
        "order_number": "IO-2024-001",
        "employee_id": 3,
        "product_batch_id": 1,
        "warehouse_id": 1
    }
}`
		expectedCode := 201
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail_invalid_json", func(t *testing.T) {
		// Arrange
		mockService := new(inbound_order.MockInboundOrderService)
		hd := NewInboundOrderHandler(mockService)

		// Invalid JSON body
		body := `{
			invalid json format
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inboundOrders", bytes.NewReader([]byte(body)))
		res := httptest.NewRecorder()

		// Act
		hd.CreateInboundOrder(res, req)

		// Assert
		expected := `{"message":"error: Validation error", "status":"Unprocessable Entity"}`
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail_validation_error", func(t *testing.T) {
		// Arrange
		mockService := new(inbound_order.MockInboundOrderService)
		hd := NewInboundOrderHandler(mockService)

		// Missing required fields
		input := models.InboundOrder{
			// OrderDate is missing
			EmployeeID: 3,
			// OrderNumber is missing
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		inputJson, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/inboundOrders", bytes.NewReader(inputJson))
		res := httptest.NewRecorder()

		// Act
		hd.CreateInboundOrder(res, req)

		// Assert
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)

		var response struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}
		err = json.NewDecoder(res.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, "Unprocessable Entity", response.Status)
		require.Contains(t, response.Message, "error:")
	})

	t.Run("create_fail 422", func(t *testing.T) {
		mockService := new(inbound_order.MockInboundOrderService)
		orderDate, err := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
		if err != nil {
			log.Fatal(err)
		}
		input := models.InboundOrder{
			OrderDate:      orderDate,
			EmployeeID:     3,
			OrderNumber:    "IO-2024-001",
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		mockService.On("Create", input).Return(models.InboundOrder{}, pkg.ServiceErrors[pkg.ErrUnprocessableEntity])
		inputJson, errInputMarshal := json.Marshal(input)
		if errInputMarshal != nil {
			log.Fatal("Error")
		}
		hd := NewInboundOrderHandler(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inboundOrders", bytes.NewReader(inputJson))
		res := httptest.NewRecorder()
		hd.CreateInboundOrder(res, req)
		expected := `{
				"message":"error: Validation error", "status":"Unprocessable Entity"
				}`
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail 409", func(t *testing.T) {
		mockService := new(inbound_order.MockInboundOrderService)
		orderDate, err := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
		if err != nil {
			log.Fatal(err)
		}
		input := models.InboundOrder{
			OrderDate:      orderDate,
			EmployeeID:     3,
			OrderNumber:    "IO-2024-001",
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		mockService.On("Create", input).Return(models.InboundOrder{}, pkg.ServiceErrors[pkg.ErrConflict])
		inputJson, errInputMarshal := json.Marshal(input)
		if errInputMarshal != nil {
			log.Fatal("Error")
		}
		hd := NewInboundOrderHandler(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inboundOrders", bytes.NewReader(inputJson))
		res := httptest.NewRecorder()
		hd.CreateInboundOrder(res, req)
		expected := `{"message":"error: Resource conflict", "status":"Conflict"}`
		expectedCode := 409
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}
