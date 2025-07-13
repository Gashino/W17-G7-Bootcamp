package service

import (
	"app/pkg/models"
)

type EmployeeService interface {
	FindAll() (map[int]models.Employee, error)
	FindById(id int) (models.Employee, error)
	Save(employee models.Employee) (models.Employee, error)
	ReportInboundOrdersCountByEmployee(id *int) ([]models.EmployeeReport, error)
	Update(employee models.Employee, id int) (models.Employee, error)
	Delete(id int) error
}
