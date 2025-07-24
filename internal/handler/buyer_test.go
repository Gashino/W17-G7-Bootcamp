package handler

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/buyer"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateBuyer(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		cardNumberID := "12345678"
		firstName := "John"
		lastName := "Doe"
		id := 1

		mockService.On("Create", models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardNumberID,
				FirstName:    firstName,
				LastName:     lastName,
			},
		}).Return(models.Buyer{
			ID: id,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardNumberID,
				FirstName:    firstName,
				LastName:     lastName,
			},
		}, nil)

		hd := NewBuyerHandler(mockService)
		body := `{
			"card_number_id": "12345678",
			"first_name": "John",
			"last_name": "Doe"
		}`

		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/buyers", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"data": {
				"id": 1,
				"card_number_id": "12345678",
				"first_name": "John",
				"last_name": "Doe"
			}
		}`
		expectedCode := 201
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail missing card_number_id", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		hd := NewBuyerHandler(mockService)
		body := `{
			"first_name": "John",
			"last_name": "Doe"
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/buyers", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"error: card_number_id is required", "status":"Unprocessable Entity"
		}`
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail missing first_name", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		hd := NewBuyerHandler(mockService)
		body := `{
			"card_number_id": "12345678",
			"last_name": "Doe"
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/buyers", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"error: first_name is required", "status":"Unprocessable Entity"
		}`
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_fail missing last_name", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		hd := NewBuyerHandler(mockService)
		body := `{
			"card_number_id": "12345678",
			"first_name": "John"
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/buyers", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{
			"message":"error: last_name is required", "status":"Unprocessable Entity"
		}`
		expectedCode := 422
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("create_conflict duplicate card_number_id", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		cardNumberID := "12345678"
		firstName := "John"
		lastName := "Doe"

		mockService.On("Create", models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardNumberID,
				FirstName:    firstName,
				LastName:     lastName,
			},
		}).Return(models.Buyer{}, pkg.ServiceErrors[pkg.ErrConflict])

		hd := NewBuyerHandler(mockService)
		body := `{
			"card_number_id": "12345678",
			"first_name": "John",
			"last_name": "Doe"
		}`

		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/buyers", reqBody)
		res := httptest.NewRecorder()
		hd.Create()(res, req)

		expected := `{"message":"Resource conflict", "status":"Conflict"}`
		expectedCode := 409
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})
}

func TestFindBuyer(t *testing.T) {
	t.Run("find_all", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		cardNumberID := "12345678"
		firstName := "John"
		lastName := "Doe"
		id := 1

		mockService.On("GetAll").Return(
			map[int]models.Buyer{
				1: {
					ID: id,
					BuyerAttributes: models.BuyerAttributes{
						CardNumberID: cardNumberID,
						FirstName:    firstName,
						LastName:     lastName,
					},
				},
			},
			nil,
		)
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers", nil)
		res := httptest.NewRecorder()
		expected := `{
			"data": [
				{
					"id": 1,
					"card_number_id": "12345678",
					"first_name": "John",
					"last_name": "Doe"
				}
			]
		}`
		expectedCode := 200
		hd.GetAll()(res, req)
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		mockService.On("GetByID", 999).Return(
			models.Buyer{},
			pkg.ServiceErrors[pkg.ErrNotFound],
		)
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers/999", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 404
		hd.GetByID()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})

	t.Run("find_by_id_existent", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		cardNumberID := "12345678"
		firstName := "John"
		lastName := "Doe"
		id := 1

		mockService.On("GetByID", 1).Return(
			models.Buyer{
				ID: id,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberID: cardNumberID,
					FirstName:    firstName,
					LastName:     lastName,
				},
			},
			nil,
		)
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		expected := `{
			"data": {
				"id": 1,
				"card_number_id": "12345678",
				"first_name": "John",
				"last_name": "Doe"
			}
		}`
		expectedCode := 200
		hd.GetByID()(res, req)
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("find_by_id_invalid_id", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers/invalid", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 400
		hd.GetByID()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})
}

func TestUpdateBuyer(t *testing.T) {
	t.Run("update_ok", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		cardNumberID := "87654321"
		firstName := "Jane"
		lastName := "Smith"
		id := 1

		mockService.On("Update", 1, models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardNumberID,
				FirstName:    firstName,
				LastName:     lastName,
			},
		}).Return(
			models.Buyer{
				ID: id,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberID: cardNumberID,
					FirstName:    firstName,
					LastName:     lastName,
				},
			},
			nil,
		)
		hd := NewBuyerHandler(mockService)

		body := `{
			"card_number_id": "87654321",
			"first_name": "Jane",
			"last_name": "Smith"
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/buyers/1", reqBody)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()
		expected := `{
			"data": {
				"id": 1,
				"card_number_id": "87654321",
				"first_name": "Jane",
				"last_name": "Smith"
			}
		}`
		expectedCode := 200
		hd.Update()(res, req)
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("update_non_existent", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		cardNumberID := "87654321"
		firstName := "Jane"
		lastName := "Smith"

		mockService.On("Update", 999, models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: cardNumberID,
				FirstName:    firstName,
				LastName:     lastName,
			},
		}).Return(
			models.Buyer{},
			pkg.ServiceErrors[pkg.ErrNotFound],
		)
		hd := NewBuyerHandler(mockService)

		body := `{
			"card_number_id": "87654321",
			"first_name": "Jane",
			"last_name": "Smith"
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/buyers/999", reqBody)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 404
		hd.Update()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})

	t.Run("update_invalid_id", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		hd := NewBuyerHandler(mockService)

		body := `{
			"first_name": "Jane"
		}`
		reqBody := bytes.NewReader([]byte(body))
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/buyers/invalid", reqBody)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 400
		hd.Update()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})
}

func TestDeleteBuyer(t *testing.T) {
	t.Run("delete_non_existent", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		mockService.On("Delete", 999).Return(
			pkg.ServiceErrors[pkg.ErrNotFound],
		)
		hd := NewBuyerHandler(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/buyers/999", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 404
		hd.Delete()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})

	t.Run("delete_ok", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		mockService.On("Delete", 1).Return(nil)
		hd := NewBuyerHandler(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 204
		hd.Delete()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})

	t.Run("delete_invalid_id", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		hd := NewBuyerHandler(mockService)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/buyers/invalid", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		res := httptest.NewRecorder()

		expectedCode := 400
		hd.Delete()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})
}

func TestGetPurchaseOrdersReport(t *testing.T) {
	t.Run("get_all_buyers_report", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		expectedReports := []models.BuyerPurchaseOrderReport{
			{
				ID:                  1,
				CardNumberID:        "12345678",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 5,
			},
			{
				ID:                  2,
				CardNumberID:        "87654321",
				FirstName:           "Jane",
				LastName:            "Smith",
				PurchaseOrdersCount: 3,
			},
		}

		mockService.On("GetPurchaseOrdersReport", (*int)(nil)).Return(expectedReports, nil)
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders", nil)
		res := httptest.NewRecorder()
		expected := `{
			"data": [
				{
					"id": 1,
					"card_number_id": "12345678",
					"first_name": "John",
					"last_name": "Doe",
					"purchase_orders_count": 5
				},
				{
					"id": 2,
					"card_number_id": "87654321",
					"first_name": "Jane",
					"last_name": "Smith",
					"purchase_orders_count": 3
				}
			]
		}`
		expectedCode := 200
		hd.GetPurchaseOrdersReport()(res, req)
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("get_specific_buyer_report", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		buyerID := 1
		expectedReport := []models.BuyerPurchaseOrderReport{
			{
				ID:                  1,
				CardNumberID:        "12345678",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 5,
			},
		}

		mockService.On("GetPurchaseOrdersReport", &buyerID).Return(expectedReport, nil)
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders?id=1", nil)
		res := httptest.NewRecorder()
		expected := `{
			"data": [
				{
					"id": 1,
					"card_number_id": "12345678",
					"first_name": "John",
					"last_name": "Doe",
					"purchase_orders_count": 5
				}
			]
		}`
		expectedCode := 200
		hd.GetPurchaseOrdersReport()(res, req)
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expected, res.Body.String())
	})

	t.Run("invalid_id_format", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders?id=invalid", nil)
		res := httptest.NewRecorder()

		expectedCode := 400
		hd.GetPurchaseOrdersReport()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})

	t.Run("buyer_not_found", func(t *testing.T) {
		mockService := new(buyer.MockBuyerService)
		buyerID := 999
		mockService.On("GetPurchaseOrdersReport", &buyerID).Return([]models.BuyerPurchaseOrderReport{}, pkg.ServiceErrors[pkg.ErrNotFound])
		hd := NewBuyerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/buyers/reportPurchaseOrders?id=999", nil)
		res := httptest.NewRecorder()

		expectedCode := 404
		hd.GetPurchaseOrdersReport()(res, req)
		require.Equal(t, expectedCode, res.Code)
	})
}
