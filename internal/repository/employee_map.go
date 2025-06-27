package repository

import (
	"app/pkg"
	"app/pkg/models"
)

type EmployeeRepositoryMap struct {
	db    map[int]models.Employee
	maxId int
}

func NewEmployeeMapRepository(db map[int]models.Employee) EmployeeRepository {
	// Calculate initial maxId
	var maxId int
	for id := range db {
		if id > maxId {
			maxId = id
		}
	}

	return &EmployeeRepositoryMap{
		db:    db,
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
				Code:         100,
				ResponseCode: 400,
				Message:      "Card ID already exists",
			}
		}
	}

	r.maxId++
	newId := r.maxId

	r.db[newId] = employee

	return employee, nil
}

func (r *EmployeeRepositoryMap) Update(employee models.Employee, id int) (models.Employee, error) {
	if _, ok := r.db[id]; !ok {
		return models.Employee{}, nil
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
