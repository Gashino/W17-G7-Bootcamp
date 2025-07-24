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
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inboundOrders", bytes.NewReader(inputJson))
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
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inboundOrders", bytes.NewReader(inputJson))
		res := httptest.NewRecorder()
		hd.CreateInboundOrder(res, req)
		expected := `{
				"message":"error: last_name is required", "status":"Unprocessable Entity"
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
		req := httptest.NewRequest(http.MethodGet, "/api/v1/inboundOrders", bytes.NewReader(inputJson))
		res := httptest.NewRecorder()
		hd.CreateInboundOrder(res, req)
		expected := `{"message":"error: Resource conflict", "status":"Conflict"}`
		expectedCode := 409
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}
