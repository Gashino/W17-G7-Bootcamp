package employee

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) FindAll() (map[int]models.Employee, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Employee), args.Error(1)
}

// ReportInboundOrdersCountByEmployee implements service.EmployeeService.
func (m *MockEmployeeService) ReportInboundOrdersCountByEmployee(id *int) ([]models.EmployeeReport, error) {
	args := m.Called(id)
	return args.Get(0).([]models.EmployeeReport), args.Error(1)
}

func (m *MockEmployeeService) FindById(id int) (models.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) Save(employee models.Employee) (models.Employee, error) {
	args := m.Called(employee)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) Update(employee models.Employee, id int) (models.Employee, error) {
	args := m.Called(employee, id)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
