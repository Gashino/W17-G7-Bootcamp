package repository

import (
	"app/pkg"
	"app/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMap_GetAll(t *testing.T) {
	// Arrange
	mockProducts := map[int]models.Product{
		1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
		2: {ID: 2, ProductAttributes: models.ProductAttributes{ProductCode: "P2", Description: "Product 2"}},
	}
	mockTypes := map[int]models.ProductType{}
	repo := NewProductMap(&mockProducts, &mockTypes)

	// Act
	result := repo.GetAll()

	// Assert
	assert.Equal(t, mockProducts, result, "GetAll should return all products")
	assert.Len(t, result, 2, "GetAll should return two products")
}

func TestProductMap_GetByID(t *testing.T) {
	// Arrange
	mockProducts := map[int]models.Product{
		1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
		2: {ID: 2, ProductAttributes: models.ProductAttributes{ProductCode: "P2", Description: "Product 2"}},
	}
	mockTypes := map[int]models.ProductType{}
	repo := NewProductMap(&mockProducts, &mockTypes)

	t.Run("existing product", func(t *testing.T) {
		// Act
		product, err := repo.GetById(1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, mockProducts[1], *product)
	})

	t.Run("Non-existing product", func(t *testing.T) {
		// Act
		product, err := repo.GetById(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})
}

func TestProductMap_Create(t *testing.T) {
	t.Run("Create new product successfully", func(t *testing.T) {
		// Arrange
		db := make(map[int]models.Product)
		dbTypes := make(map[int]models.ProductType)
		repo := NewProductMap(&db, &dbTypes)
		newProduct := models.Product{
			ProductAttributes: models.ProductAttributes{
				Description: "New Product",
				ProductCode: "NP1",
			},
		}

		// Act
		err := repo.Create(newProduct)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, len(db)) // Verify a product was added to the map
	})

	t.Run("Create product with duplicate code", func(t *testing.T) {
		// Arrange
		existingProducts := map[int]models.Product{
			1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
		}
		mockTypes := map[int]models.ProductType{}
		repo := NewProductMap(&existingProducts, &mockTypes)
		newProduct := models.Product{
			ProductAttributes: models.ProductAttributes{
				Description: "New Product",
				ProductCode: "P1", // Same code as existing product
			},
		}

		// Act
		err := repo.Create(newProduct)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Product code already exist")
		assert.Len(t, *repo.db, 1) // No new product debería ser agregado
	})
}

func TestProductMap_Update(t *testing.T) {
	t.Run("Update existing product", func(t *testing.T) {
		// Arrange
		mockProducts := map[int]models.Product{
			1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
			2: {ID: 2, ProductAttributes: models.ProductAttributes{ProductCode: "P2", Description: "Product 2"}},
		}
		mockTypes := map[int]models.ProductType{}
		repo := NewProductMap(&mockProducts, &mockTypes)
		updatedProduct := models.Product{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				Description: "Updated Product",
				ProductCode: "UP1",
			},
		}

		// Act
		err := repo.Update(1, updatedProduct)

		// Assert
		assert.NoError(t, err)
		product, _ := repo.GetById(1)
		assert.Equal(t, updatedProduct, *product)
	})

	t.Run("Update with duplicate product code", func(t *testing.T) {
		// Arrange
		mockProducts := map[int]models.Product{
			1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
			2: {ID: 2, ProductAttributes: models.ProductAttributes{ProductCode: "P2", Description: "Product 2"}},
		}
		mockTypes := map[int]models.ProductType{}
		repo := NewProductMap(&mockProducts, &mockTypes)
		updatedProduct := models.Product{
			ID: 1,
			ProductAttributes: models.ProductAttributes{
				Description: "Updated Product",
				ProductCode: "P2", // Same as product 2
			},
		}

		// Act
		err := repo.Update(1, updatedProduct)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Product code already exist")
		// Original product should remain unchanged
		product, _ := repo.GetById(1)
		assert.Equal(t, mockProducts[1], *product)
	})
}

func TestProductMap_Delete(t *testing.T) {
	t.Run("Delete existing product", func(t *testing.T) {
		// Arrange
		mockProducts := map[int]models.Product{
			1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
			2: {ID: 2, ProductAttributes: models.ProductAttributes{ProductCode: "P2", Description: "Product 2"}},
		}
		mockTypes := map[int]models.ProductType{}
		repo := NewProductMap(&mockProducts, &mockTypes)

		// Act
		err := repo.Delete(1)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, *repo.db, 1)
		_, exists := (*repo.db)[1]
		assert.False(t, exists)
	})

	t.Run("Delete non-existing product", func(t *testing.T) {
		// Arrange
		mockProducts := map[int]models.Product{
			1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
		}
		mockTypes := map[int]models.ProductType{}
		repo := NewProductMap(&mockProducts, &mockTypes)

		// Act
		err := repo.Delete(999)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		assert.Len(t, *repo.db, 1) // DB should remain unchanged
	})
}

func TestNewProductMap(t *testing.T) {
	t.Run("Create with nil db", func(t *testing.T) {
		// Act
		repo := NewProductMap(&map[int]models.Product{}, &map[int]models.ProductType{})

		// Assert
		assert.NotNil(t, repo)
		assert.NotNil(t, repo.db)
		assert.Len(t, *repo.db, 0)
		assert.Equal(t, 1, repo.lastId)
	})

	t.Run("Create with existing db", func(t *testing.T) {
		// Arrange
		mockProducts := map[int]models.Product{
			1: {ID: 1, ProductAttributes: models.ProductAttributes{ProductCode: "P1", Description: "Product 1"}},
			2: {ID: 2, ProductAttributes: models.ProductAttributes{ProductCode: "P2", Description: "Product 2"}},
		}
		mockTypes := map[int]models.ProductType{}

		// Act
		repo := NewProductMap(&mockProducts, &mockTypes)

		// Assert
		assert.NotNil(t, repo)
		assert.Equal(t, mockProducts, *repo.db)
		assert.Equal(t, 3, repo.lastId) // lastId debería ser len(db) + 1
	})
}
