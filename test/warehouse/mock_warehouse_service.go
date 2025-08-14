package warehouse

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

/*
	FindAll() (v map[int]models.Warehouse, err error)
	FindByID(id int) (v models.Warehouse, err error)
	Add(v models.WarehouseDoc) (w models.Warehouse, err error)
	Update(id int, v models.WarehouseDoc) (w models.Warehouse, err error)
	Delete(id int) (err error)
*/

type MockWarehouseService struct {
	mock.Mock
}

func (m *MockWarehouseService) FindAll() (map[int]models.Warehouse, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) FindByID(id int) (v models.Warehouse, err error) {
	args := m.Called(id)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) Add(v models.WarehouseDoc) (w models.Warehouse, err error) {
	args := m.Called(v)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) Update(id int, v models.WarehouseDoc) (w models.Warehouse, err error) {
	args := m.Called(id, v)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) Delete(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}
