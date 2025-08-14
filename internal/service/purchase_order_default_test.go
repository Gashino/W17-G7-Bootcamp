package service

import (
	"app/pkg"
	"app/pkg/models"
	"app/test/purchase_order"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePurchaseOrder(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// arrange
		mockRepo := new(purchase_order.MockPurchaseOrderRepository)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#2",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		expected := models.PurchaseOrder{
			ID: 2,
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#2",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		mockRepo.On("Create", input).Return(expected, nil)
		service := NewPurchaseOrderDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_conflict", func(t *testing.T) {
		// arrange
		mockRepo := new(purchase_order.MockPurchaseOrderRepository)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#1",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		mockRepo.On("Create", input).Return(models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrBadRequest])
		service := NewPurchaseOrderDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrBadRequest], err)
		require.Equal(t, models.PurchaseOrder{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create_repository_error", func(t *testing.T) {
		// arrange
		mockRepo := new(purchase_order.MockPurchaseOrderRepository)
		input := models.PurchaseOrder{
			PurchaseOrderAttributes: models.PurchaseOrderAttributes{
				OrderNumber:     "order#3",
				OrderDate:       "2021-04-05",
				TrackingCode:    "xyz789",
				BuyerID:         1,
				ProductRecordID: 2,
			},
		}

		mockRepo.On("Create", input).Return(models.PurchaseOrder{}, pkg.ServiceErrors[pkg.ErrInternalServer])
		service := NewPurchaseOrderDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, pkg.ServiceErrors[pkg.ErrInternalServer], err)
		require.Equal(t, models.PurchaseOrder{}, result)
		mockRepo.AssertExpectations(t)
	})
}
