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

func TestCreateBuyer(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		input := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "12345678",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		expectedID := int64(1)
		mock.ExpectExec("INSERT INTO buyers").
			WithArgs("12345678", "John", "Doe").
			WillReturnResult(sqlmock.NewResult(expectedID, 1))

		// act
		result, err := repo.Create(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "12345678", result.CardNumberID)
		require.Equal(t, "John", result.FirstName)
		require.Equal(t, "Doe", result.LastName)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create_conflict", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		input := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "12345678",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		// Simulate MySQL duplicate entry error
		mysqlErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
		mock.ExpectExec("INSERT INTO buyers").
			WithArgs("12345678", "John", "Doe").
			WillReturnError(mysqlErr)

		// act
		result, err := repo.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.Buyer{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetAllBuyers(t *testing.T) {
	t.Run("find_all", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
			AddRow(1, "12345678", "John", "Doe").
			AddRow(2, "87654321", "Jane", "Smith")

		mock.ExpectQuery("SELECT id, card_number_id, first_name, last_name FROM buyers").
			WillReturnRows(rows)

		// act
		result, err := repo.GetAll()

		// assert
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, "John", result[1].FirstName)
		require.Equal(t, "Jane", result[2].FirstName)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetBuyerByID(t *testing.T) {
	t.Run("find_by_id_existent", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		buyerID := 1
		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
			AddRow(1, "12345678", "John", "Doe")

		mock.ExpectQuery("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?").
			WithArgs(buyerID).
			WillReturnRows(rows)

		// act
		result, err := repo.GetByID(buyerID)

		// assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "12345678", result.CardNumberID)
		require.Equal(t, "John", result.FirstName)
		require.Equal(t, "Doe", result.LastName)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		buyerID := 999

		mock.ExpectQuery("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?").
			WithArgs(buyerID).
			WillReturnError(sql.ErrNoRows)

		// act
		result, err := repo.GetByID(buyerID)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.Buyer{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateBuyer(t *testing.T) {
	t.Run("update_existent", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		input := models.Buyer{
			ID: 1,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		// Mock GetByID call first (to check if buyer exists)
		existingRows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
			AddRow(1, "12345678", "John", "Doe")
		mock.ExpectQuery("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?").
			WithArgs(1).
			WillReturnRows(existingRows)

		// Mock the update operation
		mock.ExpectExec("UPDATE buyers SET.*WHERE id = ?").
			WithArgs("87654321", "Jane", "Smith", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock GetByID call after update
		updatedRows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
			AddRow(1, "87654321", "Jane", "Smith")
		mock.ExpectQuery("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?").
			WithArgs(1).
			WillReturnRows(updatedRows)

		// act
		result, err := repo.Update(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "87654321", result.CardNumberID)
		require.Equal(t, "Jane", result.FirstName)
		require.Equal(t, "Smith", result.LastName)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update_non_existent", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		input := models.Buyer{
			ID: 999,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		// Mock GetByID call (buyer doesn't exist)
		mock.ExpectQuery("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		// act
		result, err := repo.Update(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Equal(t, models.Buyer{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteBuyer(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		buyerID := 1

		mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
			WithArgs(buyerID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// act
		err = repo.Delete(buyerID)

		// assert
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("delete_non_existent", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		buyerID := 999

		mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
			WithArgs(buyerID).
			WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected

		// act
		err = repo.Delete(buyerID)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetPurchaseOrdersReport(t *testing.T) {
	t.Run("get_all_buyers_report", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(1, "12345678", "John", "Doe", 5).
			AddRow(2, "87654321", "Jane", "Smith", 3)

		mock.ExpectQuery("SELECT.*FROM buyers b.*LEFT JOIN purchase_orders po.*GROUP BY.*ORDER BY").
			WillReturnRows(rows)

		// act
		result, err := repo.GetPurchaseOrdersReport(nil)

		// assert
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, 1, result[0].ID)
		require.Equal(t, "John", result[0].FirstName)
		require.Equal(t, 5, result[0].PurchaseOrdersCount)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("get_specific_buyer_report", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		buyerID := 1
		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(1, "12345678", "John", "Doe", 5)

		mock.ExpectQuery("SELECT.*FROM buyers b.*LEFT JOIN purchase_orders po.*WHERE b.id = ?.*GROUP BY").
			WithArgs(buyerID).
			WillReturnRows(rows)

		// act
		result, err := repo.GetPurchaseOrdersReport(&buyerID)

		// assert
		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, 1, result[0].ID)
		require.Equal(t, "John", result[0].FirstName)
		require.Equal(t, 5, result[0].PurchaseOrdersCount)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("get_report_buyer_not_found", func(t *testing.T) {
		// arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)
		buyerID := 999
		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"})

		mock.ExpectQuery("SELECT.*FROM buyers b.*LEFT JOIN purchase_orders po.*WHERE b.id = ?.*GROUP BY").
			WithArgs(buyerID).
			WillReturnRows(rows)

		// act
		result, err := repo.GetPurchaseOrdersReport(&buyerID)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, err.(pkg.ServiceError).ResponseCode)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
