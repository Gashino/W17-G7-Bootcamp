package product

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

// Create implements repository.ProductRepository.
func (m MockProductRepository) Create(product models.Product) (*models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*models.Product), args.Error(1)
}

// Delete implements repository.ProductRepository.
func (m MockProductRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetAll implements repository.ProductRepository.
func (m MockProductRepository) GetAll() map[int]models.Product {
	args := m.Called()
	return args.Get(0).(map[int]models.Product)
}

// GetById implements repository.ProductRepository.
func (m MockProductRepository) GetById(id int) (*models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Product), args.Error(1)
}

// Update implements repository.ProductRepository.
func (m MockProductRepository) Update(id int, product models.Product) error {
	args := m.Called(id, product)
	return args.Error(0)
}

// GetProductRecords implements repository.ProductRepository.
func (m MockProductRepository) GetProductRecords(id *int) ([]models.ProductRecordResponse, error) {
	args := m.Called(id)
	return args.Get(0).([]models.ProductRecordResponse), args.Error(1)
}
