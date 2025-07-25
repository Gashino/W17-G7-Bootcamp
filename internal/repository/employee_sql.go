package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const (
	selectAllEmployees = `
		SELECT id, card_number_id, first_name, last_name, warehouse_id
		FROM employees
	`
	selectEmployeeById = `
		SELECT id, card_number_id, first_name, last_name, warehouse_id
		FROM employees
		WHERE id = ?
	`
	insertEmployee = "INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES (?, ?, ?, ?)"
	updateEmployee = `UPDATE employees SET 
		card_number_id = COALESCE(?, card_number_id),
		first_name = COALESCE(?, first_name),
		last_name = COALESCE(?, last_name),
		warehouse_id = COALESCE(?, warehouse_id)
		WHERE id = ?`
	deleteEmployee                     = "DELETE FROM employees WHERE id = ?"
	reportInboundOrdersCountByEmployee = `
		SELECT e.id, e.card_number_id, e.first_name, e.last_name, COUNT(io.id) AS inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders io ON e.id = io.employee_id
		GROUP BY e.id, e.card_number_id, e.first_name, e.last_name
	`
	reportInboundOrdersCountByEmployeeWithId = `
		SELECT e.id, e.card_number_id, e.first_name, e.last_name, COUNT(io.id) AS inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders io ON e.id = io.employee_id
		WHERE e.id = ?
		GROUP BY e.id, e.card_number_id, e.first_name, e.last_name
	`
)

type EmployeeRepositoryMap struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	// Calculate initial maxId
	return &EmployeeRepositoryMap{
		db: db,
	}
}

func (r *EmployeeRepositoryMap) FindAll() (map[int]models.Employee, error) {
	rows, err := r.db.Query(selectAllEmployees)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]models.Employee{}
	for rows.Next() {
		var id int
		var cardNumberID, firstName, lastName string
		var warehouseID int

		err := rows.Scan(&id, &cardNumberID, &firstName, &lastName, &warehouseID)
		if err != nil {
			return nil, err
		}

		employee := models.Employee{
			ID:           &id,
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}
		result[id] = employee
	}

	return result, nil
}

func (r *EmployeeRepositoryMap) FindById(id int) (models.Employee, error) {
	rows, err := r.db.Query(selectEmployeeById, id)
	if err != nil {
		return models.Employee{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	var employeeID int
	var cardNumberID, firstName, lastName string
	var warehouseID int

	err = rows.Scan(&employeeID, &cardNumberID, &firstName, &lastName, &warehouseID)
	if err != nil {
		return models.Employee{}, err
	}

	employee := models.Employee{
		ID:           &employeeID,
		CardNumberID: &cardNumberID,
		FirstName:    &firstName,
		LastName:     &lastName,
		WarehouseID:  &warehouseID,
	}

	return employee, nil
}

func (r *EmployeeRepositoryMap) Save(employee models.Employee) (models.Employee, error) {
	// Validate that the WarehouseID exists
	data, err := r.db.Exec(insertEmployee, employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				// Error foreign key constraint (warehouse_id no existe)
				srvError := pkg.ServiceErrors[pkg.ErrNotFound]
				srvError.InternalError = fmt.Errorf("Warehouse ID does not exist")
				return models.Employee{}, srvError
			case 1062:
				// Error unique constraint (card_number_id ya existe)
				srvError := pkg.ServiceErrors[pkg.ErrConflict]
				srvError.InternalError = fmt.Errorf("Card number ID already exists")
				return models.Employee{}, srvError
			default:
				// Otro error de MySQL
				return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
			}
		}
		// Error no es de MySQL
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	lastInsertId, err := data.LastInsertId()
	if err != nil {
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	id := int(lastInsertId)
	employee.ID = &id

	return employee, nil
}

func (r *EmployeeRepositoryMap) Update(employee models.Employee, id int) (models.Employee, error) {
	// Start transaction for atomic update
	tx, err := r.db.Begin()
	if err != nil {
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	defer tx.Rollback()

	_, err = tx.Exec(updateEmployee,
		employee.CardNumberID,
		employee.FirstName,
		employee.LastName,
		employee.WarehouseID,
		id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				srvError := pkg.ServiceErrors[pkg.ErrNotFound]
				srvError.InternalError = fmt.Errorf("Warehouse ID does not exist")
				return models.Employee{}, srvError
			case 1062:
				srvError := pkg.ServiceErrors[pkg.ErrConflict]
				srvError.InternalError = fmt.Errorf("Card number ID already exists")
				return models.Employee{}, srvError
			default:
				return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
			}
		}
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	// Get current values from database to verify record exists and return updated data
	var updatedEmployeeID int
	var updatedCardNumberID, updatedFirstName, updatedLastName string
	var updatedWarehouseID int

	err = tx.QueryRow(selectEmployeeById, id).Scan(
		&updatedEmployeeID,
		&updatedCardNumberID,
		&updatedFirstName,
		&updatedLastName,
		&updatedWarehouseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound]
		}
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	// Note: RowsAffected might be 0 if values are identical (COALESCE doesn't change anything)
	// This is normal behavior and not an error

	updatedEmployee := models.Employee{
		ID:           &updatedEmployeeID,
		CardNumberID: &updatedCardNumberID,
		FirstName:    &updatedFirstName,
		LastName:     &updatedLastName,
		WarehouseID:  &updatedWarehouseID,
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	return updatedEmployee, nil
}

func (r *EmployeeRepositoryMap) ReportInboundOrdersCountByEmployee(id *int) ([]models.EmployeeReport, error) {
	var rows *sql.Rows
	var err error

	if id != nil && *id > 0 {
		// Consulta solo para un empleado específico
		rows, err = r.db.Query(reportInboundOrdersCountByEmployeeWithId, id)
	} else {
		// Consulta para todos los empleados
		rows, err = r.db.Query(reportInboundOrdersCountByEmployee)
	}

	if err != nil {
		return nil, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	defer rows.Close()

	var reports []models.EmployeeReport
	hasResults := false

	for rows.Next() {
		hasResults = true
		var r models.EmployeeReport
		if err := rows.Scan(&r.ID, &r.CardNumberID, &r.FirstName, &r.LastName, &r.InboundOrdersCount); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	// Verificar si hubo errores durante la iteración
	if err = rows.Err(); err != nil {
		return nil, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	// Si se buscó un empleado específico y no hay resultados, es un error
	if id != nil && *id > 0 && !hasResults {
		srvError := pkg.ServiceErrors[pkg.ErrNotFound]
		srvError.InternalError = fmt.Errorf("Employee with ID %d not found", *id)
		return nil, srvError
	}

	return reports, nil
}

func (r *EmployeeRepositoryMap) Delete(id int) error {
	_, err := r.db.Exec(deleteEmployee, id)
	if err != nil {
		return pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	return nil
}
