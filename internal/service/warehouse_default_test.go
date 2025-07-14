package service

/*
import (
	"app/pkg"
	"app/pkg/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeRepository es un mock del repositorio para testing
type MockWarehouseRepository struct {
	mock.Mock
}

func (m *MockWarehouseRepository) FindAll() (map[int]models.Warehouse, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) FindByID(id int) (models.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) Add(warehouse models.Warehouse) error {
	args := m.Called(warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWarehouseRepository) FindAvailableID() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *MockWarehouseRepository) FindWarehouseByCode(code string) (models.Warehouse, error) {
	args := m.Called(code)
	return args.Get(0).(models.Warehouse), args.Error(1)
}
func TestNewWarehouseServiceDefault(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}

	// When
	service := NewWarehouseDefault(mockRepo)

	// Then
	assert.NotNil(t, service)
	assert.IsType(t, &WarehouseDefault{}, service)
}

func TestWarehouseServiceDefault_FindAll_Success(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	expectedEmployees := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "7 Calle Principal, Ciudad 1, Pa\u00eds", Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4},
		2: {ID: 2, WarehouseCode: "WH002", Address: "14 Calle Principal, Ciudad 2, Pa\u00eds", Telephone: "+1234567002", MinCapacity: 4930, MinTemperature: -24},
	}

	mockRepo.On("FindAll").Return(expectedEmployees, nil)

	// When
	result, err := service.FindAll()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedEmployees, result)
	mockRepo.AssertExpectations(t)
}

func TestWarehouseServiceDefault_FindAll_Error(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	expectedError := errors.New("database error")
	mockRepo.On("FindAll").Return(map[int]models.Warehouse{}, expectedError)

	// When
	result, err := service.FindAll()

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestWarehouseServiceDefault_FindById_Success(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	expectedWarehouse := models.Warehouse{
		ID: 1, WarehouseCode: "WH001", Address: "7 Calle Principal, Ciudad 1, Pa\u00eds", Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4,
	}

	mockRepo.On("FindByID", 1).Return(expectedWarehouse, nil)

	// When
	result, err := service.FindByID(1)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedWarehouse, result)
	mockRepo.AssertExpectations(t)
}

func TestWarehouseServiceDefault_FindById_NotFound(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	expectedError := pkg.ServiceErrors[pkg.ErrNotFound]
	mockRepo.On("FindByID", 999).Return(models.Warehouse{}, expectedError)

	// When
	result, err := service.FindByID(999)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, models.Warehouse{}, result)
	mockRepo.AssertExpectations(t)
}

func Test_WarehouseServiceDefault_Add_Success(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	inputWarehouse := models.WarehouseDoc{
		ID:             501,
		WarehouseCode:  "WH00501",
		Address:        "7 Calle Principal, Ciudad 1, País",
		Telephone:      "+1234567001",
		MinCapacity:    3694,
		MinTemperature: -4,
	}

	mockRepo.On("FindWarehouseByCode", "WH00501").Return(models.Warehouse{}, errors.New("not found"))

	mockRepo.On("FindAvailableID").Return(600, nil)

	warehouse := models.Warehouse{
		ID:             600,
		WarehouseCode:  inputWarehouse.WarehouseCode,
		Address:        inputWarehouse.Address,
		Telephone:      inputWarehouse.Telephone,
		MinCapacity:    inputWarehouse.MinCapacity,
		MinTemperature: inputWarehouse.MinTemperature,
	}
	mockRepo.On("Add", warehouse).Return(nil)

	// When
	_, err := service.Add(inputWarehouse)

	// Then
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func Test_WarehouseServiceDefault_Add_Error(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	inputWarehouse := models.WarehouseDoc{
		ID:             501,
		WarehouseCode:  "WH00501",
		Address:        "7 Calle Principal, Ciudad 1, País",
		Telephone:      "+1234567001",
		MinCapacity:    3694,
		MinTemperature: -4,
	}

	// Simulamos que el código ya existe retornando un Warehouse y nil error
	existingWarehouse := models.Warehouse{
		ID:            999,
		WarehouseCode: "WH00501",
	}
	mockRepo.On("FindWarehouseByCode", "WH00501").Return(existingWarehouse, nil)
	// When
	_, err := service.Add(inputWarehouse)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "Warehouse Code is not unique", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestWarehouseServiceDefault_Update_Success_AllFields(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	existingWarehouse := models.Warehouse{
		ID: 1, WarehouseCode: "WH001", Address: "7 Calle Principal, Ciudad 1, País",
		Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4,
	}
	updatedWarehouse := models.Warehouse{
		ID: 1, WarehouseCode: "WH001", Address: "7 Calle Principal, Ciudad 2, País",
		Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4,
	}
	warehouseDocInput := models.WarehouseDoc{
		ID: 1, WarehouseCode: "WH001", Address: "7 Calle Principal, Ciudad 2, País",
		Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4,
	}

	mockRepo.On("FindByID", 1).Return(existingWarehouse, nil)
	mockRepo.On("Add", updatedWarehouse).Return(nil)

	// When
	_, err := service.Update(1, warehouseDocInput)

	// Then
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWarehouseServiceDefault_Update_Error_WarehouseCodeNotUnique(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	existingWarehouse := models.Warehouse{
		ID: 1, WarehouseCode: "WH001", Address: "7 Calle Principal, Ciudad 1, País",
		Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4,
	}
	conflictingWarehouse := models.Warehouse{
		ID: 2, WarehouseCode: "WH002", Address: "Another Address",
		Telephone: "+9876543210", MinCapacity: 5000, MinTemperature: -10,
	}
	updateInput := models.WarehouseDoc{
		ID: 1, WarehouseCode: "WH002", Address: "7 Calle Principal, Ciudad 1, País",
		Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4,
	}

	mockRepo.On("FindByID", 1).Return(existingWarehouse, nil)
	mockRepo.On("FindWarehouseByCode", "WH002").Return(conflictingWarehouse, nil)

	// When
	_, err := service.Update(1, updateInput)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "Warehouse Code is not unique", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestWarehouseServiceDefault_Delete_Success(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	mockRepo.On("FindByID", 1).Return(models.Warehouse{
		ID: 1, WarehouseCode: "WH001", Address: "7 Calle Principal, Ciudad 1, País",
		Telephone: "+1234567001", MinCapacity: 3694, MinTemperature: -4,
	}, nil)

	mockRepo.On("Delete", 1).Return(nil)

	// When
	err := service.Delete(1)

	// Then
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWarehouseServiceDefault_Delete_Error_NotFound(t *testing.T) {
	// Given
	mockRepo := &MockWarehouseRepository{}
	service := NewWarehouseDefault(mockRepo)

	mockRepo.On("FindByID", 999).Return(models.Warehouse{}, errors.New("warehouse not found"))

	// When
	err := service.Delete(999)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "warehouse not found", err.Error())
	mockRepo.AssertExpectations(t)
}
*/
