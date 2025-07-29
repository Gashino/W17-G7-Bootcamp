package repository

import (
	"app/pkg"
	"app/pkg/models"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestInsertProductRecord(t *testing.T) {
	t.Run("insert_ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewProductRecordSqlRepository(db)
		now := time.Now()
		input := models.ProductRecord{
			LastUpdateDate: &now,
			PurchasePrice:  &[]float64{10.50}[0],
			SalePrice:      &[]float64{15.99}[0],
			ProductId:      &[]int{1}[0],
		}

		expectedID := int64(1)
		mock.ExpectExec("INSERT INTO product_records").
			WithArgs(input.LastUpdateDate, input.PurchasePrice, input.SalePrice, input.ProductId).
			WillReturnResult(sqlmock.NewResult(expectedID, 1))

		// act
		result, err := repo.Insert(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, input.LastUpdateDate, result.LastUpdateDate)
		require.Equal(t, input.PurchasePrice, result.PurchasePrice)
		require.Equal(t, input.SalePrice, result.SalePrice)
		require.Equal(t, input.ProductId, result.ProductId)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insert_exec_error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewProductRecordSqlRepository(db)
		now := time.Now()
		input := models.ProductRecord{
			LastUpdateDate: &now,
			PurchasePrice:  &[]float64{10.50}[0],
			SalePrice:      &[]float64{15.99}[0],
			ProductId:      &[]int{1}[0],
		}

		mock.ExpectExec("INSERT INTO product_records").
			WithArgs(input.LastUpdateDate, input.PurchasePrice, input.SalePrice, input.ProductId).
			WillReturnError(errors.New("database error"))

		// act
		result, err := repo.Insert(input)

		// assert
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Contains(t, err.Error(), "invalid product_id")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insert_last_insert_id_error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewProductRecordSqlRepository(db)
		now := time.Now()
		input := models.ProductRecord{
			LastUpdateDate: &now,
			PurchasePrice:  &[]float64{10.50}[0],
			SalePrice:      &[]float64{15.99}[0],
			ProductId:      &[]int{1}[0],
		}

		mock.ExpectExec("INSERT INTO product_records").
			WithArgs(input.LastUpdateDate, input.PurchasePrice, input.SalePrice, input.ProductId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId error")))

		// act
		result, err := repo.Insert(input)

		// assert
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewProductRecordSqlRepository(t *testing.T) {
	t.Run("new_product_record_sql_constructor", func(t *testing.T) {
		// arrange
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// act
		repo := NewProductRecordSqlRepository(db)

		// assert
		require.NotNil(t, repo)
		require.IsType(t, &ProductRecordSql{}, repo)
	})
}
