package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestInboundOrderSQL_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		orderDate := time.Now()
		inboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "IO001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		mock.ExpectExec("INSERT INTO inbound_orders").WithArgs(
			inboundOrder.OrderDate,
			inboundOrder.OrderNumber,
			inboundOrder.EmployeeID,
			inboundOrder.ProductBatchID,
			inboundOrder.WarehouseID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewInboundOrderSQL(db)

		// Act
		createdOrder, err := repo.Create(inboundOrder)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, createdOrder)
		require.Equal(t, inboundOrder.OrderNumber, createdOrder.OrderNumber)
		require.Equal(t, inboundOrder.EmployeeID, createdOrder.EmployeeID)
		require.Equal(t, inboundOrder.ProductBatchID, createdOrder.ProductBatchID)
		require.Equal(t, inboundOrder.WarehouseID, createdOrder.WarehouseID)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_employee_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		orderDate := time.Now()
		inboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "IO001",
			EmployeeID:     999, // Invalid employee ID
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`frescos`.`inbound_orders`, CONSTRAINT `fk_inbound_orders_employee_id` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`))",
		}

		mock.ExpectExec("INSERT INTO inbound_orders").WithArgs(
			inboundOrder.OrderDate,
			inboundOrder.OrderNumber,
			inboundOrder.EmployeeID,
			inboundOrder.ProductBatchID,
			inboundOrder.WarehouseID,
		).WillReturnError(mysqlErr)

		repo := NewInboundOrderSQL(db)

		// Act
		createdOrder, err := repo.Create(inboundOrder)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.InboundOrder{}, createdOrder)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_product_batch_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		orderDate := time.Now()
		inboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "IO001",
			EmployeeID:     1,
			ProductBatchID: 999, // Invalid product batch ID
			WarehouseID:    1,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`frescos`.`inbound_orders`, CONSTRAINT `fk_inbound_orders_product_batch_id` FOREIGN KEY (`product_batch_id`) REFERENCES `product_batches` (`id`))",
		}

		mock.ExpectExec("INSERT INTO inbound_orders").WithArgs(
			inboundOrder.OrderDate,
			inboundOrder.OrderNumber,
			inboundOrder.EmployeeID,
			inboundOrder.ProductBatchID,
			inboundOrder.WarehouseID,
		).WillReturnError(mysqlErr)

		repo := NewInboundOrderSQL(db)

		// Act
		createdOrder, err := repo.Create(inboundOrder)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.InboundOrder{}, createdOrder)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_warehouse_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		orderDate := time.Now()
		inboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "IO001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    999, // Invalid warehouse ID
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`frescos`.`inbound_orders`, CONSTRAINT `fk_inbound_orders_warehouse_id` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`))",
		}

		mock.ExpectExec("INSERT INTO inbound_orders").WithArgs(
			inboundOrder.OrderDate,
			inboundOrder.OrderNumber,
			inboundOrder.EmployeeID,
			inboundOrder.ProductBatchID,
			inboundOrder.WarehouseID,
		).WillReturnError(mysqlErr)

		repo := NewInboundOrderSQL(db)

		// Act
		createdOrder, err := repo.Create(inboundOrder)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.InboundOrder{}, createdOrder)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_order_number", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		orderDate := time.Now()
		inboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "IO001", // Duplicate order number
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'IO001' for key 'order_number'",
		}

		mock.ExpectExec("INSERT INTO inbound_orders").WithArgs(
			inboundOrder.OrderDate,
			inboundOrder.OrderNumber,
			inboundOrder.EmployeeID,
			inboundOrder.ProductBatchID,
			inboundOrder.WarehouseID,
		).WillReturnError(mysqlErr)

		repo := NewInboundOrderSQL(db)

		// Act
		createdOrder, err := repo.Create(inboundOrder)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.InboundOrder{}, createdOrder)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		orderDate := time.Now()
		inboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "IO001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		mock.ExpectExec("INSERT INTO inbound_orders").WithArgs(
			inboundOrder.OrderDate,
			inboundOrder.OrderNumber,
			inboundOrder.EmployeeID,
			inboundOrder.ProductBatchID,
			inboundOrder.WarehouseID,
		).WillReturnError(sql.ErrConnDone)

		repo := NewInboundOrderSQL(db)

		// Act
		createdOrder, err := repo.Create(inboundOrder)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.InboundOrder{}, createdOrder)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
