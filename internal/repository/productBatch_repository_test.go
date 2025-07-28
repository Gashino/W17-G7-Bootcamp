package repository

import (
	"app/pkg"
	"app/pkg/models"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductBatchRepository_InsertProductBatch(t *testing.T) {
	t.Run("Se inserta el productBatch correctamente", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

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

		mock.ExpectExec("INSERT INTO product_batches.*VALUES.*").
			WithArgs(1, 10, 10.0, "10/10/2025", "10/10/2025", 12, 12.0, 10, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rp := NewProductBatchSqlRepository(db)

		// Act
		result, err := rp.InsertProductBatch(productBatch)

		// Assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, 1, result.BatchNumber)
		require.Equal(t, 10, result.CurrentQuantity)
		require.Equal(t, 10.0, result.CurrentTemperature)
		require.Equal(t, "10/10/2025", result.DueDate)
		require.Equal(t, 10, result.InitialQuantity)
		require.Equal(t, "10/10/2025", result.ManufacturingDate)
		require.Equal(t, 12, result.ManufacturingHour)
		require.Equal(t, 12.0, result.MinumumTemperature)
		require.Equal(t, 1, result.ProductId)
		require.Equal(t, 1, result.SectionId)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Devuelve error conflict cuando batch_number ya existe", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

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

		mock.ExpectExec("INSERT INTO product_batches.*VALUES.*").
			WithArgs(1, 10, 10.0, "10/10/2025", "10/10/2025", 12, 12.0, 10, 1, 1).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

		rp := NewProductBatchSqlRepository(db)

		// Act
		result, err := rp.InsertProductBatch(productBatch)

		// Assert
		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrConflict]
		svcErr.InternalError = fmt.Errorf("batch_number duplicado")
		require.Equal(t, svcErr, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Devuelve error conflict cuando product_id o section_id no existen", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

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
				ProductId:          999,
				SectionId:          999,
			},
		}

		mock.ExpectExec("INSERT INTO product_batches.*VALUES.*").
			WithArgs(1, 10, 10.0, "10/10/2025", "10/10/2025", 12, 12.0, 10, 999, 999).
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row"})

		rp := NewProductBatchSqlRepository(db)

		// Act
		result, err := rp.InsertProductBatch(productBatch)

		// Assert
		require.Error(t, err)
		svcErr := pkg.ServiceErrors[pkg.ErrConflict]
		svcErr.InternalError = fmt.Errorf("product_id o section_id no existen")
		require.Equal(t, svcErr, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Devuelve error interno si falla la query", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

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

		mock.ExpectExec("INSERT INTO product_batches.*VALUES.*").
			WithArgs(1, 10, 10.0, "10/10/2025", "10/10/2025", 12, 12.0, 10, 1, 1).
			WillReturnError(assert.AnError)

		rp := NewProductBatchSqlRepository(db)

		// Act
		result, err := rp.InsertProductBatch(productBatch)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Devuelve error interno si falla LastInsertId", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

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

		mock.ExpectExec("INSERT INTO product_batches.*VALUES.*").
			WithArgs(1, 10, 10.0, "10/10/2025", "10/10/2025", 12, 12.0, 10, 1, 1).
			WillReturnResult(sqlmock.NewErrorResult(assert.AnError))

		rp := NewProductBatchSqlRepository(db)

		// Act
		result, err := rp.InsertProductBatch(productBatch)

		// Assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
