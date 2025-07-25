package service

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/locality"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalityService_Create(t *testing.T) {
	t.Run("Cuando el repository crea la locality exitosamente el service la retorna", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		inputLocality := models.Locality{
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "Buenos Aires",
				ProvinceName: "Buenos Aires",
				CountryName:  "Argentina",
			},
		}
		expectedLocality := inputLocality
		expectedLocality.ID = 1

		mockRepo.On("Create", inputLocality).Return(expectedLocality, nil)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.Create(inputLocality)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedLocality, result)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "Buenos Aires", result.LocalityName)
		require.Equal(t, "Buenos Aires", result.ProvinceName)
		require.Equal(t, "Argentina", result.CountryName)
		mockRepo.AssertCalled(t, "Create", inputLocality)
	})

	t.Run("Cuando el repository devuelve error de conflicto el service propaga el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		inputLocality := models.Locality{
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "Buenos Aires",
				ProvinceName: "Buenos Aires",
				CountryName:  "Argentina",
			},
		}
		expectedError := pkg.ServiceErrors[pkg.ErrConflict]

		mockRepo.On("Create", inputLocality).Return(models.Locality{}, expectedError)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.Create(inputLocality)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "Create", inputLocality)
	})

	t.Run("Cuando hay error de validación el service propaga el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		inputLocality := models.Locality{
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "", // Nombre vacío
				ProvinceName: "Buenos Aires",
				CountryName:  "Argentina",
			},
		}
		expectedError := pkg.ServiceErrors[pkg.ErrUnprocessableEntity]

		mockRepo.On("Create", inputLocality).Return(models.Locality{}, expectedError)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.Create(inputLocality)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "Create", inputLocality)
	})
}

func TestLocalityService_GetById(t *testing.T) {
	t.Run("Cuando el repository encuentra la locality el service la retorna", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		expectedLocality := models.Locality{
			ID: 1,
			LocalitiesAttributes: models.LocalitiesAttributes{
				LocalityName: "Buenos Aires",
				ProvinceName: "Buenos Aires",
				CountryName:  "Argentina",
			},
		}
		mockRepo.On("GetById", 1).Return(expectedLocality, nil)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.GetById(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedLocality, result)
		require.Equal(t, 1, result.ID)
		require.Equal(t, "Buenos Aires", result.LocalityName)
		require.Equal(t, "Buenos Aires", result.ProvinceName)
		require.Equal(t, "Argentina", result.CountryName)
		mockRepo.AssertCalled(t, "GetById", 1)
	})

	t.Run("Cuando el repository no encuentra la locality el service devuelve el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		expectedError := pkg.ServiceErrors[pkg.ErrNotFound]
		mockRepo.On("GetById", 999).Return(models.Locality{}, expectedError)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.GetById(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "GetById", 999)
	})

	t.Run("Cuando hay error de conexión a base de datos el service propaga el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		expectedError := fmt.Errorf("database connection error")
		mockRepo.On("GetById", 1).Return(models.Locality{}, expectedError)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.GetById(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.Locality{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "GetById", 1)
	})
}

func TestLocalityService_GetCantSellersByLocality(t *testing.T) {
	t.Run("Cuando el repository encuentra la locality y sellers el service retorna el reporte", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		expectedResponse := models.LocalityBySellerResponse{
			ID:           1,
			LocalityName: models.StringPtr("Buenos Aires"),
			SellerCount:  models.StringPtr("5"),
		}
		mockRepo.On("GetCantSellersByLocality", 1).Return(expectedResponse, nil)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.GetCantSellersByLocality(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedResponse, result)
		require.Equal(t, 1, result.ID)
		require.NotNil(t, result.LocalityName)
		require.Equal(t, "Buenos Aires", *result.LocalityName)
		require.NotNil(t, result.SellerCount)
		require.Equal(t, "5", *result.SellerCount)
		mockRepo.AssertCalled(t, "GetCantSellersByLocality", 1)
	})

	t.Run("Cuando la locality no existe el service devuelve error not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		expectedError := pkg.ServiceErrors[pkg.ErrNotFound]
		mockRepo.On("GetCantSellersByLocality", 999).Return(models.LocalityBySellerResponse{}, expectedError)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.GetCantSellersByLocality(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.LocalityBySellerResponse{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "GetCantSellersByLocality", 999)
	})

	t.Run("Cuando hay error en la consulta el service propaga el error", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		expectedError := fmt.Errorf("query execution error")
		mockRepo.On("GetCantSellersByLocality", 1).Return(models.LocalityBySellerResponse{}, expectedError)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.GetCantSellersByLocality(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, models.LocalityBySellerResponse{}, result)
		require.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "GetCantSellersByLocality", 1)
	})

	t.Run("Cuando la locality existe pero no tiene sellers el service retorna cero sellers", func(t *testing.T) {
		// Arrange
		mockRepo := new(locality.MockLocalityRepository)
		expectedResponse := models.LocalityBySellerResponse{
			ID:           1,
			LocalityName: models.StringPtr("Córdoba"),
			SellerCount:  models.StringPtr("0"),
		}
		mockRepo.On("GetCantSellersByLocality", 1).Return(expectedResponse, nil)

		service := NewLocalityDefault(mockRepo)

		// Act
		result, err := service.GetCantSellersByLocality(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedResponse, result)
		require.Equal(t, 1, result.ID)
		require.NotNil(t, result.LocalityName)
		require.Equal(t, "Córdoba", *result.LocalityName)
		require.NotNil(t, result.SellerCount)
		require.Equal(t, "0", *result.SellerCount)
		mockRepo.AssertCalled(t, "GetCantSellersByLocality", 1)
	})
}
