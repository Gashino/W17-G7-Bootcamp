package service

import (
	"app/pkg/models"
	"app/test/inbound_order"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateInboundOrder(t *testing.T) {
	t.Run("create_ok", func(t *testing.T) {
		// arrange
		mockRepo := new(inbound_order.MockInboundOrderRepository)
		orderDate := time.Now()
		input := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORDER001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		expected := models.InboundOrder{
			ID:             1,
			OrderDate:      orderDate,
			OrderNumber:    "ORDER001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		mockRepo.On("Create", input).Return(expected, nil)
		service := NewInboundOrderServiceDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("create_fail_invalid_order_date", func(t *testing.T) {
		// arrange
		mockRepo := new(inbound_order.MockInboundOrderRepository)
		input := models.InboundOrder{
			OrderDate:      time.Time{}, // Zero time
			OrderNumber:    "ORDER001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		service := NewInboundOrderServiceDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, "order_date is required", err.Error())
		require.Equal(t, models.InboundOrder{}, result)
	})

	t.Run("create_fail_future_date", func(t *testing.T) {
		// arrange
		mockRepo := new(inbound_order.MockInboundOrderRepository)
		input := models.InboundOrder{
			OrderDate:      time.Now().Add(24 * time.Hour), // Future date
			OrderNumber:    "ORDER001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		service := NewInboundOrderServiceDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, "order_date cannot be in the future", err.Error())
		require.Equal(t, models.InboundOrder{}, result)
	})

	t.Run("create_fail_empty_order_number", func(t *testing.T) {
		// arrange
		mockRepo := new(inbound_order.MockInboundOrderRepository)
		input := models.InboundOrder{
			OrderDate:      time.Now(),
			OrderNumber:    "",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		service := NewInboundOrderServiceDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, "order_number is required", err.Error())
		require.Equal(t, models.InboundOrder{}, result)
	})

	t.Run("create_fail_invalid_employee_id", func(t *testing.T) {
		// arrange
		mockRepo := new(inbound_order.MockInboundOrderRepository)
		input := models.InboundOrder{
			OrderDate:      time.Now(),
			OrderNumber:    "ORDER001",
			EmployeeID:     0,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		service := NewInboundOrderServiceDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, "employee_id must be a valid positive number", err.Error())
		require.Equal(t, models.InboundOrder{}, result)
	})

	t.Run("create_fail_invalid_product_batch_id", func(t *testing.T) {
		// arrange
		mockRepo := new(inbound_order.MockInboundOrderRepository)
		input := models.InboundOrder{
			OrderDate:      time.Now(),
			OrderNumber:    "ORDER001",
			EmployeeID:     1,
			ProductBatchID: 0,
			WarehouseID:    1,
		}

		service := NewInboundOrderServiceDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, "product_batch_id must be a valid positive number", err.Error())
		require.Equal(t, models.InboundOrder{}, result)
	})

	t.Run("create_fail_invalid_warehouse_id", func(t *testing.T) {
		// arrange
		mockRepo := new(inbound_order.MockInboundOrderRepository)
		input := models.InboundOrder{
			OrderDate:      time.Now(),
			OrderNumber:    "ORDER001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    0,
		}

		service := NewInboundOrderServiceDefault(mockRepo)

		// act
		result, err := service.Create(input)

		// assert
		require.Error(t, err)
		require.Equal(t, "warehouse_id must be a valid positive number", err.Error())
		require.Equal(t, models.InboundOrder{}, result)
	})
}
