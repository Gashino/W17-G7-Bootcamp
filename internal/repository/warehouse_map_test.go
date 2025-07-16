package repository

/*
import (
	"app/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestRepo() *WarehouseMap {
	warehouseDb := map[int]models.Warehouse{
		1: {ID: 1, WarehouseCode: "WH001", Address: "Address 1"},
	}
	sectionDb := map[int]models.Section{}
	employeeDb := map[int]models.Employee{}

	return NewWarehouseMap(&sectionDb, &employeeDb, &warehouseDb)
}

func TestFindAll(t *testing.T) {
	repo := setupTestRepo()

	warehouses, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, warehouses, 1)
	assert.Equal(t, "WH001", warehouses[1].WarehouseCode)
}

func TestFindByID_Success(t *testing.T) {
	repo := setupTestRepo()

	warehouse, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, warehouse.ID)
	assert.Equal(t, "WH001", warehouse.WarehouseCode)
}

func TestFindByID_NotFound(t *testing.T) {
	repo := setupTestRepo()

	_, err := repo.FindByID(999)

	assert.Error(t, err)
	assert.Equal(t, "warehouse not found", err.Error())
}

func TestAdd(t *testing.T) {
	repo := setupTestRepo()

	newWarehouse := models.Warehouse{ID: 2, WarehouseCode: "WH002", Address: "Address 2"}
	err := repo.Add(newWarehouse)

	assert.NoError(t, err)
	warehouse, _ := repo.FindByID(2)
	assert.Equal(t, "WH002", warehouse.WarehouseCode)
}

func TestFindWarehouseByCode_Success(t *testing.T) {
	repo := setupTestRepo()

	warehouse, err := repo.FindWarehouseByCode("WH001")

	assert.NoError(t, err)
	assert.Equal(t, 1, warehouse.ID)
}

func TestFindWarehouseByCode_NotFound(t *testing.T) {
	repo := setupTestRepo()

	_, err := repo.FindWarehouseByCode("UNKNOWN")

	assert.Error(t, err)
	assert.Equal(t, "warehouse not found", err.Error())
}

func TestFindAvailableID(t *testing.T) {
	repo := setupTestRepo()

	id, err := repo.FindAvailableID()

	assert.NoError(t, err)
	assert.Equal(t, 2, id)
}

func TestDelete_Success(t *testing.T) {
	repo := setupTestRepo()

	err := repo.Delete(1)

	assert.NoError(t, err)
	_, err = repo.FindByID(1)
	assert.Error(t, err)
}
*/
