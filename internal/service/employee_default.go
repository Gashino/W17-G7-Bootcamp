package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

type EmployeeServiceDefault struct {
	repository repository.EmployeeRepository
}

func NewEmployeeServiceDefault(repository repository.EmployeeRepository) EmployeeService {
	return &EmployeeServiceDefault{
		repository: repository,
	}
}

func (s *EmployeeServiceDefault) FindAll() (map[int]models.Employee, error) {
	return s.repository.FindAll()
}

func (s *EmployeeServiceDefault) FindById(id int) (models.Employee, error) {
	return s.repository.FindById(id)
}

func (s *EmployeeServiceDefault) Save(employee models.Employee) (models.Employee, error) {
	return s.repository.Save(employee)
}

func (s *EmployeeServiceDefault) Update(employee models.Employee, id int) (models.Employee, error) {
	e, err := s.repository.FindById(id)
	if err != nil {
		return models.Employee{}, err
	}
	if employee.WarehouseID != 0 {
		e.WarehouseID = employee.WarehouseID
	}
	if employee.FirstName != "" {
		e.FirstName = employee.FirstName
	}
	if employee.LastName != "" {
		e.LastName = employee.LastName
	}
	if e.CardNumberID != employee.CardNumberID {
		e.CardNumberID = employee.CardNumberID
	}
	return s.repository.Update(e, id)
}

func (s *EmployeeServiceDefault) Delete(id int) error {
	return s.repository.Delete(id)
}
