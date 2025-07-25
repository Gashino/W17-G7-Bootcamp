package service

import (
	"app/pkg/models"
	"app/test/warehouse"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
FindAll() (v map[int]models.Warehouse, err error)
FindByID(id int) (v models.Warehouse, err error)
Add(v models.WarehouseDoc) (w models.Warehouse, err error)
Update(id int, v models.WarehouseDoc) (w models.Warehouse, err error)
Delete(id int) (err error)
*/
func TestWarehouseService_FindAll(t *testing.T) {
	t.Run("Si la lista posee n elementos devolverá un cantidad de los elementos totales", func(t *testing.T) {
		// Arrange
		expectedWarehouses := map[int]models.Warehouse{
			1: {
				ID:             1,
				WarehouseCode:  "WH001",
				Address:        "123 Main St",
				Telephone:      "555-1234",
				MinCapacity:    100,
				MinTemperature: -10,
			},
			2: {
				ID:             2,
				WarehouseCode:  "WH002",
				Address:        "456 Oak Ave",
				Telephone:      "555-5678",
				MinCapacity:    200,
				MinTemperature: -20,
			},
		}

		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindAll").Return(expectedWarehouses, nil)

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.FindAll()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouses, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindAll").Return(nil, errors.New("database error"))

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.FindAll()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestWarehouseService_FindByID(t *testing.T) {
	t.Run("Si el elemento buscado por id existe devolverá la información del elemento solicitado", func(t *testing.T) {
		// Arrange
		expectedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-1234",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindByID", 1).Return(expectedWarehouse, nil)

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.FindByID(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si el elemento buscado por id no existe retorna nulo", func(t *testing.T) {
		// Arrange
		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindByID", 999).Return(models.Warehouse{}, errors.New("warehouse not found"))

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.FindByID(999)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Warehouse{}, result)
		assert.Equal(t, "warehouse not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestWarehouseService_Add(t *testing.T) {
	t.Run("Si contiene los campos necesarios se creará", func(t *testing.T) {
		// Arrange
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-1234",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		expectedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-1234",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		mockRepo := &warehouse.MockWarehouseRepository{}
		// El código no existe, por lo que devuelve error
		mockRepo.On("FindWarehouseByCode", "WH001").Return(models.Warehouse{}, errors.New("warehouse not found"))

		expectedWarehouseToAdd := models.Warehouse{
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-1234",
			MinCapacity:    100,
			MinTemperature: -10,
		}
		mockRepo.On("Add", expectedWarehouseToAdd).Return(expectedWarehouse, nil)

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.Add(warehouseDoc)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si el warehouse_code ya existe no podrá ser creado", func(t *testing.T) {
		// Arrange
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-1234",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		existingWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001",
			Address:        "456 Oak Ave",
			Telephone:      "555-5678",
			MinCapacity:    200,
			MinTemperature: -20,
		}

		mockRepo := &warehouse.MockWarehouseRepository{}
		// El código ya existe
		mockRepo.On("FindWarehouseByCode", "WH001").Return(existingWarehouse, nil)

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.Add(warehouseDoc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Warehouse{}, result)
		assert.Equal(t, "Warehouse Code is not unique", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestWarehouseService_Update(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá el warehouse con la información actualizada", func(t *testing.T) {
		// Arrange
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "WH001-UPDATED",
			Address:        "123 Updated St",
			Telephone:      "555-9999",
			MinCapacity:    150,
			MinTemperature: -15,
		}

		existingWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-1234",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		updatedWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001-UPDATED",
			Address:        "123 Updated St",
			Telephone:      "555-9999",
			MinCapacity:    150,
			MinTemperature: -15,
		}

		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindByID", 1).Return(existingWarehouse, nil)
		mockRepo.On("Update", 1, updatedWarehouse).Return(updatedWarehouse, nil)

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.Update(1, warehouseDoc)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updatedWarehouse, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si el warehouse que se desea actualizar no existe se retorna null.", func(t *testing.T) {
		// Arrange
		warehouseDoc := models.WarehouseDoc{
			WarehouseCode:  "WH001-UPDATED",
			Address:        "123 Updated St",
			Telephone:      "555-9999",
			MinCapacity:    150,
			MinTemperature: -15,
		}

		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindByID", 999).Return(models.Warehouse{}, errors.New("warehouse not found"))

		service := NewWarehouseDefault(mockRepo)

		// Act
		result, err := service.Update(999, warehouseDoc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Warehouse{}, result)
		assert.Equal(t, "warehouse not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestWarehouseService_Delete(t *testing.T) {
	t.Run("Cuando la eliminación de datos sea exitosa no se devolverá error", func(t *testing.T) {
		// Arrange
		existingWarehouse := models.Warehouse{
			ID:             1,
			WarehouseCode:  "WH001",
			Address:        "123 Main St",
			Telephone:      "555-1234",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindByID", 1).Return(existingWarehouse, nil)
		mockRepo.On("Delete", 1).Return(nil)

		service := NewWarehouseDefault(mockRepo)

		// Act
		err := service.Delete(1)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Si el warehouse que se desea eliminar no existe se retorna error.", func(t *testing.T) {
		// Arrange
		mockRepo := &warehouse.MockWarehouseRepository{}
		mockRepo.On("FindByID", 999).Return(models.Warehouse{}, errors.New("warehouse not found"))

		service := NewWarehouseDefault(mockRepo)

		// Act
		err := service.Delete(999)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "warehouse not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
