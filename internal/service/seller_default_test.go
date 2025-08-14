package service

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/seller"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSellerService_FindAll(t *testing.T) {
	t.Run("Cuando el repository devuelve sellers exitosamente el service los retorna", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		expectedSellers := map[int]models.Seller{
			1: {ID: 1, SellerAttributes: models.SellerAttributes{CId: "12345", CompanyName: "Test Company 1", Address: "Test Address 1", Telephone: "123456789", LocalityID: 1}},
			2: {ID: 2, SellerAttributes: models.SellerAttributes{CId: "67890", CompanyName: "Test Company 2", Address: "Test Address 2", Telephone: "987654321", LocalityID: 2}},
		}
		mockRepo.On("FindAll").Return(expectedSellers, nil)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.FindAll()

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSellers, result)
		require.Len(t, result, 2)
		mockRepo.AssertCalled(t, "FindAll")
	})

	t.Run("Cuando el repository devuelve error el service propaga el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		expectedError := fmt.Errorf("database connection error")
		mockRepo.On("FindAll").Return(map[int]models.Seller{}, expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.FindAll()

		// Assert
		require.Error(t, err)
		require.Empty(t, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "FindAll")
	})
}

func TestSellerService_GetById(t *testing.T) {
	t.Run("Cuando el repository encuentra el seller el service lo retorna", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		expectedSeller := models.Seller{
			ID: 1,
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}
		mockRepo.On("GetById", 1).Return(expectedSeller, nil)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.GetById(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSeller, result)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "12345", result.CId)
		mockRepo.AssertCalled(t, "GetById", 1)
	})

	t.Run("Cuando el repository no encuentra el seller el service devuelve el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		expectedError := pkg.ServiceErrors[pkg.ErrNotFound]
		mockRepo.On("GetById", 999).Return(models.Seller{}, expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.GetById(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "GetById", 999)
	})
}

func TestSellerService_Create(t *testing.T) {
	t.Run("Cuando el repository crea el seller exitosamente el service lo retorna", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "New Company",
				Address:     "New Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}
		expectedSeller := inputSeller
		expectedSeller.ID = 1

		mockRepo.On("Create", inputSeller).Return(expectedSeller, nil)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.Create(inputSeller)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSeller, result)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "12345", result.CId)
		mockRepo.AssertCalled(t, "Create", inputSeller)
	})

	t.Run("Cuando el repository devuelve error de conflicto el service propaga el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Duplicate Company",
				Address:     "Some Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}
		expectedError := pkg.ServiceErrors[pkg.ErrConflict]

		mockRepo.On("Create", inputSeller).Return(models.Seller{}, expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.Create(inputSeller)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "Create", inputSeller)
	})

	t.Run("Cuando la locality_id no existe el service devuelve error not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		inputSeller := models.Seller{
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "123456789",
				LocalityID:  999, // ID inexistente
			},
		}
		expectedError := pkg.ServiceErrors[pkg.ErrNotFound]

		mockRepo.On("Create", inputSeller).Return(models.Seller{}, expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.Create(inputSeller)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "Create", inputSeller)
	})
}

func TestSellerService_UpdateFields(t *testing.T) {
	t.Run("Cuando el repository actualiza exitosamente el service retorna el seller actualizado", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		updateData := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
			Address:     models.StringPtr("Updated Address"),
		}
		expectedSeller := models.Seller{
			ID: 1,
			SellerAttributes: models.SellerAttributes{
				CId:         "12345",
				CompanyName: "Updated Company",
				Address:     "Updated Address",
				Telephone:   "123456789",
				LocalityID:  1,
			},
		}

		mockRepo.On("UpdateFields", 1, updateData).Return(expectedSeller, nil)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.UpdateFields(1, updateData)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSeller, result)
		require.Equal(t, "Updated Company", result.CompanyName)
		require.Equal(t, "Updated Address", result.Address)
		mockRepo.AssertCalled(t, "UpdateFields", 1, updateData)
	})

	t.Run("Cuando el seller no existe el service devuelve error not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		updateData := models.SellerCreateRequest{
			CompanyName: models.StringPtr("Updated Company"),
		}
		expectedError := pkg.ServiceErrors[pkg.ErrNotFound]

		mockRepo.On("UpdateFields", 999, updateData).Return(models.Seller{}, expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.UpdateFields(999, updateData)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "UpdateFields", 999, updateData)
	})

	t.Run("Cuando se intenta actualizar con CId duplicado el service devuelve error de conflicto", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		updateData := models.SellerCreateRequest{
			CId: models.StringPtr("67890"), // CId que ya existe
		}
		expectedError := pkg.ServiceErrors[pkg.ErrConflict]

		mockRepo.On("UpdateFields", 1, updateData).Return(models.Seller{}, expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		result, err := service.UpdateFields(1, updateData)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Seller{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "UpdateFields", 1, updateData)
	})
}

func TestSellerService_DeleteSeller(t *testing.T) {
	t.Run("Cuando el repository elimina exitosamente el service no retorna error", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		mockRepo.On("DeleteSeller", 1).Return(nil)

		service := NewSellerDefault(mockRepo)

		// Act
		err := service.DeleteSeller(1)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertCalled(t, "DeleteSeller", 1)
	})

	t.Run("Cuando el seller no existe el service devuelve error not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		expectedError := pkg.ServiceErrors[pkg.ErrNotFound]
		mockRepo.On("DeleteSeller", 999).Return(expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		err := service.DeleteSeller(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "DeleteSeller", 999)
	})

	t.Run("Cuando hay referencias en otras tablas el service devuelve error de conflicto", func(t *testing.T) {
		// Arrange
		mockRepo := new(seller.MockSellerRepository)
		expectedError := pkg.ServiceErrors[pkg.ErrConflict]
		mockRepo.On("DeleteSeller", 1).Return(expectedError)

		service := NewSellerDefault(mockRepo)

		// Act
		err := service.DeleteSeller(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "DeleteSeller", 1)
	})
}
