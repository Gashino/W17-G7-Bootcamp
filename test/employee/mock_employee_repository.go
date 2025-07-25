package employee

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) FindAll() (map[int]models.Employee, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) FindById(id int) (models.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Save(employee models.Employee) (models.Employee, error) {
	args := m.Called(employee)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Update(employee models.Employee, id int) (models.Employee, error) {
	args := m.Called(employee, id)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) ReportInboundOrdersCountByEmployee(id *int) ([]models.EmployeeReport, error) {
	panic("unimplemented")
}
