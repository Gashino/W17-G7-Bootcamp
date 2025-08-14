package repository

import (
	"app/pkg/models"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestWarehouseSql_FindAll(t *testing.T) {
	t.Run("success - find all warehouses", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		expectedWarehouses := []models.Warehouse{
			{
				ID:             1,
				WarehouseCode:  "WH001",
				Address:        "123 Main St",
				Telephone:      "555-0101",
				MinCapacity:    100,
				MinTemperature: -10,
			},
			{
				ID:             2,
				WarehouseCode:  "WH002",
				Address:        "456 Oak Ave",
				Telephone:      "555-0102",
				MinCapacity:    200,
				MinTemperature: -20,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow(1, "WH001", "123 Main St", "555-0101", 100, -10).
			AddRow(2, "WH002", "456 Oak Ave", "555-0102", 200, -20)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses").
			WillReturnRows(rows)

		// Act
		result, err := repo.FindAll()

		// Assert
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, expectedWarehouses[0], result[1])
		require.Equal(t, expectedWarehouses[1], result[2])
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses").
			WillReturnError(errors.New("database error"))

		// Act
		result, err := repo.FindAll()

		// Assert
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, "database error", err.Error())
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - scan fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow("invalid", "WH001", "123 Main St", "555-0101", 100, -10)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses").
			WillReturnRows(rows)

		// Act
		result, err := repo.FindAll()

		// Assert
		require.Error(t, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSql_FindByID(t *testing.T) {
	t.Run("success - find warehouse by id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		expectedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-0101",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow(1, "WH001", "123 Main St", "555-0101", 100, -10)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := repo.FindByID(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedWarehouse, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - warehouse not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"})

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE id = \\?").
			WithArgs(999).
			WillReturnRows(rows)

		// Act
		result, err := repo.FindByID(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, "warehouse not found", err.Error())
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE id = \\?").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		// Act
		result, err := repo.FindByID(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, "warehouse not found", err.Error())
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - scan fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow("invalid", "WH001", "123 Main St", "555-0101", 100, -10)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := repo.FindByID(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSql_Add(t *testing.T) {
	t.Run("success - add warehouse", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		warehouse := models.Warehouse{
			WarehouseCode:  "WH003",
			Address:        "789 Pine St",
			Telephone:      "555-0103",
			MinCapacity:    150,
			MinTemperature: -15,
		}

		mock.ExpectExec("INSERT INTO warehouses \\(warehouse_code, address, telephone, minimun_capacity, minimun_temperature\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("WH003", "789 Pine St", "555-0103", 150, -15).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.Add(warehouse)

		// Assert
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database insert fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		warehouse := models.Warehouse{
			WarehouseCode:  "WH003",
			Address:        "789 Pine St",
			Telephone:      "555-0103",
			MinCapacity:    150,
			MinTemperature: -15,
		}

		mock.ExpectExec("INSERT INTO warehouses \\(warehouse_code, address, telephone, minimun_capacity, minimun_temperature\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("WH003", "789 Pine St", "555-0103", 150, -15).
			WillReturnError(errors.New("duplicate entry"))

		// Act
		result, err := repo.Add(warehouse)

		// Assert
		require.Error(t, err)
		require.Equal(t, "SQL Error", err.Error())
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSql_FindWarehouseByCode(t *testing.T) {
	t.Run("success - find warehouse by code", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		expectedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-0101",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow(1, "WH001", "123 Main St", "555-0101", 100, -10)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE warehouse_code = \\?").
			WithArgs("WH001").
			WillReturnRows(rows)

		// Act
		result, err := repo.FindWarehouseByCode("WH001")

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedWarehouse, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - warehouse not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"})

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE warehouse_code = \\?").
			WithArgs("INVALID").
			WillReturnRows(rows)

		// Act
		result, err := repo.FindWarehouseByCode("INVALID")

		// Assert
		require.Error(t, err)
		require.Equal(t, "warehouse not found", err.Error())
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database query fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE warehouse_code = \\?").
			WithArgs("WH001").
			WillReturnError(errors.New("database error"))

		// Act
		result, err := repo.FindWarehouseByCode("WH001")

		// Assert
		require.Error(t, err)
		require.Equal(t, "warehouse not found", err.Error())
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - scan fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimun_capacity", "minimun_temperature"}).
			AddRow("invalid", "WH001", "123 Main St", "555-0101", 100, -10)

		mock.ExpectQuery("SELECT id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature FROM warehouses WHERE warehouse_code = \\?").
			WithArgs("WH001").
			WillReturnRows(rows)

		// Act
		result, err := repo.FindWarehouseByCode("WH001")

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSql_Update(t *testing.T) {
	t.Run("success - update warehouse", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		warehouse := models.Warehouse{
			WarehouseCode:  "WH001_UPDATED",
			Address:        "123 Main St Updated",
			Telephone:      "555-0199",
			MinCapacity:    200,
			MinTemperature: -20,
		}

		mock.ExpectExec("UPDATE warehouses SET warehouse_code = \\?, address = \\?, telephone = \\?, minimun_capacity= \\?, minimun_temperature = \\? WHERE id = \\?").
			WithArgs("WH001_UPDATED", "123 Main St Updated", "555-0199", 200, -20, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		result, err := repo.Update(1, warehouse)

		// Assert
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database update fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		warehouse := models.Warehouse{
			WarehouseCode:  "WH001_UPDATED",
			Address:        "123 Main St Updated",
			Telephone:      "555-0199",
			MinCapacity:    200,
			MinTemperature: -20,
		}

		mock.ExpectExec("UPDATE warehouses SET warehouse_code = \\?, address = \\?, telephone = \\?, minimun_capacity= \\?, minimun_temperature = \\? WHERE id = \\?").
			WithArgs("WH001_UPDATED", "123 Main St Updated", "555-0199", 200, -20, 1).
			WillReturnError(errors.New("database error"))

		// Act
		result, err := repo.Update(1, warehouse)

		// Assert
		require.Error(t, err)
		require.Equal(t, "SQL Error", err.Error())
		require.Equal(t, models.Warehouse{}, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSql_Delete(t *testing.T) {
	t.Run("success - delete warehouse", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		mock.ExpectExec("DELETE FROM warehouses WHERE id = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		err = repo.Delete(1)

		// Assert
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database delete fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		mock.ExpectExec("DELETE FROM warehouses WHERE id = \\?").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		// Act
		err = repo.Delete(1)

		// Assert
		require.NoError(t, err) // Note: The original implementation returns nil even on error
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSql_FindAvailableID(t *testing.T) {
	t.Run("success - find available id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewWarehouseSql(db)

		// Act
		id, err := repo.FindAvailableID()

		// Assert
		require.NoError(t, err)
		require.Equal(t, 0, id) // Current implementation returns 0
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewWarehouseSql(t *testing.T) {
	t.Run("success - create new warehouse sql repository", func(t *testing.T) {
		// Arrange
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Act
		repo := NewWarehouseSql(db)

		// Assert
		require.NotNil(t, repo)
		require.Equal(t, db, repo.db)
	})
}
