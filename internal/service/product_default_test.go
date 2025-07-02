package service

import (
	"app/pkg"
	"app/pkg/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Util para *string
func ptr(s string) *string { return &s }

// Mock implementation of ProductRepository for testing
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll() map[int]models.Product {
	args := m.Called()
	return args.Get(0).(map[int]models.Product)
}
func (m *MockProductRepository) GetById(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}
func (m *MockProductRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockProductRepository) Create(product models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductRepository) Update(id int, product models.Product) error {
	args := m.Called(id, product)
	return args.Error(0)
}
func TestNewProductDefault(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	// Act
	service := NewProductDefault(mockRepo)
	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.rp)
}
func TestProductService_GetAll(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	mockProducts := map[int]models.Product{
		1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P1"), Description: ptr("Product 1")}},
		2: {ID: 2, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P2"), Description: ptr("Product 2")}},
	}
	mockRepo.On("GetAll").Return(mockProducts)
	service := NewProductDefault(mockRepo)
	// Act
	result, err := service.GetAll()
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockProducts, result)
	mockRepo.AssertExpectations(t)
}
func TestProductService_GetById(t *testing.T) {
	t.Run("Existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		product := &models.Product{ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P1"), Description: ptr("Product 1")}}
		mockRepo.On("GetById", 1).Return(product, nil)
		service := NewProductDefault(mockRepo)
		// Act
		result, err := service.GetById(1)
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, product, result)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Non-existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		mockRepo.On("GetById", 999).Return(nil, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewProductDefault(mockRepo)
		// Act
		result, err := service.GetById(999)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})
}
func TestProductService_Delete(t *testing.T) {
	t.Run("Delete existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		mockRepo.On("Delete", 1).Return(nil)
		service := NewProductDefault(mockRepo)
		// Act
		err := service.Delete(1)
		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Delete non-existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		mockRepo.On("Delete", 999).Return(pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewProductDefault(mockRepo)
		// Act
		err := service.Delete(999)
		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})
}
func TestProductService_Create(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	product := models.Product{ProductAttributes: models.ProductAttributes{Description: ptr("New Product"), ProductCode: ptr("NP1")}}
	mockRepo.On("Create", product).Return(nil)
	service := NewProductDefault(mockRepo)
	// Act
	err := service.Create(product)
	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
func TestProductService_Update(t *testing.T) {
	t.Run("Update existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		existingProduct := &models.Product{ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P1"), Description: ptr("Product 1")}}
		// Create a ProductDoc with updated description
		description := "Updated Product"
		productDoc := models.Product{ProductAttributes: models.ProductAttributes{Description: &description}}
		// After mapping the doc to the struct
		updatedProduct := &models.Product{ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P1"), Description: ptr("Updated Product")}}
		mockRepo.On("GetById", 1).Return(existingProduct, nil)
		mockRepo.On("Update", 1, *updatedProduct).Return(nil)
		service := NewProductDefault(mockRepo)
		// Act
		result, err := service.Update(1, productDoc)
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updatedProduct, result)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Update non-existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		description := "Updated Product"
		productDoc := models.Product{ProductAttributes: models.ProductAttributes{Description: &description}}
		mockRepo.On("GetById", 999).Return(nil, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewProductDefault(mockRepo)
		// Act
		result, err := service.Update(999, productDoc)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Update with repository error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		existingProduct := &models.Product{ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P1"), Description: ptr("Product 1")}}
		// Create a ProductDoc with updated description
		description := "Updated Product"
		productDoc := models.Product{ProductAttributes: models.ProductAttributes{Description: &description}}
		// After mapping the doc to the struct
		updatedProduct := &models.Product{ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P1"), Description: ptr("Updated Product")}}
		mockRepo.On("GetById", 1).Return(existingProduct, nil)
		mockRepo.On("Update", 1, *updatedProduct).Return(errors.New("update error"))
		service := NewProductDefault(mockRepo)
		// Act
		result, err := service.Update(1, productDoc)
		// Assert
		assert.Error(t, err)
		assert.Equal(t, updatedProduct, result)
		assert.Equal(t, "update error", err.Error())
		mockRepo.AssertExpectations(t)
	})
	t.Run("Update with seller ID", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		existingProduct := &models.Product{ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: ptr("P1"), Description: ptr("Product 1")}}
		// Create a ProductDoc with updated seller ID
		sellerId := 5
		productDoc := models.Product{ProductAttributes: models.ProductAttributes{SellerId: &sellerId}}
		// After mapping the doc to the struct
		updatedProduct := &models.Product{ID: 1, ProductAttributes: models.ProductAttributes{
			ProductCode: ptr("P1"),
			Description: ptr("Product 1"),
			SellerId:    &sellerId, // Now correctly sets SellerId
		}}
		mockRepo.On("GetById", 1).Return(existingProduct, nil)
		mockRepo.On("Update", 1, *updatedProduct).Return(nil)
		service := NewProductDefault(mockRepo)
		// Act
		result, err := service.Update(1, productDoc)
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updatedProduct, result)
		mockRepo.AssertExpectations(t)
	})
}
