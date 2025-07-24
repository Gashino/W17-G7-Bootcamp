package service

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/employee"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEmployee(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 21

		input := models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		expected := models.Employee{
			ID:           &id,
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		mockRepo.On("Save", input).Return(expected, nil)
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		result, err := service.Save(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("create_conflict", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3

		input := models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		mockRepo.On("Save", input).Return(models.Employee{}, pkg.ServiceErrors[pkg.ErrConflict])
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		result, err := service.Save(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
		require.Equal(t, models.Employee{}, result)
	})
}

func TestFindEmployee(t *testing.T) {
	t.Run("find_all", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 21

		expected := map[int]models.Employee{
			1: {
				ID:           &id,
				CardNumberID: &cardNumberID,
				FirstName:    &firstName,
				LastName:     &lastName,
				WarehouseID:  &warehouseID,
			},
		}

		mockRepo.On("FindAll").Return(expected, nil)
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		result, err := service.FindAll()

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		require.Len(t, result, 1)
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		mockRepo.On("FindById", 3).Return(models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		result, err := service.FindById(3)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.Equal(t, models.Employee{}, result)
	})

	t.Run("find_by_id_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 3

		expected := models.Employee{
			ID:           &id,
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		mockRepo.On("FindById", 3).Return(expected, nil)
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		result, err := service.FindById(3)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("update_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3
		id := 3

		input := models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		expected := models.Employee{
			ID:           &id,
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		mockRepo.On("Update", input, 3).Return(expected, nil)
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		result, err := service.Update(input, 3)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("update_non_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		cardNumberID := "11223342"
		firstName := "Carlos"
		lastName := "López"
		warehouseID := 3

		input := models.Employee{
			CardNumberID: &cardNumberID,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		mockRepo.On("Update", input, 3).Return(models.Employee{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		result, err := service.Update(input, 3)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.Equal(t, models.Employee{}, result)
	})
}

func TestDeleteEmployee(t *testing.T) {
	t.Run("delete_non_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		mockRepo.On("Delete", 3).Return(pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		err := service.Delete(3)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})

	t.Run("delete_ok", func(t *testing.T) {
		// arrange
		mockRepo := new(employee.MockEmployeeRepository)
		mockRepo.On("Delete", 3).Return(nil)
		service := NewEmployeeServiceDefault(mockRepo)

		// act
		err := service.Delete(3)

		// assert
		require.NoError(t, err)
	})
}
