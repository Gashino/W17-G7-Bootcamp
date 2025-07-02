package repository

import (
	"app/pkg"
	"app/pkg/models"
)

type EmployeeRepositoryMap struct {
	db          *map[int]models.Employee
	maxId       int
	warehouseDb *map[int]models.Warehouse
}

func NewEmployeeMapRepository(employeeDb *map[int]models.Employee, warehouseDb *map[int]models.Warehouse) EmployeeRepository {
	// Calculate initial maxId
	var maxId int
	for _, e := range *employeeDb {
		if e.ID > maxId {
			maxId = e.ID
		}
	}

	return &EmployeeRepositoryMap{
		db:          employeeDb,
		warehouseDb: warehouseDb,
		maxId:       maxId,
	}
}

func (r *EmployeeRepositoryMap) FindAll() (map[int]models.Employee, error) {
	result := make(map[int]models.Employee)
	for k, v := range *r.db {
		result[k] = v
	}
	return result, nil
}

func (r *EmployeeRepositoryMap) FindById(id int) (models.Employee, error) {
	if employee, ok := (*r.db)[id]; ok {
		return employee, nil
	}
	return models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound]
}

func (r *EmployeeRepositoryMap) Save(employee models.Employee) (models.Employee, error) {
	// Validate that the WarehouseID exists
	if _, exists := (*r.warehouseDb)[employee.WarehouseID]; !exists {
		return models.Employee{}, pkg.ServiceError{
			Code:         400,
			ResponseCode: 400,
			Message:      "Warehouse ID does not exist",
		}
	}

	for _, existingEmployee := range *r.db {
		if existingEmployee.CardNumberID == employee.CardNumberID {
			return models.Employee{}, pkg.ServiceError{
				Code:         400,
				ResponseCode: 400,
				Message:      "Card ID already exists",
			}
		}
	}

	r.maxId++
	newId := r.maxId
	employee.ID = newId
	(*r.db)[newId] = employee

	return employee, nil
}

func (r *EmployeeRepositoryMap) Update(employee models.Employee, id int) (models.Employee, error) {
	if _, ok := (*r.db)[id]; !ok {
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	// Validate that the WarehouseID exists
	if _, exists := (*r.warehouseDb)[employee.WarehouseID]; !exists {
		return models.Employee{}, pkg.ServiceError{
			Code:         400,
			ResponseCode: 400,
			Message:      "Warehouse ID does not exist",
		}
	}

	currentEmployee := (*r.db)[id]

	if currentEmployee.CardNumberID != employee.CardNumberID {

		for existingID, existingEmployee := range *r.db {
			if existingID != id && existingEmployee.CardNumberID == employee.CardNumberID {
				return models.Employee{}, pkg.ServiceErrors[pkg.ErrConflict]
			}
		}
	}

	(*r.db)[id] = employee
	return employee, nil
}

func (r *EmployeeRepositoryMap) Delete(id int) error {
	if _, ok := (*r.db)[id]; ok {
		delete(*r.db, id)
	}
	return nil
}
