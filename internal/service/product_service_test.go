package service

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/product"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper functions for creating pointers
func float64Ptr(f float64) *float64 { return &f }

// Function to create a test product
func createTestProduct() models.Product {
	return models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			ProductCode:                    models.StringPtr("TEST001"),
			Description:                    models.StringPtr("Test Product"),
			NetWeight:                      float64Ptr(10.5),
			ExpirationRate:                 models.IntPtr(30),
			RecommendedFreezingTemperature: float64Ptr(-18.0),
			FreezingRate:                   models.IntPtr(10),
			ProductTypeId:                  models.IntPtr(101),
			SellerId:                       models.IntPtr(1),
		},
		Dimensions: models.Dimensions{
			Width:  float64Ptr(10.0),
			Height: float64Ptr(20.0),
			Length: float64Ptr(30.0),
		},
	}
}

// Tests
func TestProductService_GetAll(t *testing.T) {
	// Arrange
	mockRepo := new(product.MockProductRepository)
	productService := NewProductDefault(mockRepo)

	expectedProducts := map[int]models.Product{
		1: createTestProduct(),
	}

	mockRepo.On("GetAll").Return(expectedProducts)

	// Act
	products, err := productService.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertCalled(t, "GetAll")
}

func TestProductService_GetById(t *testing.T) {
	t.Run("find_by_id_existent", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		expectedProduct := createTestProduct()
		mockRepo.On("GetById", 1).Return(&expectedProduct, nil)

		// Act
		product, err := productService.GetById(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &expectedProduct, product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		mockRepo.On("GetById", 999).Return(nil, pkg.ServiceErrors[pkg.ErrNotFound])

		// Act
		product, err := productService.GetById(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertCalled(t, "GetById", 999)
	})
}

func TestProductService_Delete(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		mockRepo.On("Delete", 1).Return(nil)

		// Act
		err := productService.Delete(1)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", 1)
	})

	t.Run("delete_non_existent", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		mockRepo.On("Delete", 999).Return(pkg.ServiceErrors[pkg.ErrNotFound])

		// Act
		err := productService.Delete(999)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertCalled(t, "Delete", 999)
	})
}

func TestProductService_Create(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		newProduct := createTestProduct()
		newProduct.ID = 0 // ID should be assigned by the repository

		expectedProduct := createTestProduct() // ID = 1

		mockRepo.On("Create", newProduct).Return(&expectedProduct, nil)

		// Act
		product, err := productService.Create(newProduct)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &expectedProduct, product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_conflict", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		newProduct := createTestProduct()
		newProduct.ID = 0

		mockRepo.On("Create", newProduct).Return(nil, pkg.ServiceErrors[pkg.ErrConflict])

		// Act
		product, err := productService.Create(newProduct)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
		mockRepo.AssertCalled(t, "Create", newProduct)
	})
}

func TestProductService_Update(t *testing.T) {
	t.Run("update_existent", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		existingProduct := createTestProduct()
		updatedProduct := models.Product{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				Description: models.StringPtr("Updated Description"),
				NetWeight:   float64Ptr(15.0),
			},
		}

		// The product after the update
		expectedProduct := createTestProduct()
		*expectedProduct.ProductAttributes.Description = "Updated Description"
		*expectedProduct.ProductAttributes.NetWeight = 15.0

		mockRepo.On("GetById", 1).Return(&existingProduct, nil)
		mockRepo.On("Update", 1, expectedProduct).Return(nil)

		// Act
		result, err := productService.Update(1, updatedProduct)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &expectedProduct, result)
		mockRepo.AssertCalled(t, "GetById", 1)
		mockRepo.AssertCalled(t, "Update", 1, expectedProduct)
	})

	t.Run("update_non_existent", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		updatedProduct := models.Product{
			ID: 999,
			ProductAttributes: models.ProductAttributes{
				Description: models.StringPtr("Updated Description"),
			},
		}

		mockRepo.On("GetById", 999).Return(nil, pkg.ServiceErrors[pkg.ErrNotFound])

		// Act
		result, err := productService.Update(999, updatedProduct)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertCalled(t, "GetById", 999)
	})

	t.Run("update error", func(t *testing.T) {
		// Arrange
		mockRepo := new(product.MockProductRepository)
		productService := NewProductDefault(mockRepo)

		existingProduct := createTestProduct()
		updatedProduct := models.Product{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				ProductCode: models.StringPtr("DUPLICATE"),
			},
		}

		// The product after the update
		expectedProduct := createTestProduct()
		*expectedProduct.ProductAttributes.ProductCode = "DUPLICATE"

		mockRepo.On("GetById", 1).Return(&existingProduct, nil)
		mockRepo.On("Update", 1, expectedProduct).Return(errors.New("duplicate product code"))

		// Act
		result, err := productService.Update(1, updatedProduct)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "duplicate product code", err.Error())
		mockRepo.AssertCalled(t, "GetById", 1)
		mockRepo.AssertCalled(t, "Update", 1, expectedProduct)
	})
}
