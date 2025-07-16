package repository

/*
import (
	"app/pkg"
	"app/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmployeeMapRepository(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
		2: {ID: 2, CardNumberID: "87654321", FirstName: "Jane", LastName: "Smith", WarehouseID: 2},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
		2: {ID: 2, WarehouseCode: "WH002", Address: "Address 2", Telephone: "987654321"},
	}

	// When
	repo := NewEmployeeMapRepository(&employees, &warehouses)

	// Then
	assert.NotNil(t, repo)
	mapRepo := repo.(*EmployeeRepositoryMap)
	assert.Equal(t, 2, len(*mapRepo.db))
	assert.Equal(t, 2, mapRepo.maxId)
}

func TestEmployeeRepositoryMap_FindAll(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
		2: {ID: 2, CardNumberID: "87654321", FirstName: "Jane", LastName: "Smith", WarehouseID: 2},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
		2: {ID: 2, WarehouseCode: "WH002", Address: "Address 2", Telephone: "987654321"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)

	// When
	result, err := repo.FindAll()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "John", result[1].FirstName)
	assert.Equal(t, "Jane", result[2].FirstName)
}

func TestEmployeeRepositoryMap_FindById_Success(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)

	// When
	result, err := repo.FindById(1)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "12345678", result.CardNumberID)
	assert.Equal(t, 1, result.WarehouseID)
}

func TestEmployeeRepositoryMap_FindById_NotFound(t *testing.T) {
	// Given
	employees := map[int]models.Employee{}
	warehouses := map[int]models.Warehouse{}
	repo := NewEmployeeMapRepository(&employees, &warehouses)

	// When
	result, err := repo.FindById(999)

	// Then
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	assert.Equal(t, models.Employee{}, result)
}

func TestEmployeeRepositoryMap_Save_Success(t *testing.T) {
	// Given
	employees := map[int]models.Employee{}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)
	newEmployee := models.Employee{
		CardNumberID: "12345678",
		FirstName:    "John",
		LastName:     "Doe",
		WarehouseID:  1,
	}

	// When
	result, err := repo.Save(newEmployee)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "12345678", result.CardNumberID)
	assert.Equal(t, 1, result.WarehouseID)

	// Verify it was actually saved
	saved, err := repo.FindById(1)
	assert.NoError(t, err)
	assert.Equal(t, result, saved)
}

func TestEmployeeRepositoryMap_Save_DuplicateCardID(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
		2: {ID: 2, WarehouseCode: "WH002", Address: "Address 2", Telephone: "987654321"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)
	newEmployee := models.Employee{
		CardNumberID: "12345678", // Same as existing
		FirstName:    "Jane",
		LastName:     "Smith",
		WarehouseID:  2,
	}

	// When
	result, err := repo.Save(newEmployee)

	// Then
	assert.Error(t, err)
	assert.Equal(t, models.Employee{}, result)
	serviceErr, ok := err.(pkg.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, 400, serviceErr.Code)
	assert.Equal(t, "Card ID already exists", serviceErr.Message)
}

func TestEmployeeRepositoryMap_Update_Success(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
		2: {ID: 2, WarehouseCode: "WH002", Address: "Address 2", Telephone: "987654321"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)
	updatedEmployee := models.Employee{
		ID:           1,
		CardNumberID: "87654321",
		FirstName:    "Jane",
		LastName:     "Smith",
		WarehouseID:  2,
	}

	// When
	result, err := repo.Update(updatedEmployee, 1)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, updatedEmployee, result)

	// Verify it was actually updated
	updated, err := repo.FindById(1)
	assert.NoError(t, err)
	assert.Equal(t, "Jane", updated.FirstName)
	assert.Equal(t, "Smith", updated.LastName)
	assert.Equal(t, "87654321", updated.CardNumberID)
}

func TestEmployeeRepositoryMap_Update_NotFound(t *testing.T) {
	// Given
	employees := map[int]models.Employee{}
	warehouses := map[int]models.Warehouse{
		2: {ID: 2, WarehouseCode: "WH002", Address: "Address 2", Telephone: "987654321"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)
	updatedEmployee := models.Employee{
		ID:           999,
		CardNumberID: "87654321",
		FirstName:    "Jane",
		LastName:     "Smith",
		WarehouseID:  2,
	}

	// When
	result, err := repo.Update(updatedEmployee, 999)

	// Then
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	assert.Equal(t, models.Employee{}, result)
}

func TestEmployeeRepositoryMap_Update_DuplicateCardID(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
		2: {ID: 2, CardNumberID: "87654321", FirstName: "Jane", LastName: "Smith", WarehouseID: 2},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
		2: {ID: 2, WarehouseCode: "WH002", Address: "Address 2", Telephone: "987654321"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)
	updatedEmployee := models.Employee{
		ID:           1,
		CardNumberID: "87654321", // Same as employee ID 2
		FirstName:    "John",
		LastName:     "Doe",
		WarehouseID:  1,
	}

	// When
	result, err := repo.Update(updatedEmployee, 1)

	// Then
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
	assert.Equal(t, models.Employee{}, result)
}

func TestEmployeeRepositoryMap_Delete_Success(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)

	// When
	err := repo.Delete(1)

	// Then
	assert.NoError(t, err)

	// Verify it was actually deleted
	_, err = repo.FindById(1)
	assert.Error(t, err)
	assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
}

func TestEmployeeRepositoryMap_Delete_NotFound(t *testing.T) {
	// Given
	employees := map[int]models.Employee{}
	warehouses := map[int]models.Warehouse{}
	repo := NewEmployeeMapRepository(&employees, &warehouses)

	// When
	err := repo.Delete(999)

	// Then
	assert.NoError(t, err) // Delete should not error even if ID doesn't exist
}

func TestEmployeeRepositoryMap_Save_InvalidWarehouseID(t *testing.T) {
	// Given
	employees := map[int]models.Employee{}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)
	newEmployee := models.Employee{
		CardNumberID: "12345678",
		FirstName:    "John",
		LastName:     "Doe",
		WarehouseID:  999, // Invalid warehouse ID
	}

	// When
	result, err := repo.Save(newEmployee)

	// Then
	assert.Error(t, err)
	assert.Equal(t, models.Employee{}, result)
	serviceErr, ok := err.(pkg.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, 400, serviceErr.Code)
	assert.Equal(t, "Warehouse ID does not exist", serviceErr.Message)
}

func TestEmployeeRepositoryMap_Update_InvalidWarehouseID(t *testing.T) {
	// Given
	employees := map[int]models.Employee{
		1: {ID: 1, CardNumberID: "12345678", FirstName: "John", LastName: "Doe", WarehouseID: 1},
	}
	warehouses := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1", Telephone: "123456789"},
	}
	repo := NewEmployeeMapRepository(&employees, &warehouses)
	updatedEmployee := models.Employee{
		ID:           1,
		CardNumberID: "87654321",
		FirstName:    "Jane",
		LastName:     "Smith",
		WarehouseID:  999, // Invalid warehouse ID
	}

	// When
	result, err := repo.Update(updatedEmployee, 1)

	// Then
	assert.Error(t, err)
	assert.Equal(t, models.Employee{}, result)
	serviceErr, ok := err.(pkg.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, 400, serviceErr.Code)
	assert.Equal(t, "Warehouse ID does not exist", serviceErr.Message)
}
*/
