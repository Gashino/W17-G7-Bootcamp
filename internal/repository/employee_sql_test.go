package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestEmployeeRepositoryMap_FindAll(t *testing.T) {
	t.Run("success_with_multiple_employees", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id",
		}).
			AddRow(1, "E001", "John", "Doe", 1).
			AddRow(2, "E002", "Jane", "Smith", 2)

		mock.ExpectQuery("SELECT (.+) FROM employees").WillReturnRows(rows)

		repo := NewEmployeeRepository(db)

		// Act
		employees, err := repo.FindAll()

		// Assert
		require.NoError(t, err)
		require.NotNil(t, employees)
		require.Equal(t, 2, len(employees))

		// Check first employee
		employee1, exists := employees[1]
		require.True(t, exists)
		require.Equal(t, 1, *employee1.ID)
		require.Equal(t, "E001", *employee1.CardNumberID)
		require.Equal(t, "John", *employee1.FirstName)
		require.Equal(t, "Doe", *employee1.LastName)
		require.Equal(t, 1, *employee1.WarehouseID)

		// Check second employee
		employee2, exists := employees[2]
		require.True(t, exists)
		require.Equal(t, 2, *employee2.ID)
		require.Equal(t, "E002", *employee2.CardNumberID)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM employees").WillReturnError(sql.ErrConnDone)

		repo := NewEmployeeRepository(db)

		// Act
		employees, err := repo.FindAll()

		// Assert
		require.Error(t, err)
		require.Nil(t, employees)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepositoryMap_FindById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id",
		}).AddRow(employeeId, "E001", "John", "Doe", 1)

		mock.ExpectQuery("SELECT (.+) FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnRows(rows)

		repo := NewEmployeeRepository(db)

		// Act
		employee, err := repo.FindById(employeeId)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, employee)
		require.Equal(t, employeeId, *employee.ID)
		require.Equal(t, "E001", *employee.CardNumberID)
		require.Equal(t, "John", *employee.FirstName)
		require.Equal(t, "Doe", *employee.LastName)
		require.Equal(t, 1, *employee.WarehouseID)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 999
		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id",
		})

		mock.ExpectQuery("SELECT (.+) FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnRows(rows)

		repo := NewEmployeeRepository(db)

		// Act
		employee, err := repo.FindById(employeeId)

		// Assert
		require.Error(t, err)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, serviceErr.ResponseCode)
		require.Equal(t, models.Employee{}, employee)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		mock.ExpectQuery("SELECT (.+) FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnError(sql.ErrConnDone)

		repo := NewEmployeeRepository(db)

		// Act
		employee, err := repo.FindById(employeeId)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, employee)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepositoryMap_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		cardNumberID := "E001"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		mock.ExpectExec("INSERT INTO employees").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewEmployeeRepository(db)

		// Act
		savedEmployee, err := repo.Save(employee)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, savedEmployee)
		require.Equal(t, *employee.CardNumberID, *savedEmployee.CardNumberID)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_card_number_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		cardNumberID := "E001"
		employee := models.Employee{
			CardNumberID: &cardNumberID,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'E001' for key 'card_number_id'",
		}

		mock.ExpectExec("INSERT INTO employees").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
		).WillReturnError(mysqlErr)

		repo := NewEmployeeRepository(db)

		// Act
		savedEmployee, err := repo.Save(employee)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, savedEmployee)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_warehouse_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		warehouseID := 999
		employee := models.Employee{
			WarehouseID: &warehouseID,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails",
		}

		mock.ExpectExec("INSERT INTO employees").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
		).WillReturnError(mysqlErr)

		repo := NewEmployeeRepository(db)

		// Act
		savedEmployee, err := repo.Save(employee)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, savedEmployee)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("mysql_error_other", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		cardNumberID := "E001"
		employee := models.Employee{
			CardNumberID: &cardNumberID,
		}

		// MySQL error diferente a 1452 y 1062
		mysqlErr := &mysql.MySQLError{
			Number:  1146,
			Message: "Table doesn't exist",
		}

		mock.ExpectExec("INSERT INTO employees").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
		).WillReturnError(mysqlErr)

		repo := NewEmployeeRepository(db)

		// Act
		savedEmployee, err := repo.Save(employee)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, savedEmployee)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("non_mysql_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		cardNumberID := "E001"
		employee := models.Employee{
			CardNumberID: &cardNumberID,
		}

		mock.ExpectExec("INSERT INTO employees").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
		).WillReturnError(fmt.Errorf("generic database error"))

		repo := NewEmployeeRepository(db)

		// Act
		savedEmployee, err := repo.Save(employee)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, savedEmployee)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("last_insert_id_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		cardNumberID := "E001"
		employee := models.Employee{
			CardNumberID: &cardNumberID,
		}

		mock.ExpectExec("INSERT INTO employees").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
		).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("last insert id error")))

		repo := NewEmployeeRepository(db)

		// Act
		savedEmployee, err := repo.Save(employee)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, savedEmployee)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepositoryMap_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		mock.ExpectExec("DELETE FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnResult(sqlmock.NewResult(0, 1))

		repo := NewEmployeeRepository(db)

		// Act
		err = repo.Delete(employeeId)

		// Assert
		require.NoError(t, err)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		mock.ExpectExec("DELETE FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnError(sql.ErrConnDone)

		repo := NewEmployeeRepository(db)

		// Act
		err = repo.Delete(employeeId)

		// Assert
		require.Error(t, err)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepositoryMap_ReportInboundOrdersCountByEmployee(t *testing.T) {
	t.Run("success_with_specific_employee", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "inbound_orders_count",
		}).AddRow(employeeId, "E001", "John", "Doe", 5)

		mock.ExpectQuery("SELECT (.+) FROM employees e LEFT JOIN inbound_orders").WithArgs(employeeId).WillReturnRows(rows)

		repo := NewEmployeeRepository(db)

		// Act
		reports, err := repo.ReportInboundOrdersCountByEmployee(&employeeId)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, reports)
		require.Equal(t, 1, len(reports))
		require.Equal(t, employeeId, reports[0].ID)
		require.Equal(t, 5, reports[0].InboundOrdersCount)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_with_all_employees", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "inbound_orders_count",
		}).
			AddRow(1, "E001", "John", "Doe", 5).
			AddRow(2, "E002", "Jane", "Smith", 3)

		mock.ExpectQuery("SELECT (.+) FROM employees e LEFT JOIN inbound_orders").WillReturnRows(rows)

		repo := NewEmployeeRepository(db)

		// Act
		reports, err := repo.ReportInboundOrdersCountByEmployee(nil)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, reports)
		require.Equal(t, 2, len(reports))

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("employee_not_found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 999
		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "inbound_orders_count",
		})

		mock.ExpectQuery("SELECT (.+) FROM employees e LEFT JOIN inbound_orders").WithArgs(employeeId).WillReturnRows(rows)

		repo := NewEmployeeRepository(db)

		// Act
		reports, err := repo.ReportInboundOrdersCountByEmployee(&employeeId)

		// Assert
		require.Error(t, err)
		require.Nil(t, reports)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].ResponseCode, serviceErr.ResponseCode)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		mock.ExpectQuery("SELECT (.+) FROM employees e LEFT JOIN inbound_orders").WithArgs(employeeId).WillReturnError(sql.ErrConnDone)

		repo := NewEmployeeRepository(db)

		// Act
		reports, err := repo.ReportInboundOrdersCountByEmployee(&employeeId)

		// Assert
		require.Error(t, err)
		require.Nil(t, reports)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].Code, serviceErr.Code)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].ResponseCode, serviceErr.ResponseCode)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEmployeeRepositoryMap_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		cardNumberID := "E001"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 2

		employee := models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		// Expect the update query
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE employees SET").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
			employeeId,
		).WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect the select query to get updated values
		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id",
		}).AddRow(employeeId, cardNumberID, firstName, lastName, warehouseID)

		mock.ExpectQuery("SELECT (.+) FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnRows(rows)
		mock.ExpectCommit()

		repo := NewEmployeeRepository(db)

		// Act
		updatedEmployee, err := repo.Update(employee, employeeId)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, updatedEmployee)
		require.Equal(t, employeeId, *updatedEmployee.ID)
		require.Equal(t, cardNumberID, *updatedEmployee.CardNumberID)
		require.Equal(t, firstName, *updatedEmployee.FirstName)
		require.Equal(t, lastName, *updatedEmployee.LastName)
		require.Equal(t, warehouseID, *updatedEmployee.WarehouseID)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid_warehouse_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		warehouseID := 999 // Invalid warehouse ID
		employee := models.Employee{
			WarehouseID: &warehouseID,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE employees SET").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
			employeeId,
		).WillReturnError(mysqlErr)
		mock.ExpectRollback()

		repo := NewEmployeeRepository(db)

		// Act
		updatedEmployee, err := repo.Update(employee, employeeId)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, updatedEmployee)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)
		require.Equal(t, "Warehouse ID does not exist", serviceErr.InternalError.Error())

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_card_number_id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		cardNumberID := "E002" // Already exists for another employee
		employee := models.Employee{
			CardNumberID: &cardNumberID,
		}

		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'E002' for key 'card_number_id'",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE employees SET").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
			employeeId,
		).WillReturnError(mysqlErr)
		mock.ExpectRollback()

		repo := NewEmployeeRepository(db)

		// Act
		updatedEmployee, err := repo.Update(employee, employeeId)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, updatedEmployee)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].Code, serviceErr.Code)
		require.Equal(t, "Card number ID already exists", serviceErr.InternalError.Error())

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("employee_not_found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 999
		cardNumberID := "E001"
		employee := models.Employee{
			CardNumberID: &cardNumberID,
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE employees SET").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
			employeeId,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		// Expect the select query to return no rows
		mock.ExpectQuery("SELECT (.+) FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		repo := NewEmployeeRepository(db)

		// Act
		updatedEmployee, err := repo.Update(employee, employeeId)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, updatedEmployee)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("transaction_begin_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		cardNumberID := "E001"
		employee := models.Employee{
			CardNumberID: &cardNumberID,
		}

		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		repo := NewEmployeeRepository(db)

		// Act
		updatedEmployee, err := repo.Update(employee, employeeId)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, updatedEmployee)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("transaction_commit_error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		employeeId := 1
		cardNumberID := "E001"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 2

		employee := models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE employees SET").WithArgs(
			employee.CardNumberID,
			employee.FirstName,
			employee.LastName,
			employee.WarehouseID,
			employeeId,
		).WillReturnResult(sqlmock.NewResult(0, 1))

		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id",
		}).AddRow(employeeId, cardNumberID, firstName, lastName, warehouseID)

		mock.ExpectQuery("SELECT (.+) FROM employees WHERE id = \\?").WithArgs(employeeId).WillReturnRows(rows)
		mock.ExpectCommit().WillReturnError(sql.ErrTxDone)
		// Eliminamos la expectativa de rollback ya que el código no lo llama después de un error de commit

		repo := NewEmployeeRepository(db)

		// Act
		updatedEmployee, err := repo.Update(employee, employeeId)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Employee{}, updatedEmployee)
		serviceErr, ok := err.(pkg.ServiceError)
		require.True(t, ok)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer].Code, serviceErr.Code)

		// Ensure all expectations were met
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
