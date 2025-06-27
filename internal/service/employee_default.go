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
	return s.repository.Update(employee, id)
}

func (s *EmployeeServiceDefault) Delete(id int) error {
	return s.repository.Delete(id)
}
