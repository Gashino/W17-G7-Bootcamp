package repository

import (
	"app/pkg"
	"app/pkg/models"
)

type EmployeeRepositoryMap struct {
	db    map[int]models.Employee
	maxId int
}

func NewEmployeeMapRepository(db []models.EmployeeDocument) EmployeeRepository {
	// Calculate initial maxId
	var maxId int
	dbMap := make(map[int]models.Employee)
	for _, e := range db {
		if e.ID > maxId {
			maxId = e.ID
		}
		dbMap[e.ID] = models.Employee{
			ID:           e.ID,
			CardNumberID: e.CardNumberID,
			FirstName:    e.FirstName,
			LastName:     e.LastName,
			WarehouseID:  e.WarehouseID,
		}
	}

	return &EmployeeRepositoryMap{
		db:    dbMap,
		maxId: maxId,
	}
}

func (r *EmployeeRepositoryMap) FindAll() (map[int]models.Employee, error) {
	result := make(map[int]models.Employee)
	for k, v := range r.db {
		result[k] = v
	}
	return result, nil
}

func (r *EmployeeRepositoryMap) FindById(id int) (models.Employee, error) {
	if employee, ok := r.db[id]; ok {
		return employee, nil
	}
	return models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound]
}

func (r *EmployeeRepositoryMap) Save(employee models.Employee) (models.Employee, error) {

	for _, existingEmployee := range r.db {
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
	r.db[newId] = employee

	return employee, nil
}

func (r *EmployeeRepositoryMap) Update(employee models.Employee, id int) (models.Employee, error) {
	if _, ok := r.db[id]; !ok {
		return models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound]
	}

	currentEmployee := r.db[id]

	if currentEmployee.CardNumberID != employee.CardNumberID {

		for existingID, existingEmployee := range r.db {
			if existingID != id && existingEmployee.CardNumberID == employee.CardNumberID {
				return models.Employee{}, pkg.ServiceErrors[pkg.ErrConflict]
			}
		}
	}

	r.db[id] = employee
	return employee, nil
}

func (r *EmployeeRepositoryMap) Delete(id int) error {
	if _, ok := r.db[id]; ok {
		delete(r.db, id)
	}
	return nil
}
