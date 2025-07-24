package service

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/buyer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateBuyer(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		input := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "12345678",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		expected := models.Buyer{
			ID: 1,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "12345678",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		mockRepo.On("Create", input).Return(expected, nil)
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_conflict", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		input := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "12345678",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		mockRepo.On("Create", input).Return(models.Buyer{}, pkg.ServiceErrors[pkg.ErrConflict])
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrConflict], err)
		require.Equal(t, models.Buyer{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestFindAllBuyers(t *testing.T) {
	t.Run("find_all", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		expected := map[int]models.Buyer{
			1: {
				ID: 1,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberID: "12345678",
					FirstName:    "John",
					LastName:     "Doe",
				},
			},
			2: {
				ID: 2,
				BuyerAttributes: models.BuyerAttributes{
					CardNumberID: "87654321",
					FirstName:    "Jane",
					LastName:     "Smith",
				},
			},
		}

		mockRepo.On("GetAll").Return(expected, nil)
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.GetAll()

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		require.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

func TestFindBuyerByID(t *testing.T) {
	t.Run("find_by_id_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 1
		expected := models.Buyer{
			ID: 1,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "12345678",
				FirstName:    "John",
				LastName:     "Doe",
			},
		}

		mockRepo.On("GetByID", buyerID).Return(expected, nil)
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.GetByID(buyerID)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 999

		mockRepo.On("GetByID", buyerID).Return(models.Buyer{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.GetByID(buyerID)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.Equal(t, models.Buyer{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateBuyer(t *testing.T) {
	t.Run("update_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 1
		input := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		inputWithID := models.Buyer{
			ID: buyerID,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		expected := models.Buyer{
			ID: buyerID,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		mockRepo.On("Update", inputWithID).Return(expected, nil)
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.Update(buyerID, input)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update_non_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 999
		input := models.Buyer{
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		inputWithID := models.Buyer{
			ID: buyerID,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: "87654321",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		mockRepo.On("Update", inputWithID).Return(models.Buyer{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.Update(buyerID, input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.Equal(t, models.Buyer{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteBuyer(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 1

		mockRepo.On("Delete", buyerID).Return(nil)
		service := NewBuyerDefault(mockRepo)

		// act
		err := service.Delete(buyerID)

		// assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete_non_existent", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 999

		mockRepo.On("Delete", buyerID).Return(pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewBuyerDefault(mockRepo)

		// act
		err := service.Delete(buyerID)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPurchaseOrdersReport(t *testing.T) {
	t.Run("get_all_buyers_report", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		expected := []models.BuyerPurchaseOrderReport{
			{
				ID:                  1,
				CardNumberID:        "12345678",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 5,
			},
			{
				ID:                  2,
				CardNumberID:        "87654321",
				FirstName:           "Jane",
				LastName:            "Smith",
				PurchaseOrdersCount: 3,
			},
		}

		mockRepo.On("GetPurchaseOrdersReport", (*int)(nil)).Return(expected, nil)
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.GetPurchaseOrdersReport(nil)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		require.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get_specific_buyer_report", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 1
		expected := []models.BuyerPurchaseOrderReport{
			{
				ID:                  1,
				CardNumberID:        "12345678",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 5,
			},
		}

		mockRepo.On("GetPurchaseOrdersReport", &buyerID).Return(expected, nil)
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.GetPurchaseOrdersReport(&buyerID)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		require.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get_report_buyer_not_found", func(t *testing.T) {
		// arrange
		mockRepo := new(buyer.MockBuyerRepository)
		buyerID := 999

		mockRepo.On("GetPurchaseOrdersReport", &buyerID).Return([]models.BuyerPurchaseOrderReport{}, pkg.ServiceErrors[pkg.ErrNotFound])
		service := NewBuyerDefault(mockRepo)

		// act
		result, err := service.GetPurchaseOrdersReport(&buyerID)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
		require.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}
