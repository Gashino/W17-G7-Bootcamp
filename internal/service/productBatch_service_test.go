package service

import (
	"app/pkg"
	"app/pkg/models"
	productBatchMock "app/test/productBatch"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProductBatchService_PostProductBatch(t *testing.T) {
	t.Run("Crea correctamente un productBatch", func(t *testing.T) {
		// Arrange
		rp := new(productBatchMock.MockProductBatchRepository)
		sv := NewProductBatchDefault(rp)
		productBatch := models.ProductBatch{
			ID: 0,
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        1,
				CurrentQuantity:    10,
				CurrentTemperature: 10.0,
				DueDate:            "10/10/2025",
				InitialQuantity:    10,
				ManufacturingDate:  "10/10/2025",
				ManufacturingHour:  12,
				MinumumTemperature: 12.0,
				ProductId:          1,
				SectionId:          1,
			},
		}
		expectedProductBatch := models.ProductBatch{
			ID: 1,
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        1,
				CurrentQuantity:    10,
				CurrentTemperature: 10.0,
				DueDate:            "10/10/2025",
				InitialQuantity:    10,
				ManufacturingDate:  "10/10/2025",
				ManufacturingHour:  12,
				MinumumTemperature: 12.0,
				ProductId:          1,
				SectionId:          1,
			},
		}
		rp.On("InsertProductBatch", productBatch).Return(expectedProductBatch, nil)

		// Act
		result, err := sv.PostProductBatch(productBatch)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedProductBatch, result)
		rp.AssertCalled(t, "InsertProductBatch", productBatch)
	})

	t.Run("Devuelve error cuando el repository falla", func(t *testing.T) {
		// Arrange
		rp := new(productBatchMock.MockProductBatchRepository)
		sv := NewProductBatchDefault(rp)
		productBatch := models.ProductBatch{
			ID: 0,
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        1,
				CurrentQuantity:    10,
				CurrentTemperature: 10.0,
				DueDate:            "10/10/2025",
				InitialQuantity:    10,
				ManufacturingDate:  "10/10/2025",
				ManufacturingHour:  12,
				MinumumTemperature: 12.0,
				ProductId:          1,
				SectionId:          1,
			},
		}
		expectedError := fmt.Errorf("repository error")
		rp.On("InsertProductBatch", productBatch).Return(models.ProductBatch{}, expectedError)

		// Act
		result, err := sv.PostProductBatch(productBatch)

		// Assert
		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Empty(t, result)
		rp.AssertCalled(t, "InsertProductBatch", productBatch)
	})

	t.Run("Devuelve error de conflicto cuando hay batch_number duplicado", func(t *testing.T) {
		// Arrange
		rp := new(productBatchMock.MockProductBatchRepository)
		sv := NewProductBatchDefault(rp)
		productBatch := models.ProductBatch{
			ID: 0,
			ProductBatchAttributes: models.ProductBatchAttributes{
				BatchNumber:        1,
				CurrentQuantity:    10,
				CurrentTemperature: 10.0,
				DueDate:            "10/10/2025",
				InitialQuantity:    10,
				ManufacturingDate:  "10/10/2025",
				ManufacturingHour:  12,
				MinumumTemperature: 12.0,
				ProductId:          1,
				SectionId:          1,
			},
		}
		svcErr := pkg.ServiceErrors[pkg.ErrConflict]
		svcErr.InternalError = fmt.Errorf("batch_number ya existe")
		rp.On("InsertProductBatch", productBatch).Return(models.ProductBatch{}, svcErr)

		// Act
		result, err := sv.PostProductBatch(productBatch)

		// Assert
		require.Error(t, err)
		require.Equal(t, svcErr, err)
		require.Empty(t, result)
		rp.AssertCalled(t, "InsertProductBatch", productBatch)
	})
}
