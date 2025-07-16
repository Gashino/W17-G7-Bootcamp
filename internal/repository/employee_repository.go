package repository

import "app/pkg/models"

type EmployeeRepository interface {
	FindAll() (map[int]models.Employee, error)
	FindById(id int) (models.Employee, error)
	ReportInboundOrdersCountByEmployee(id *int) ([]models.EmployeeReport, error)
	Save(employee models.Employee) (models.Employee, error)
	Update(employee models.Employee, id int) (models.Employee, error)
	Delete(id int) error
}
