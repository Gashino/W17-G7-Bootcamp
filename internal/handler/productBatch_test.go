package handler

import (
	"app/pkg"
	"app/pkg/models"
	productbatch "app/test/productBatch"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProductBatchDefault_CreateBatch(t *testing.T) {
	t.Run("Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado", func(t *testing.T) {
		// Arrange
		mockService := new(productbatch.MockProductBatchService)
		requestBody := `{"batch_number": 123, "current_quantity": 100, "current_temperature": 5.0, "due_date": "2024-12-31", "initial_quantity": 150, "manufacturing_date": "2024-01-01", "manufacturing_hour": 8, "minumum_temperature": 2.0, "product_id": 1, "section_id": 1}`
		expectedBatch := models.ProductBatch{
			ID: 1,
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}
		mockService.On("PostProductBatch", mock.AnythingOfType("models.ProductBatch")).Return(expectedBatch, nil)
		hd := NewProductBatchDefault(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", strings.NewReader(requestBody))
		res := httptest.NewRecorder()
		// Act
		hd.CreateBatch()(res, req)
		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{"data":{"id":1,"batch_number":123,"current_quantity":100,"current_temperature":5,"due_date":"2024-12-31","initial_quantity":150,"manufacturing_date":"2024-01-01","manufacturing_hour":8,"minumum_temperature":2,"product_id":1,"section_id":1}}`
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "PostProductBatch", mock.AnythingOfType("models.ProductBatch"))
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el JSON es inválido se devolverá un código 400", func(t *testing.T) {
		// Arrange
		mockService := new(productbatch.MockProductBatchService)
		requestBody := `{"invalid json`

		hd := NewProductBatchDefault(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		// Act
		hd.CreateBatch()(res, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{
							"status": "Bad Request",
							"message": "error: Bad request"
						}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si la validación falla con error service se devuelve el código del error", func(t *testing.T) {
		// Arrange
		mockService := new(productbatch.MockProductBatchService)
		requestBody := `{"batch_number": 0, "current_quantity": 100, "current_temperature": 5.0, "due_date": "2024-12-31", "initial_quantity": 150, "manufacturing_date": "2024-01-01", "manufacturing_hour": 8, "minumum_temperature": 2.0, "product_id": 1, "section_id": 1}`

		hd := NewProductBatchDefault(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		// Act
		hd.CreateBatch()(res, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{
							"status": "Unprocessable Entity",
							"message": "error: product batch is not valid"
						}`
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el service devuelve error service se devuelve el código del error", func(t *testing.T) {
		// Arrange
		mockService := new(productbatch.MockProductBatchService)
		requestBody := `{"batch_number": 123, "current_quantity": 100, "current_temperature": 5.0, "due_date": "2024-12-31", "initial_quantity": 150, "manufacturing_date": "2024-01-01", "manufacturing_hour": 8, "minumum_temperature": 2.0, "product_id": 1, "section_id": 1}`
		svcErr := pkg.ServiceErrors[pkg.ErrConflict]
		svcErr.InternalError = fmt.Errorf("batch_number ya existe")
		mockService.On("PostProductBatch", mock.AnythingOfType("models.ProductBatch")).Return(models.ProductBatch{}, svcErr)

		hd := NewProductBatchDefault(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		// Act
		hd.CreateBatch()(res, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{
							"status": "Conflict",
							"message": "error: batch_number ya existe"
						}`
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "PostProductBatch", mock.AnythingOfType("models.ProductBatch"))
		require.JSONEq(t, expectedBody, res.Body.String())
	})

	t.Run("Si el service devuelve error no-service se devuelve código 500", func(t *testing.T) {
		// Arrange
		mockService := new(productbatch.MockProductBatchService)
		requestBody := `{"batch_number": 123, "current_quantity": 100, "current_temperature": 5.0, "due_date": "2024-12-31", "initial_quantity": 150, "manufacturing_date": "2024-01-01", "manufacturing_hour": 8, "minumum_temperature": 2.0, "product_id": 1, "section_id": 1}`
		mockService.On("PostProductBatch", mock.AnythingOfType("models.ProductBatch")).Return(models.ProductBatch{}, fmt.Errorf("some internal error"))

		hd := NewProductBatchDefault(mockService)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		// Act
		hd.CreateBatch()(res, req)

		// Assert
		expectedCode := http.StatusInternalServerError
		expectedBody := `{
							"status": "Internal Server Error",
							"message": "error: Internal server error"
						}`
		require.Equal(t, expectedCode, res.Code)
		mockService.AssertCalled(t, "PostProductBatch", mock.AnythingOfType("models.ProductBatch"))
		require.JSONEq(t, expectedBody, res.Body.String())
	})
}

func TestProductBatchDefault_ValidatePostProductBatch(t *testing.T) {
	hd := &ProductBatchDefault{}

	t.Run("Devuelve error cuando batch_number es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        0,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando current_quantity es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    0,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando current_temperature es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 0.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando due_date es vacío", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando initial_quantity es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    0,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando manufacturing_date es vacío", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando manufacturing_hour es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  0,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando minimum_temperature es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 0.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando product_id es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          0,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando section_id es 0", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          0,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve error cuando current_temperature < minimum_temperature", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 1.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]
		svcErr.InternalError = fmt.Errorf("product batch is not valid")
		require.Equal(t, svcErr, err)
	})

	t.Run("Devuelve nil cuando el product batch es válido", func(t *testing.T) {
		batch := models.ProductBatch{
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        123,
				CurrentQuantity:    100,
				CurrentTemperature: 5.0,
				DueDate:            "2024-12-31",
				InitialQuantity:    150,
				ManufacturingDate:  "2024-01-01",
				ManufacturingHour:  8,
				MinumumTemperature: 2.0,
				ProductId:          1,
				SectionId:          1,
			},
		}

		err := hd.ValidatePostProductBatch(batch)

		require.NoError(t, err)
	})
}
