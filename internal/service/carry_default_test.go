package service

import (
	"app/pkg/models"
	"app/test/carry"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
SearchByLocality(locality_id int) (c map[string]models.CarryByLocality, err error)
Create(v models.Carry) (w models.Carry, err error)
*/

func TestCarryService_SearchByLocality(t *testing.T) {
	t.Run("Si la búsqueda por localidad es exitosa devolverá un mapa con los carries de esa localidad", func(t *testing.T) {
		// Arrange
		localityId := 1
		expectedCarries := map[string]models.CarryByLocality{
			"1": {
				LocalityId:   "1",
				LocalityName: "Buenos Aires",
				CarriesCount: 3,
			},
		}

		mockRepo := &carry.MockCarryRepository{}
		mockRepo.On("SearchByLocality", localityId).Return(expectedCarries, nil)

		service := NewCarryDefault(mockRepo)

		// Act
		result, err := service.SearchByLocality(localityId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCarries, result)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si ocurre un error en la búsqueda por localidad se retorna error", func(t *testing.T) {
		// Arrange
		localityId := 999

		mockRepo := &carry.MockCarryRepository{}
		mockRepo.On("SearchByLocality", localityId).Return(map[string]models.CarryByLocality{}, errors.New("database error"))

		service := NewCarryDefault(mockRepo)

		// Act
		result, err := service.SearchByLocality(localityId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, map[string]models.CarryByLocality{}, result)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si no se encuentran carries para la localidad devolverá un mapa vacío", func(t *testing.T) {
		// Arrange
		localityId := 2
		expectedCarries := map[string]models.CarryByLocality{}

		mockRepo := &carry.MockCarryRepository{}
		mockRepo.On("SearchByLocality", localityId).Return(expectedCarries, nil)

		service := NewCarryDefault(mockRepo)

		// Act
		result, err := service.SearchByLocality(localityId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCarries, result)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestCarryService_Create(t *testing.T) {
	t.Run("Si contiene los campos necesarios se creará el carry exitosamente", func(t *testing.T) {
		// Arrange
		carryInput := models.Carry{
			Cid:         "C001",
			CompanyName: "Express Delivery Co",
			Address:     "123 Main St",
			Telephone:   "555-1234",
			LocalityId:  1,
		}

		expectedCarry := models.Carry{
			ID:          1,
			Cid:         "C001",
			CompanyName: "Express Delivery Co",
			Address:     "123 Main St",
			Telephone:   "555-1234",
			LocalityId:  1,
		}

		mockRepo := &carry.MockCarryRepository{}
		mockRepo.On("Create", carryInput).Return(expectedCarry, nil)

		service := NewCarryDefault(mockRepo)

		// Act
		result, err := service.Create(carryInput)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCarry, result)
		assert.Equal(t, expectedCarry.ID, result.ID)
		assert.Equal(t, expectedCarry.Cid, result.Cid)
		assert.Equal(t, expectedCarry.CompanyName, result.CompanyName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si ocurre un error durante la creación se retorna error", func(t *testing.T) {
		// Arrange
		carryInput := models.Carry{
			Cid:         "C001",
			CompanyName: "Express Delivery Co",
			Address:     "123 Main St",
			Telephone:   "555-1234",
			LocalityId:  1,
		}

		mockRepo := &carry.MockCarryRepository{}
		mockRepo.On("Create", carryInput).Return(models.Carry{}, errors.New("database error"))

		service := NewCarryDefault(mockRepo)

		// Act
		result, err := service.Create(carryInput)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Carry{}, result)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si el CID ya existe se retorna error de duplicado", func(t *testing.T) {
		// Arrange
		carryInput := models.Carry{
			Cid:         "C001",
			CompanyName: "Express Delivery Co",
			Address:     "123 Main St",
			Telephone:   "555-1234",
			LocalityId:  1,
		}

		mockRepo := &carry.MockCarryRepository{}
		mockRepo.On("Create", carryInput).Return(models.Carry{}, errors.New("carry CID already exists"))

		service := NewCarryDefault(mockRepo)

		// Act
		result, err := service.Create(carryInput)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Carry{}, result)
		assert.Equal(t, "carry CID already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
