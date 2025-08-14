package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/purchase_order"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePurchaseOrder(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)
		orderNumber := "order#2"
		orderDate := "2021-04-05"
		trackingCode := "xyz789"
		buyerID := 1
		productRecordID := 2
		id := 2

		mockService.On("Create", models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     orderNumber,
				OrderDate:       orderDate,
				TrackingCode:    trackingCode,
				BuyerID:         buyerID,
				ProductRecordID: productRecordID,
			},
		}).Return(models.PurchaseOrder{
			ID: id,
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     orderNumber,
				OrderDate:       orderDate,
				TrackingCode:    trackingCode,
				BuyerID:         buyerID,
				ProductRecordID: productRecordID,
			},
		}, nil)

		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_number": "order#2",
			"order_date": "2021-04-05",
			"tracking_code": "xyz789",
			"buyer_id": 1,
			"product_record_id": 2
		}`

		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"data": {
				"id": 2,
				"order_number": "order#2",
				"order_date": "2021-04-05",
				"tracking_code": "xyz789",
				"buyer_id": 1,
				"product_record_id": 2
			}
		}`
		expectedCode := 201
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("create_fail missing order_number", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)
		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_date": "2021-04-05",
			"tracking_code": "xyz789",
			"buyer_id": 1,
			"product_record_id": 2
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"order_number is required", "status":"Bad Request"
		}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail missing order_date", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)
		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_number": "order#2",
			"tracking_code": "xyz789",
			"buyer_id": 1,
			"product_record_id": 2
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"order_date is required", "status":"Bad Request"
		}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail missing tracking_code", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)
		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_number": "order#2",
			"order_date": "2021-04-05",
			"buyer_id": 1,
			"product_record_id": 2
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"tracking_code is required", "status":"Bad Request"
		}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail missing buyer_id", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)
		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_number": "order#2",
			"order_date": "2021-04-05",
			"tracking_code": "xyz789",
			"product_record_id": 2
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"buyer_id is required", "status":"Bad Request"
		}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail missing product_record_id", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)
		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_number": "order#2",
			"order_date": "2021-04-05",
			"tracking_code": "xyz789",
			"buyer_id": 1
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"product_record_id is required", "status":"Bad Request"
		}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail invalid JSON", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)
		hd := NewPurchaseOrderHandler(mockService)

		reqBody := bytes.NewReader([]byte("invalid json"))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"Invalid JSON format", "status":"Bad Request"
		}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail duplicate order_number", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)

		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#1",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		mockService.On("Create", input).Return(models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrBadRequest])

		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_number": "order#1",
			"order_date": "2021-04-05",
			"tracking_code": "xyz789",
			"buyer_id": 1,
			"product_record_id": 2
		}`

		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"Bad request", "status":"Bad Request"
		}`
		expectedCode := 400
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("create_fail service error", func(t *testing.T) {
		mockService := new(purchase_order.MockPurchaseOrderService)

		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#2",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		mockService.On("Create", input).Return(models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrInternalServer])

		hd := NewPurchaseOrderHandler(mockService)
		body := `{
			"order_number": "order#2",
			"order_date": "2021-04-05",
			"tracking_code": "xyz789",
			"buyer_id": 1,
			"product_record_id": 2
		}`

		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/purchaseOrders", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"Internal server error", "status":"Internal Server Error"
		}`
		expectedCode := 500
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
		mockService.AssertExpectations(t)
	})
}
