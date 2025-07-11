package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
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
	rows, err := r.db.Query("SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]models.Employee{}
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(&employee.ID, &employee.CardNumberID, &employee.FirstName, &employee.LastName, &employee.WarehouseID)
		if err != nil {
			return nil, err
		}
		result[employee.ID] = employee
	}

	return result, nil
}

func (r *EmployeeRepositoryMap) FindById(id int) (models.Employee, error) {
	rows, err := r.db.Query("SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?", id)
	if err != nil {
		return models.Employee{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	var employee models.Employee
	err = rows.Scan(&employee.ID, &employee.CardNumberID, &employee.FirstName, &employee.LastName, &employee.WarehouseID)
	if err != nil {
		return models.Employee{}, err
	}

	return employee, nil
}

func (r *EmployeeRepositoryMap) Save(employee models.Employee) (models.Employee, error) {
	// Validate that the WarehouseID exists
	_, err := r.db.Exec(
		"INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES (?, ?, ?, ?)",
		employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID,
	)
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

	return employee, nil
}

func (r *EmployeeRepositoryMap) Update(employee models.Employee, id int) (models.Employee, error) {
	_, err := r.db.Exec("UPDATE employees SET card_number_id = ?, first_name = ?, last_name = ?, warehouse_id = ? WHERE id = ?", employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID, id)
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

	return employee, nil
}

func (r *EmployeeRepositoryMap) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		return pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	return nil
}
