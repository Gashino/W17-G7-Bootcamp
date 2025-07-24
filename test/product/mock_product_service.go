package product

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

// Create implements service.IProductService.
func (m *MockProductService) Create(product models.Product) (*models.Product, error) {
	args := m.Called(product)
	prod := args.Get(0).(*models.Product)
	prod.ID = 1
	return prod, args.Error(1)
}

// Delete implements service.IProductService.
func (m *MockProductService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetAll implements service.IProductService.
func (m *MockProductService) GetAll() (map[int]models.Product, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Product), args.Error(1)
}

// GetById implements service.IProductService.
func (m *MockProductService) GetById(id int) (*models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Product), args.Error(1)
}

// GetProductRecords implements service.IProductService.
func (m *MockProductService) GetProductRecords(id *int) ([]models.ProductRecordResponse, error) {
	args := m.Called(id)
	return args.Get(0).([]models.ProductRecordResponse), args.Error(1)
}

// Update implements service.IProductService.
func (m *MockProductService) Update(id int, product models.Product) (*models.Product, error) {
	args := m.Called(id, product)
	return args.Get(0).(*models.Product), args.Error(1)
}
