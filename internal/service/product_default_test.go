package service

import (
	"app/pkg"
	"app/pkg/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Funciones auxiliares para crear punteros
func strPtr(s string) *string       { return &s }
func intPtr(i int) *int             { return &i }
func float64Ptr(f float64) *float64 { return &f }

// Mock del repositorio de productos
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

func (m *MockProductRepository) Create(product models.Product) (*models.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(id int, product models.Product) error {
	args := m.Called(id, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductRecords(id *int) ([]models.ProductRecordResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ProductRecordResponse), args.Error(1)
}

// Función para crear un producto de prueba
func createTestProduct() models.Product {
	return models.Product{
		ID: 1,
		ProductAttributes: models.ProductAttributes{
			ProductCode:                    strPtr("TEST001"),
			Description:                    strPtr("Test Product"),
			NetWeight:                      float64Ptr(10.5),
			ExpirationRate:                 intPtr(30),
			RecommendedFreezingTemperature: float64Ptr(-18.0),
			FreezingRate:                   intPtr(10),
			ProductTypeId:                  intPtr(101),
			SellerId:                       intPtr(1),
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
	mockRepo := new(MockProductRepository)
	service := NewProductDefault(mockRepo)

	expectedProducts := map[int]models.Product{
		1: createTestProduct(),
	}

	mockRepo.On("GetAll").Return(expectedProducts)

	// Act
	products, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetById(t *testing.T) {
	t.Run("existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		expectedProduct := createTestProduct()
		mockRepo.On("GetById", 1).Return(&expectedProduct, nil)

		// Act
		product, err := service.GetById(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &expectedProduct, product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("non-existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		mockRepo.On("GetById", 999).Return(nil, pkg.ServiceErrors[pkg.ErrNotFound])

		// Act
		product, err := service.GetById(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Delete(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		mockRepo.On("Delete", 1).Return(nil)

		// Act
		err := service.Delete(1)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deletion error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		mockRepo.On("Delete", 999).Return(pkg.ServiceErrors[pkg.ErrNotFound])

		// Act
		err := service.Delete(999)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		productToCreate := createTestProduct()
		productToCreate.ID = 0 // ID should be assigned by the repository

		expectedProduct := createTestProduct() // With ID=1

		mockRepo.On("Create", productToCreate).Return(&expectedProduct, nil)

		// Act
		createdProduct, err := service.Create(productToCreate)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &expectedProduct, createdProduct)
		mockRepo.AssertExpectations(t)
	})

	t.Run("creation error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		productToCreate := createTestProduct()
		productToCreate.ID = 0

		mockRepo.On("Create", productToCreate).Return(nil, errors.New("duplicate product code"))

		// Act
		createdProduct, err := service.Create(productToCreate)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, createdProduct)
		assert.Equal(t, "duplicate product code", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		existingProduct := createTestProduct()
		updatedProduct := models.Product{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				Description: strPtr("Updated Description"),
				NetWeight:   float64Ptr(15.0),
			},
		}

		// El producto después de la actualización
		expectedProduct := createTestProduct()
		*expectedProduct.ProductAttributes.Description = "Updated Description"
		*expectedProduct.ProductAttributes.NetWeight = 15.0

		mockRepo.On("GetById", 1).Return(&existingProduct, nil)
		mockRepo.On("Update", 1, expectedProduct).Return(nil)

		// Act
		result, err := service.Update(1, updatedProduct)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &expectedProduct, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update non-existing product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		updatedProduct := models.Product{
			ID: 999,
			ProductAttributes: models.ProductAttributes{
				Description: strPtr("Updated Description"),
			},
		}

		mockRepo.On("GetById", 999).Return(nil, pkg.ServiceErrors[pkg.ErrNotFound])

		// Act
		result, err := service.Update(999, updatedProduct)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		existingProduct := createTestProduct()
		updatedProduct := models.Product{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				ProductCode: strPtr("DUPLICATE"),
			},
		}

		// El producto después de la actualización
		expectedProduct := createTestProduct()
		*expectedProduct.ProductAttributes.ProductCode = "DUPLICATE"

		mockRepo.On("GetById", 1).Return(&existingProduct, nil)
		mockRepo.On("Update", 1, expectedProduct).Return(errors.New("duplicate product code"))

		// Act
		result, err := service.Update(1, updatedProduct)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "duplicate product code", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetProductRecords(t *testing.T) {
	t.Run("get all product records", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		expectedRecords := []models.ProductRecordResponse{
			{ProductId: 1, Description: "Product 1", RecordsCount: 5},
			{ProductId: 2, Description: "Product 2", RecordsCount: 3},
		}

		mockRepo.On("GetProductRecords", (*int)(nil)).Return(expectedRecords, nil)

		// Act
		records, err := service.GetProductRecords(nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedRecords, records)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get product records for specific product", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		productId := 1
		expectedRecords := []models.ProductRecordResponse{
			{ProductId: 1, Description: "Product 1", RecordsCount: 5},
		}

		mockRepo.On("GetProductRecords", &productId).Return(expectedRecords, nil)

		// Act
		records, err := service.GetProductRecords(&productId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedRecords, records)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting product records", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockProductRepository)
		service := NewProductDefault(mockRepo)

		productId := 999
		mockRepo.On("GetProductRecords", &productId).Return(nil, pkg.ServiceErrors[pkg.ErrNotFound])

		// Act
		records, err := service.GetProductRecords(&productId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, records)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})
}
