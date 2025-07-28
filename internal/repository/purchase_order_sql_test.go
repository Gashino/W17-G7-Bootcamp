package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestCreatePurchaseOrder(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPurchaseOrderSQL(db)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#2",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		expectedID := int64(2)
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#2", "2021-04-05", "xyz789", 1, 2).
			WillReturnResult(sqlmock.NewResult(expectedID, 1))

		// act
		result, err := repo.Create(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, 2, result.ID)
		require.Equal(t, "order#2", result.OrderNumber)
		require.Equal(t, "2021-04-05", result.OrderDate)
		require.Equal(t, "xyz789", result.TrackingCode)
		require.Equal(t, 1, result.BuyerID)
		require.Equal(t, 2, result.ProductRecordID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create_conflict_duplicate_order_number", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPurchaseOrderSQL(db)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#1",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		// Simulate MySQL duplicate entry error (1062)
		mysqlErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'order#1' for key 'order_number'"}
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#1", "2021-04-05", "xyz789", 1, 2).
			WillReturnError(mysqlErr)

		// act
		result, err := repo.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.PurchaseOrder{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create_conflict_foreign_key_buyer_id", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPurchaseOrderSQL(db)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#3",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         999, // non-existent buyer
				ProductRecordID: 2,
			},
		}

		// Simulate MySQL foreign key constraint error (1452)
		mysqlErr := &mysql.MySQLError{Number: 1452, Message: "Cannot add or update a child row: a foreign key constraint fails"}
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#3", "2021-04-05", "xyz789", 999, 2).
			WillReturnError(mysqlErr)

		// act
		result, err := repo.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrBadRequest].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.PurchaseOrder{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create_mysql_error_other", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPurchaseOrderSQL(db)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#4",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		// Simulate other MySQL error
		mysqlErr := &mysql.MySQLError{Number: 1054, Message: "Unknown column 'unknown_column' in 'field list'"}
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#4", "2021-04-05", "xyz789", 1, 2).
			WillReturnError(mysqlErr)

		// act
		result, err := repo.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.PurchaseOrder{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create_non_mysql_error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPurchaseOrderSQL(db)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#5",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		// Simulate non-MySQL error
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#5", "2021-04-05", "xyz789", 1, 2).
			WillReturnError(sql.ErrConnDone)

		// act
		result, err := repo.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.PurchaseOrder{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create_last_insert_id_error", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPurchaseOrderSQL(db)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#6",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		// Simulate successful insert but failed LastInsertId
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#6", "2021-04-05", "xyz789", 1, 2).
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		// act
		result, err := repo.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.PurchaseOrder{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
