package warehouse

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

/*
	FindAll() (v map[int]models.Warehouse, err error)
	FindByID(id int) (v models.Warehouse, err error)
	Add(v models.Warehouse) (w models.Warehouse, err error)
	FindAvailableID() (id int, err error)
	FindWarehouseByCode(code string) (v models.Warehouse, err error)
	Delete(id int) (err error)
	Update(id int, v models.Warehouse) (w models.Warehouse, err error)
*/

type MockWarehouseRepository struct {
	mock.Mock
}

func (m *MockWarehouseRepository) FindAll() (v map[int]models.Warehouse, err error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) FindByID(id int) (v models.Warehouse, err error) {
	args := m.Called(id)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) Add(v models.Warehouse) (w models.Warehouse, err error) {
	args := m.Called(v)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) Update(id int, v models.Warehouse) (w models.Warehouse, err error) {
	args := m.Called(id, v)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) Delete(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}
